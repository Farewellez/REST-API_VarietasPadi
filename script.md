# 🎤 **SCRIPT PRESENTASI – PRESENTER 1 (Lapisan Domain – Functional Programming Mindset)**

*(Durasi: 3–5 menit)*

---

## 🎯 **Pembukaan**

Halo semuanya, perkenalkan, saya akan memulai presentasi dengan membahas fondasi paling penting dari proyek kita: **Lapisan Domain**.
Lapisan inilah yang menentukan *bahasa*, *aturan*, dan *kebenaran logika* dalam sistem kita.

Dan menariknya, di proyek REST API Go untuk *Klasifikasi Varietas Padi* ini, kita merancang domain dengan pendekatan **Functional Programming**, terutama fokus pada **immutability**, **kontrak yang jelas**, dan **integritas data**.

Mari kita mulai masuk ke struktur inti domain-nya.

---

## 🧩 **1. Definisi Data: VarietasPadi sebagai Pure Value Object**

Di Functional Programming, data idealnya bersifat **immutable**—tidak diubah-ubah seenaknya.
Kenapa ini penting? Karena data yang tidak berubah memberikan prediktabilitas, keamanan, dan mencegah bug halus pada sistem concurrent seperti Go.

Berikut struktur domain kita:

### <bagian kode>

```go
package domain

type VarietasPadi struct {
    IDPadi          string
    VarietasKelas   string
    Warna           string
    PanjangBijiMM   float64
    WaktuPembuatan  int
}
```

### </bagian kode>

### <penjelasan>

Struct **VarietasPadi** kita posisikan sebagai **Value Object** murni.
Tidak ada method yang mengubah dirinya sendiri, tidak ada state tersembunyi, tidak ada behavior berbahaya.

Dengan menjadikan struct ini *immutable secara desain* (tidak ada setter, tidak ada pointer ke dalam), maka:

* Data lebih aman dari perubahan tidak sengaja
* Operasi CRUD menjadi lebih dapat diprediksi
* Integritas data lebih terjaga, terutama ketika API ini digunakan secara concurrent

Ini mengikuti prinsip FP bahwa “**data adalah nilai, bukan entitas yang berubah-ubah**”.

### </penjelasan>

---

## 🏗️ **2. Kontrak Domain: Interface VarietasRepository**

Setelah datanya stabil, kita butuh jembatan yang menghubungkan domain dengan dunia luar—database, API handler, dan sebagainya.
Dalam konsep Functional Programming dan Clean Architecture, ini disebut **contract**.

### <bagian kode>

```go
type VarietasRepository interface {
    Create(v VarietasPadi) error
    GetAll() ([]VarietasPadi, error)
    FindByID(id string) (VarietasPadi, error)
    Update(v VarietasPadi) error
    Delete(id string) error
}
```

### </bagian kode>

### <penjelasan>

Interface ini adalah **kontrak murni**—sebuah perjanjian.
Ia tidak peduli apakah data berasal dari PostgreSQL, file JSON, atau API lain.

Dengan memisahkan **aksi** (CRUD) dari **implementasi**, kita mendapatkan:

* Fleksibilitas tinggi: repository bisa diganti kapan saja
* Testability meningkat: kita bisa menulis unit test hanya dengan mock
* Keamanan arsitektur: domain tidak pernah “tercemar” detail teknis database

Inilah inti filosofi FP: **fokus pada fungsi, bukan perangkat keras atau tempat data disimpan**.

### </penjelasan>

---

## 🛡️ **3. Integritas & Keamanan Data Melalui Prinsip FP**

Functional Programming bukan hanya soal gaya—dia membantu menjaga **kehandalan**, terutama dalam proyek backend seperti ini.

### 💡 Dengan **Immutability**:

* Data tidak akan dimodifikasi secara tidak sengaja oleh goroutine lain
* Struktur data menjadi lebih mudah diverifikasi
* Risiko *race condition* lebih rendah

### 💡 Dengan **Contracts**:

* Akses data lebih terstruktur
* Validasi alur data menjadi konsisten
* Domain bebas dari manipulasi ilegal atau bypassing logika

Menggabungkan keduanya membuat domain kita **aman, konsisten, dan tahan terhadap kesalahan implementasi**.

---

## 🧭 **Penutup**

Jadi, lapisan domain kita dibangun dengan tiga pondasi besar:

1. **Struktur data yang murni** — VarietasPadi sebagai immutable value.
2. **Kontrak fungsional** — repository sebagai interface yang memisahkan aturan dengan implementasi.
3. **Integritas data yang kuat** melalui prinsip-prinsip Functional Programming.

Lapisan domain ini menjadi pusat stabilitas seluruh REST API kita. Tanpa domain yang bersih, lapisan lain—handler, service, repository—tidak akan pernah kokoh.

Terima kasih. Setelah ini, Presentasi akan dilanjutkan ke lapisan berikutnya.

---

# 🎤 **SCRIPT PRESENTASI – PRESENTER 2 (Lapisan Service – Functional Programming Approach)**

*(Durasi 3–5 menit)*

---

## 🌱 **Transisi dari Presenter 1**

Terima kasih kepada Presenter 1 yang sudah memaparkan konsep domain dengan sangat jelas—mulai dari *immutable value object* hingga kontrak repository.
Sekarang, saya akan melanjutkan pembahasan ke lapisan berikutnya, yaitu **Lapisan Service**, yang menjadi pusat logika bisnis dalam arsitektur kita.

Ini adalah lapisan yang bertugas untuk menggabungkan prinsip Functional Programming dengan kebutuhan praktis dari sebuah REST API backend.

---

# ⚙️ **1. Gambaran Umum Lapisan Service**

Lapisan **Service** adalah jembatan antara domain yang bersifat murni dan lapisan infrastruktur seperti database dan handler.
Tugas utamanya bukan melakukan penyimpanan atau pengambilan data, melainkan:

* menjalankan **pure logic**,
* mengatur **urutan side effect**,
* dan memastikan **aturan domain terimplementasi dengan benar**.

Pada bagian ini, kita menggunakan pendekatan Functional Programming untuk memisahkan **logika murni** dari **aksi I/O**, sehingga kode tetap bersih, teruji, dan aman.

---

# 🧩 **2. Dependency Injection – Bergantung pada Interface, Bukan Implementasi**

Mari lihat definisi struct `VarietasService`:

### <bagian kode>

```go
type VarietasService struct {
	repo domain.VarietasRepository // Hanya bergantung pada INTERFACE
}
```

### </bagian kode>

### <penjelasan>

Di sini kita menerapkan **Dependency Injection (DI)** dengan cara yang sangat bersih.
Service **hanya mengetahui interface**, bukan implementasinya.

Artinya:

* Layer service tidak tahu atau peduli apakah repository memakai PostgreSQL, MongoDB, atau in-memory.
* Ketika diuji, kita bisa mengganti repository dengan *mock* tanpa memodifikasi service.
* Kode menjadi lebih modular, loosely coupled, dan sangat sesuai dengan prinsip FP:
  **fungsi tidak boleh bergantung pada detail efek samping**.

### </penjelasan>

---

# 🧪 **3. Pure Functions – Logika Murni Tanpa Efek Samping**

Sekarang mari kita fokus pada method utama: `TambahkanData`.

### <bagian kode>

```go
if data.VarietasKelas == "" || data.PanjangBijiMM <= 0 {
	return domain.VarietasPadi{}, errors.New("data varietas tidak valid: kelas dan panjang biji harus diisi")
}
```

### </bagian kode>

### <penjelasan>

Bagian ini adalah **pure function**:

* Hanya menerima input.
* Hanya menghasilkan output.
* Tidak melakukan I/O.
* Tidak mengubah state global.

Validasi seperti ini sangat cocok untuk pendekatan FP karena:

* mudah ditest,
* predictable,
* bebas side-effect.

Dengan memisahkan pure logic dari tindakan I/O, kita memastikan service tetap ringan dan aman.

### </penjelasan>

---

# 🔄 **4. Side-Effect Orchestration – Peran Utama Service Layer**

Setelah pure logic selesai, barulah service mengatur side effect. Mari lihat dua bagian side effect berikut:

### <bagian kode>

```go
// SIDE EFFECT 1: Memberikan nilai WaktuPembuatan
data.WaktuPembuatan = time.Now()
```

### </bagian kode>

### <penjelasan>

Mengambil waktu saat ini (`time.Now()`) adalah efek samping karena melibatkan akses ke *system clock*.
Di FP, efek samping harus **dikontrol dan diisolasi**, dan layanan (service) adalah tempat yang tepat.

### </penjelasan>

---

### <bagian kode>

```go
// SIDE EFFECT 2: Memanggil Repository untuk I/O
return s.repo.Create(ctx, data)
```

### </bagian kode>

### <penjelasan>

Memanggil repository adalah efek samping I/O yang paling jelas.
Service tidak mengerjakan penyimpanan data secara langsung, melainkan hanya:

* memanggil kontrak,
* mengatur urutan eksekusi,
* dan menangani error.

Dengan menempatkan I/O hanya di satu titik terkontrol seperti ini, kita menjaga arsitektur tetap bersih dan konsisten.

### </penjelasan>

---

# 🧪 **5. Testability – Kekuatan Kombinasi FP + DI**

Pendekatan Functional Programming memberikan keuntungan besar pada **pengujian**:

1. **Pure logic** di service dapat diuji tanpa database.
2. Karena dependency injection, kita bisa menggunakan **mock repository** untuk mensimulasikan hasil I/O.
3. Dengan memisahkan logic & side-effect, setiap bagian lebih sederhana dan stabil.

Artinya, service ini mendukung *unit test* yang cepat, presisi, dan mudah ditulis.

---

# 🎯 **Penutup**

Jadi, pada lapisan Service kita menerapkan prinsip Functional Programming dengan cara:

* Menggunakan **dependency injection** untuk memisahkan logika dari implementasi I/O.
* Menempatkan **pure functions** untuk logika validasi yang aman dan mudah diuji.
* Mengatur **side-effect orchestration** agar efek samping tetap terisolasi dan terprediksi.
* Meningkatkan **testability** sistem secara keseluruhan.

Lapisan service inilah yang memastikan domain dijalankan dengan benar dan data yang masuk ke repository tetap terjaga integritasnya.

Terima kasih, dan presentasi akan dilanjutkan ke lapisan berikutnya.

# 🎤 **SCRIPT PRESENTASI – PRESENTER 3 (Lapisan Repository – FP & Cybersecurity Focus)**

*(Durasi: 3–5 menit)*

---

## 🔗 **Transisi dari Presenter 2**

Terima kasih kepada Presenter 2 yang telah menjelaskan lapisan Service, termasuk peran **dependency injection**, **pure function**, dan **side-effect orchestration**.
Sekarang, saya akan membahas **Lapisan Repository**, tempat semua **side effects** dikontrol dan dieksekusi secara aman.

Di sinilah service kita tetap murni, sementara I/O—baik database maupun network—dikelola secara terisolasi, sesuai prinsip Functional Programming. Selain itu, repository juga menjadi titik kunci dalam penerapan **praktik cybersecurity**.

---

# ⚙️ **1. Repository: Isolasi Side Effects**

Di Functional Programming, kita ingin **meminimalkan efek samping**.
Repository bertindak sebagai **satu-satunya layer yang melakukan I/O**:

### <bagian kode>

```go
type VarietasRepositoryImpl struct {
	DB *sql.DB // Dependency koneksi database
}
```

### </bagian kode>

### <penjelasan>

`VarietasRepositoryImpl` hanya mengetahui koneksi database (`*sql.DB`) dan mengimplementasikan kontrak `VarietasRepository`.
Ini artinya:

* Semua operasi Service tetap murni dan tidak tahu detil database.
* Repository menjadi satu-satunya titik eksekusi **side effects**, menjaga integritas domain dan logika bisnis.
* Kode lebih mudah diuji karena Service bisa diuji tanpa database, cukup dengan mock repository.

### </penjelasan>

---

# ⏱️ **2. Context Management: Mengelola State Non-Lokal**

FP mengajarkan bahwa state harus **dioperasikan secara eksplisit**, termasuk saat melakukan I/O.
Go menyediakan `context.Context` untuk menangani timeout, cancelation, dan propagation metadata.

### <bagian kode>

```go
func (r *VarietasRepositoryImpl) FindByID(ctx context.Context, id int64) (domain.VarietasPadi, error) {
	var data domain.VarietasPadi
	
	query := "SELECT id_padi, varietas_kelas, warna, panjang_biji_mm, tekstur_permukaan, bentuk_ujung_daun, waktu_pembuatan FROM data_pengamatan_padi WHERE id_padi = $1"
	
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&data.IDPadi, 
		&data.VarietasKelas, 
		&data.Warna, 
		&data.PanjangBijiMM,
		&data.TeksturPermukaan,
		&data.BentukUjungDaun,
		&data.WaktuPembuatan,
	)
```

### </bagian kode>

### <penjelasan>

Dengan `ctx`:

* Kita bisa mengatur **timeout** agar query tidak menggantung.
* Membatalkan query jika client sudah menutup koneksi.
* Menjaga **efisiensi dan keamanan** server.

Ini adalah contoh **pengelolaan state non-lokal ala FP**, di mana efek samping I/O dikontrol secara eksplisit dan tidak tersembunyi.

### </penjelasan>

---

# 🛡️ **3. Security Best Practices – Prepared Statements**

Repository juga bertanggung jawab menjaga keamanan saat berinteraksi dengan database.
Prepared statements dan parameterized queries menjadi praktik wajib untuk mencegah SQL Injection.

### <bagian kode>

```go
query := "SELECT ... FROM data_pengamatan_padi WHERE id_padi = $1"
err := r.DB.QueryRowContext(ctx, query, id).Scan(
	&data.IDPadi, 
	&data.VarietasKelas, 
	&data.Warna, 
	// ...
)
```

### </bagian kode>

### <penjelasan>

* `$1` adalah placeholder parameter, bukan string konkatenasi.
* Data user tidak pernah langsung dimasukkan ke query string.
* Ini **mengurangi risiko serangan SQL Injection**, menjaga integritas dan keamanan sistem.

Kombinasi FP + praktik keamanan membuat lapisan repository:

* Aman
* Terisolasi
* Dapat diprediksi

### </penjelasan>

---

# 🔄 **4. Mengembalikan Data Murni ke Service**

Setelah side effect selesai:

### <bagian kode>

```go
if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.VarietasPadi{}, errors.New("data tidak ditemukan")
	}
	return domain.VarietasPadi{}, err
}
return data, nil
```

### </bagian kode>

### <penjelasan>

Repository mengembalikan **struct VarietasPadi** yang bersifat murni ke service:

* Service menerima data siap pakai, tanpa mengetahui query atau I/O.
* Kode tetap modular, testable, dan predictable.
* Integritas data dijaga, karena service tidak memanipulasi hasil query secara langsung.

### </penjelasan>

---

# 🎯 **5. Penutup**

Kesimpulan dari lapisan repository:

1. **Side-effect isolation:** Semua I/O hanya terjadi di sini, menjaga service tetap murni.
2. **Context management:** State non-lokal dikontrol secara eksplisit menggunakan `context.Context`.
3. **Security:** Prepared statements mencegah SQL Injection.
4. **FP Compliance:** Mengembalikan data murni ke service, memisahkan logika bisnis dan I/O.

Dengan pendekatan ini, REST API kita aman, modular, mudah diuji, dan tetap sesuai prinsip Functional Programming.

Terima kasih. Presenter berikutnya akan membahas **Lapisan Handler / Controller** untuk menghubungkan API dengan dunia luar.

# 🎤 **SCRIPT PRESENTASI – PRESENTER 4 (Integrasi & Deployment – FP & DevOps Focus)**

*(Durasi: 3–5 menit)*

---

## 🔗 **Transisi dari Presenter 3**

Terima kasih kepada Presenter 3 yang telah membahas lapisan Repository, termasuk **isolasi side effects**, **context management**, dan **praktik keamanan**.
Sekarang kita sampai pada **lapisan paling atas**: Integrasi seluruh sistem dan deployment-nya.

Di sinilah **semua lapisan domain, service, dan repository dikomposisikan secara fungsional** dan kemudian **di-deploy** menggunakan paradigma **declarative infrastructure** ala DevOps.

---

# ⚙️ **1. Wiring / Functional Composition di main.go**

`cmd/server/main.go` bertindak sebagai titik komposisi utama.
Di sini, kita menyatukan semua lapisan dengan **dependency injection** dan menyiapkan side effect akhir—yaitu HTTP server.

### <bagian kode>

```go
func main() {
	// 1. Dependency Initialization (The State)
	config, err := config.LoadConfig() 
	db, err := database.NewPostgresDB(config.DBURL)

	// 2. Wiring Dependencies (Functional Composition)
	varietasRepo := repository.NewVarietasRepository(db)
	varietasService := service.NewVarietasService(varietasRepo)
	varietasHandler := handler.NewVarietasHandler(varietasService)

	// 3. Start Side Effect (HTTP Server)
	router := http.NewRouter(varietasHandler)
	log.Printf("🚀 Server siap dijalankan pada port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
```

### </bagian kode>

### <penjelasan>

* **Step 1 – Initialization:** Semua **state** eksternal seperti konfigurasi dan koneksi database dibuat secara eksplisit.
* **Step 2 – Wiring / Composition:**

  * Repository, Service, dan Handler dikomposisikan.
  * Ini adalah **composition function ala FP**: setiap lapisan adalah input bagi lapisan berikutnya.
  * Dependency Injection menjaga **modularitas dan testability**.
* **Step 3 – Final Side Effect:**

  * `http.ListenAndServe` adalah efek samping terakhir yang memulai aplikasi.
  * Karena semua logika bisnis sudah dipisahkan, side effect ini bersih, sederhana, dan terprediksi.

### </penjelasan>

---

# 📦 **2. Declarative Infrastructure dengan Docker Compose**

Untuk deployment, kita menerapkan **paradigma deklaratif**: mendeklarasikan *apa yang diinginkan*, bukan *bagaimana memprosesnya*.

Ini selaras dengan prinsip FP: menyatakan **state** dan **komposisi** daripada mengubahnya secara imperatif.

### <bagian kode>

```yaml
# docker-compose.yml (Declarative Infrastructure)

services:
  api:
    build: . 
    environment:
      DB_URL: postgres://${NEON_USER}:...
    ports:
      - "8080:8080"
```

### </bagian kode>

### <penjelasan>

* Kita mendeklarasikan **service API**, environment variable, dan port mapping.
* Docker Compose menangani proses build dan run container.
* Kita tidak menulis skrip imperatif untuk memulai server atau mengatur container—semua **terdeskripsikan**.
* Ini memudahkan **reproduksi lingkungan**, **scaling**, dan **continuous deployment**.

### </penjelasan>

---

# 🔄 **3. Hubungan FP dengan DevOps**

* **Functional Composition** → modularitas tinggi, setiap layer saling terhubung tanpa mengubah state global.
* **Declarative Infrastructure** → lingkungan dideklarasikan, bukan diprogram imperatif.
* Keduanya menekankan **prediktabilitas**, **reproducibility**, dan **kontrol penuh atas side effect**.

Hasilnya, sistem backend kita:

* **Aman**: side effect dikontrol di repository dan server.
* **Testable**: service bisa diuji secara murni.
* **Deployable**: lingkungan container bisa direplikasi dan dikelola oleh DevOps.

---

# 🎯 **4. Penutup**

Dengan `main.go` dan `docker-compose.yml`, kita telah menunjukkan:

1. **Wiring secara fungsional**: semua layer disusun secara modular dan composable.
2. **Side effect final**: HTTP server dijalankan secara terkontrol.
3. **Declarative Deployment**: Docker Compose menyatakan lingkungan tanpa imperatif.

Ini adalah contoh nyata bagaimana **Functional Programming dan DevOps** saling melengkapi dalam pengembangan REST API modern.

Terima kasih, dan ini menutup sesi presentasi tentang **Integrasi dan Deployment**.
