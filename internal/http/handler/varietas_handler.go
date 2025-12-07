package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux" // Contoh router untuk mengambil path parameter

	"github.com/Farewellez/REST-API_VarietasPadi/internal/domain"
	// TIDAK PERLU IMPORT internal/repository lagi!
)

// VarietasHandler menerima VarietasService Interface
type VarietasHandler struct {
	service domain.VarietasService // Menggunakan Interface Service
}

// NewVarietasHandler adalah constructor Handler Layer.
func NewVarietasHandler(service domain.VarietasService) *VarietasHandler {
	return &VarietasHandler{service: service}
}

// Helper function untuk mengirim response JSON
func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// --- FUNGSI HANDLER CRUD ---

// GetAll: GET /varietas
func (h *VarietasHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Panggil Service Layer
	// Kita ambil context dari request untuk diteruskan ke Service
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.DapatkanSemuaData(ctx) // Panggil Service, BUKAN Repository
	if err != nil {
		// Logika Service Error
		respondJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": "Gagal mengambil data: " + err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"total":   len(data),
		"data":    data,
	})
}

// Create: POST /varietas
func (h *VarietasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var varietas domain.VarietasPadi
	// 1. Parsing Request Body
	if err := json.NewDecoder(r.Body).Decode(&varietas); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "Format data JSON tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 2. Panggil Service Layer (Validasi terjadi di Service)
	newVarietas, err := h.service.TambahkanData(ctx, varietas) // Panggil Service Create
	if err != nil {
		// Asumsi error dari Service adalah karena Validasi/Bisnis Logic
		respondJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "Gagal membuat data: " + err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{"success": true, "data": newVarietas})
}

// ReadByID: GET /varietas/{id}
func (h *VarietasHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "ID varietas tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.DapatkanDataByID(ctx, id)
	if err != nil {
		// Cek jika errornya adalah 'not found' (dari Service Layer)
		if err.Error() == "data varietas tidak ditemukan" {
			respondJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": "Gagal mengambil data: " + err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"success": true, "data": data})
}

// Update: PUT /varietas/{id}
func (h *VarietasHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "ID varietas tidak valid"})
		return
	}

	var varietas domain.VarietasPadi
	if err := json.NewDecoder(r.Body).Decode(&varietas); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "Format data JSON tidak valid"})
		return
	}

	// Pastikan ID dari path digunakan, bukan dari body JSON
	varietas.ID = id

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	updatedData, err := h.service.UbahData(ctx, varietas)
	if err != nil {
		// Cek jika errornya adalah not found atau validation
		if err.Error() == "data varietas tidak ditemukan" {
			respondJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": err.Error()})
			return
		}
		respondJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "Gagal mengubah data: " + err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"success": true, "data": updatedData})
}

// Delete: DELETE /varietas/{id}
func (h *VarietasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "ID varietas tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.service.HapusData(ctx, id)
	if err != nil {
		// Cek jika errornya adalah 'not found'
		if err.Error() == "data varietas yang akan dihapus tidak ditemukan" {
			respondJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": "Gagal menghapus data: " + err.Error()})
		return
	}

	// Menggunakan Status 204 No Content untuk operasi DELETE yang sukses tanpa mengembalikan body
	respondJSON(w, http.StatusNoContent, nil)
}
