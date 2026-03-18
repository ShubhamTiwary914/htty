package htty

type HttpType struct {
	Path string
 	Method string 	
	Headers map[string]string
	Body string
}
