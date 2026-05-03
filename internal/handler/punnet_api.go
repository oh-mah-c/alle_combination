package handler

import (
	"alle_combination/internal/genetics"
	"encoding/json"
	"net/http"
	"strings"
)

type CalculateRequest struct {
	Parent1 string `json:"parent_1"`
	Parent2 string `json:"parent_2"`
}

func CalculatePunnet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Please use POST", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload format", http.StatusBadRequest)
	}

	if len(req.Parent1) != len(req.Parent2) {
		http.Error(w, "P1 and P2 must have the same length", http.StatusBadRequest)
	}

	if len(req.Parent1)%2 != 0 {
		http.Error(w, "Gen must be even (Ex: AaBb)", http.StatusBadRequest)
	}

	result := genetics.CalculatePolyHybrid(req.Parent1, req.Parent2)

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

type LinkedRequest struct {
	Parent1 string  `json:"parent_1"`
	Parent2 string  `json:"parent_2"`
	F       float64 `json:"f"`
}

func CalculateLinkedPunnet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Please use POST", http.StatusMethodNotAllowed)
		return
	}

	var req LinkedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload format", http.StatusBadRequest)
	}

	if !strings.Contains(req.Parent1, "/") || !strings.Contains(req.Parent2, "/") {
		http.Error(w, "INvalid phase Format, Must use '/' (ex: AB/ab)", http.StatusBadRequest)
		return
	}

	if req.F < 0.0 || req.F > 0.5 {
		http.Error(w, "Recombination frequency (F) must be between 0 - 0.5", http.StatusBadRequest)
		return
	}

	result := genetics.CalculateLinkedCross(req.Parent1, req.Parent2, req.F)

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
