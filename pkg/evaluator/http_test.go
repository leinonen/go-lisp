package evaluator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

func TestHttpGet(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()

	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test http-get
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "http-get"},
			&types.StringExpr{Value: server.URL},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("http-get failed: %v", err)
	}

	responseMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("expected HashMapValue, got %T", result)
	}

	// Check status
	status, exists := responseMap.Elements[":status"]
	if !exists {
		t.Fatal("response missing status field")
	}
	if status != types.NumberValue(200) {
		t.Errorf("expected status 200, got %v", status)
	}

	// Check body
	body, exists := responseMap.Elements[":body"]
	if !exists {
		t.Fatal("response missing body field")
	}
	expectedBody := `{"message": "Hello, World!"}`
	if body != types.StringValue(expectedBody) {
		t.Errorf("expected body %s, got %v", expectedBody, body)
	}

	// Check headers
	headers, exists := responseMap.Elements[":headers"]
	if !exists {
		t.Fatal("response missing headers field")
	}
	headersMap, ok := headers.(*types.HashMapValue)
	if !ok {
		t.Fatalf("expected headers to be HashMapValue, got %T", headers)
	}
	contentType, exists := headersMap.Elements["Content-Type"]
	if !exists {
		t.Fatal("response headers missing Content-Type")
	}
	if contentType != types.StringValue("application/json") {
		t.Errorf("expected Content-Type application/json, got %v", contentType)
	}
}

func TestHttpPost(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Echo back the request body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	}))
	defer server.Close()

	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test http-post
	requestBody := `{"name": "test"}`
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "http-post"},
			&types.StringExpr{Value: server.URL},
			&types.StringExpr{Value: requestBody},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("http-post failed: %v", err)
	}

	responseMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("expected HashMapValue, got %T", result)
	}

	// Check status
	status, exists := responseMap.Elements[":status"]
	if !exists {
		t.Fatal("response missing status field")
	}
	if status != types.NumberValue(201) {
		t.Errorf("expected status 201, got %v", status)
	}

	// Check body (should echo back our request)
	body, exists := responseMap.Elements[":body"]
	if !exists {
		t.Fatal("response missing body field")
	}
	if body != types.StringValue(requestBody) {
		t.Errorf("expected body %s, got %v", requestBody, body)
	}
}

func TestHttpPostWithHeaders(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check custom header
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header custom-value, got %s", r.Header.Get("X-Custom-Header"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Create headers hash map
	headersExpr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "hash-map"},
			&types.StringExpr{Value: "X-Custom-Header"},
			&types.StringExpr{Value: "custom-value"},
			&types.StringExpr{Value: "Content-Type"},
			&types.StringExpr{Value: "text/plain"},
		},
	}

	// Test http-post with headers
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "http-post"},
			&types.StringExpr{Value: server.URL},
			&types.StringExpr{Value: "test body"},
			headersExpr,
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("http-post with headers failed: %v", err)
	}

	responseMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("expected HashMapValue, got %T", result)
	}

	// Check status
	status, exists := responseMap.Elements[":status"]
	if !exists {
		t.Fatal("response missing status field")
	}
	if status != types.NumberValue(200) {
		t.Errorf("expected status 200, got %v", status)
	}
}

func TestHttpPut(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Updated"))
	}))
	defer server.Close()

	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test http-put
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "http-put"},
			&types.StringExpr{Value: server.URL},
			&types.StringExpr{Value: `{"id": 1, "name": "updated"}`},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("http-put failed: %v", err)
	}

	responseMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("expected HashMapValue, got %T", result)
	}

	// Check status
	status, exists := responseMap.Elements[":status"]
	if !exists {
		t.Fatal("response missing status field")
	}
	if status != types.NumberValue(200) {
		t.Errorf("expected status 200, got %v", status)
	}
}

func TestHttpDelete(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	// Test http-delete
	expr := &types.ListExpr{
		Elements: []types.Expr{
			&types.SymbolExpr{Name: "http-delete"},
			&types.StringExpr{Value: server.URL},
		},
	}

	result, err := evaluator.Eval(expr)
	if err != nil {
		t.Fatalf("http-delete failed: %v", err)
	}

	responseMap, ok := result.(*types.HashMapValue)
	if !ok {
		t.Fatalf("expected HashMapValue, got %T", result)
	}

	// Check status
	status, exists := responseMap.Elements[":status"]
	if !exists {
		t.Fatal("response missing status field")
	}
	if status != types.NumberValue(204) {
		t.Errorf("expected status 204, got %v", status)
	}
}

func TestHttpErrors(t *testing.T) {
	env := NewEnvironment()
	evaluator := NewEvaluator(env)

	tests := []struct {
		name string
		expr types.Expr
	}{
		{
			name: "http-get with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "http-get"},
				},
			},
		},
		{
			name: "http-get with non-string URL",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "http-get"},
					&types.NumberExpr{Value: 123},
				},
			},
		},
		{
			name: "http-post with wrong number of arguments",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "http-post"},
				},
			},
		},
		{
			name: "http-post with non-string body",
			expr: &types.ListExpr{
				Elements: []types.Expr{
					&types.SymbolExpr{Name: "http-post"},
					&types.StringExpr{Value: "http://example.com"},
					&types.NumberExpr{Value: 123},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluator.Eval(tt.expr)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
