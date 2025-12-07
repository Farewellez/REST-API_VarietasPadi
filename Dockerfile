# ==========================================================
# STAGE 1: BUILDER (MENGKOMPILASI APLIKASI GO)
# ==========================================================
# Kita gunakan image golang yang full-featured untuk kompilasi
FROM golang:1.25.4-alpine AS builder

# Atur Environment Variable untuk kompilasi statis (penting untuk container minimal)
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set working directory di dalam container builder
WORKDIR /app

# Copy go.mod dan go.sum untuk mendownload dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
# Perintah ini akan menggunakan cache, yang mempercepat build berikutnya
RUN go mod download

# Copy semua kode sumber ke dalam container
COPY . .

# Kompilasi aplikasi. Output binary ditaruh di /app/main
# Kita kompilasi dari cmd/server/main.go
RUN go build -ldflags "-s -w" -o /app/main ./cmd/server

# ==========================================================
# STAGE 2: RUNNER (IMAGE RUMTIME MINIMAL)
# ==========================================================
# Kita gunakan image yang sangat kecil (Scratch, Alpine) atau image dasar yang minimal.
# Alpine ideal karena kecil dan memiliki dasar shell/library yang cukup.
FROM alpine:latest

# Tambahkan sertifikat SSL/TLS root (penting untuk koneksi keluar, misal ke NeonDB)
# Ini memastikan koneksi HTTPS/SSL ke database eksternal berfungsi
RUN apk add --no-cache ca-certificates

# Set user non-root untuk keamanan (best practice)
# RUN adduser -D nonroot
# USER nonroot 

# Set working directory
WORKDIR /root/

# Copy binary yang sudah dikompilasi dari STAGE 1 ke STAGE 2
COPY --from=builder /app/main .

COPY views ./views

# Expose port yang digunakan server Go (default: 8080)
EXPOSE 8080

# Command utama yang dijalankan saat container di-start
# Program akan dijalankan dengan binary /root/main
CMD ["./main"]