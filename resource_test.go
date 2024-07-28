package resource_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iwanhae/resource"
)

// MockResource is a simple struct that implements the Validator interface
type MockResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (m MockResource) ValidateCreate(ctx resource.Context) error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func (m MockResource) ValidateUpdate(ctx resource.Context, id string) error {
	return m.ValidateCreate(ctx)
}

// Mock functions for CRUD operations
func mockList(ctx resource.Context, offset int, limit int) ([]MockResource, error) {
	return []MockResource{{ID: "1", Name: "Test"}}, nil
}

func mockCreate(ctx resource.Context, resource MockResource) (MockResource, error) {
	resource.ID = "2"
	return resource, nil
}

func mockGet(ctx resource.Context, id string) (MockResource, error) {
	return MockResource{ID: id, Name: "Test"}, nil
}

func mockUpdate(ctx resource.Context, id string, resource MockResource) (MockResource, error) {
	resource.ID = id
	return resource, nil
}

func mockDelete(ctx resource.Context, id string) error {
	return nil
}

func TestResourceRegistration(t *testing.T) {
	r := resource.New[MockResource]().
		Name("mock").
		Plural("mocks").
		List(mockList).
		Create(mockCreate).
		Get(mockGet).
		Update(mockUpdate).
		Delete(mockDelete)

	mux := http.NewServeMux()
	r.RegisterMux(mux)

	// Test each endpoint
	testCases := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{"List", "GET", "/mocks", "", http.StatusOK},
		{"Create", "POST", "/mocks", `{"name":"New Mock"}`, http.StatusCreated},
		{"Get", "GET", "/mocks/1", "", http.StatusOK},
		{"Update", "PUT", "/mocks/1", `{"name":"Updated Mock"}`, http.StatusOK},
		{"Delete", "DELETE", "/mocks/1", "", http.StatusNoContent},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedStatus)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	r := resource.New[MockResource]().
		Name("mock").
		Plural("mocks").
		Create(mockCreate)

	mux := http.NewServeMux()
	r.RegisterMux(mux)

	// Test validation failure
	req := httptest.NewRequest("POST", "/mocks", strings.NewReader(`{"name":""}`))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}

	var errResp resource.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	if err != nil {
		t.Fatalf("Could not decode error response: %v", err)
	}

	if errResp.Message != "name is required" {
		t.Errorf("Unexpected error message: got %v want %v",
			errResp.Message, "name is required")
	}
}
