package schema

// Posts ... root data element for Delicious links list
type Posts struct {
	PostsList []Post `xml:"post"`
	User      string `xml:"user,attr"`
	Update    string `xml:"update,attr"`
	Tag       string `xml:"tag,attr"`
	Total     int    `xml:"total,attr"`
}

// Post ... link's data layout
type Post struct {
	Href        string `xml:"href,attr"`
	Description string `xml:"description,attr"`
	Tags        string `xml:"tag,attr"`
	Time        string `xml:"time,attr"`
	Extended    string `xml:"extended,attr"`
	Hash        string `xml:"hash,attr"`
	Meta        string `xml:"meta,attr"`
}

// PostMetadata ... metadata for post
type PostMetadata struct {
	Description     string
	TimeAddedOrigin string
	HashOrigin      string
	MetaOrigin      string
	From            string
}

// ImportError ... information about import error(s)
type ImportError struct {
	Link  string
	Error string
}

// AddResponse ... JSON-RPC response for Add method
type AddResponse struct {
	Message string
}
