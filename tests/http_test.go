package htty_test

import (
	http "htty/http"
	types "htty/types"
	"testing"
	"encoding/json"
)

func TestHTTPMethod(tt *testing.T){
	method := "PUT"
	var correctMethod bool = http.AssertHTTPMethodType(method)
	if !correctMethod {
		tt.Errorf("%s is not valid", method)
	}
}

func TestHTTPCall_SimpleGet(tt *testing.T){
	req := types.HttpType{
		Path:   "http://example.com",
		Method: "GET",
	}
	_, status, err := http.HTTPCaller(req)
	if err != nil {
		tt.Fatalf("unexpected error: %v", err)
	}
	if status != 200 {
		tt.Fatalf("expected status 200, got %d", status)
	}
}


// curl equivalent:
// curl -X POST https://httpbin.org/post -H "Content-Type: application/json" -H "X-Custom-Test: htty" \
//   -d '{"name":"htty","role":"dev"}'
func TestHTTPCaller_PostWithHeadersAndBody(tt *testing.T) {
	req := types.HttpType{
		Path:   "https://httpbin.org/post",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Custom-Test": "htty",
		},
		Body: []string{
			`{"name":"htty","role":"dev"}`,
		},
	}
	body, status, err := http.HTTPCaller(req)
	if err != nil {
		tt.Fatalf("unexpected error: %v", err)
	}

	if status != 200 {
		tt.Fatalf("expected status 200, got %d", status)
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		tt.Fatalf("invalid json response: %v", err)
	}

	// validate echoed header
	headers := parsed["headers"].(map[string]any)
	if headers["X-Custom-Test"] != "htty" {
		tt.Fatalf("header not echoed correctly")
	}

	// validate echoed body
	data := parsed["data"].(string)
	if data != `{"name":"htty","role":"dev"}` {
		tt.Fatalf("body not echoed correctly")
	}
}
