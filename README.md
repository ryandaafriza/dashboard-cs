# 📊 Dashboard CS

Sistem monitoring dan manajemen tiket Customer Service — dibangun dengan **Go** (backend) dan **React.js** (frontend).

---

## Tech Stack

| Layer     | Teknologi                        |
|-----------|----------------------------------|
| Backend   | Go                               |
| Frontend  | React.js, Vite                   |
| Database  | MySQL 8.0                        |
| Excel     | excelize v2                      |

---

## Fitur

- 📈 Dashboard real-time — SLA, CSAT, tren tiket harian & per jam
- 🚨 Manajemen incident aktif (buat, pantau, resolve)
- 📥 Import data tiket massal dari file `.xlsx`
- 📤 Export laporan multi-sheet ke file `.xlsx`

---

## Menjalankan Proyek

### Backend

```bash
# 1. Salin konfigurasi
cp .env.example .env

# 2. Sesuaikan variabel di .env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=dashboard_cs
PORT=8080

# 3. Jalankan
go mod tidy
go run .
```

### Frontend

```bash
cp .env.example .env
# Set VITE_API_BASE_URL=http://localhost:8080

npm install
npm run dev
```
---
