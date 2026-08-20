package tests

import (
	"net/http/httptest"
	"newproject/internal/handlers"
	"newproject/internal/storage"
	"strings"
	"testing"
)

func TestRootHandlerPost(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost:5432/urlshortener_test?sslmode=disable"
	store, err := storage.New(dsn)
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer store.Close()

	h := handlers.New(store, "http://localhost:8080")

	req := httptest.NewRequest("POST", "/", strings.NewReader("https://google.com"))
	rr := httptest.NewRecorder()
	h.HandleTextShorten(rr, req)
	if rr.Code != 201 {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestRootHandlerEmptyPost(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost:5432/urlshortener_test?sslmode=disable"
	store, err := storage.New(dsn)
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer store.Close()

	h := handlers.New(store, "http://localhost:8080")

	req := httptest.NewRequest("POST", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	h.HandleTextShorten(rr, req)
	if rr.Code != 400 {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRootHandlerGet(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost:5432/urlshortener_test?sslmode=disable"
	store, err := storage.New(dsn)
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer store.Close()

	h := handlers.New(store, "http://localhost:8080")

	postReq := httptest.NewRequest("POST", "/", strings.NewReader("https://google.com"))
	postRr := httptest.NewRecorder()
	h.HandleTextShorten(postRr, postReq)

	shortURL := postRr.Body.String()
	id := shortURL[strings.LastIndex(shortURL, "/")+1:]

	getReq := httptest.NewRequest("GET", "/"+id, nil)
	getRr := httptest.NewRecorder()
	h.HandleGet(getRr, getReq)

	if getRr.Code != 307 {
		t.Errorf("expected 307, got %d", getRr.Code)
	}
}

func TestRootHandlerGetEmpty(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost:5432/urlshortener_test?sslmode=disable"
	store, err := storage.New(dsn)
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer store.Close()

	h := handlers.New(store, "http://localhost:8080")

	getReq := httptest.NewRequest("GET", "/", nil)
	getRr := httptest.NewRecorder()
	h.HandleGet(getRr, getReq)

	if getRr.Code != 400 {
		t.Errorf("expected 400, got %d", getRr.Code)
		return
	}
}
