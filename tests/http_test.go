package htty_test

import (
	http "htty/utils"
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
		Body: `{"name":"htty","role":"dev"}`,
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



func TestHeaderKVparser_Basic(tt *testing.T) {
	raw := `
Authorization: Bearer Token
Content-Type: application/json
X-Test: 123
`
	headers := http.HeaderKVparser(raw)
	if len(headers) != 3 {
		tt.Fatalf("expected 3 headers, got %d", len(headers))
	}
	if headers["Authorization"] != "Bearer Token" {
		tt.Fatalf("unexpected Authorization value: %s", headers["Authorization"])
	}
	if headers["Content-Type"] != "application/json" {
		tt.Fatalf("unexpected Content-Type value: %s", headers["Content-Type"])
	}
	if headers["X-Test"] != "123" {
		tt.Fatalf("unexpected X-Test value: %s", headers["X-Test"])
	}
}

func TestHeaderKVparser_IgnoreInvalidAndEmpty(tt *testing.T) {
	raw := `
InvalidLineWithoutColon
Another-Bad-Line

Accept: application/json
`
	headers := http.HeaderKVparser(raw)
	if len(headers) != 1 {
		tt.Fatalf("expected 1 valid header, got %d", len(headers))
	}
	if headers["Accept"] != "application/json" {
		tt.Fatalf("unexpected Accept value: %s", headers["Accept"])
	}
}
