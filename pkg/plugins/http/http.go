// Package http implements HTTP functions as a plugin
package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/leinonen/lisp-interpreter/pkg/functions"
	"github.com/leinonen/lisp-interpreter/pkg/plugins"
	"github.com/leinonen/lisp-interpreter/pkg/registry"
	"github.com/leinonen/lisp-interpreter/pkg/types"
)

// HTTPPlugin implements HTTP functions
type HTTPPlugin struct {
	*plugins.BasePlugin
	client *http.Client
}

// NewHTTPPlugin creates a new HTTP plugin
func NewHTTPPlugin() *HTTPPlugin {
	return &HTTPPlugin{
		BasePlugin: plugins.NewBasePlugin(
			"http",
			"1.0.0",
			"HTTP client functions (http-get, http-post, http-put, http-delete)",
			[]string{}, // No dependencies
		),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Functions returns the list of functions provided by this plugin
func (p *HTTPPlugin) Functions() []string {
	return []string{
		"http-get", "http-post", "http-put", "http-delete",
	}
}

// RegisterFunctions registers all HTTP functions with the registry
func (p *HTTPPlugin) RegisterFunctions(reg registry.FunctionRegistry) error {
	// http-get function
	httpGetFunc := functions.NewFunction(
		"http-get",
		registry.CategoryHTTP,
		1,
		"Perform HTTP GET request: (http-get \"https://api.example.com\")",
		p.evalHttpGet,
	)
	if err := reg.Register(httpGetFunc); err != nil {
		return err
	}

	// http-post function
	httpPostFunc := functions.NewFunction(
		"http-post",
		registry.CategoryHTTP,
		-1, // Variable arguments (url, body, optional headers)
		"Perform HTTP POST request: (http-post \"url\" \"body\") or (http-post \"url\" \"body\" headers)",
		p.evalHttpPost,
	)
	if err := reg.Register(httpPostFunc); err != nil {
		return err
	}

	// http-put function
	httpPutFunc := functions.NewFunction(
		"http-put",
		registry.CategoryHTTP,
		-1, // Variable arguments (url, body, optional headers)
		"Perform HTTP PUT request: (http-put \"url\" \"body\") or (http-put \"url\" \"body\" headers)",
		p.evalHttpPut,
	)
	if err := reg.Register(httpPutFunc); err != nil {
		return err
	}

	// http-delete function
	httpDeleteFunc := functions.NewFunction(
		"http-delete",
		registry.CategoryHTTP,
		-1, // Variable arguments (url, optional headers)
		"Perform HTTP DELETE request: (http-delete \"url\") or (http-delete \"url\" headers)",
		p.evalHttpDelete,
	)
	if err := reg.Register(httpDeleteFunc); err != nil {
		return err
	}

	return nil
}

// Helper function to create HTTP response hash map
func (p *HTTPPlugin) createResponseMap(resp *http.Response, body []byte) types.Value {
	// Create headers map
	headersMap := make(map[string]types.Value)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headersMap[key] = types.StringValue(values[0])
		}
	}

	responseMap := map[string]types.Value{
		":status":      types.NumberValue(float64(resp.StatusCode)),
		":status-text": types.StringValue(resp.Status),
		":body":        types.StringValue(string(body)),
		":headers":     &types.HashMapValue{Elements: headersMap},
	}

	return &types.HashMapValue{Elements: responseMap}
}

// Helper function to extract headers from hash map
func (p *HTTPPlugin) extractHeaders(headersValue types.Value) (map[string]string, error) {
	hashMap, ok := headersValue.(*types.HashMapValue)
	if !ok {
		return nil, fmt.Errorf("headers must be a hash map, got %T", headersValue)
	}

	headers := make(map[string]string)
	for key, value := range hashMap.Elements {
		stringValue, ok := value.(types.StringValue)
		if !ok {
			return nil, fmt.Errorf("header value for %s must be a string, got %T", key, value)
		}
		headers[key] = string(stringValue)
	}

	return headers, nil
}

// Helper function to perform HTTP request
func (p *HTTPPlugin) performRequest(method, url string, body io.Reader, headers map[string]string) (types.Value, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s request: %v", method, err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set default content type for POST/PUT if not specified
	if (method == "POST" || method == "PUT") && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %v", method, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return p.createResponseMap(resp, responseBody), nil
}

// evalHttpGet performs an HTTP GET request
func (p *HTTPPlugin) evalHttpGet(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("http-get requires exactly 1 argument, got %d", len(args))
	}

	urlValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlString, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-get URL must be a string, got %T", urlValue)
	}

	return p.performRequest("GET", string(urlString), nil, nil)
}

// evalHttpPost performs an HTTP POST request
func (p *HTTPPlugin) evalHttpPost(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("http-post requires 2 or 3 arguments, got %d", len(args))
	}

	urlValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlString, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-post URL must be a string, got %T", urlValue)
	}

	bodyValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	bodyString, ok := bodyValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-post body must be a string, got %T", bodyValue)
	}

	var headers map[string]string
	if len(args) == 3 {
		headersValue, err := evaluator.Eval(args[2])
		if err != nil {
			return nil, err
		}

		headers, err = p.extractHeaders(headersValue)
		if err != nil {
			return nil, fmt.Errorf("http-post headers: %v", err)
		}
	}

	return p.performRequest("POST", string(urlString), strings.NewReader(string(bodyString)), headers)
}

// evalHttpPut performs an HTTP PUT request
func (p *HTTPPlugin) evalHttpPut(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("http-put requires 2 or 3 arguments, got %d", len(args))
	}

	urlValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlString, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-put URL must be a string, got %T", urlValue)
	}

	bodyValue, err := evaluator.Eval(args[1])
	if err != nil {
		return nil, err
	}

	bodyString, ok := bodyValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-put body must be a string, got %T", bodyValue)
	}

	var headers map[string]string
	if len(args) == 3 {
		headersValue, err := evaluator.Eval(args[2])
		if err != nil {
			return nil, err
		}

		headers, err = p.extractHeaders(headersValue)
		if err != nil {
			return nil, fmt.Errorf("http-put headers: %v", err)
		}
	}

	return p.performRequest("PUT", string(urlString), strings.NewReader(string(bodyString)), headers)
}

// evalHttpDelete performs an HTTP DELETE request
func (p *HTTPPlugin) evalHttpDelete(evaluator registry.Evaluator, args []types.Expr) (types.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("http-delete requires 1 or 2 arguments, got %d", len(args))
	}

	urlValue, err := evaluator.Eval(args[0])
	if err != nil {
		return nil, err
	}

	urlString, ok := urlValue.(types.StringValue)
	if !ok {
		return nil, fmt.Errorf("http-delete URL must be a string, got %T", urlValue)
	}

	var headers map[string]string
	if len(args) == 2 {
		headersValue, err := evaluator.Eval(args[1])
		if err != nil {
			return nil, err
		}

		headers, err = p.extractHeaders(headersValue)
		if err != nil {
			return nil, fmt.Errorf("http-delete headers: %v", err)
		}
	}

	return p.performRequest("DELETE", string(urlString), nil, headers)
}
