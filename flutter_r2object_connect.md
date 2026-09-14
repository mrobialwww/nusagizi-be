# Mekanisme Upload Gambar ke Cloudflare R2 dari Sisi Flutter (Clean Architecture)

### Pelengkap dari `r2object_connect.md` (backend Go + Gin) dan `flutter_clean_architecture.md`

---

## 1. Konfirmasi: Apakah Ini Best Practice?

**Ya.** Pola presigned URL (presign → upload langsung ke R2 → confirm) adalah pola resmi yang direkomendasikan Cloudflare untuk kasus upload dari client (mobile/web), termasuk untuk kebutuhan seperti foto profil, foto produk/resep, dan fitur bergaya feed (foto harian anak) dengan volume upload tinggi.

Dokumentasi resmi Cloudflare R2 menyatakan presigned URL cocok untuk memberi akses sementara ke objek tertentu — misalnya mengizinkan pengguna mengupload file langsung ke R2 — dan bisa dipakai langsung dari browser, aplikasi mobile, atau HTTP client apa pun.

### ⚠️ Koreksi Penting: Presigned PUT vs Presigned POST

Banyak tutorial Flutter + S3 di internet memakai pola **presigned POST** (form fields: `policy`, `signature`, `AWSAccessKeyId`, dikirim via `FormData`/`MultipartFile`). **Pola ini TIDAK berlaku untuk R2.**

Dokumentasi resmi R2 menyatakan R2 hanya mendukung presigned URL untuk method **GET, HEAD, PUT, dan DELETE** — method **POST (multipart form upload) tidak didukung**.

**Implikasi ke Flutter:** upload ke R2 wajib pakai `dio.put()` dengan body berupa raw bytes/stream file — **bukan** `FormData`/`MultipartFile` seperti kebanyakan contoh S3-generic di internet.

Dua aturan lain dari dokumentasi resmi R2 yang relevan ke sisi Flutter:
- **Content-Type harus identik** antara saat generate presign (di backend) dan saat request PUT (di Flutter) — signature menyertakan header ini, mismatch → `403 SignatureDoesNotMatch`.
- **CORS hanya relevan untuk Flutter Web** (browser). Untuk Flutter mobile native (Android/iOS), konfigurasi CORS di bucket R2 tidak diperlukan.

---

## 2. Struktur Folder (Feature Baru: `image_upload`)

Dibuat sebagai satu fitur reusable (bukan per use case), mengikuti struktur wajib di `flutter_clean_architecture.md`:

```
lib/features/image_upload/
├── data/
│   ├── datasources/
│   │   ├── image_api_service.dart        # ke backend Go (presign + confirm)
│   │   └── image_storage_service.dart    # PUT langsung ke R2
│   ├── models/
│   │   └── presign_upload_model.dart
│   └── repositories/
│       └── image_upload_repository_impl.dart
├── domain/
│   ├── entities/
│   │   └── presign_upload_entity.dart
│   ├── repositories/
│   │   └── image_upload_repository.dart
│   └── usecases/
│       └── upload_image_usecase.dart
└── presentation/
    ├── cubit/
    │   ├── image_upload_cubit.dart
    │   └── image_upload_state.dart
    └── pages/    # page pemanggil: foto profil / produk / foto harian anak
```

**Kenapa dua data source terpisah (`image_api_service` & `image_storage_service`)?**
Karena secara fisik ini dua server berbeda — backend Go vs `*.r2.cloudflarestorage.com` — dengan kebutuhan header yang berbeda (backend butuh auth token, R2 justru tidak boleh menerima header asing). Ini analog dengan alasan "Tahap 1 vs 2 wajib terpisah" di `r2object_connect.md`.

> Catatan Aturan Wajib #1 (`flutter_clean_architecture.md`): `dart:io File` aman dipakai di domain layer karena bagian dari Dart SDK inti, bukan package eksternal/Flutter — setara dengan `dartz` yang sudah diizinkan di dokumen tersebut.

---

## 3. Alur Lengkap

```
[Page] user pilih foto (image_picker) → File
      |
      v
[Cubit] emit(InProgress) → panggil UploadImageUseCase(params)
      |
      v
[UseCase] orkestrasi 3 tahap (business logic, sesuai Aturan Wajib #7):
      1. repository.requestPresignUrl()   → backend Go
      2. repository.putFileToStorage()    → PUT langsung ke R2
      3. repository.confirmUpload()       → backend Go
      |
      v
[RepositoryImpl] gabungkan 2 data source (ImageApiService + ImageStorageService)
      |
      v
[Cubit] result.fold() → emit Error / Success
      |
      v
[Page] BlocConsumer: listener → SnackBar, builder → progress bar
```

---

## 4. Kode Lengkap per Layer

### Domain Layer

```dart
// domain/entities/presign_upload_entity.dart
class PresignUploadEntity {
  final String uploadUrl;
  final String objectKey;
  const PresignUploadEntity({required this.uploadUrl, required this.objectKey});
}
```

```dart
// domain/repositories/image_upload_repository.dart
import 'dart:io';
import 'package:dartz/dartz.dart';
import 'package:nusagizi/core/error/failures.dart';
import 'presign_upload_entity.dart';

typedef UploadProgressCallback = void Function(int sent, int total);

abstract class ImageUploadRepository {
  Future<Either<Failure, PresignUploadEntity>> requestPresignUrl({
    required String category,      // 'profile' | 'product' | 'recipe' | 'child-daily'
    required String contentType,
    String? ownerId,
  });

  Future<Either<Failure, void>> putFileToStorage({
    required String uploadUrl,
    required File file,
    required String contentType,
    UploadProgressCallback? onProgress,
  });

  Future<Either<Failure, void>> confirmUpload({
    required String objectKey,
    required String category,
    String? ownerId,
    Map<String, dynamic>? metadata,
  });
}
```

```dart
// domain/usecases/upload_image_usecase.dart
import 'dart:io';
import 'package:dartz/dartz.dart';
import 'package:nusagizi/core/error/failures.dart';
import 'package:nusagizi/core/usecase/usecase.dart';
import '../repositories/image_upload_repository.dart';

class UploadImageParams {
  final File file;
  final String category;
  final String contentType;
  final String? ownerId;
  final Map<String, dynamic>? metadata;
  final UploadProgressCallback? onProgress;

  UploadImageParams({
    required this.file,
    required this.category,
    required this.contentType,
    this.ownerId,
    this.metadata,
    this.onProgress,
  });
}

class UploadImageUseCase implements UseCase<void, UploadImageParams> {
  final ImageUploadRepository repository;
  UploadImageUseCase({required this.repository});

  @override
  Future<Either<Failure, void>> call(UploadImageParams params) async {
    final presignResult = await repository.requestPresignUrl(
      category: params.category,
      contentType: params.contentType,
      ownerId: params.ownerId,
    );

    return presignResult.fold(
      (failure) async => Left(failure),
      (presign) async {
        final uploadResult = await repository.putFileToStorage(
          uploadUrl: presign.uploadUrl,
          file: params.file,
          contentType: params.contentType,
          onProgress: params.onProgress,
        );

        return uploadResult.fold(
          (failure) async => Left(failure),
          (_) => repository.confirmUpload(
            objectKey: presign.objectKey,
            category: params.category,
            ownerId: params.ownerId,
            metadata: params.metadata,
          ),
        );
      },
    );
  }
}
```

### Data Layer

```dart
// data/models/presign_upload_model.dart
import 'package:nusagizi/features/image_upload/domain/entities/presign_upload_entity.dart';

class PresignUploadModel extends PresignUploadEntity {
  const PresignUploadModel({required super.uploadUrl, required super.objectKey});

  factory PresignUploadModel.fromJson(Map<String, dynamic> json) => PresignUploadModel(
        uploadUrl: json['upload_url'] as String,
        objectKey: json['object_key'] as String,
      );
}
```

```dart
// data/datasources/image_api_service.dart — ke backend Go
import 'package:dio/dio.dart';
import 'package:nusagizi/core/error/exceptions.dart';
import '../models/presign_upload_model.dart';

abstract class ImageApiService {
  Future<PresignUploadModel> presignUpload({
    required String category, required String contentType, String? ownerId,
  });
  Future<void> confirmUpload({
    required String objectKey, required String category,
    String? ownerId, Map<String, dynamic>? metadata,
  });
}

class ImageApiServiceImpl implements ImageApiService {
  final Dio dio; // Dio ber-baseUrl ke backend Go + interceptor Authorization
  ImageApiServiceImpl({required this.dio});

  @override
  Future<PresignUploadModel> presignUpload({
    required String category, required String contentType, String? ownerId,
  }) async {
    try {
      final res = await dio.post('/images/presign-upload', data: {
        'category': category,
        'content_type': contentType,
        if (ownerId != null) 'owner_id': ownerId,
      });
      return PresignUploadModel.fromJson(res.data);
    } on DioException catch (e) {
      throw ServerException(message: e.message ?? 'Gagal meminta presigned URL');
    }
  }

  @override
  Future<void> confirmUpload({
    required String objectKey, required String category,
    String? ownerId, Map<String, dynamic>? metadata,
  }) async {
    try {
      await dio.post('/images/confirm', data: {
        'object_key': objectKey,
        'category': category,
        if (ownerId != null) 'owner_id': ownerId,
        if (metadata != null) ...metadata,
      });
    } on DioException catch (e) {
      throw ServerException(message: e.message ?? 'Gagal konfirmasi upload');
    }
  }
}
```

```dart
// data/datasources/image_storage_service.dart — PUT langsung ke R2
import 'dart:io';
import 'package:dio/dio.dart';
import 'package:nusagizi/core/error/exceptions.dart';

abstract class ImageStorageService {
  Future<void> putFile({
    required String uploadUrl, required File file, required String contentType,
    void Function(int sent, int total)? onProgress,
  });
}

class ImageStorageServiceImpl implements ImageStorageService {
  // PENTING: Dio instance ini TERPISAH dari Dio backend Go.
  // Tidak boleh punya baseUrl / interceptor Authorization — request ini
  // menuju *.r2.cloudflarestorage.com, bukan server kita.
  final Dio storageDio;
  ImageStorageServiceImpl({required this.storageDio});

  @override
  Future<void> putFile({
    required String uploadUrl, required File file, required String contentType,
    void Function(int sent, int total)? onProgress,
  }) async {
    try {
      final length = await file.length();
      await storageDio.put(
        uploadUrl,
        data: file.openRead(),           // stream → hemat memory utk file besar
        options: Options(headers: {
          Headers.contentLengthHeader: length, // wajib, agar onSendProgress akurat
          'Content-Type': contentType,          // HARUS = Content-Type saat presign
        }),
        onSendProgress: onProgress,
      );
    } on DioException catch (e) {
      if (e.response?.statusCode == 403) {
        throw ServerException(
          message: 'Upload ditolak R2 (403) — Content-Type tidak cocok atau URL expired',
        );
      }
      throw ServerException(message: e.message ?? 'Upload ke storage gagal');
    }
  }
}
```

```dart
// data/repositories/image_upload_repository_impl.dart
import 'dart:io';
import 'package:dartz/dartz.dart';
import 'package:nusagizi/core/error/exceptions.dart';
import 'package:nusagizi/core/error/failures.dart';
import '../datasources/image_api_service.dart';
import '../datasources/image_storage_service.dart';
import 'package:nusagizi/features/image_upload/domain/entities/presign_upload_entity.dart';
import 'package:nusagizi/features/image_upload/domain/repositories/image_upload_repository.dart';

class ImageUploadRepositoryImpl implements ImageUploadRepository {
  final ImageApiService apiService;
  final ImageStorageService storageService;
  ImageUploadRepositoryImpl({required this.apiService, required this.storageService});

  @override
  Future<Either<Failure, PresignUploadEntity>> requestPresignUrl({
    required String category, required String contentType, String? ownerId,
  }) async {
    try {
      final model = await apiService.presignUpload(
        category: category, contentType: contentType, ownerId: ownerId);
      return Right(model);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    }
  }

  @override
  Future<Either<Failure, void>> putFileToStorage({
    required String uploadUrl, required File file, required String contentType,
    UploadProgressCallback? onProgress,
  }) async {
    try {
      await storageService.putFile(
        uploadUrl: uploadUrl, file: file, contentType: contentType, onProgress: onProgress);
      return const Right(null);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    }
  }

  @override
  Future<Either<Failure, void>> confirmUpload({
    required String objectKey, required String category,
    String? ownerId, Map<String, dynamic>? metadata,
  }) async {
    try {
      await apiService.confirmUpload(
        objectKey: objectKey, category: category, ownerId: ownerId, metadata: metadata);
      return const Right(null);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message));
    }
  }
}
```

### Presentation Layer

```dart
// presentation/cubit/image_upload_state.dart
abstract class ImageUploadState extends Equatable {
  const ImageUploadState();
  @override
  List<Object> get props => [];
}
class ImageUploadInitial extends ImageUploadState {}
class ImageUploadInProgress extends ImageUploadState {
  final double progress;
  const ImageUploadInProgress({required this.progress});
  @override
  List<Object> get props => [progress];
}
class ImageUploadSuccess extends ImageUploadState {}
class ImageUploadError extends ImageUploadState {
  final String message;
  const ImageUploadError({required this.message});
  @override
  List<Object> get props => [message];
}
```

```dart
// presentation/cubit/image_upload_cubit.dart
import 'dart:io';
import 'package:mime/mime.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/usecases/upload_image_usecase.dart';
import 'image_upload_state.dart';

class ImageUploadCubit extends Cubit<ImageUploadState> {
  final UploadImageUseCase uploadImageUseCase;
  ImageUploadCubit({required this.uploadImageUseCase}) : super(ImageUploadInitial());

  Future<void> upload({
    required File file,
    required String category,
    String? ownerId,
    Map<String, dynamic>? metadata,
  }) async {
    emit(const ImageUploadInProgress(progress: 0));

    // Deteksi content-type SEKALI, dipakai konsisten di tahap presign & PUT
    final contentType = lookupMimeType(file.path) ?? 'application/octet-stream';

    final result = await uploadImageUseCase(UploadImageParams(
      file: file,
      category: category,
      contentType: contentType,
      ownerId: ownerId,
      metadata: metadata,
      onProgress: (sent, total) {
        if (total > 0) emit(ImageUploadInProgress(progress: sent / total));
      },
    ));

    result.fold(
      (failure) => emit(ImageUploadError(message: failure.message)),
      (_) => emit(ImageUploadSuccess()),
    );
  }
}
```

Pemilihan file di Page: pakai `image_picker` (kamera/galeri), convert `XFile` → `File(xfile.path)`, lalu panggil `cubit.upload(...)`.

### Dependency Injection

```dart
// core/di/service_locator.dart (tambahan)

// Dio khusus backend Go — baseUrl + interceptor auth
sl.registerLazySingleton<Dio>(
  () => Dio(BaseOptions(baseUrl: 'https://api.nusagizi.com'))
    ..interceptors.add(AuthInterceptor()),
  instanceName: 'apiDio',
);

// Dio khusus upload ke R2 — polos, TANPA baseUrl/interceptor
sl.registerLazySingleton<Dio>(() => Dio(), instanceName: 'storageDio');

sl.registerLazySingleton<ImageApiService>(
  () => ImageApiServiceImpl(dio: sl<Dio>(instanceName: 'apiDio')));
sl.registerLazySingleton<ImageStorageService>(
  () => ImageStorageServiceImpl(storageDio: sl<Dio>(instanceName: 'storageDio')));
sl.registerLazySingleton<ImageUploadRepository>(
  () => ImageUploadRepositoryImpl(
    apiService: sl<ImageApiService>(),
    storageService: sl<ImageStorageService>(),
  ));
sl.registerLazySingleton<UploadImageUseCase>(
  () => UploadImageUseCase(repository: sl<ImageUploadRepository>()));
sl.registerFactory<ImageUploadCubit>(
  () => ImageUploadCubit(uploadImageUseCase: sl<UploadImageUseCase>()));
```

`GetIt` mendukung `instanceName` untuk mendaftarkan dua instance `Dio` berbeda dalam satu tipe — cara paling praktis memisahkan Dio backend vs Dio storage tanpa wrapper class tambahan.

---

## 5. Detail Teknis yang Sering Terlewat

1. **Content-Type harus identik** di dua tempat: payload `presign-upload` (dipakai backend Go generate signature) dan header `Content-Type` saat `dio.put()`. Deteksi sekali pakai package `mime` (`lookupMimeType`), simpan, pakai ulang.
2. **`Content-Length` wajib diset manual** agar `onSendProgress` akurat — dikonfirmasi langsung oleh dokumentasi resmi `dio`: content-length harus diset jika ingin subscribe ke sending progress.
3. **Gunakan `file.openRead()` (stream), bukan `file.readAsBytes()`** — supaya file besar tidak dimuat penuh ke memory sebelum dikirim.
4. **Retry**: presigned URL bisa dipakai berulang selama belum expired, jadi retry cukup panggil ulang `putFileToStorage` dengan `uploadUrl` yang sama. Tapi `file.openRead()` harus dipanggil ulang tiap percobaan — Stream yang sudah dikonsumsi tidak bisa dipakai lagi.
5. **Dio terpisah untuk R2**: jangan pakai Dio yang sama dengan interceptor `Authorization`. Header asing yang ikut terkirim ke R2 berisiko membuat request tidak bersih dan menyulitkan debugging.
6. **Presigned PUT ≠ Presigned POST**: R2 tidak mendukung presigned POST (form upload). Selalu pakai `dio.put()` dengan body stream/bytes, bukan `FormData`.

---

## 6. Referensi Dokumentasi Resmi

- Cloudflare R2 — Presigned URLs: https://developers.cloudflare.com/r2/api/s3/presigned-urls/
- Cloudflare R2 — Get Started: https://developers.cloudflare.com/r2/get-started/
- Cloudflare R2 — S3 API Compatibility: https://developers.cloudflare.com/r2/api/s3/api/
- Cloudflare R2 — Konfigurasi CORS (untuk Flutter Web): https://developers.cloudflare.com/r2/buckets/cors/
- Cloudflare R2 — Contoh aws-sdk-go: https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/
- Cloudflare R2 — Public Buckets: https://developers.cloudflare.com/r2/buckets/public-buckets/
- Package `dio` (dart.dev/pub.dev) — dokumentasi resmi package: https://pub.dev/packages/dio