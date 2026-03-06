package htty

type HttyConfig struct {
	Common HttyCommonConfig
	Panels HttyPanelsConfig
	Log    HttyLogConfig
	Key    HttyKeyConfig
}

type HttyBorder struct {
	Top    bool
	Left   bool
	Right  bool
	Bottom bool
}

type HttyCommonConfig struct {
	Border             HttyBorder
	Focus_border_color string
	Unfocus_border_color string
	Default_bg         string
	Default_fg         string
	Default_text_color string
}

type HttyPanel struct {
	ID        int
	Width     int
	Height    int
	BG        string
	FG        string
	Text_color string
	Border    HttyBorder
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
	Quit        string
	Nextpanel   string
	Sendapicall string
}
