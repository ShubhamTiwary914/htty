package htty


type HttyConfig struct {
	Common HttyCommonConfig
	Panels HttyPanelsConfig
	Log    HttyLogConfig
	Key    HttyKeyConfig
}

type HttyCommonConfig struct {
	Focus_border_color string
	Background_color string
	Unfocus_border_color string
	Textoptions_selection_color string
}

type HttyPanel struct {
	ID        int
	Width     int
	Height    int
	Margin    [2]int //[x,y] 
	Layer     int
}

type HttyPanelsConfig struct {
	Main           HttyPanel
	Side           HttyPanel
	Main_req       HttyPanel
	Main_req_method HttyPanel
	Main_req_url    HttyPanel
	Main_req_headers HttyPanel
	Main_req_body   HttyPanel
	Main_res        HttyPanel
}

type HttyLogConfig struct {
	Allow bool
	File  string
	Level string
}

type HttyKeyConfig struct {
	Quit        	string
	Nextpanel   	string
	Sendapicall 	string
	CompleteText 	string
	Textoptions_prev string
	Textoptions_next string
}
