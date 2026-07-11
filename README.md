# New Website Lelang API

Boilerplate backend Go dengan Gin sebagai HTTP framework, GORM sebagai ORM, dan SQLite sebagai database lokal. Struktur aplikasi memakai pendekatan Domain-Driven Design (DDD) sederhana dan dependency inversion.

## Struktur

```text
cmd/api/                         composition root dan HTTP server
internal/domain/reference/       entity, repository port, dan domain service
internal/infrastructure/database/ koneksi SQLite dan repository GORM
internal/infrastructure/test/     seluruh unit dan black-box test
internal/interfaces/httpapi/      Gin handler, DTO, mapper, dan router
```

## Menjalankan

Membutuhkan Go 1.25 atau lebih baru.

```bash
go run -buildvcs=false ./cmd/api
```

Atau dengan Docker:

```bash
docker build -t lelang-api .
docker run --rm -p 8080:8080 lelang-api
```

Port default untuk local adalah `80`. Image Docker menggunakan port `8080`.

Konfigurasi opsional:

```bash
PORT=80
SQLITE_PATH=lelang.db
DATABASE_URL=jdbc:oracle:thin:@//localhost:1521/FREEPDB1
DATABASE_USERNAME=system
DATABASE_PASSWORD=your-password
RUN_MIGRATIONS=false
MIGRATION_SCHEMA=CMS
```

## Deploy: build source di server dengan Docker Compose

Server hanya membutuhkan Git, Docker, dan Docker Compose. Go tidak perlu di-install
di server karena proses `go build` dijalankan oleh stage builder di `Dockerfile`.

Project ini sudah membaca konfigurasi Oracle dari environment variable. Redis tidak
ditambahkan karena aplikasi saat ini tidak menggunakannya.

1. Pastikan Oracle sudah berjalan dan tergabung ke external network yang sama.

   ```bash
   docker network create shared-network
   ```

   Perintah tersebut cukup dijalankan sekali. Jika network sudah ada, Docker akan
   menampilkan pesan bahwa network sudah tersedia.

2. Clone source di server dan buat `.env` produksi dari contoh. Jangan commit `.env`.

   ```bash
   cd /opt
   git clone URL_REPOSITORY new-web-lelang
   cd new-web-lelang
   cp .env.example .env
   nano .env
   ```

   Isi `DATABASE_URL` dengan hostname container Oracle pada `shared-network`, bukan
   `localhost`. Contoh: `jdbc:oracle:thin:@//shared-oracle:1521/XEPDB1`.

3. Validasi, build image dari source, dan jalankan container.

   ```bash
   docker-compose -f docker-compose.yaml config
   docker-compose -f docker-compose.yaml up -d --build
   docker logs -f new-web-lelang-api
   ```

   Compose memetakan port server `8081` ke port aplikasi `8080`. Tes dengan:

   ```bash
   curl http://localhost:8081/health
   ```

4. Untuk deployment berikutnya:

   ```bash
   cd /opt/new-web-lelang
   git pull origin main
   docker-compose -f docker-compose.yaml up -d --build
   docker logs -f new-web-lelang-api
   ```

Jika server memakai Compose v2, perintah yang sama dapat ditulis sebagai
`docker compose` (tanpa tanda hubung). `RUN_MIGRATIONS` sebaiknya tetap `false`
untuk startup normal dan hanya diaktifkan secara terkontrol saat migration memang
akan dijalankan.

## Database migration

Migration memakai GORM dan menyimpan riwayat pada tabel `GORM_SCHEMA_MIGRATIONS`. Migration hanya dijalankan ketika `RUN_MIGRATIONS=true`.

Contoh migration `001` membuat tabel `GORM_MIGRATION_EXAMPLE` pada schema yang ditentukan oleh `MIGRATION_SCHEMA`. User koneksi harus memiliki izin membuat object pada schema tersebut. Jalankan sekali dengan:

```powershell
$env:RUN_MIGRATIONS="true"
go run -buildvcs=false ./cmd/api
```

File migration berada di `internal/infrastructure/database/migration` dan hanya berisi SQL. Format nama file adalah `V001__description.sql`. Migration yang sudah tercatat akan dilewati. Jangan mengubah file migration lama; tambahkan file SQL dengan versi berikutnya.

## Endpoint

```bash
curl http://localhost/health
curl http://localhost/api/v1/reference-data
curl "http://localhost/api/v1/assets?search=rumah&page=1&limit=10"
curl http://localhost/api/v1/awards
```

Endpoint reference data menghasilkan data `kategori`, `tipe_aset`, `provinsi`, `metode_penjualan`, dan `kpknl`. Endpoint assets sementara menghasilkan data hardcode dan sudah menerima query filter serta pagination. Endpoint awards membaca data aktif (`IS_DELETED = 0`) dari tabel Oracle `CMS.MST_AWARDS`.

## Test

```bash
go test ./...
```
