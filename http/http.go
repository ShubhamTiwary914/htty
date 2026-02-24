package htty
	
import (
	"net/http"
	"bytes"
	"io"
	"time"
	"strings"
	types "htty/types"
	utils "htty/utils"
)  

// sends an HTTP request with the given method, URL, headers, and body.
// returns the response body bytes, status code, and an error if any.
func HTTPCaller(httpObj types.HttpType) ([]byte, int, error) {
	if !AssertHTTPMethodType(httpObj.Method) {
		utils.Errorf("HTTP method type not valid: %s", httpObj.Method)			
	}
	client := &http.Client{
		Timeout: time.Duration(types.REQUEST_TIMEOUT) * time.Second,
	}
	//encapsulate request
	var bodyBuffer io.Reader
	if len(httpObj.Body) > 0 {
		bodyBuffer = bytes.NewBufferString(strings.Join(httpObj.Body, "\n"))
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

func AssertHTTPMethodType(method string) bool{
	_, found := types.HTTP_METHOD[method]
	if found {
		return true;
	}
	return false;
}
