package htty
	
import (
	"net/http"
	"io"
	"time"
	"strings"
	global "htty/globals"
	types "htty/types"
)  


/* 
	Sends an HTTP request with the given method, URL, headers, and body.
 	
	Returns: (response body bytes, status code, error if any).
*/
func HTTPCaller(httpObj types.HttpType) ([]byte, map[string]string, int, error) {
	if !AssertHTTPMethodType(httpObj.Method) {
		Errorf("HTTP method type not valid: %s", httpObj.Method)			
	}
	client := &http.Client{
		Timeout: time.Duration(global.REQUEST_TIMEOUT) * time.Second,
	}
	//encapsulate request
	var bodyBuffer io.Reader
	if httpObj.Body != "" {
		bodyBuffer = strings.NewReader(httpObj.Body)
	}
	request, err := http.NewRequest(httpObj.Method, httpObj.Path, bodyBuffer)
	if err != nil {
		return nil, nil, 0, err
	}	
	for header, val := range httpObj.Headers {
		request.Header.Set(header, val)
	}
	//make request
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()
	responseHeaders := getHTTPHeadersMap(resp.Header)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, resp.StatusCode, err
	}
	return respBody, responseHeaders, resp.StatusCode, nil
}


/*
	Takes headers string of k/v pairs seperated by lines and compose header map(key -> val)
	Assumes key value are delimited by ':', example:

	Authorization: Bearer Token
	Content-Type: application/json
	...
*/
func HeaderKVparser(rawHeaders string) map[string]string{
	headers := make(map[string]string)
	for _, line := range strings.Split(rawHeaders, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		headers[key] = val
	}
	return headers;
}

//check if http method is of standard type(GET/POST/PUT/etc..)
func AssertHTTPMethodType(method string) bool{
	_, found := global.HTTP_METHOD[method]
	if found {
		return true;
	}
	return false;
}


//extract's request/response headers as map(key->val)
func getHTTPHeadersMap(headers http.Header) map[string]string{
	respHeaders := make(map[string]string, len(headers))
	for key, val := range headers {
		respHeaders[key] = strings.Join(val, ", ")
	}
	return respHeaders
}

