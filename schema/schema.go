package schema

import (
	"time"
)

// Posts ... root data element for Delicious links list
type Posts struct { // TODO: think of renaming (to more semantically descriptive name)
	PostsList []Post `xml:"post"`
	User      string `xml:"user,attr"`
	Update    string `xml:"update,attr"`
	Tag       string `xml:"tag,attr"`
	Total     int    `xml:"total,attr"`
}

// Post ... link's data layout
type Post struct { // TODO: think of renaming (to more semantically descriptive name)
	Href        string `xml:"href,attr"`
	Description string `xml:"description,attr"`
	Tags        string `xml:"tag,attr"`
	Time        string `xml:"time,attr"`
	Extended    string `xml:"extended,attr"`
	Hash        string `xml:"hash,attr"`
	Meta        string `xml:"meta,attr"`
}

// PostMetadata ... metadata for post
type PostMetadata struct { // TODO: think of renaming (to more semantically descriptive name)
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

// Entry ... db entry data layout, for search results representation
type Entry struct {
	Content    string
	Type       string
	Tags       []string
	Scheduled  time.Time
	Deadline   time.Time
	AddedAt    time.Time
	ModifiedAt time.Time
	Priority   string
	TodoStatus string
	Metadata   PostMetadata
	UUID       string
}

// SearchResponse ... JSON-RPC response for Search method
type SearchResponse struct {
	Error   string
	Count   int
	Entries []Entry
}
