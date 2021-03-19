package bye

import (
	"html/template"

	"bye/pkg/bye/repr"
)

type Representator struct {
	files map[*File]*repr.File
	pages map[*Page]*repr.Page
}

func NewRepresentator() *Representator {
	return &Representator{
		files: map[*File]*repr.File{},
		pages: map[*Page]*repr.Page{},
	}
}

func (r *Representator) PageCommonInfo(page *PageCommonInfo) *repr.PageCommonInfo {
	return &repr.PageCommonInfo{
		HeadTitle:   page.HeadTitle,
		HeaderTitle: page.HeaderTitle,
		SiteRoot:    template.URL(page.SiteRoot),
		StylePath:   template.URL(page.StylePath),
		GithubLink:  template.URL(page.GithubLink),
		Theme:       page.Theme,
		Dst:         page.Dst,
	}
}

func (r *Representator) IndexPage(page *IndexPage) *repr.IndexPage {
	pages := make([]*repr.Page, 0, len(page.pages))
	for _, page := range page.pages {
		pages = append(pages, r.Page(page))
	}

	return &repr.IndexPage{
		Common:      r.PageCommonInfo(page.common),
		Title:       page.title,
		Description: template.HTML(page.description),
		Pages:       pages,
	}
}

func (r *Representator) TagPage(page *TagPage) *repr.TagPage {
	pages := make([]*repr.Page, 0, len(page.pages))
	for _, page := range page.pages {
		pages = append(pages, r.Page(page))
	}

	return &repr.TagPage{
		Common: r.PageCommonInfo(page.common),
		Title:  page.title,
		Pages:  pages,
	}
}

func (r *Representator) Page(page *Page) *repr.Page {
	if page == nil {
		return nil
	}

	p, has := r.pages[page]
	if has {
		return p
	}

	p = &repr.Page{
		Common:        r.PageCommonInfo(page.common),
		Title:         page.title,
		Tags:          page.tags,
		GeneratedTime: page.generatedTime,
		Link:          page.link,
	}
	r.pages[page] = p

	parts := make([]*repr.PagePart, 0, len(page.parts))
	for _, part := range page.parts {
		parts = append(parts, r.PagePart(part))
	}

	nextPage := r.Page(page.nextPage)
	prevPage := r.Page(page.prevPage)

	p.Parts = parts
	p.NextPage = nextPage
	p.PrevPage = prevPage

	return p
}

func (r *Representator) File(file *File) *repr.File {
	f, has := r.files[file]
	if has {
		return f
	}

	f = &repr.File{
		Name:   file.name,
		Ext:    file.ext,
		Code:   file.code,
		Parts:  file.parts,
		Parent: r.Page(file.parent),
	}

	r.files[file] = f

	return f
}

func (r *Representator) PagePart(part *PagePart) *repr.PagePart {
	p := &repr.PagePart{
		Text:    part.Text,
		TextRaw: part.TextRaw,
		IsFile:  part.IsFile,
	}

	if part.IsFile {
		p.File = r.File(part.File)
	}

	return p
}
