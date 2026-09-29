package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"newproject/internal/ports"
)

type Handler struct {
	service ports.Service
	baseURL string
}

type ShortenResponse struct {
	ShortURL string `json:"result"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func New(service ports.Service, baseURL string) *Handler {
	return &Handler{
		service: service,
		baseURL: baseURL,
	}
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return parsedURL.Host != ""
}

func (h *Handler) HandleAPIShorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ShortenRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !isValidURL(req.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(ctx, req.URL)
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{
		ShortURL: h.baseURL + "/" + id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) HandleTextShorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	rawURL := string(body)

	if !isValidURL(rawURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(ctx, rawURL)
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	shortURL := h.baseURL + "/" + id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(shortURL))
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Empty id", http.StatusBadRequest)
		return
	}

	longURL, err := h.service.Get(ctx, id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) BatchShorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	var reqs []BatchRequest

	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(reqs) == 0 {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	urls := make([]string, 0, len(reqs))

	for _, req := range reqs {
		if !isValidURL(req.OriginalURL) {
			http.Error(w, "Invalid URL in request", http.StatusBadRequest)
			return
		}

		urls = append(urls, req.OriginalURL)
	}

	ids, err := h.service.CreateBatch(ctx, urls)
	if err != nil {
		http.Error(w, "Failed to save URLs", http.StatusInternalServerError)
		return
	}

	resp := make([]BatchResponse, 0, len(reqs))

	for i, req := range reqs {
		resp = append(resp, BatchResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      h.baseURL + "/" + ids[i],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(resp)
}
