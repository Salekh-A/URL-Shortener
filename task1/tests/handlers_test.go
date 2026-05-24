package tests

import (
	"net/http/httptest"
	"newproject/internal/handlers"
	"newproject/internal/storage"
	"strings"
	"testing"
)

func TestRootHandlerPost(t *testing.T) {
	store := storage.New()
	h := handlers.New(store, "http://localhost:8080")

	req := httptest.NewRequest("POST", "/", strings.NewReader("https://google.com"))
	rr := httptest.NewRecorder()
	h.Root(rr, req)
	if rr.Code != 201 {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestRootHandlerEmptyPost(t *testing.T) {
	store := storage.New()
	h := handlers.New(store, "http://localhost:8080")

	req := httptest.NewRequest("POST", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	h.Root(rr, req)
	if rr.Code != 400 {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRootHandlerGet(t *testing.T) {
	store := storage.New()
	h := handlers.New(store, "http://localhost:8080")

	postReq := httptest.NewRequest("POST", "/", strings.NewReader("https://google.com"))
	postRr := httptest.NewRecorder()
	h.Root(postRr, postReq)

	getReq := httptest.NewRequest("GET", "/0", nil)
	getRr := httptest.NewRecorder()
	h.Root(getRr, getReq)

	if getRr.Code != 307 {
		t.Errorf("expected 307, got %d", getRr.Code)
	}
}

func TestRootHandlerGetEmpty(t *testing.T) {
	store := storage.New()
	h := handlers.New(store, "http://localhost:8080")

	getReq := httptest.NewRequest("GET", "/", nil)
	getRr := httptest.NewRecorder()
	h.Root(getRr, getReq)

	if getRr.Code != 400 {
		t.Errorf("expected 400, got %d", getRr.Code)
		return
	}
}
