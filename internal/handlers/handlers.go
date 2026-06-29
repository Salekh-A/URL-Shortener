package handlers

import (
	"io"
	"net/http"
	storage "newproject/internal/storage"
)

type Handler struct {
	storage *storage.Storage
	baseURL string
}

func New(store *storage.Storage, baseURL string) *Handler {
	return &Handler{
		storage: store,
		baseURL: baseURL,
	}
}
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		if string(body) == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		longURL := string(body)
		id, err := h.storage.Save(longURL)
		if err != nil {
			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		}

		shortURL := h.baseURL + "/" + id
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not found"))
		return
	}

}
