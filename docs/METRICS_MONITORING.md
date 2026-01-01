# Lightweight Metrics & Monitoring

## Overview
Sistem monitoring yang **super ringan** menggunakan built-in Go packages tanpa dependencies eksternal seperti Prometheus.

## Features

### ✅ Zero External Dependencies
- Tidak menggunakan Prometheus (menghemat ~5.2MB)
- Tidak menggunakan Fiber Monitor (menghemat ~1.6MB)
- Hanya menggunakan Go stdlib: `sync/atomic` dan `runtime`

### 📊 Metrics Yang Tersedia

#### 1. HTTP Metrics
- **Total Requests**: Jumlah total request yang diterima
- **Success Requests**: Request dengan status < 400
- **Error Requests**: Request dengan status >= 400
- **In-Flight Requests**: Request yang sedang diproses
- **Success Rate**: Persentase request yang berhasil
- **Average Latency**: Rata-rata waktu response (ms)

#### 2. System Metrics
- **Memory Allocation**: Memory yang dialokasikan saat ini (MB)
- **Total Memory**: Total memory yang pernah dialokasikan (MB)
- **System Memory**: Memory yang direserve dari OS (MB)
- **GC Runs**: Jumlah garbage collection yang sudah dijalankan
- **Goroutines**: Jumlah goroutine yang aktif

#### 3. Performance Metrics
- **Uptime**: Waktu server berjalan (detik)
- **Total Duration**: Total waktu pemrosesan semua request (detik)

## Usage

### 1. Endpoint
```bash
GET /metrics
```

### 2. Response Example
```json
{
  "uptime_seconds": 123.45,
  "requests": {
    "total": 1000,
    "success": 950,
    "errors": 50,
    "in_flight": 3,
    "success_rate": 95.0
  },
  "performance": {
    "avg_latency_ms": 45.2,
    "total_duration_seconds": 45.2
  },
  "system": {
    "memory_alloc_mb": 12.5,
    "memory_total_mb": 50.3,
    "memory_sys_mb": 25.7,
    "gc_runs": 15,
    "goroutines": 25
  }
}
```

### 3. cURL Example
```bash
curl http://localhost:3000/metrics | jq .
```

## Implementation Details

### Thread-Safe Counters
Menggunakan `sync/atomic` untuk operasi thread-safe tanpa mutex:

```go
atomic.AddUint64(&totalRequests, 1)
atomic.LoadUint64(&totalRequests)
```

### Middleware Integration
Metrics otomatis direcord oleh `MetricsMiddleware()`:

```go
// config/app.go
app.Use(middleware.MetricsMiddleware())
```

### Skipped Paths
Endpoint berikut **tidak** tercatat dalam metrics:
- `/metrics` - Endpoint metrics itu sendiri
- `/health` - Health check endpoint

## Advantages

### 🚀 Performance
- **No External Libraries**: Tidak ada overhead dari Prometheus client
- **Atomic Operations**: Fast & lock-free operations
- **Minimal Memory**: Hanya menggunakan beberapa uint64 counters

### 📦 Size Comparison
| Solution | Package Size | Dependencies |
|----------|-------------|--------------|
| Prometheus | ~5.2 MB | 88 packages |
| Fiber Monitor | ~1.6 MB | 58 packages (requires Fiber v3) |
| **Custom Metrics** | **0 MB** | **0 packages** ✅ |

### 🔧 Maintenance
- Tidak ada dependency yang perlu diupdate
- Tidak ada breaking changes dari external packages
- Code yang simple dan mudah dimodifikasi

## Monitoring

### Production Best Practices

1. **Export to External System**
   Untuk production, export metrics ke external monitoring:
   ```bash
   # Poll metrics setiap 30 detik
   while true; do
     curl -s http://localhost:3000/metrics >> /var/log/app-metrics.log
     sleep 30
   done
   ```

2. **Integration dengan Grafana/Prometheus**
   Jika diperlukan, bisa tambahkan endpoint Prometheus format:
   ```go
   // Konversi ke Prometheus text format
   func PrometheusFormat(c *fiber.Ctx) error {
     metrics := middleware.GetMetrics()
     // Format ke Prometheus exposition format
     return c.SendString(formatToPrometheus(metrics))
   }
   ```

3. **Alerting**
   Setup alerting berdasarkan threshold:
   ```bash
   # Check error rate
   ERROR_RATE=$(curl -s http://localhost:3000/metrics | jq '.requests.errors / .requests.total * 100')
   if (( $(echo "$ERROR_RATE > 10" | bc -l) )); then
     # Send alert
   fi
   ```

## Future Enhancements

Jika dibutuhkan metrics tambahan, bisa ditambahkan di:
- `helper/metrics.go` - Custom business metrics (auth, DB, dll)
- `middleware/prometheus.go` - HTTP metrics yang sudah ada

Semua menggunakan atomic operations untuk thread-safety tanpa dependencies tambahan!

## FAQ

**Q: Kenapa tidak pakai Prometheus?**
A: Prometheus client library berat (~5.2MB) dan menambah 88 dependencies. Untuk aplikasi starter, metrics sederhana sudah cukup.

**Q: Kenapa tidak pakai Fiber Monitor?**
A: Fiber Monitor membutuhkan Fiber v3 (masih RC/beta) dan menambah 58 dependencies. Lebih baik tunggu Fiber v3 stable.

**Q: Apakah thread-safe?**
A: Ya! Menggunakan `sync/atomic` yang thread-safe tanpa perlu mutex.

**Q: Bagaimana cara add custom metrics?**
A: Tambahkan atomic counter di `middleware/prometheus.go` atau `helper/metrics.go`:
```go
var myCounter uint64
atomic.AddUint64(&myCounter, 1)
```

**Q: Bisa export ke Prometheus?**
A: Bisa! Tinggal tambahkan endpoint yang convert JSON ke Prometheus exposition format.

---

**Last Updated**: 2026-01-01
**Version**: 1.0.0 (Lightweight Metrics)
