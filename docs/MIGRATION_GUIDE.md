# Atlas Migration Quick Start Guide

Panduan lengkap untuk menggunakan Atlas migrations di Go Fiber Starter Kit.

## 🎯 Apa itu Atlas?

Atlas adalah tool modern untuk database schema management yang **otomatis generate migrations dari GORM entities**. Tidak perlu menulis SQL manual lagi!

### Keuntungan Atlas

✅ **Auto-generate** - Migrations dibuat otomatis dari GORM models
✅ **Type-safe** - Migrations dari Go code, bukan SQL manual
✅ **Version control** - Migration files di-commit ke git
✅ **Multi-database** - Support MySQL, PostgreSQL, SQL Server
✅ **Smart diffing** - Hanya generate perubahan yang diperlukan
✅ **CI/CD ready** - CLI commands untuk automation

## 📦 Instalasi

### 1. Install Atlas CLI

**Linux/macOS (curl)**:
```bash
curl -sSf https://atlasgo.sh | sh
```

**macOS (Homebrew)**:
```bash
brew install ariga/tap/atlas
```

**Go Install**:
```bash
go install ariga.io/atlas@latest
```

**Using Project Makefile**:
```bash
make atlas-install
```

**Verify Installation**:
```bash
atlas version
```

### 2. Install Go Dependencies

```bash
go get -u ariga.io/atlas-provider-gorm/gormschema
```

### 3. Setup Database URL

Tambahkan ke `.env`:

```bash
# Development Database
ATLAS_DEV_DB_URL=mysql://root:password@localhost:3306/myapp_dev

# Production Database (untuk deployment)
ATLAS_DB_URL=mysql://root:password@localhost:3306/myapp
```

**PostgreSQL**:
```bash
ATLAS_DEV_DB_URL=postgres://user:pass@localhost:5432/myapp_dev?sslmode=disable
ATLAS_DB_URL=postgres://user:pass@localhost:5432/myapp?sslmode=disable
```

### 4. Run Setup Script

```bash
./scripts/atlas-setup.sh
```

Script ini akan:
- Check instalasi Atlas
- Validate project structure
- Check konfigurasi environment
- Test Atlas configuration

## 🚀 Quick Start - 3 Langkah Mudah

### Step 1: Buat/Ubah GORM Entity

```go
// internal/domain/product/entity.go
package product

import "gorm.io/gorm"

type Product struct {
    ID          uint    `gorm:"primaryKey"`
    Name        string  `gorm:"type:varchar(200);not null"`
    Description string  `gorm:"type:text"`
    Price       float64 `gorm:"type:decimal(10,2);not null"`
    Stock       int     `gorm:"default:0"`
    gorm.Model
}
```

### Step 2: Register Model

Edit `internal/config/database.go`, tambahkan model baru di function `GetModelsForMigration()`:

```go
func GetModelsForMigration() []interface{} {
    models := []interface{}{
        &user.User{},
        &post.Post{},
        &product.Product{}, // ← Tambahkan ini
        // ... models lainnya
    }
    return models
}
```

**Keuntungan**: Hanya perlu tambah 1x! Function ini digunakan untuk:
- ✅ AutoMigrate (dev mode dengan `DB_GEN=true`)
- ✅ Atlas migrations (production mode)

### Step 3: Generate & Apply Migration

```bash
# Generate migration
make atlas-diff

# Review migration file di: migrations/

# Apply migration
make atlas-apply
```

Selesai! ✅ Database schema sudah terupdate!

## 📝 Common Commands

### Generate Migrations

```bash
# Auto-generate dengan nama otomatis
make atlas-diff

# Dengan custom name
atlas migrate diff create_products_table --env dev

# Untuk production environment
make atlas-diff-prod
```

### Apply Migrations

```bash
# Apply ke development
make atlas-apply

# Apply ke production (dengan konfirmasi)
make atlas-apply-prod

# Dry run (preview tanpa apply)
atlas migrate apply --env dev --dry-run
```

### Check Status

```bash
# Lihat status migrations
make atlas-status

# Validate migration files
make atlas-validate

# Inspect database schema
make atlas-inspect
```

### Full Workflow

```bash
# Generate + Validate + Apply dalam 1 command
make atlas-migrate
```

## 🔄 Development Workflow

### Mode 1: AutoMigrate (Prototyping Cepat)

Untuk development yang cepat, gunakan GORM AutoMigrate:

**Set di `.env`**:
```bash
DB_GEN=true
```

Server akan otomatis sync schema saat startup.

⚠️ **Hanya untuk lokal development!** Jangan gunakan di production.

### Mode 2: Atlas Migrations (Production Ready)

Untuk production dan team collaboration:

**Set di `.env`**:
```bash
DB_GEN=false
```

**Workflow**:
1. Ubah GORM entity
2. `make atlas-diff` - Generate migration
3. Review SQL yang di-generate
4. `make atlas-apply` - Apply migration
5. Commit migration files ke git

## 🎓 Tutorial: Membuat Fitur Baru

### Contoh: Menambahkan Tabel Categories

**1. Buat Domain Structure**

```bash
mkdir -p internal/domain/category
```

**2. Buat Entity**

```go
// internal/domain/category/entity.go
package category

import "gorm.io/gorm"

type Category struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"type:varchar(100);not null;uniqueIndex"`
    Slug        string `gorm:"type:varchar(100);not null;uniqueIndex"`
    Description string `gorm:"type:text"`
    gorm.Model
}
```

**3. Register di Config Database**

```go
// internal/config/database.go
import "starter-gofiber/internal/domain/category"

func GetModelsForMigration() []interface{} {
    models := []interface{}{
        &user.User{},
        &post.Post{},
        &category.Category{}, // Tambahkan ini
        // ... models lainnya
    }
    return models
}
```

**4. Generate Migration**

```bash
atlas migrate diff create_categories_table --env dev
```

**Output**:
```
Analyzing GORM models...
Comparing with database schema...
✓ Generated: migrations/20260107120000_create_categories_table.sql
```

**5. Review Generated SQL**

```bash
cat migrations/20260107120000_create_categories_table.sql
```

```sql
-- Create "categories" table
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `slug` varchar(100) NOT NULL,
  `description` text,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_categories_name` (`name`),
  UNIQUE INDEX `idx_categories_slug` (`slug`),
  INDEX `idx_categories_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

**6. Apply Migration**

```bash
make atlas-apply
```

**Output**:
```
✓ Applied migration 20260107120000_create_categories_table.sql
Migration complete!
```

**7. Verify**

```bash
make atlas-status
```

```
Migration Status: OK
Current Version: 20260107120000
Pending Migrations: 0
```

### Contoh: Menambahkan Kolom ke Tabel Existing

**1. Update Entity**

```go
// internal/domain/user/entity.go
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"type:varchar(200);not null"`
    Email    string `gorm:"type:varchar(200);uniqueIndex;not null"`
    Bio      string `gorm:"type:text"` // ← Field baru
    Website  string `gorm:"type:varchar(200)"` // ← Field baru
    Password string `gorm:"type:varchar(150);not null"`
    gorm.Model
}
```

**2. Generate Migration**

```bash
atlas migrate diff add_user_bio_and_website --env dev
```

**3. Review & Apply**

```bash
# Review
cat migrations/20260107130000_add_user_bio_and_website.sql

# Apply
make atlas-apply
```

## 🌍 Production Deployment

### Step 1: Development (Local)

```bash
# 1. Buat/ubah GORM entities
# 2. Generate migrations
make atlas-diff-dev

# 3. Test migrations
make atlas-apply-dev
make atlas-status
```

### Step 2: Commit Changes

```bash
git add migrations/
git add internal/domain/
git add internal/config/database.go
git commit -m "feat: add categories table"
git push origin main
```

### Step 3: Production Deployment

**Manual**:
```bash
# SSH ke production server
ssh user@production-server

# Pull latest code
cd /path/to/app
git pull origin main

# Set production DB URL
export ATLAS_DB_URL="mysql://user:pass@prod-db:3306/myapp"

# Apply migrations
make atlas-apply-prod
```

**Using CI/CD** (GitHub Actions):

```yaml
# .github/workflows/migrate.yml
name: Database Migration

on:
  push:
    branches: [main]
    paths:
      - 'migrations/**'

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Atlas
        run: curl -sSf https://atlasgo.sh | sh
      
      - name: Apply Migrations
        run: atlas migrate apply --env prod
        env:
          ATLAS_DB_URL: ${{ secrets.PROD_DB_URL }}
```

## 🔧 Advanced Usage

### Custom Migration (Data Migration)

Jika perlu migration dengan logic khusus:

```bash
# Generate empty migration
atlas migrate new custom_data_migration --env dev
```

Edit file dan tambahkan SQL custom:

```sql
-- migrations/20260107140000_custom_data_migration.sql

-- Migrate existing data
UPDATE users SET role = 'user' WHERE role IS NULL;

-- Add constraint
ALTER TABLE users MODIFY COLUMN role ENUM('admin', 'user') NOT NULL;
```

### Multiple Databases

Jika project menggunakan multiple databases:

```hcl
// atlas.hcl
env "db1" {
  src = data.external_schema.gorm.url
  url = var.db1_url
  migration {
    dir = "file://migrations/db1"
  }
}

env "db2" {
  src = data.external_schema.gorm.url
  url = var.db2_url
  migration {
    dir = "file://migrations/db2"
  }
}
```

Generate migrations:
```bash
atlas migrate diff --env db1
atlas migrate diff --env db2
```

### Migration Policies

Prevent destructive changes di production:

```hcl
// atlas.hcl
env "prod" {
  lint {
    destructive {
      error = true  // Block DROP operations
    }
  }
  
  diff {
    skip {
      drop_schema = true
      drop_table  = true
    }
  }
}
```

## ❗ Troubleshooting

### "Failed to connect to database"

```bash
# Check database URL
echo $ATLAS_DEV_DB_URL

# Test connection
atlas schema inspect --env dev
```

### "Schema drift detected"

Database schema berbeda dengan migrations:

```bash
# Inspect current state
atlas schema inspect --env dev > current_schema.sql

# Generate migration untuk fix drift
atlas migrate diff fix_schema_drift --env dev
```

### Model Tidak Terdeteksi

Pastikan model sudah di-register di `internal/config/database.go`:

```go
func GetModelsForMigration() []interface{} {
    models := []interface{}{
        &user.User{},
        &yourmodel.YourModel{}, // Harus ada di sini!
        // ... models lainnya
    }
    return models
}
```

Ini adalah single source of truth untuk semua models.

### Migration Failed

Jika migration gagal di tengah jalan:

```bash
# Check status
atlas migrate status --env dev

# Force ke version yang diketahui baik (hati-hati!)
atlas migrate set 20260107120000 --env dev
```

## 📚 Best Practices

### ✅ DO

- ✅ Review generated migrations sebelum apply
- ✅ Test migrations di development dulu
- ✅ Commit migration files ke version control
- ✅ Gunakan descriptive migration names
- ✅ Keep migrations small dan focused
- ✅ Run `atlas migrate validate` sebelum deploy
- ✅ Backup database sebelum apply di production

### ❌ DON'T

- ❌ Edit migration files yang sudah di-apply
- ❌ Delete migration files
- ❌ Apply untested migrations ke production
- ❌ Skip migration review
- ❌ Gunakan `DB_GEN=true` di production
- ❌ Force version tanpa memahami konsekuensinya

## 🔗 Resources

### Documentation

- **Project Docs**: [docs/DATABASE.md](docs/DATABASE.md)
- **Migration README**: [migrations/README.md](migrations/README.md)
- **Atlas Official**: https://atlasgo.io/
- **GORM Provider**: https://github.com/ariga/atlas-provider-gorm

### Available Commands

```bash
make help              # Lihat semua commands
make atlas-install     # Install Atlas CLI
make atlas-diff        # Generate migration
make atlas-apply       # Apply migrations
make atlas-status      # Check status
make atlas-validate    # Validate migrations
make atlas-inspect     # Inspect database
make atlas-migrate     # Full workflow
```

## 💡 Tips & Tricks

1. **Gunakan Descriptive Names**: 
   ```bash
   atlas migrate diff add_user_avatar --env dev
   # Lebih baik daripada:
   atlas migrate diff update_users --env dev
   ```

2. **Always Review Before Apply**:
   ```bash
   # Generate
   make atlas-diff
   
   # Review file
   cat migrations/[latest-file].sql
   
   # Apply jika OK
   make atlas-apply
   ```

3. **Test Rollback**:
   ```bash
   # Atlas otomatis generate down migrations
   # Test dengan apply lalu rollback
   ```

4. **Use Dry Run**:
   ```bash
   atlas migrate apply --env dev --dry-run
   ```

5. **Monitor Schema**:
   ```bash
   # Inspect schema secara berkala
   make atlas-inspect > schema_$(date +%Y%m%d).sql
   ```

## 🆘 Getting Help

Jika mengalami masalah:

1. Check [Troubleshooting](#troubleshooting) section
2. Review [docs/DATABASE.md](docs/DATABASE.md)
3. Check Atlas docs: https://atlasgo.io/
4. Run: `atlas migrate lint --env dev`

---

**Happy Migrating! 🚀**

Last Updated: January 2026
