package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Georgy-collab/Black-Box/internal/server"
)

func chdirRepoRoot(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := wd
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			if err := os.Chdir(root); err != nil {
				t.Fatal(err)
			}
			return
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("go.mod not found")
		}
		root = parent
	}
}

func TestUserCRUDFlow(t *testing.T) {
	chdirRepoRoot(t)
	handler := server.New()

	createBody := `{"name":" Alice ","email":"Alice@Example.com","age":25}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body=%s", rr.Code, rr.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created["name"] != "Alice" || created["email"] != "alice@example.com" {
		t.Fatalf("unexpected create response: %#v", created)
	}

	// duplicate email
	req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(createBody))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("duplicate email status = %d", rr.Code)
	}

	// get missing
	req = httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("missing user status = %d", rr.Code)
	}

	// create second and filter
	req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"Bob","email":"bob@example.com","age":40}`))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("second create status = %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/users?min_age=30", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list status = %d", rr.Code)
	}
	var listed []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0]["name"] != "Bob" {
		t.Fatalf("unexpected list: %#v", listed)
	}

	// update conflict
	req = httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString(`{"email":"bob@example.com"}`))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("update conflict status = %d", rr.Code)
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d", rr.Code)
	}
}
