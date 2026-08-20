package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	storage "newproject/internal/storage"
	"strings"
)

type Handler struct {
	storage *storage.Storage
	baseURL string
}

type ShortenResponse struct {
	ShortURL string `json:"result"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

func New(store *storage.Storage, baseURL string) *Handler {
	return &Handler{
		storage: store,
		baseURL: baseURL,
	}
}
func (h *Handler) HandleAPIShorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ShortenRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := h.storage.Save(ctx, req.URL)
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{
		ShortURL: h.baseURL + "/" + id,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h *Handler) HandleTextShorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	if string(body) == "" {
		http.Error(w, "Empty request", http.StatusBadRequest)
		return
	}

	id, err := h.storage.Save(ctx, string(body))
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	shortURL := h.baseURL + "/" + id
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Empty id", http.StatusBadRequest)
		return
	}

	longURL, err := h.storage.Load(ctx, id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
