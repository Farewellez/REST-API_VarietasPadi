# 1. varietas.go

```
// internal/domain/varietas.go
package domain

import "time"

type VarietasPadi struct {
	ID               int       `json:"id_padi"`
	VarietasKelas    string    `json:"varietas_kelas"`
	Warna            string    `json:"warna"`
	PanjangBijiMM    float64   `json:"panjang_biji_mm"`
	TeksturPermukaan string    `json:"tekstur_permukaan"`
	BentukUjungDaun  string    `json:"bentuk_ujung_daun"`
	WaktuPembuatan   time.Time `json:"waktu_pembuatan"`
}

type VarietasRepository interface {
	Create(data VarietasPadi) (VarietasPadi, error)
	FindByID(id int) (VarietasPadi, error)
	FindAll() ([]VarietasPadi, error)
	Update(data VarietasPadi) (VarietasPadi, error)
	Delete(id int) error
}

type VarietasService interface {
	TambahkanData(data VarietasPadi) (VarietasPadi, error)
	DapatkanDataByID(id int) (VarietasPadi, error)
	DapatkanSemuaData() ([]VarietasPadi, error)
	UbahData(data VarietasPadi) (VarietasPadi, error)
	HapusData(id int) error
}
```

VarietasRepository interface di sini berfungsi untuk menentukan apa yang harus dilakukan oleh database tanpa peduli bagaimana cara melakukannya. bagian ini juga yang akan menyediakan fitur CRUD data ke sumber eksternal. interface kedua yaitu VarietasService

# 2. varietas_repository.go

File ini berisi **implementasi** konkret dari kontrak `domain.VarietasRepository`. Lapisan ini bertanggung jawab penuh untuk berinteraksi dengan dunia luar (database) dan mengelola semua **Side Effects** I/O. 

### Struktur Repository

```
// internal/repository/varietas_repository.go
package repository

import (
	"context"
	"database/sql"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/domain"
)

type VarietasRepositoryImpl struct {
	DB *sql.DB // Koneksi database sebagai dependensi
}

func NewVarietasRepository(db *sql.DB) domain.VarietasRepository {
	return &VarietasRepositoryImpl{DB: db}
}
```

Struct `VarietasRepositoryImpl` menerima dependensi `*sql.DB`. Ini memastikan bahwa *repository* ini hanyalah **adapter** untuk database, memenuhi kontrak `domain.VarietasRepository`.

### Implementasi: FindByID (Security Focus)
Metode ini menunjukkan praktik terbaik dalam *data access*, terutama dalam hal **keamanan** dan **konkurensi** (aspek penting FP/Go).

```
func (r *VarietasRepositoryImpl) FindByID(ctx context.Context, id int) (domain.VarietasPadi, error) {
	var data domain.VarietasPadi
	
	// Prepared Statement: Mencegah SQL Injection
	query := "SELECT id_padi, varietas_kelas, warna, ... FROM data_pengamatan_padi WHERE id_padi = $1"
	
	// Eksekusi I/O dengan Context
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&data.IDPadi, 
		&data.VarietasKelas, 
        // ... scanning field lain
	)
    // ... Error handling
	return data, nil
}
```

1.  **`context.Context`**: Digunakan untuk mengelola batas waktu dan pembatalan permintaan. Ini adalah mekanisme kunci Golang untuk mengelola *state* non-lokal di operasi I/O yang rentan terhadap latensi (sejalan dengan prinsip FP).
2.  **Prepared Statement (`id_padi = $1`)**: Dengan memisahkan *query* dan parameter (`id`), kita secara aktif mencegah serangan **SQL Injection**, sebuah praktik wajib dalam **Cybersecurity**.
3.  **Side-Effect Execution**: Baris `r.DB.QueryRowContext(...)` adalah satu-satunya tempat di aplikasi yang secara langsung memicu *side effect* (interaksi database) yang diisolasi dari lapisan di atasnya.

# 3. varietas_service.go

Lapisan Service adalah inti dari logika bisnis aplikasi. File ini mengimplementasikan kontrak `domain.VarietasService` dan bertanggung jawab untuk **validasi data** dan **orkestrasi** alur kerja. Kunci di sini adalah menjaga fungsi-fungsi se-**Murni** mungkin, sesuai prinsip **Pemrograman Fungsional (FP)**.

### Struktur Service

```
// internal/service/varietas_service.go
package service

import "github.com/Farewellez/REST-API_VarietasPadi/internal/domain"

type VarietasServiceImpl struct {
	repo domain.VarietasRepository // Hanya bergantung pada Interface!
}

func NewVarietasService(repo domain.VarietasRepository) domain.VarietasService {
	return &VarietasServiceImpl{
		repo: repo,
	}
}
```

Struct `VarietasServiceImpl` menerima `domain.VarietasRepository` sebagai dependensi (Dependency Injection). Ini penting karena Service **tidak tahu** dan **tidak peduli** bagaimana data disimpan; ia hanya tahu cara memanggil kontrak yang sudah disepakati. Ini menjaga Service **terisolasi** dari dunia I/O (database).

### Implementasi: TambahkanData (Pure Functions)

Metode ini menunjukkan bagaimana **Logika Murni (Pure Logic)** dijalankan sebelum memicu *side effect* I/O.

```
// internal/service/varietas_service.go

func (s *VarietasServiceImpl) TambahkanData(data domain.VarietasPadi) (domain.VarietasPadi, error) {
	
	// PURE FUNCTION / VALIDASI: Tidak ada I/O, hanya memeriksa input.
	if data.VarietasKelas == "" || data.PanjangBijiMM <= 0 {
		return domain.VarietasPadi{}, errors.New("data varietas tidak valid: kelas dan panjang biji harus diisi")
	}
	
    // SIDE EFFECT MINOR: Memberikan timestamp
	data.WaktuPembuatan = time.Now()

	// SIDE EFFECT MAYOR: Memanggil Repository untuk I/O
	return s.repo.Create(data) 
}
```

1.  **Pure Logic**: Bagian *validasi* adalah fungsi yang murni. Untuk *input* yang sama, ia selalu menghasilkan *output* yang sama (error atau tidak), tanpa mengakses *state* global atau I/O.
2.  **Orkestrasi Side Effects**: Setelah validasi murni, Service memicu *side effect* utama dengan memanggil `s.repo.Create(data)`. Service bertindak sebagai orkestrator yang memastikan aturan bisnis dipenuhi *sebelum* interaksi database dilakukan.

# 4. varietas_handler.go

File ini berisi **HTTP Handler** yang berfungsi sebagai **Lapisan Eksternal** atau **Adapter**. Tugas utamanya adalah menerjemahkan permintaan HTTP (misalnya, body JSON) menjadi *method call* ke Service Layer, dan menerjemahkan hasilnya kembali menjadi respons HTTP (status code dan body JSON).

### Struktur Handler

```
// internal/http/handler/varietas_handler.go
package handler

import (
	"net/http"
    "encoding/json"
	"github.com/Farewellez/REST-API_VarietasPadi/internal/domain"
)

type VarietasHandler struct {
    service domain.VarietasService // Hanya bergantung pada Service Interface
}

func NewVarietasHandler(svc domain.VarietasService) *VarietasHandler {
    return &VarietasHandler{service: svc}
}
```
Handler menerima `domain.VarietasService` sebagai dependensi. Hal ini sesuai dengan aturan Clean Architecture: Lapisan Luar (Handler) harus bergantung pada Kontrak dari Lapisan Dalam (Service Interface).

### Implementasi: CreateVarietas (Adaptasi Permintaan)

Metode ini menunjukkan tiga langkah utama Handler: **Adaptasi Input**, **Pemanggilan Logika**, dan **Adaptasi Output**.

```
// internal/http/handler/varietas_handler.go

func (h *VarietasHandler) CreateVarietas(w http.ResponseWriter, r *http.Request) {
    var inputData domain.VarietasPadi
    
    // 1. ADAPTASI INPUT: Decode JSON dari HTTP Request
    if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest) // Error Eksternal (HTTP)
        return
    }

    // 2. PEMANGGILAN LOGIKA: Memanggil Service Layer (Logika Murni)
    createdData, err := h.service.TambahkanData(inputData) 

    if err != nil {
        // Service mengembalikan Error Bisnis/Validasi. Handler menentukan Status Code.
        http.Error(w, err.Error(), http.StatusUnprocessableEntity) 
        return
    }
    
    // 3. ADAPTASI OUTPUT: Encode JSON ke HTTP Response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // Status 201 Created
    json.NewEncoder(w).Encode(createdData)
}
```

* **Fokus Adapter**: Handler tidak melakukan validasi data bisnis (misalnya, apakah `PanjangBijiMM` > 0); tugas itu didelegasikan sepenuhnya ke Service Layer.
* **Status Code**: Handler bertanggung jawab menerjemahkan hasil (sukses/gagal) dari Service menjadi Status Code HTTP standar (misalnya, `201 Created` atau `422 Unprocessable Entity`), memastikan API berkomunikasi dengan benar ke klien.

# 5. cmd/server/main.go

File `main.go` adalah **titik komposit** aplikasi. Tugasnya bukanlah menjalankan logika bisnis (itu tugas Service), melainkan **merangkai (wiring)** semua dependensi (Repository, Service, Handler) dan memulai server.

Lapisan ini sepenuhnya bergantung pada prinsip **Dependency Injection (DI)** dan merupakan implementasi dari **Komposisi Fungsional**—membangun aplikasi besar dari unit-unit fungsi murni yang lebih kecil.

### Alur Wiring Dependencies

```
// cmd/server/main.go
package main

func main() {
    // 1. Initialization (Config dan Koneksi Database)
    config := config.LoadConfig() 
    db := database.NewPostgresDB(config.DBURL)
    
    // 2. WIRING: Komposisi Lapisan dari Dalam ke Luar
    // Repo dibuat dengan dependensi DB
    varietasRepo := repository.NewVarietasRepository(db) 
    
    // Service dibuat dengan dependensi Repo (Interface)
    varietasService := service.NewVarietasService(varietasRepo) 
    
    // Handler dibuat dengan dependensi Service (Interface)
    varietasHandler := handler.NewVarietasHandler(varietasService) 

    // 3. FINAL SIDE EFFECT: Menjalankan HTTP Server
    router := http.NewRouter(varietasHandler)
    log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
```
1.  **DI dan Komposisi**: Setiap lapisan dibuat secara berurutan, dengan lapisan luar menerima implementasi lapisan di bawahnya. Contoh: `VarietasService` tidak dibuat sampai `VarietasRepository` selesai, dan `VarietasHandler` tidak dibuat sampai `VarietasService` selesai. Hal ini menciptakan grafik dependensi yang jelas dan membuat setiap komponen mudah diuji secara independen (*unit testing*).
2.  **Context Side Effects**: Langkah inisiasi (`config.LoadConfig` dan `database.NewPostgresDB`) mengelola *initial side effects* (membaca *environment variables* dan membuat koneksi eksternal).
3.  **Final Side Effect**: `http.ListenAndServe` adalah *side effect* terakhir, yaitu menjalankan server HTTP yang akan menerima permintaan. Setelah titik ini, aplikasi berjalan dan menunggu I/O.

Terima kasih sudah memberitahu\! Maafkan kalau ada *glitch* di respon sebelumnya. Mari kita fokus kembali dan buatkan penjelasan yang jelas dan terstruktur untuk file **Setup & Deployment (DevOps)**, mengikuti format yang kamu minta.


# 6\. Setup & Deployment (DevOps)

File **`Dockerfile`** dan **`docker-compose.yml`** adalah implementasi dari prinsip **Infrastruktur sebagai Kode (IaC)**. Tugasnya adalah mendefinisikan dan mengotomatisasi *packaging* serta menjalankan aplikasi. Aspek ini selaras dengan prinsip **Deklaratif** (mirip FP), di mana kita mendefinisikan *state* akhir yang diinginkan, bukan langkah-langkah prosedural instalasi.

## A. Dockerfile (The Multi-Stage Build)

`Dockerfile` menggunakan pendekatan **Multi-Stage Build** untuk menghasilkan *binary* Go yang optimal, fokus pada **efisiensi** dan **keamanan** (*Cybersecurity*).

### Alur Multi-Stage Build


```dockerfile
# STAGE 1: BUILDER
FROM golang:1.25.4-alpine AS builder 

WORKDIR /app
COPY . .

# Build aplikasi menjadi binary statis

RUN go build -o /varietas\_api ./cmd/server/main.go

# STAGE 2: FINAL

FROM alpine:latest

WORKDIR /root/

# Copy hanya binary dan sertifikat SSL/TLS

COPY --from=builder /varietas\_api .

# Perintah menjalankan binary final

CMD ["./varietas\_api"]

```

1.  **Isolasi**: Tahap **Builder** mengandung semua *tools* pengembangan. Tahap **Final** *hanya* mengandung *binary* aplikasi yang sudah dikompilasi (`/varietas_api`).
2.  **Keamanan (Attack Surface Reduction)**: *Image* akhir yang ringan (berbasis `alpine:latest`) **tidak memiliki** *compiler*, *source code* Go, atau *tools* pengembangan lainnya. Ini mengurangi **Permukaan Serangan (Attack Surface)** aplikasi secara drastis, sebuah praktik penting dalam *hardening* aplikasi.
3.  **Final Side Effect**: Perintah `CMD ["./varietas_api"]` adalah *final side effect* yang dieksekusi oleh *container*, yaitu menjalankan aplikasi yang telah di-*wire* di `main.go`.

## B. docker-compose.yml (Declarative Configuration)

File ini mendefinisikan **layanan** dan **konfigurasi lingkungan** yang dibutuhkan aplikasi, memetakan *port*, dan menyuntikkan variabel penting dari file `.env`.

### Alur Konfigurasi Layanan

```yaml
services:
  api:
    build: . 
    # Mengambil environment variables dari file .env
    env_file:
      - .env
    ports:
      - "8080:8080"
    restart: always
````


1.  **Deklaratif**: Konfigurasi ini mendefinisikan **APA** yang harus dilakukan (*build* dari direktori ini, gunakan *port* 8080) tanpa menjelaskan **BAGAIMANA** prosesnya. Ini adalah cerminan dari prinsip deklaratif (mirip FP).
2.  **Manajemen Side Effects (Environment)**: `env_file: .env` adalah cara aman untuk menyuntikkan variabel sensitif (`DB_URL`, dll.) yang akan dibaca oleh `main.go` saat *initial side effects* (koneksi DB) dilakukan.
3.  **Port Mapping**: Memastikan *side effect* HTTP yang dimulai oleh `main.go` di *port* 8080 dapat diakses dari *host* lokal.


# 7\. database.go

File **`internal/database/database.go`** adalah *factory function* yang bertanggung jawab atas **Side Effect** kritis: **inisiasi koneksi database** PostgreSQL. Tugasnya adalah mengisolasi detail koneksi dari lapisan Repository.

```
// internal/database/database.go
package database

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

// NewPostgresDB membuat koneksi baru ke database
func NewPostgresDB(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	// Ping untuk memastikan koneksi aktif
	if err = db.Ping(); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}

	log.Println("Koneksi database PostgreSQL (NeonDB) berhasil!")
	return db
}
```

Fungsi `NewPostgresDB` di sini berfungsi sebagai **satu-satunya tempat** di mana koneksi I/O database dibuka, menjadikannya **Inisiasi Side Effect** yang penting. Objek `*sql.DB` yang dikembalikan kemudian disuntikkan (**Dependency Injection**) ke `VarietasRepositoryImpl`. Penggunaan **`db.Ping()`** dan `log.Fatalf` memastikan prinsip **DevOps** *fail fast* diterapkan; jika database tidak tersedia, aplikasi akan segera berhenti.

Tentu\! Saya akan memastikan format dan struktur penjelasannya sama persis seperti contoh yang kamu berikan, termasuk gaya *code block* dan penjelasannya.

Langkah logis berikutnya adalah mendokumentasikan file konfigurasi yang bertanggung jawab membaca *environment variables*, yaitu **`internal/config/config.go`**.


# 8\. config.go

File **`internal/config/config.go`** menangani **Side Effect** pertama di fase inisiasi: **membaca *environment variables*** (seperti `DB_URL` dari file `.env`). Tugasnya adalah mengisolasi detail pembacaan konfigurasi dari lapisan aplikasi lainnya.

```
// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config merepresentasikan konfigurasi aplikasi
type Config struct {
	DBURL string
	Port  string
}

// LoadConfig memuat konfigurasi dari environment variables atau file .env
func LoadConfig() *Config {
	// Side Effect: Membaca file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("Tidak ada file .env yang ditemukan. Menggunakan environment variable sistem.")
	}

	// Membaca variabel environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL environment variable tidak diatur.")
	}
	
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DBURL: dbURL,
		Port:  port,
	}
}
```

Fungsi `LoadConfig` di sini berfungsi untuk **mengisolasi** proses pembacaan variabel lingkungan dari *main function*. Ini adalah **Side Effect** yang dilakukan di awal *startup* aplikasi. *Struct* `Config` yang dikembalikan kemudian disuntikkan (**DI**) ke fungsi `NewPostgresDB` di lapisan **Database** (seperti yang terlihat di `main.go`), memastikan seluruh aplikasi menerima konfigurasi yang sudah diolah.


