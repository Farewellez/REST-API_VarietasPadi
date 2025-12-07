# **Bagian 1: Quick Start (Instruksi DevOps)**

## 1. Pendahuluan

Bagian ini dirancang agar Anda dapat menjalankan REST API berbasis Go secepat mungkin menggunakan Docker. Fokusnya sederhana: siapkan konfigurasi, jalankan satu perintah, dan API Anda langsung aktif tanpa perlu konfigurasi manual tambahan.

## 2. Prasyarat (Dependencies)

Sebelum memulai, pastikan lingkungan Anda sudah memiliki:

* **Docker** versi terbaru
* **Docker Compose**
* Akses ke **database PostgreSQL NeonDB** (atau koneksi PostgreSQL kompatibel lainnya)
* Sistem operasi apa pun yang mendukung Docker (Linux, macOS, Windows)

Pastikan seluruh dependency berjalan dengan baik karena Docker akan mengurus sisanya.

## 3. Konfigurasi Lingkungan (.env)

Aplikasi ini menggunakan file `.env` untuk menyimpan kredensial koneksi database dan konfigurasi runtime.
Letakkan file ini di direktori root proyek.

Penting untuk menjaga **keamanan kredensial**, terutama jika Anda bekerja dalam workflow kolaboratif. Jangan pernah commit file `.env` ke repositori publik.

Berikut template `.env` yang harus Anda isi:

```env
# Template File .env
# Kredensial untuk koneksi PostgreSQL (NeonDB)

NEON_USER=user_anda
NEON_PASSWORD=password_anda
NEON_HOST=ep-tight-field-5910.ap-southeast-1.aws.neon.tech
NEON_DBNAME=varietaspadi
NEON_PORT=5432

# Full Connection String (digunakan oleh aplikasi Go)
DB_URL="postgres://${NEON_USER}:${NEON_PASSWORD}@${NEON_HOST}:${NEON_PORT}/${NEON_DBNAME}?sslmode=require"

# Port aplikasi Go berjalan
API_PORT=8080
```

Catatan penting:
Gunakan `sslmode=require` pada `DB_URL` untuk memastikan koneksi ke NeonDB tetap terenkripsi dan aman.

## 4. Build dan Run (Satu Perintah)

Setelah `.env` siap, jalankan aplikasi dalam satu langkah menggunakan Docker Compose.

### **Build & Run**

```sh
docker compose up --build
```

Perintah ini akan:

* Membangun image Go API
* Mengonfigurasi environment sesuai `.env`
* Menjalankan container beserta dependensinya

Setelah proses selesai, REST API Anda dapat diakses pada port yang telah Anda tentukan.

### **Hentikan**

```sh
docker compose down
```

Perintah ini menghentikan seluruh container dan merapikan resource yang digunakan.

---

# **Bagian 2: Endpoint API (Fungsionalitas)**

## 1. Pendahuluan

Bagian ini menjelaskan seluruh fungsionalitas REST API CRUD untuk entitas *Varietas Padi*. Informasi disusun agar pengembang maupun penguji dapat memahami cara menggunakan setiap endpoint, struktur data yang dibutuhkan, serta status respons yang diharapkan. Semua endpoint berjalan dari **base URL**:

```
http://localhost:8080
```

(Atau URL tunneling seperti **ngrok** jika digunakan.)

---

## 2. Status dan UI Endpoints

Endpoint berikut digunakan untuk mengakses halaman antarmuka HTML sederhana dan mengambil seluruh data dalam format JSON.

| Method  | Endpoint        | Deskripsi                                             |
| :------ | :-------------- | :---------------------------------------------------- |
| **GET** | `/`             | Menyajikan halaman UI HTML (Demo CRUD).               |
| **GET** | `/api/varietas` | Mengambil semua data Varietas Padi dalam format JSON. |

---

## 3. Operasi CRUD

### Tabel Operasi CRUD

| Method     | Endpoint             | Deskripsi                        | Status Code Sukses |
| :--------- | :------------------- | :------------------------------- | :----------------- |
| **POST**   | `/api/varietas`      | Membuat data baru.               | **201 Created**    |
| **GET**    | `/api/varietas/{id}` | Membaca data berdasarkan ID.     | **200 OK**         |
| **PUT**    | `/api/varietas/{id}` | Memperbarui data yang sudah ada. | **200 OK**         |
| **DELETE** | `/api/varietas/{id}` | Menghapus data.                  | **204 No Content** |

---

### Contoh Body JSON

Contoh ini digunakan untuk **POST** dan **PUT** pada endpoint `/api/varietas` dan `/api/varietas/{id}`.

```json
// Contoh Body JSON untuk POST /api/varietas dan PUT /api/varietas/{id}
{
    "varietas_kelas": "Japonica",
    "warna": "Putih Susu",
    "panjang_biji_mm": 5.8,
    "tekstur_permukaan": "Halus",
    "bentuk_ujung_daun": "Tumpul"
}
```

**Catatan penting:**

* Field **id_padi** dan **waktu_pembuatan** **tidak perlu dikirimkan** pada POST atau PUT.
* Keduanya dibuat dan diisi otomatis oleh backend/database saat data disimpan.

---

# **Bagian 3: Penjelasan Kode & Arsitektur (The Why)**

## 1. Pendahuluan

Arsitektur proyek ini dibangun dengan memadukan *Clean Architecture* dan prinsip *Functional Programming* (FP). Keduanya berfokus pada pemisahan concerns, isolasi efek samping, serta kontrol yang jelas terhadap aliran data. Pendekatan ini memungkinkan sistem tetap mudah di-*test*, stabil, dan dapat berkembang tanpa menambah beban teknis jangka panjang. Dalam konteks API CRUD, integrasi FP membantu menjaga logika bisnis tetap murni, dapat diprediksi, dan bebas dari dependensi eksternal.

---

## 2. Struktur Proyek

Struktur direktori dirancang agar setiap bagian memiliki tanggung jawab jelas. Berikut struktur utama proyek:

```
C:\...\REST-API_VARIETASPADI
│   .env
│   Dockerfile
│   go.mod
│   docker-compose.yml
│
+---cmd
│   \---server
│           main.go
│
+---internal
    +---domain                (Kontrak Data dan Interface)
    +---service               (Logika Bisnis / Pure Functions)
    +---repository            (Side Effects I/O / Database)
    \---http                  (Handler / External Layer)
            views
                index.html
```

**Ringkasan Fungsi Folder:**

* **cmd/server**: Titik masuk aplikasi yang hanya menginisialisasi komponen.
* **internal/domain**: Struktur data, interface repository, dan kontrak layanan yang stabil.
* **internal/service**: Tempat logika bisnis murni disusun menggunakan prinsip FP.
* **internal/repository**: Akses database, berisi semua *side effects* yang tidak boleh masuk ke logika murni.
* **internal/http**: Handler, routing, dan integrasi ke luar (browser, HTTP client).
* **views**: Berisi HTML untuk antarmuka demo.

---

## 3. Clean Architecture dan FP

Clean Architecture memberikan batasan yang tegas antara domain murni dan lapisan yang berurusan dengan side effects. Prinsip FP kemudian memperkuat batasan tersebut melalui konsep *pure functions*, *immutability*, dan *isolation*.

### **Lapisan Domain (Entities & Contracts)**

Lapisan ini mendefinisikan struktur data dan interface. Tidak ada logika bisnis kompleks maupun ketergantungan pada database. Ini adalah pusat dari prinsip *immutability* karena objek domain hanya digunakan sebagai nilai, bukan sebagai entitas yang dimodifikasi berantai.

### **Lapisan Service (Pure Functions)**

Lapisan service memuat logika yang bersifat deterministik. Selama input sama dan dependency repository dipisahkan via interface, hasilnya dapat diprediksi.

Berikut contoh implementasi berupa *pure logic*:

```go
// internal/service/varietas_service.go (Pure Logic Example)

func (s *VarietasService) TambahkanData(ctx context.Context, data domain.VarietasPadi) (domain.VarietasPadi, error) {
	
	// VALIDASI adalah Pure Function
	if data.VarietasKelas == "" || data.PanjangBijiMM <= 0 {
		return domain.VarietasPadi{}, errors.New("data varietas tidak valid")
	}
	
	// Memanggil Side Effect
	return s.repo.Create(ctx, data) 
}
```

<penjelasan>  
Blok validasi di atas termasuk *pure function* karena tidak berinteraksi dengan I/O, tidak memodifikasi state di luar fungsi, dan hanya bergantung pada input. Ketika akhirnya memanggil repository, barulah terjadi side effect yang sengaja diisolasi. Strategi ini menjaga prediktabilitas dan memudahkan pengujian unit.
</penjelasan>

### **Lapisan Repository (Side Effects & Security Layer)**

Repository bertanggung jawab pada segala bentuk interaksi dengan database. Di sinilah aturan keamanan diterapkan: penggunaan prepared statements, sanitasi input, dan pemetaan data.

Contoh implementasi:

```go
// internal/repository/varietas_repository.go (Security Example)

func (r *VarietasRepositoryImpl) FindByID(ctx context.Context, id int64) (domain.VarietasPadi, error) {
	query := "SELECT ... FROM data_pengamatan_padi WHERE id_padi = $1"
	
	// Menggunakan Prepared Statement (QueryRowContext)
	err := r.DB.QueryRowContext(ctx, query, id).Scan(...)
    // ...
}
```

<penjelasan>  
Prepared statement (`$1`) melindungi aplikasi dari SQL Injection. Lapisan ini mengisolasi efek samping sehingga logika service tetap bersifat deterministik. Prinsip *isolation* dalam FP memastikan hanya lapisan inilah yang boleh menyentuh infrastruktur.  
</penjelasan>

### **Lapisan HTTP (Handlers / Delivery Layer)**

Handler bertugas menerima permintaan HTTP, memvalidasi input, lalu meneruskan ke service. Lapisan ini tidak boleh berisi logika bisnis utama. Handler hanya menjadi *translator* antara dunia eksternal dan dunia internal.

---

## 4. Deployment & DevOps

Aplikasi menggunakan pendekatan **Docker Multi-Stage Build** untuk meningkatkan efisiensi dan keamanan.

**Alasan pendekatan ini penting:**

* **Binary akhir jauh lebih kecil** karena hanya membawa dependensi yang diperlukan.
* **Mengurangi attack surface** karena layer build terpisah dari layer runtime.
* **Mempercepat build CI/CD** melalui caching layer.
* **Menstabilkan environment** sehingga API berjalan konsisten di berbagai platform.

Prinsip isolasi pada FP juga tercermin pada DevOps: runtime hanya menjalankan binary yang bersih, tanpa embel-embel compiler atau tool bawaan build environment. Ini selaras dengan Clean Architecture yang membatasi dependency eksternal agar tidak mencemari logika inti.

---
