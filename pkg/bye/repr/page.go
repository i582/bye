package repr

import "html/template"

type PageCommonInfo struct {
	HeadTitle   string
	HeaderTitle string
	SiteRoot    template.URL
	StylePath   template.URL
	GithubLink  template.URL
	Theme       string
	Dst         string
}

type PagePart struct {
	Text    template.HTML
	TextRaw string

	File   *File
	IsFile bool
}

type Page struct {
	// Common information for the page.
	Common *PageCommonInfo

	// Page title.
	Title string

	// Tags for the page.
	Tags []string

	// Page generation time.
	GeneratedTime string

	Parts []*PagePart

	// Next page.
	//
	// If nil, it means there is no next page.
	NextPage *Page

	// Previous page.
	//
	// If nil, it means there is no previous page.
	PrevPage *Page

	// Link to the page relative to the root path.
	Link string
}
