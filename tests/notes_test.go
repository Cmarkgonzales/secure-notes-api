package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Response structs
type LoginResponse struct {
	Token string `json:"token"`
}

type NoteResponse struct {
	ID      uint   `json:"ID"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Main test workflow
func TestFullWorkflow(t *testing.T) {
	InitTestSuite()
	ClearDB()

	token := registerAndLogin(t)
	noteID := createNote(t, token)
	getNotes(t, token)
	updateNote(t, token, noteID)
	deleteNote(t, token, noteID)
}

// REGISTER and LOGIN
func registerAndLogin(t *testing.T) string {
	// REGISTER
	registerBody := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	resp := performRequest(t, "POST", "/api/v1/register", registerBody, "")

	assertStatus(t, resp, http.StatusCreated, "Register failed")

	// LOGIN
	loginBody := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	resp = performRequest(t, "POST", "/api/v1/login", loginBody, "")

	assertStatus(t, resp, http.StatusOK, "Login failed")

	var loginResp LoginResponse
	decodeJSON(t, resp.Body.Bytes(), &loginResp)
	log.Printf("[LOGIN] JWT Token: %s\n", loginResp.Token)

	return loginResp.Token
}

// CREATE NOTE
func createNote(t *testing.T, token string) uint {
	noteBody := map[string]string{
		"title":   "Encrypted Title",
		"content": "Encrypted Content",
	}
	resp := performRequest(t, "POST", "/api/v1/notes", noteBody, token)

	assertStatus(t, resp, http.StatusCreated, "Create note failed")

	var noteResp NoteResponse
	decodeJSON(t, resp.Body.Bytes(), &noteResp)
	log.Printf("[CREATE_NOTE] Created Note ID: %d\n", noteResp.ID)
	return noteResp.ID
}

// GET NOTES
func getNotes(t *testing.T, token string) {
	resp := performRequest(t, "GET", "/api/v1/notes", nil, token)

	assertStatus(t, resp, http.StatusOK, "Get notes failed")

	var notes []NoteResponse
	decodeJSON(t, resp.Body.Bytes(), &notes)

	if len(notes) != 1 {
		t.Fatalf("Expected 1 note, got: %d", len(notes))
	}
	log.Printf("[GET_NOTES] Found Note ID: %d\n", notes[0].ID)
}

// UPDATE NOTE
func updateNote(t *testing.T, token string, noteID uint) {
	updateBody := map[string]string{
		"title":   "Updated Title",
		"content": "Updated Content",
	}
	path := fmt.Sprintf("/api/v1/notes/%d", noteID)
	resp := performRequest(t, "PUT", path, updateBody, token)

	assertStatus(t, resp, http.StatusOK, "Update note failed")
	log.Println("[UPDATE_NOTE] Success")
}

// DELETE NOTE
func deleteNote(t *testing.T, token string, noteID uint) {
	path := fmt.Sprintf("/api/v1/notes/%d", noteID)
	resp := performRequest(t, "DELETE", path, nil, token)

	assertStatus(t, resp, http.StatusOK, "Delete note failed")
	log.Println("[DELETE_NOTE] Success")
}

// === Helper Functions ===

// Perform API request
func performRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("Failed to encode body: %v", err)
		}
	}

	req, err := http.NewRequest(method, path, &buf)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)
	return w
}

// Assert HTTP status code
func assertStatus(t *testing.T, resp *httptest.ResponseRecorder, expected int, message string) {
	if resp.Code != expected {
		t.Fatalf("%s: expected %d, got %d, body: %s", message, expected, resp.Code, resp.Body.String())
	}
}

// Decode JSON response
func decodeJSON(t *testing.T, data []byte, v interface{}) {
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}
}
