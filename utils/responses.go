/*
Takes HTTP response's "Content-Type" header to format response as html, text, json, ...

(ref for Content-Type in HTTP - https://httpwg.org/specs/rfc9110.html#field.content-type)
*/
package utils 

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"golang.org/x/net/html"
)


func ResponseParser_main(body []byte, headers map[string]string, status int, isVerbose bool) string{
	contentType := GetResponse_ContentMimeType(headers["Content-Type"])
	if contentType == ""  {
		return defaultFormatter(body, headers, status, isVerbose)
	}
	var output strings.Builder
	if isVerbose {
		headersFmt := formatHeaders_verbose(headers, status)
		fmt.Fprintf(&output, "%s\n", headersFmt)
	}
	var bodyFmt string
	switch(contentType){
		case "text/plain":
			bodyFmt = string(body) 
		case "text/html":
			bodyFmt = formatBody_html(body)
		case "application/json":	
			bodyFmt = formatBody_json(body)
		default:
			return defaultFormatter(body, headers, status, isVerbose)
	}
	fmt.Fprintf(&output, "\n%s", bodyFmt)
	return output.String()
}

//for a response's Content-Type, what will be its file extension
//(example -> json for application/json, html for )
func GetResponseContent_FileExtension(respContentType string) string{
	respContentMimeType := GetResponse_ContentMimeType(respContentType)
	switch(respContentMimeType){
		case "text/plain":
			return "txt"
		case "text/html":
			return "html"
		case "application/json":	
			return "json"
		default:
			return "txt"
	}
}

/// ---------------------------------------
//Formatters ------------------------
func defaultFormatter(body []byte, headers map[string]string, status int, isVerbose bool) string{
	var output strings.Builder
	if isVerbose {
		fmt.Fprintf(&output, "%s", formatHeaders_verbose(headers, status))
	}
	fmt.Fprintf(&output, "\n%s", string(body))
	return output.String()
}

func formatHeaders_verbose(headers map[string]string, status int) string{
	var output strings.Builder
	fmt.Fprintf(&output, "Status: %d\n\n", status)
	// headers 
	if len(headers) > 0 {
		output.WriteString("Headers:\n")
		for key, value := range headers {
			fmt.Fprintf(&output, "  %s: %s\n", key, value)
		}
	}
	return output.String()
}

func formatBody_json(body []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return fmt.Sprintf("[Invalid JSON - showing raw]\n%s", string(body))
	}
	return prettyJSON.String()
}

func formatBody_html(body []byte) string {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return fmt.Sprintf("[Invalid HTML - showing raw]\n%s", string(body))
	}
	var buf bytes.Buffer
	
	var render func(*html.Node, int)
	render = func(n *html.Node, depth int) {
		indent := bytes.Repeat([]byte("  "), depth)
		
		switch n.Type {
		case html.ElementNode:
			buf.Write(indent)
			buf.WriteString("<" + n.Data)
			for _, attr := range n.Attr {
				fmt.Fprintf(&buf, " %s=\"%s\"", attr.Key, attr.Val)
			}
			buf.WriteString(">\n")
			
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				render(c, depth+1)
			}
			
			buf.Write(indent)
			buf.WriteString("</" + n.Data + ">\n")
			
		case html.TextNode:
			text := bytes.TrimSpace([]byte(n.Data))
			if len(text) > 0 {
				buf.Write(indent)
				buf.Write(text)
				buf.WriteString("\n")
			}
			
		case html.DocumentNode:
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				render(c, depth)
			}
		}
	}
	
	render(doc, 0)
	return buf.String()
}

//--------------------------------
//Other utilities

// extract just the MIME type for response's Content-Type 
func GetResponse_ContentMimeType(respContentTypeHeader string) string{
	contentType := strings.TrimSpace(strings.Split(respContentTypeHeader, ";")[0])
	return contentType
}
