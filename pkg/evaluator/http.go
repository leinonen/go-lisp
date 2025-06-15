// Package evaluator_http contains HTTP request functionality
package evaluator

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// HTTP client with reasonable timeout
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// evalHttpGet performs an HTTP GET request
func (e *Evaluator) evalHttpGet(args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("http-get requires exactly 1 argument (URL), got %d", len(args))
	}

	// Evaluate the URL
	urlValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlStr, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-get URL must be a string, got %T", urlValue)
	}

	// Make the HTTP GET request
	resp, err := httpClient.Get(string(urlStr))
	if err != nil {
		return nil, fmt.Errorf("http-get request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http-get failed to read response body: %v", err)
	}

	// Create response hash map
	responseMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			":status":      types.NumberValue(resp.StatusCode),
			":status-text": types.StringValue(resp.Status),
			":body":        types.StringValue(string(body)),
			":headers":     createHeadersHashMap(resp.Header),
		},
	}

	return responseMap, nil
}

// evalHttpPost performs an HTTP POST request
func (e *Evaluator) evalHttpPost(args []types.Expr) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("http-post requires 2-3 arguments (URL, body, optional headers), got %d", len(args))
	}

	// Evaluate the URL
	urlValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlStr, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-post URL must be a string, got %T", urlValue)
	}

	// Evaluate the body
	bodyValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	bodyStr, ok := bodyValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-post body must be a string, got %T", bodyValue)
	}

	// Create the request
	req, err := http.NewRequest("POST", string(urlStr), bytes.NewBufferString(string(bodyStr)))
	if err != nil {
		return nil, fmt.Errorf("http-post failed to create request: %v", err)
	}

	// Set default content type
	req.Header.Set("Content-Type", "application/json")

	// Handle optional headers
	if len(args) == 3 {
		headersValue, err := e.Eval(args[2])
		if err != nil {
			return nil, err
		}

		headersMap, ok := headersValue.(*types.HashMapValue)
		if !ok {
			return nil, fmt.Errorf("http-post headers must be a hash map, got %T", headersValue)
		}

		// Set headers from the hash map
		for key, value := range headersMap.Elements {
			if strValue, ok := value.(types.StringValue); ok {
				req.Header.Set(key, string(strValue))
			}
		}
	}

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http-post request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http-post failed to read response body: %v", err)
	}

	// Create response hash map
	responseMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			":status":      types.NumberValue(resp.StatusCode),
			":status-text": types.StringValue(resp.Status),
			":body":        types.StringValue(string(body)),
			":headers":     createHeadersHashMap(resp.Header),
		},
	}

	return responseMap, nil
}

// evalHttpPut performs an HTTP PUT request
func (e *Evaluator) evalHttpPut(args []types.Expr) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("http-put requires 2-3 arguments (URL, body, optional headers), got %d", len(args))
	}

	// Evaluate the URL
	urlValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlStr, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-put URL must be a string, got %T", urlValue)
	}

	// Evaluate the body
	bodyValue, err := e.Eval(args[1])
	if err != nil {
		return nil, err
	}

	bodyStr, ok := bodyValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-put body must be a string, got %T", bodyValue)
	}

	// Create the request
	req, err := http.NewRequest("PUT", string(urlStr), bytes.NewBufferString(string(bodyStr)))
	if err != nil {
		return nil, fmt.Errorf("http-put failed to create request: %v", err)
	}

	// Set default content type
	req.Header.Set("Content-Type", "application/json")

	// Handle optional headers
	if len(args) == 3 {
		headersValue, err := e.Eval(args[2])
		if err != nil {
			return nil, err
		}

		headersMap, ok := headersValue.(*types.HashMapValue)
		if !ok {
			return nil, fmt.Errorf("http-put headers must be a hash map, got %T", headersValue)
		}

		// Set headers from the hash map
		for key, value := range headersMap.Elements {
			if strValue, ok := value.(types.StringValue); ok {
				req.Header.Set(key, string(strValue))
			}
		}
	}

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http-put request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http-put failed to read response body: %v", err)
	}

	// Create response hash map
	responseMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			":status":      types.NumberValue(resp.StatusCode),
			":status-text": types.StringValue(resp.Status),
			":body":        types.StringValue(string(body)),
			":headers":     createHeadersHashMap(resp.Header),
		},
	}

	return responseMap, nil
}

// evalHttpDelete performs an HTTP DELETE request
func (e *Evaluator) evalHttpDelete(args []types.Expr) (types.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("http-delete requires 1-2 arguments (URL, optional headers), got %d", len(args))
	}

	// Evaluate the URL
	urlValue, err := e.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlStr, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-delete URL must be a string, got %T", urlValue)
	}

	// Create the request
	req, err := http.NewRequest("DELETE", string(urlStr), nil)
	if err != nil {
		return nil, fmt.Errorf("http-delete failed to create request: %v", err)
	}

	// Handle optional headers
	if len(args) == 2 {
		headersValue, err := e.Eval(args[1])
		if err != nil {
			return nil, err
		}

		headersMap, ok := headersValue.(*types.HashMapValue)
		if !ok {
			return nil, fmt.Errorf("http-delete headers must be a hash map, got %T", headersValue)
		}

		// Set headers from the hash map
		for key, value := range headersMap.Elements {
			if strValue, ok := value.(types.StringValue); ok {
				req.Header.Set(key, string(strValue))
			}
		}
	}

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http-delete request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http-delete failed to read response body: %v", err)
	}

	// Create response hash map
	responseMap := &types.HashMapValue{
		Elements: map[string]types.Value{
			":status":      types.NumberValue(resp.StatusCode),
			":status-text": types.StringValue(resp.Status),
			":body":        types.StringValue(string(body)),
			":headers":     createHeadersHashMap(resp.Header),
		},
	}

	return responseMap, nil
}

// createHeadersHashMap converts HTTP headers to a hash map
func createHeadersHashMap(headers http.Header) *types.HashMapValue {
	headersMap := &types.HashMapValue{
		Elements: make(map[string]types.Value),
	}

	for key, values := range headers {
		if len(values) == 1 {
			headersMap.Elements[key] = types.StringValue(values[0])
		} else {
			// Multiple values for the same header - create a list
			listValues := make([]types.Value, len(values))
			for i, value := range values {
				listValues[i] = types.StringValue(value)
			}
			headersMap.Elements[key] = &types.ListValue{Elements: listValues}
		}
	}

	return headersMap
}
