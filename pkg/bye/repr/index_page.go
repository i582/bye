package repr

import "html/template"

type IndexPage struct {
	Common *PageCommonInfo

	Title       string
	Description template.HTML
	Pages       []*Page
}
