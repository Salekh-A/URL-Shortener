package handlers

import (
	"encoding/json"
	"net/http"
	storage "newproject/internal/storage"
)

type Handler struct {
	storage *storage.Storage
	baseURL string
}

type ShortRequest struct {
	URL string `json:"url"`
}

type ShortResponse struct {
	ShortUrl string `json:"short_url"`
}

func New(store *storage.Storage, baseURL string) *Handler {
	return &Handler{
		storage: store,
		baseURL: baseURL,
	}
}
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var req ShortRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		longURL := req.URL
		id, err := h.storage.Save(longURL)
		if err != nil {
			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
			return
		}

		shortURL := ShortResponse{
			ShortUrl: h.baseURL + "/" + id,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(shortURL)
		return
	}

	if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if longURL, ok := h.storage.Load(id); ok {
			w.Header().Set("Location", longURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed"))
}
