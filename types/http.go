package htty

const (
	REQUEST_TIMEOUT int = 10
)

var HTTP_METHOD = map[string]int{
	"GET": 1,
	"HEAD": 2,
	"POST": 3,
	"PUT": 4,
	"PATCH": 5,
	"DELETE": 6,
	"CONNECT": 7,
	"OPTIONS": 8,
	"TRACE": 9,
}

type HttpType struct {
	Path string
 	Method string 	
	Headers map[string]string
	Body []string
}
