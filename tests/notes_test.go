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

func registerAndLogin(t *testing.T) string {
	// REGISTER
	registerBody := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	resp := performRequest(t, "POST", "/register", registerBody, "")

	if resp.Code != http.StatusCreated {
		t.Fatalf("Register failed: %d, %s", resp.Code, resp.Body.String())
	}

	// LOGIN
	loginBody := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	resp = performRequest(t, "POST", "/login", loginBody, "")

	if resp.Code != http.StatusOK {
		t.Fatalf("Login failed: %d, %s", resp.Code, resp.Body.String())
	}

	var loginResp LoginResponse
	err := json.Unmarshal(resp.Body.Bytes(), &loginResp)
	if err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}
	log.Printf("[LOGIN] JWT Token: %s\n", loginResp.Token)
	return loginResp.Token
}

func createNote(t *testing.T, token string) uint {
	noteBody := map[string]string{
		"title":   "Encrypted Title",
		"content": "Encrypted Content",
	}
	resp := performRequest(t, "POST", "/notes", noteBody, token)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Create note failed: %d, %s", resp.Code, resp.Body.String())
	}

	var noteResp NoteResponse
	err := json.Unmarshal(resp.Body.Bytes(), &noteResp)
	if err != nil {
		t.Fatalf("Failed to parse create note response: %v", err)
	}
	log.Printf("[CREATE_NOTE] Created Note ID: %d\n", noteResp.ID)
	return noteResp.ID
}

func getNotes(t *testing.T, token string) {
	resp := performRequest(t, "GET", "/notes", nil, token)

	if resp.Code != http.StatusOK {
		t.Fatalf("Get notes failed: %d, %s", resp.Code, resp.Body.String())
	}

	var notes []NoteResponse
	err := json.Unmarshal(resp.Body.Bytes(), &notes)
	if err != nil {
		t.Fatalf("Failed to parse notes response: %v", err)
	}

	if len(notes) != 1 {
		t.Fatalf("Expected 1 note, got: %d", len(notes))
	}
	log.Printf("[GET_NOTES] Found Note ID: %d\n", notes[0].ID)
}

func updateNote(t *testing.T, token string, noteID uint) {
	updateBody := map[string]string{
		"title":   "Updated Title",
		"content": "Updated Content",
	}
	resp := performRequest(t, "PUT", "/notes/"+uintToString(noteID), updateBody, token)

	if resp.Code != http.StatusOK {
		t.Fatalf("Update note failed: %d, %s", resp.Code, resp.Body.String())
	}
	log.Println("[UPDATE_NOTE] Success")
}

func deleteNote(t *testing.T, token string, noteID uint) {
	resp := performRequest(t, "DELETE", "/notes/"+uintToString(noteID), nil, token)

	if resp.Code != http.StatusOK {
		t.Fatalf("Delete note failed: %d, %s", resp.Code, resp.Body.String())
	}
	log.Println("[DELETE_NOTE] Success")
}

// Helper functions

func performRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
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

func uintToString(val uint) string {
	return fmt.Sprintf("%d", val)
}
