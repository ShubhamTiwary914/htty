package types 

type HttpRespState struct {
    Output  string
    Raw     string
    Headers map[string]string
    Status  int
}

type HttyState struct {
	Method 		string	
	URL 	 	string
	ReqHeaders 	map[string]string
	ReqBody     	string
	Response 	HttpRespState
}

