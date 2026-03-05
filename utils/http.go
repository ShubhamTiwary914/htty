package htty
	
import (
	"net/http"
	"io"
	"time"
	"strings"
	types "htty/types"
)  


/* 
	Sends an HTTP request with the given method, URL, headers, and body.
 	
	Returns: (response body bytes, status code, error if any).
*/
func HTTPCaller(httpObj types.HttpType) ([]byte, int, error) {
	if !AssertHTTPMethodType(httpObj.Method) {
		Errorf("HTTP method type not valid: %s", httpObj.Method)			
	}
	client := &http.Client{
		Timeout: time.Duration(types.REQUEST_TIMEOUT) * time.Second,
	}
	//encapsulate request
	var bodyBuffer io.Reader
	if httpObj.Body != "" {
		bodyBuffer = strings.NewReader(httpObj.Body)
	}
	request, err := http.NewRequest(httpObj.Method, httpObj.Path, bodyBuffer)
	if err != nil {
		return nil, 0, err
	}	
	for header, val := range httpObj.Headers {
		request.Header.Set(header, val)
	}
	//make request
	resp, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respBody, resp.StatusCode, nil
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
	_, found := types.HTTP_METHOD[method]
	if found {
		return true;
	}
	return false;
}
