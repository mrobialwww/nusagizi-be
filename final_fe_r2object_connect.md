# Implementasi: Upload Gambar ke Cloudflare R2 (Flutter)

**Endpoint yang dicakup (sesi ini):**
- `PATCH /users` — update profil user (termasuk foto)
- `POST /children` — tambah profil anak + foto
- `PATCH /children/{child_id}` — update profil anak + foto

> **[REVISI]** Versi sebelumnya dokumen ini memakai dua pola berbeda: foto user lewat `POST /images/confirm` (backend yang menyimpan ke DB), sedangkan foto anak lewat pengiriman `photo_url` langsung di body. Setelah disamakan dengan `be_r2object_connect` (yang menambahkan field `PhotoURL` ke `UpdateUserInput` untuk `PATCH /users`) dan `endpoint.md` (endpoint 1 sudah direvisi untuk menerima `photo_url`), **ketiga endpoint sekarang memakai satu pola yang sama** — lihat di bawah.

---

## Klarifikasi Teknis

> [!NOTE]
> **Flutter-lah yang melakukan PUT langsung ke R2.** Backend Go hanya berperan sebagai penjaga gerbang: generate presigned URL (`POST /images/presign-upload`). Proses transfer binary file terjadi dari Flutter → R2 menggunakan `Dio.put()` ke presigned URL. Tidak ada SDK cloudflare/aws yang dibutuhkan di Flutter — cukup `Dio` (sudah ada) dan `mime` (package baru) untuk deteksi Content-Type.

---

## Pola Seragam untuk Ketiga Endpoint

Tidak ada lagi pembedaan flow antar endpoint — semuanya 2 tahap yang sama:

| Endpoint | Ada `photo_url` di request body? | Cara simpan foto |
|----------|----------------------------------|------------------|
| `PATCH /users` | **YA** (`photo_url` opsional, field baru) | Kirim `objectKey` sebagai `photo_url` di body |
| `POST /children` | **YA** (`photo_url` opsional) | Kirim `objectKey` sebagai `photo_url` di body |
| `PATCH /children/{child_id}` | **YA** (sama seperti POST) | Kirim `objectKey` sebagai `photo_url` di body |

**Flow yang sama untuk ketiganya:** `presign` → `PUT ke R2` → `objectKey` disertakan langsung sebagai `photo_url` di body endpoint domain terkait (`PATCH /users` / `POST /children` / `PATCH /children/{child_id}`).

`POST /images/confirm` **tidak dipakai** dalam flow utama sama sekali — sesuai `be_r2object_connect`, endpoint itu sifatnya opsional (hanya validasi keberadaan object di R2, tidak menyimpan apa pun ke DB). Kalau suatu saat ingin dipakai sebagai validasi sebelum submit, itu perubahan terpisah, bukan bagian dari flow simpan data.

---

## Proposed Changes

### 0. Package Dependencies

#### [MODIFY] [pubspec.yaml](file:///f:/Project/ProjectLomba/nuzagizi/pubspec.yaml)

```diff
  image_picker: ^1.2.2
+ mime: ^2.0.0
```

---

### 1. Fitur Baru: `lib/features/image_upload/`

#### [NEW] `domain/entities/presign_upload_entity.dart`

```dart
class PresignUploadEntity {
  final String uploadUrl;  // URL presigned ke R2 untuk di-PUT
  final String objectKey;  // Key objek di R2 bucket
  const PresignUploadEntity({required this.uploadUrl, required this.objectKey});
}
```

#### [NEW] `domain/repositories/image_upload_repository.dart`

```dart
typedef UploadProgressCallback = void Function(int sent, int total);

abstract class ImageUploadRepository {
  // Tahap 1: minta presigned URL ke backend Go
  Future<Either<Failure, PresignUploadEntity>> requestPresignUrl({
    required String category,     // 'profile' | 'child-profile'
    required String contentType,  // misal 'image/jpeg'
    String? ownerId,
  });

  // Tahap 2: Flutter PUT raw bytes langsung ke R2
  // objectKey hasil tahap 1 langsung dipakai caller sebagai photo_url
  // di body endpoint domain (PATCH /users, POST/PATCH /children) — tidak ada tahap konfirmasi terpisah.
  Future<Either<Failure, void>> putFileToStorage({
    required String uploadUrl,
    required File file,
    required String contentType,
    UploadProgressCallback? onProgress,
  });
}
```

#### [NEW] `domain/usecases/upload_image_usecase.dart`

Satu alur seragam (2 tahap) untuk semua kategori foto — tidak ada percabangan mode lagi:

```dart
class UploadImageParams {
  final File file;
  final String category;      // 'profile' | 'child-profile'
  final String contentType;
  final String? ownerId;
  final UploadProgressCallback? onProgress;
}

class UploadImageUseCase implements UseCase<String, UploadImageParams> {
  // Return Either<Failure, String> — String adalah objectKey,
  // caller (Cubit halaman terkait) yang menyisipkannya ke field `photo_url`
  // saat memanggil PATCH /users / POST /children / PATCH /children/{id}.
  Future<Either<Failure, String>> call(UploadImageParams params) async {
    // Tahap 1: presign
    final presignResult = await repository.requestPresignUrl(
      category: params.category,
      contentType: params.contentType,
      ownerId: params.ownerId,
    );
    return presignResult.fold(
      (f) async => Left(f),
      (presign) async {
        // Tahap 2: PUT ke R2
        final putResult = await repository.putFileToStorage(
          uploadUrl: presign.uploadUrl,
          file: params.file,
          contentType: params.contentType,
          onProgress: params.onProgress,
        );
        return putResult.fold(
          (f) => Left(f),
          (_) => Right(presign.objectKey),
        );
      },
    );
  }
}
```

#### [NEW] `data/models/presign_upload_model.dart`

```dart
class PresignUploadModel extends PresignUploadEntity {
  const PresignUploadModel({required super.uploadUrl, required super.objectKey});

  factory PresignUploadModel.fromJson(Map<String, dynamic> json) =>
      PresignUploadModel(
        uploadUrl: json['upload_url'] as String,
        objectKey: json['object_key'] as String,
      );
}
```

#### [NEW] `data/datasources/image_api_service.dart`

Menggunakan Dio **ber-auth** — hanya call ke backend Go:

```dart
abstract class ImageApiService {
  Future<PresignUploadModel> presignUpload({
    required String category, required String contentType, String? ownerId,
  });
}

class ImageApiServiceImpl implements ImageApiService {
  final Dio dio; // Dio utama dengan baseUrl + Authorization interceptor

  // POST /images/presign-upload
  // Request body: { category, content_type, owner_id? }
  // Response: { upload_url, object_key }
}
```

> **Catatan**: `POST /images/confirm` **tidak dipakai** oleh `ImageApiService` di flow utama ini, karena per `be_r2object_connect` endpoint tersebut hanya validasi opsional (HEAD-check ke R2), bukan penyimpan data. `objectKey` hasil presign langsung dipakai caller sebagai `photo_url` saat memanggil endpoint domain (`PATCH /users` / `POST/PATCH /children`).

#### [NEW] `data/datasources/image_storage_service.dart`

Menggunakan Dio **POLOS** — PUT langsung ke domain R2, bukan backend kita:

```dart
class ImageStorageServiceImpl implements ImageStorageService {
  final Dio storageDio; // Dio TANPA baseUrl, TANPA interceptor Authorization

  // storageDio.put(
  //   uploadUrl,                               // full URL ke R2 (dari presign)
  //   data: file.openRead(),                   // stream bytes, hemat memory
  //   options: Options(headers: {
  //     Headers.contentLengthHeader: fileLen,  // wajib agar onSendProgress akurat
  //     'Content-Type': contentType,           // HARUS identik dgn saat presign
  //   }),
  //   onSendProgress: onProgress,
  // )
  // HTTP 403 → Content-Type mismatch atau URL expired
}
```

#### [NEW] `data/repositories/image_upload_repository_impl.dart`

Gabungkan `ImageApiService` + `ImageStorageService`. Catch `ServerException` → `Left(ServerFailure)`.

#### [NEW] `presentation/cubit/image_upload_state.dart`

```dart
class ImageUploadInitial extends ImageUploadState {}
class ImageUploadInProgress extends ImageUploadState {
  final double progress; // 0.0–1.0
}
class ImageUploadSuccess extends ImageUploadState {
  final String objectKey; // objectKey tersedia di state Success agar UI bisa pakai
}
class ImageUploadError extends ImageUploadState {
  final String message;
}
```

#### [NEW] `presentation/cubit/image_upload_cubit.dart`

```dart
class ImageUploadCubit extends Cubit<ImageUploadState> {
  Future<void> upload({
    required File file,
    required String category,   // 'profile' | 'child-profile'
    String? ownerId,             // user_id / child_id — kosongkan jika resource belum ada (mis. create child, lihat endpoint.md endpoint 7)
  }) async {
    emit(ImageUploadInProgress(progress: 0));
    // Deteksi contentType SEKALI → dipakai konsisten di presign & PUT
    final contentType = lookupMimeType(file.path) ?? 'application/octet-stream';

    final result = await uploadImageUseCase(UploadImageParams(
      file: file, category: category, contentType: contentType,
      ownerId: ownerId,
      onProgress: (sent, total) {
        if (total > 0) emit(ImageUploadInProgress(progress: sent / total));
      },
    ));

    result.fold(
      (failure) => emit(ImageUploadError(message: failure.message)),
      // objectKey diteruskan ke caller — caller yang menyisipkannya sebagai
      // `photo_url` saat memanggil PATCH /users / POST/PATCH /children.
      (objectKey) => emit(ImageUploadSuccess(objectKey: objectKey)),
    );
  }
}
```

---

### 2. Dependency Injection

#### [MODIFY] [service_locator.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/core/di/service_locator.dart)

Tambah imports image_upload feature, lalu di dalam [initDependencies()](file:///f:/Project/ProjectLomba/nuzagizi/lib/core/di/service_locator.dart#179-711) setelah blok Dio utama:

```dart
// Dio polos — PUT langsung ke R2 (TANPA baseUrl / interceptor auth)
sl.registerLazySingleton<Dio>(
  () => Dio(),
  instanceName: 'storageDio',
);

// Image Upload
sl.registerLazySingleton<ImageApiService>(
  () => ImageApiServiceImpl(dio: sl<Dio>()),  // Dio ber-auth → backend Go
);
sl.registerLazySingleton<ImageStorageService>(
  () => ImageStorageServiceImpl(
    storageDio: sl<Dio>(instanceName: 'storageDio'),
  ),
);
sl.registerLazySingleton<ImageUploadRepository>(
  () => ImageUploadRepositoryImpl(
    apiService: sl<ImageApiService>(),
    storageService: sl<ImageStorageService>(),
  ),
);
sl.registerLazySingleton<UploadImageUseCase>(
  () => UploadImageUseCase(repository: sl<ImageUploadRepository>()),
);
sl.registerFactory<ImageUploadCubit>(
  () => ImageUploadCubit(uploadImageUseCase: sl<UploadImageUseCase>()),
);
```

---

### 3. Model Updates

#### [MODIFY] [update_user_profile_request_model.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/data/models/update_user_profile_request_model.dart)

**[REVISI] Perlu tambah `photoUrl`** — `endpoint.md` sudah direvisi, `PATCH /users` sekarang menerima field `photo_url` opsional (diisi `object_key`), konsisten dengan `be_r2object_connect` yang menambahkan `PhotoURL *string` ke `UpdateUserInput`.

```diff
 class UpdateUserProfileRequestModel {
   final String? fullName;
   final String? email;
   final String? phoneNumber;
   final String? gender;
+  final String? photoUrl;

   Map<String, dynamic> toJson() => {
     if (fullName != null) 'full_name': fullName,
     if (email != null) 'email': email,
     if (phoneNumber != null) 'phone_number': phoneNumber,
     if (gender != null) 'gender': gender,
+    if (photoUrl != null) 'photo_url': photoUrl,
   };
 }
```

#### [MODIFY] [update_caregiver_profile_request_model.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/caregiver/profile/data/models/update_caregiver_profile_request_model.dart)

Sama seperti di atas — tambah `photoUrl`. Foto caregiver juga disimpan di `users.photo_url` (caregiver tidak punya endpoint PATCH sendiri; `caregiver_profiles` hanya join ke `users`, lihat endpoint 55 di `endpoint.md`), jadi halaman caregiver tetap memanggil `PATCH /users` di baliknya.

[ChildProfileRequestModel](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/data/models/child_profile_request_model.dart#1-41) sudah punya `final String? photoUrl` dan entry di [toJson()](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/data/models/child_profile_request_model.dart#26-40) — **tidak perlu diubah**.

---

### 4. UI: Edit Profil Mother (`PATCH /users`)

**Flow foto user (disamakan dengan flow foto anak):** pick → upload (presign + PUT, tanpa confirm) → `objectKey` disertakan langsung sebagai `photo_url` di body `PATCH /users`.

#### [MODIFY] [edit_profile_page.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/edit_profile_page.dart)

Widget ini sudah `StatefulWidget`. Perubahan:

**Tambah state & method:**
```dart
File? _pickedFile;

Future<void> _pickImage() async {
  final xfile = await ImagePicker().pickImage(
    source: ImageSource.gallery, imageQuality: 80,
  );
  if (xfile != null) setState(() => _pickedFile = File(xfile.path));
}
```

**Widget foto** — ganti dummy `CircleAvatar`:
```dart
GestureDetector(
  onTap: _pickImage,
  child: CircleAvatar(
    radius: 50.r,
    backgroundImage: _pickedFile != null
        ? FileImage(_pickedFile!) as ImageProvider
        : (widget.initialProfile.photoUrl != null
            ? NetworkImage(widget.initialProfile.photoUrl!)
            : null),
    child: (_pickedFile == null && widget.initialProfile.photoUrl == null)
        ? Icon(Icons.person, size: 40.w, color: Colors.grey)
        : null,
  ),
)
```

**[_saveProfile()](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/edit_profile_page.dart#46-55) → async, upload dulu jika ada foto baru, lalu sisipkan objectKey ke `photo_url`:**
```dart
Future<void> _saveProfile() async {
  // Upload foto terlebih dahulu jika ada file baru (presign + PUT saja, tanpa confirm)
  String? objectKey;
  if (_pickedFile != null) {
    final uploadCubit = context.read<ImageUploadCubit>();
    await uploadCubit.upload(
      file: _pickedFile!,
      category: 'profile',
      ownerId: widget.initialProfile.id, // user_id — sudah ada karena ini edit, bukan create
    );
    if (!mounted) return;
    final st = uploadCubit.state;
    if (st is ImageUploadError) return; // error sudah di SnackBar
    if (st is ImageUploadSuccess) objectKey = st.objectKey;
  }

  // PATCH /users — photo_url diisi objectKey baru (jika ganti foto) atau tetap null (jika tidak)
  final request = UpdateUserProfileRequestModel(
    fullName: _fullNameController.text,
    email: _emailController.text,
    phoneNumber: _phoneController.text,
    gender: _selectedGender,
    photoUrl: objectKey, // null jika user tidak ganti foto → field ini tidak dikirim (lihat toJson di atas)
  );
  if (mounted) context.read<UserProfileCubit>().updateProfile(request);
}
```

**[build()](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/edit_profile_page.dart#57-226) — wrap dengan BlocProvider + BlocListener:**
```dart
BlocProvider<ImageUploadCubit>(
  create: (_) => sl<ImageUploadCubit>(),
  child: BlocListener<ImageUploadCubit, ImageUploadState>(
    listener: (context, state) {
      if (state is ImageUploadError) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(state.message), backgroundColor: Colors.red),
        );
      }
    },
    child: Scaffold(
      // ... isi Scaffold yang sudah ada ...
      // Tambah progress bar di atas body:
      body: Column(children: [
        BlocBuilder<ImageUploadCubit, ImageUploadState>(
          builder: (_, state) => state is ImageUploadInProgress
            ? LinearProgressIndicator(value: state.progress, color: Color(0xFF00A735))
            : const SizedBox.shrink(),
        ),
        Expanded(child: /* existing body content */),
      ]),
    ),
  ),
)
```

---

### 5. UI: Edit Profil Caregiver

#### [MODIFY] [caregiver_edit_profile_page.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/caregiver/profile/presentation/pages/caregiver_edit_profile_page.dart)

Perubahan **identik** dengan [edit_profile_page.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/edit_profile_page.dart) di atas (arsitektur caregiver sengaja duplikasi dari mother, dan sama-sama memanggil `PATCH /users` di baliknya karena `caregiver_profiles` tidak punya kolom foto sendiri):
- State `File? _pickedFile` + `_pickImage()`
- Widget foto real
- [_saveProfile()](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/edit_profile_page.dart#46-55) async, upload (presign + PUT saja), lalu sisipkan `objectKey` ke `photoUrl` pada `UpdateCaregiverProfileRequestModel` (model terpisah, tapi tetap memanggil endpoint `PATCH /users` yang sama di baliknya, karena `caregiver_profiles` tidak punya kolom foto sendiri)
- `BlocProvider<ImageUploadCubit>` + `BlocListener`

---

### 6. UI: Tambah/Edit Profil Anak (`POST /children` & `PATCH /children/{child_id}`)

**Flow foto anak:** pick → upload (presign + PUT) → objectKey disertakan di `POST/PATCH /children` sebagai `photo_url` — pola yang sama dengan profil user di atas.

#### [MODIFY] [add_edit_child_profile_step1.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/widgets/add_edit_child_profile_step1.dart)

**Masalah:** widget ini `StatelessWidget`. Perlu local state (`File? _pickedFile`) → ubah ke `StatefulWidget`.

```dart
class AddEditChildProfileStep1 extends StatefulWidget {
  final bool isAddMode;
  const AddEditChildProfileStep1({super.key, required this.isAddMode});

  @override
  State<AddEditChildProfileStep1> createState() => AddEditChildProfileStep1State();
  // ← State class dibuat PUBLIC (tanpa underscore) agar parent bisa akses _pickedFile via GlobalKey
}

class AddEditChildProfileStep1State extends State<AddEditChildProfileStep1> {
  File? pickedFile; // PUBLIC agar parent akses via GlobalKey

  Future<void> _pickImage() async {
    final xfile = await ImagePicker().pickImage(
      source: ImageSource.gallery, imageQuality: 80,
    );
    if (xfile != null) setState(() => pickedFile = File(xfile.path));
  }

  @override
  Widget build(BuildContext context) {
    final cubit = context.read<AddEditProfileCubit>();
    final initialState = cubit.state;

    return SingleChildScrollView(
      padding: EdgeInsets.symmetric(horizontal: 20.w),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // ... header ...
          Center(
            child: GestureDetector(
              onTap: _pickImage,
              child: Column(children: [
                CircleAvatar(
                  radius: 45.r,
                  backgroundImage: pickedFile != null
                      ? FileImage(pickedFile!) as ImageProvider
                      : (initialState.photoUrl != null
                          ? NetworkImage(initialState.photoUrl!)
                          : null),
                  child: (pickedFile == null && initialState.photoUrl == null)
                      ? Icon(Icons.person, size: 36.w, color: Colors.grey)
                      : null,
                ),
                SizedBox(height: 8.h),
                Text('Unggah Foto', style: GoogleFonts.outfit(color: Color(0xFF00A735))),
              ]),
            ),
          ),
          // ... field lain tidak berubah ...
        ],
      ),
    );
  }
}
```

#### [MODIFY] [add_or_edit_child_profile.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/add_or_edit_child_profile.dart)

```dart
class _AddOrEditChildProfileState extends State<AddOrEditChildProfile> {
  int _currentStep = 0;
  bool _hasLoadedDetail = false;

  // GlobalKey untuk akses pickedFile dari Step1
  final _step1Key = GlobalKey<AddEditChildProfileStep1State>();
```

Gunakan key di [_buildCurrentStep()](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/add_or_edit_child_profile.dart#162-174):
```dart
case 0:
  return AddEditChildProfileStep1(key: _step1Key, isAddMode: widget.childData == null);
```

Tambah `BlocProvider<ImageUploadCubit>` sebagai parent terluar:
```dart
return MultiBlocProvider(
  providers: [
    BlocProvider(create: (_) => sl<ImageUploadCubit>()),
    BlocProvider(create: (_) => sl<ChildProfileFormCubit>()),
  ],
  child: MultiBlocListener(
    listeners: [
      BlocListener<ImageUploadCubit, ImageUploadState>(
        listener: (context, state) {
          if (state is ImageUploadError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(state.message), backgroundColor: Colors.red),
            );
          }
        },
      ),
      BlocListener<ChildProfileFormCubit, ChildProfileFormState>(
        listener: (context, state) { /* existing listener */ },
      ),
    ],
    child: /* existing Scaffold */,
  ),
);
```

[_submitProfile()](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/add_or_edit_child_profile.dart#221-251) → diubah menjadi `async`, upload foto sebelum submit:
```dart
Future<void> _submitProfile(BuildContext context) async {
  final formState = context.read<AddEditProfileCubit>().state;

  // Upload foto anak jika ada file baru yang dipilih
  String? objectKey;
  final pickedFile = _step1Key.currentState?.pickedFile;
  if (pickedFile != null) {
    final uploadCubit = context.read<ImageUploadCubit>();
    await uploadCubit.upload(
      file: pickedFile,
      category: 'child-profile',
      // Mode add: child_id belum ada → ownerId dikosongkan (lihat catatan endpoint 7 di endpoint.md).
      // Mode edit: child_id sudah ada → wajib diisi (lihat catatan endpoint 8).
      ownerId: widget.childData?.id,
    );
    if (!mounted) return;
    final st = uploadCubit.state;
    if (st is ImageUploadError) return; // error sudah di SnackBar
    if (st is ImageUploadSuccess) objectKey = st.objectKey;
  }

  // Build request — photo_url diisi objectKey baru (add/edit) atau URL lama (edit tanpa ganti foto)
  final request = ChildProfileRequestModel(
    fullName: formState.fullName,
    birthDate: "${formState.birthDate!.day.toString().padLeft(2,'0')}-"
        "${formState.birthDate!.month.toString().padLeft(2,'0')}-"
        "${formState.birthDate!.year}",
    gender: formState.gender,
    photoUrl: objectKey ?? formState.photoUrl, // ← objectKey baru || URL lama di edit mode
    allergies: formState.allergies,
    chronicDiseases: formState.chronicDiseases,
    diets: formState.diets,
    favoriteFoods: formState.favoriteFoods,
    favoriteTextures: formState.favoriteTextures,
    notes: formState.notes,
  );

  if (!mounted) return;
  if (widget.childData != null) {
    context.read<ChildProfileFormCubit>().updateProfile(widget.childData!.id, request);
  } else {
    context.read<ChildProfileFormCubit>().submitProfile(
      request, weightKg: formState.weightKg,
      heightCm: formState.heightCm, headCircumferenceCm: formState.headCircumferenceCm,
    );
  }
}
```

[_onNextStep()](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/add_or_edit_child_profile.dart#210-220) — ubah signature agar bisa await [_submitProfile](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/add_or_edit_child_profile.dart#221-251):
```dart
Future<void> _onNextStep(BuildContext context) async {
  if (!_validateCurrentStep(context)) return;
  if (_currentStep < 2) {
    setState(() => _currentStep++);
  } else {
    await _submitProfile(context);
  }
}
```

Button `onPressed` harus ikut async:
```dart
onPressed: state is ChildProfileFormLoading || state is ImageUploadInProgress
    ? null
    : () => _onNextStep(context),
```

---

## File Summary

| # | File | Status | Keterangan |
|---|------|--------|------------|
| 1 | [pubspec.yaml](file:///f:/Project/ProjectLomba/nuzagizi/pubspec.yaml) | MODIFY | `+mime: ^2.0.0` |
| 2 | `image_upload/domain/entities/presign_upload_entity.dart` | NEW | |
| 3 | `image_upload/domain/repositories/image_upload_repository.dart` | NEW | 2 method saja: `requestPresignUrl` + `putFileToStorage` |
| 4 | `image_upload/domain/usecases/upload_image_usecase.dart` | NEW | Return `objectKey`, tanpa percabangan mode |
| 5 | `image_upload/data/models/presign_upload_model.dart` | NEW | |
| 6 | `image_upload/data/datasources/image_api_service.dart` | NEW | Dio ber-auth, ke backend Go. Hanya `presignUpload` (tidak wire `confirmUpload` ke flow utama) |
| 7 | `image_upload/data/datasources/image_storage_service.dart` | NEW | storageDio polos, PUT ke R2 |
| 8 | `image_upload/data/repositories/image_upload_repository_impl.dart` | NEW | |
| 9 | `image_upload/presentation/cubit/image_upload_state.dart` | NEW | Success bawa `objectKey` |
| 10 | `image_upload/presentation/cubit/image_upload_cubit.dart` | NEW | `upload({file, category, ownerId})` — tanpa `withConfirm` |
| 11 | [core/di/service_locator.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/core/di/service_locator.dart) | MODIFY | +`storageDio` + 5 registrasi |
| 12 | [update_user_profile_request_model.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/data/models/update_user_profile_request_model.dart) | **MODIFY** | **[REVISI]** tambah field `photoUrl` — `PATCH /users` sekarang menerima `photo_url` |
| 13 | [update_caregiver_profile_request_model.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/caregiver/profile/data/models/update_caregiver_profile_request_model.dart) | **MODIFY** | **[REVISI]** idem, tambah `photoUrl` |
| 14 | [child_profile_request_model.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/data/models/child_profile_request_model.dart) | **TIDAK BERUBAH** | sudah ada `photoUrl` |
| 15 | [edit_profile_page.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/edit_profile_page.dart) (mother) | MODIFY | pick foto, [_saveProfile](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/edit_profile_page.dart#46-55) async, `photoUrl: objectKey` disisipkan langsung ke `PATCH /users` |
| 16 | [caregiver_edit_profile_page.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/caregiver/profile/presentation/pages/caregiver_edit_profile_page.dart) | MODIFY | identik dengan mother |
| 17 | [add_edit_child_profile_step1.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/widgets/add_edit_child_profile_step1.dart) | MODIFY | StatelessWidget→StatefulWidget, pick & preview |
| 18 | [add_or_edit_child_profile.dart](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/add_or_edit_child_profile.dart) | MODIFY | GlobalKey, BlocProvider upload, [_submitProfile](file:///f:/Project/ProjectLomba/nuzagizi/lib/features/mother/profile/presentation/pages/add_or_edit_child_profile.dart#221-251) async |

---

## Verification Plan

### Backend yang Harus Sudah Ada
- `POST /images/presign-upload` → `{ upload_url, object_key }`
- `PATCH /users`, `POST /children`, `PATCH /children/{id}` → menerima field `photo_url` diisi `object_key` (lihat `endpoint.md` yang sudah direvisi)
- `POST /images/confirm` bersifat opsional — tidak dipanggil di flow utama manapun

### Manual Test

| # | Skenario | Flow R2 | Hasil yang diharapkan |
|---|----------|---------|----------------------|
| 1 | Edit profil mother, ganti foto | presign→PUT→`photo_url` langsung di body `PATCH /users` | `users.photo_url` tersimpan sebagai `object_key` |
| 2 | Edit profil caregiver, ganti foto | identik test 1 (tetap lewat `PATCH /users`) | idem |
| 3 | Tambah profil anak, dengan foto | presign (ownerId kosong)→PUT→`photo_url` di body `POST /children` | objectKey masuk ke `POST /children` body sebagai `photo_url` |
| 4 | Edit profil anak, ganti foto | presign (ownerId = child_id)→PUT→`photo_url` di body `PATCH /children/{id}` | objectKey masuk ke `PATCH /children/{id}` |
| 5 | Edit/tambah profil (user/anak), tanpa pilih foto | tidak ada upload | `photo_url` tidak dikirim (user) / `null` (child add) / tetap URL lama (child edit) |
| 6 | Backend mati saat upload | error di Tahap 1 | SnackBar merah, flow berhenti, profil tidak tersimpan |
