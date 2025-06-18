package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leinonen/go-lisp/pkg/evaluator"
	"github.com/leinonen/go-lisp/pkg/registry"
	"github.com/leinonen/go-lisp/pkg/types"
)

// Mock evaluator for testing
type mockEvaluator struct {
	env *evaluator.Environment
}

func newMockEvaluator() *mockEvaluator {
	return &mockEvaluator{
		env: evaluator.NewEnvironment(),
	}
}

func (me *mockEvaluator) Eval(expr types.Expr) (types.Value, error) {
	switch e := expr.(type) {
	case *types.NumberExpr:
		return types.NumberValue(e.Value), nil
	case *types.StringExpr:
		return types.StringValue(e.Value), nil
	case *types.BooleanExpr:
		return types.BooleanValue(e.Value), nil
	default:
		if ve, ok := expr.(valueExpr); ok {
			return ve.value, nil
		}
		if val, ok := expr.(types.Value); ok {
			return val, nil
		}
		return nil, nil
	}
}

func (me *mockEvaluator) CallFunction(funcValue types.Value, args []types.Expr) (types.Value, error) {
	return nil, nil // Not needed for HTTP tests
}

func (me *mockEvaluator) EvalWithBindings(expr types.Expr, bindings map[string]types.Value) (types.Value, error) {
	// For testing purposes, just call regular Eval
	// In a real implementation, this would use the bindings
	return me.Eval(expr)
}

func wrapValue(value types.Value) types.Expr {
	return valueExpr{value}
}

type valueExpr struct {
	value types.Value
}

func (ve valueExpr) String() string {
	return ve.value.String()
}

func (ve valueExpr) GetPosition() types.Position {
	return types.Position{Line: 1, Column: 1}
}

func TestHTTPPlugin_RegisterFunctions(t *testing.T) {
	plugin := NewHTTPPlugin()
	reg := registry.NewRegistry()

	err := plugin.RegisterFunctions(reg)
	if err != nil {
		t.Fatalf("Failed to register functions: %v", err)
	}

	expectedFunctions := []string{"http-get", "http-post", "http-put", "http-delete"}

	for _, fnName := range expectedFunctions {
		if !reg.Has(fnName) {
			t.Errorf("Function %s was not registered", fnName)
		}
	}
}

func TestHTTPPlugin_HttpGet(t *testing.T) {
	plugin := NewHTTPPlugin()
	evaluator := newMockEvaluator()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()

	// Test GET request
	args := []types.Expr{&types.StringExpr{Value: server.URL}}
	result, err := plugin.evalHttpGet(evaluator, args)
	if err != nil {
		t.Fatalf("evalHttpGet failed: %v", err)
	}

	// Result should be a hash map with :status, :headers, and :body
	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected hash map, got %T", result)
	}

	// Check status
	if status, exists := hashMap.Elements[":status"]; exists {
		if status.String() != "200" {
			t.Errorf("Expected status 200, got %s", status.String())
		}
	} else {
		t.Error("Response should contain :status field")
	}

	// Check body
	if body, exists := hashMap.Elements[":body"]; exists {
		expectedBody := `{"message": "Hello, World!"}`
		if body.String() != expectedBody {
			t.Errorf("Expected body %s, got %s", expectedBody, body.String())
		}
	} else {
		t.Error("Response should contain :body field")
	}
}

func TestHTTPPlugin_HttpPost(t *testing.T) {
	plugin := NewHTTPPlugin()
	evaluator := newMockEvaluator()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}))
	defer server.Close()

	// Test POST request
	args := []types.Expr{
		&types.StringExpr{Value: server.URL},
		&types.StringExpr{Value: "test data"},
	}
	result, err := plugin.evalHttpPost(evaluator, args)
	if err != nil {
		t.Fatalf("evalHttpPost failed: %v", err)
	}

	// Result should be a hash map
	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected hash map, got %T", result)
	}

	// Check status
	if status, exists := hashMap.Elements[":status"]; exists {
		if status.String() != "201" {
			t.Errorf("Expected status 201, got %s", status.String())
		}
	} else {
		t.Error("Response should contain :status field")
	}
}

func TestHTTPPlugin_HttpPut(t *testing.T) {
	plugin := NewHTTPPlugin()
	evaluator := newMockEvaluator()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Updated"))
	}))
	defer server.Close()

	// Test PUT request
	args := []types.Expr{
		&types.StringExpr{Value: server.URL},
		&types.StringExpr{Value: "updated data"},
	}
	result, err := plugin.evalHttpPut(evaluator, args)
	if err != nil {
		t.Fatalf("evalHttpPut failed: %v", err)
	}

	// Result should be a hash map
	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected hash map, got %T", result)
	}

	// Check status
	if status, exists := hashMap.Elements[":status"]; exists {
		if status.String() != "200" {
			t.Errorf("Expected status 200, got %s", status.String())
		}
	} else {
		t.Error("Response should contain :status field")
	}
}

func TestHTTPPlugin_HttpDelete(t *testing.T) {
	plugin := NewHTTPPlugin()
	evaluator := newMockEvaluator()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Test DELETE request
	args := []types.Expr{&types.StringExpr{Value: server.URL}}
	result, err := plugin.evalHttpDelete(evaluator, args)
	if err != nil {
		t.Fatalf("evalHttpDelete failed: %v", err)
	}

	// Result should be a hash map
	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected hash map, got %T", result)
	}

	// Check status
	if status, exists := hashMap.Elements[":status"]; exists {
		if status.String() != "204" {
			t.Errorf("Expected status 204, got %s", status.String())
		}
	} else {
		t.Error("Response should contain :status field")
	}
}

func TestHTTPPlugin_WithHeaders(t *testing.T) {
	plugin := NewHTTPPlugin()
	evaluator := newMockEvaluator()

	// Create a test server that checks headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			http.Error(w, "Missing or incorrect header", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Create headers hash map
	headers := &types.HashMapValue{Elements: make(map[string]types.Value)}
	headers.Elements["X-Custom-Header"] = types.StringValue("test-value")

	// Test POST request with headers
	args := []types.Expr{
		&types.StringExpr{Value: server.URL},
		&types.StringExpr{Value: "test data"},
		wrapValue(headers),
	}
	result, err := plugin.evalHttpPost(evaluator, args)
	if err != nil {
		t.Fatalf("evalHttpPost with headers failed: %v", err)
	}

	// Result should be a hash map
	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected hash map, got %T", result)
	}

	// Check status
	if status, exists := hashMap.Elements[":status"]; exists {
		if status.String() != "200" {
			t.Errorf("Expected status 200, got %s", status.String())
		}
	} else {
		t.Error("Response should contain :status field")
	}
}

func TestHTTPPlugin_ErrorCases(t *testing.T) {
	plugin := NewHTTPPlugin()
	evaluator := newMockEvaluator()

	// Test with invalid URL
	args := []types.Expr{&types.StringExpr{Value: "invalid-url"}}
	_, err := plugin.evalHttpGet(evaluator, args)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	// Test with wrong argument count
	_, err = plugin.evalHttpGet(evaluator, []types.Expr{})
	if err == nil {
		t.Error("Expected error for missing arguments")
	}

	// Test with wrong argument type
	args = []types.Expr{&types.NumberExpr{Value: 42}}
	_, err = plugin.evalHttpGet(evaluator, args)
	if err == nil {
		t.Error("Expected error for non-string URL")
	}
}

func TestHTTPPlugin_ResponseStructure(t *testing.T) {
	plugin := NewHTTPPlugin()
	evaluator := newMockEvaluator()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": true}`))
	}))
	defer server.Close()

	// Test GET request
	args := []types.Expr{&types.StringExpr{Value: server.URL}}
	result, err := plugin.evalHttpGet(evaluator, args)
	if err != nil {
		t.Fatalf("evalHttpGet failed: %v", err)
	}

	// Result should be a hash map with expected fields
	hashMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("Expected hash map, got %T", result)
	}

	// Check that all expected fields are present
	expectedFields := []string{":status", ":headers", ":body"}
	for _, field := range expectedFields {
		if _, exists := hashMap.Elements[field]; !exists {
			t.Errorf("Response should contain %s field", field)
		}
	}

	// Check headers structure
	if headersValue, exists := hashMap.Elements[":headers"]; exists {
		if headersMap, ok := headersValue.(*types.HashMapValue); ok {
			if contentType, exists := headersMap.Elements["Content-Type"]; exists {
				expectedContentType := "application/json"
				if contentType.String() != expectedContentType {
					t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType.String())
				}
			} else {
				t.Error("Headers should contain Content-Type")
			}
		} else {
			t.Error("Headers should be a hash map")
		}
	}
}
