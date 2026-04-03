package types

type PostImage struct {
	Alt   string `json:"Alt"`
	Image string `json:"Image"`
}

type PostLink struct {
	Uri         string `json:"Uri"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Thumb       string `json:"Thumb"`
}

type PostFacet struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type Post struct {
	URL           string        `json:"URL"`
	AtURI         string        `json:"AtURI"`
	Text          string        `json:"Text"`
	CreatedAt     string        `json:"CreatedAt"`
	Labels        []string      `json:"Labels"`
	Langs         []string      `json:"Langs"`
	Tags          []string      `json:"Tags"`
	Images        []PostImage   `json:"Images"`
	Links         []PostLink    `json:"Links"`
	Facets        []PostFacet   `json:"Facets"`
	AuthorName    string        `json:"AuthorName"`
	ReplyRoot     string        `json:"ReplyRoot"`
	ReplyParent   string        `json:"ReplyParent"`
	ValueAnalysis ValueAnalysis `json:"ValueAnalysis"`
}

type PostURLs struct {
	Urls []string `json:"urls"`
}
