/*
Takes HTTP response's "Content-Type" header to format response as html, text, json, ...

(ref for Content-Type in HTTP - https://httpwg.org/specs/rfc9110.html#field.content-type)
*/
package htty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)


func ResponseParser_main(body []byte, headers map[string]string, status int) string{
	contentTypeFull := headers["Content-Type"]
	//if empty -> default plain text 
	if contentTypeFull == "" {
		return defaultFormatter(body, headers, status)
	}
	// extract just the MIME type 
	contentType := strings.TrimSpace(strings.Split(contentTypeFull, ";")[0])	
	Debugf("Content-Type: %s", contentType)
	var output strings.Builder
	//TODO: verbose toggle condition to limit whether to show header/what headers to show, ...
	headersFmt := formatHeaders_verbose(headers, status)
	fmt.Fprintf(&output, "%s\n", headersFmt)
	var bodyFmt string
	switch(contentType){
		case "text/plain":
			bodyFmt = string(body) 
		case "text/html":
			bodyFmt = formatBody_html(body)
		case "application/json":	
			bodyFmt = formatBody_json(body)
		default:
			return defaultFormatter(body, headers, status)
	}
	fmt.Fprintf(&output, "Body:\n%s", bodyFmt)
	return output.String()
}

/// ---------------------------------------
//Formatters ------------------------
func defaultFormatter(body []byte, headers map[string]string, status int) string{
	var output strings.Builder
	fmt.Fprintf(&output, "%s", formatHeaders_verbose(headers, status))
	fmt.Fprintf(&output, "\nBody:\n%s", string(body))
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
