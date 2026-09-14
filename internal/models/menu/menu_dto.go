package menu

// -----------------------------------------------------------------------------
// Payload Request ke AI Engine (all_child.json format)
// -----------------------------------------------------------------------------

type FoodEnginePayload []FoodEnginePayloadChild

type FoodEnginePayloadChild struct {
	ID                   string                       `json:"id"`
	Nama                 string                       `json:"nama"`
	Sex                  string                       `json:"sex"`
	UsiaBulan            int                          `json:"usia_bulan"`
	BeratKg              float64                      `json:"berat_kg"`
	PanjangCm            float64                      `json:"panjang_cm"`
	PreferensiHarga      string                       `json:"preferensi_harga"`
	JumlahProteinPerHari int                          `json:"jumlah_protein_per_hari"`
	Alergi               []string                     `json:"alergi"`
	LingkarKepalaCm      float64                      `json:"lingkar_kepala_cm"`
	RiwayatBerat         [][]any                      `json:"riwayat_berat"`
	ResepDokter          FoodEnginePayloadResepDokter `json:"resep_dokter"`
	Kondisi              []FoodEnginePayloadKondisi   `json:"kondisi"`
}

type FoodEnginePayloadResepDokter struct {
	Tingkat         *string                    `json:"tingkat"`
	NamaDokter      string                     `json:"nama_dokter"`
	NomorStr        *string                    `json:"nomor_str"`
	KomposisiPerKg  FoodEnginePayloadKomposisi `json:"komposisi_per_kg"`
	CatatanTambahan string                     `json:"catatan_tambahan"`
}

type FoodEnginePayloadKomposisi struct {
	EnergiKkal    float64 `json:"energi_kkal"`
	ProteinG      float64 `json:"protein_g"`
	NatriumMgMaks float64 `json:"natrium_mg_maks"`
}

type FoodEnginePayloadKondisi struct {
	ID          string `json:"id"`
	HariKe      int    `json:"hari_ke"`
	SesakNapas  bool   `json:"sesak_napas"`
	NapasCepat  bool   `json:"napas_cepat"`
	DemamTinggi bool   `json:"demam_tinggi"`
}

// -----------------------------------------------------------------------------
// Response dari AI Engine (output.json format)
// -----------------------------------------------------------------------------

type FoodEngineResponse struct {
	Meta    FoodEngineMeta    `json:"meta"`
	Anak    []FoodEngineAnak  `json:"anak"`
	Catatan FoodEngineCatatan `json:"catatan"`
}

type FoodEngineCatatan struct {
	Peringatan []string `json:"peringatan"`
}

type FoodEngineMeta struct {
	HariKe     int    `json:"hari_ke"`
	HariLabel  string `json:"hari_label"`
	Tanggal    string `json:"tanggal"`
	JumlahSesi int    `json:"jumlah_sesi"`
}

type FoodEngineAnak struct {
	ID               string                     `json:"id"`
	Nama             string                     `json:"nama"`
	Pertumbuhan      FoodEnginePertumbuhan      `json:"pertumbuhan"`
	Kebutuhan        FoodEngineKebutuhan        `json:"kebutuhan"`
	KondisiKesehatan FoodEngineKondisiKesehatan `json:"kondisi_kesehatan"`
	Menu             FoodEngineMenu             `json:"menu"`
}

type FoodEnginePertumbuhan struct {
	Nama      string             `json:"nama"`
	UsiaBulan int                `json:"usia_bulan"`
	BeratKg   float64            `json:"berat_kg"`
	PanjangCm float64            `json:"panjang_cm"`
	ZScores   map[string]float64 `json:"z_scores"`
	Status    map[string]string  `json:"status"`
	Label     map[string]string  `json:"label"`
	Action    string             `json:"action"`
	Pesan     []string           `json:"pesan"`
}

type FoodEngineKebutuhan struct {
	Sumber  string                  `json:"sumber"`
	Energi  float64                 `json:"energi"`
	Protein float64                 `json:"protein"`
	Lemak   float64                 `json:"lemak"`
	Karbo   float64                 `json:"karbo"`
	PerSesi FoodEngineMakroSelingan `json:"per_sesi"`
	Pesan   []string                `json:"pesan"`
}

type FoodEngineKondisiKesehatan struct {
	RujukDarurat    bool                           `json:"rujuk_darurat"`
	RujukKonsultasi bool                           `json:"rujuk_konsultasi"`
	PesanRujuk      []string                       `json:"pesan_rujuk"`
	Penyesuaian     FoodEngineKesehatanPenyesuaian `json:"penyesuaian"`
	AturanPerSlot   map[string]string              `json:"aturan_per_slot"`
	CatatanOrangTua []string                       `json:"catatan_orang_tua"`
	KondisiAktif    []string                       `json:"kondisi_aktif"`
}

type FoodEngineKesehatanPenyesuaian struct {
	TargetFaktor           map[string]float64 `json:"target_faktor"`
	BatasSeratPer100g      *float64           `json:"batas_serat_per_100g"`
	SeratMinimumPer100g    *float64           `json:"serat_minimum_per_100g"`
	LemakFaktorGlobal      *float64           `json:"lemak_faktor_global"`
	KapasitasLambungFaktor *float64           `json:"kapasitas_lambung_faktor"`
	TeksturOverride        *string            `json:"tekstur_override"`
	PreferHangat           bool               `json:"prefer_hangat"`
	PreferDingin           bool               `json:"prefer_dingin"`
	TambahSesi             bool               `json:"tambah_sesi"`
}

type FoodEngineMenu struct {
	Sesi     map[string]FoodEngineSesi     `json:"sesi"`
	Selingan map[string]FoodEngineSelingan `json:"selingan"`
}

type FoodEngineSesi struct {
	Bahan        []FoodEngineBahan              `json:"bahan"`
	PanduanMasak FoodEnginePanduanMasak         `json:"panduan_masak"`
	Makro        map[string]FoodEngineMakroItem `json:"makro"`
	Pesan        []string                       `json:"pesan"`
}

type FoodEngineSelingan struct {
	Tipe       string                  `json:"tipe"`
	Bahan      []FoodEngineBahan       `json:"bahan"`
	Makro      FoodEngineMakroSelingan `json:"makro"`
	Resep      FoodEngineResepSelingan `json:"resep"`
	FokusMikro []string                `json:"fokus_mikro"`
	Pesan      []string                `json:"pesan"`
}

type FoodEngineBahan struct {
	Slot        string                 `json:"slot"`
	Kode        string                 `json:"kode"`
	GramMentah  float64                `json:"gram_mentah"`
	GramMatang  float64                `json:"gram_matang"`
	Satuan      string                 `json:"satuan"`
	Substitutes []FoodEngineSubstitute `json:"pengganti"`
}

type FoodEngineSubstitute struct {
	Kode    string  `json:"kode"`
	Satuan  string  `json:"satuan"`
	Energi  float64 `json:"energi"`
	Protein float64 `json:"protein"`
	Lemak   float64 `json:"lemak"`
	Karbo   float64 `json:"karbo"`
	Jarak   float64 `json:"jarak"`
}

type FoodEnginePanduanMasak struct {
	UsiaLabel    string            `json:"usia_label"`
	Tekstur      string            `json:"tekstur"`
	JenisMasakan string            `json:"jenis_masakan"`
	ResepNama    string            `json:"resep_nama"`
	WaktuMasak   string            `json:"waktu_masak"`
	Catatan      string            `json:"catatan"`
	Bumbu        []FoodEngineBumbu `json:"bumbu"`
	Langkah      []string          `json:"langkah"`
}

type FoodEngineBumbu struct {
	Nama   string `json:"nama"`
	Satuan string `json:"satuan"`
}

type FoodEngineResepSelingan struct {
	Nama    string   `json:"nama"`
	Langkah []string `json:"langkah"`
}

type FoodEngineMakroItem struct {
	Dapat  float64 `json:"dapat"`
	Target float64 `json:"target"`
	Persen float64 `json:"persen"`
}

type FoodEngineMakroSelingan struct {
	Energi  float64 `json:"energi"`
	Protein float64 `json:"protein"`
	Lemak   float64 `json:"lemak"`
	Karbo   float64 `json:"karbo"`
}
