package repr

type TagPage struct {
	Common *PageCommonInfo

	Title string
	Pages []*Page
}
