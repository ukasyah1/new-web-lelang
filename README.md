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

Port default untuk local adalah `80`. Docker tetap menggunakan port `8080`. Database SQLite otomatis dibuat sebagai `lelang.db` dan tabel serta data awal dibuat saat aplikasi pertama kali berjalan.

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
