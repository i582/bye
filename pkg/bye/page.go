package bye

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	mhtml "github.com/gomarkdown/markdown/html"

	"bye/internal/parser/page"
	"bye/pkg/bye/config"
)

type PageCommonInfo struct {
	HeadTitle   string
	HeaderTitle string
	SiteRoot    string
	StylePath   string
	GithubLink  string
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
	common *PageCommonInfo

	// Page title.
	title string

	// Tags for the page.
	tags []string

	// Page generation time. The format is determined by the dateLayout field.
	generatedTime string
	dateLayout    string

	// Files to be on the current page.
	parts []*PagePart

	// Next page.
	//
	// If nil, it means there is no next page.
	// Set in the Pages.Prepare() method.
	nextPage *Page

	// Previous page.
	//
	// If nil, it means there is no previous page.
	// Set in the Pages.Prepare() method.
	prevPage *Page

	// Link to the page relative to the root path.
	link string

	// The folder in which the generated files will be located.
	dir string

	// Template for the page.
	tmpl *template.Template

	repr     *Representator
	warnings []string
}

func NewPage(common *PageCommonInfo, filename string, desc *config.PageDescription, tmpl *template.Template, dir string, repr *Representator) (*Page, error) {
	p, err := page.NewParser(filename)
	if err != nil {
		return nil, err
	}

	p.Parse()

	pageParts := make([]*PagePart, 0, len(p.Parts()))
	pg := &Page{
		common:        common,
		title:         desc.Title,
		tags:          desc.Tags,
		generatedTime: time.Now().Format("January 2, 2006"),
		dateLayout:    "January 2, 2006",
		nextPage:      nil,
		prevPage:      nil,
		link:          filepath.ToSlash(filepath.Join(dir, "index.html")),
		dir:           dir,
		tmpl:          tmpl,
		repr:          repr,
		warnings:      p.Warnings(),
	}

	for _, part := range p.Parts() {
		if part.IsFile {
			pageParts = append(pageParts, &PagePart{
				File:   NewFile(part, pg, repr).WithHighlights().WithMarkdown(),
				IsFile: true,
			})
			continue
		}

		pageParts = append(pageParts, &PagePart{
			Text:    template.HTML(part.TextRaw),
			TextRaw: part.TextRaw,
		})
	}

	pg.parts = pageParts
	pg.withMarkdown()

	return pg, nil
}

func (p *Page) Warnings() []string {
	return p.warnings
}

func (p *Page) Tags() []string {
	return p.tags
}

func (p *Page) withMarkdown() {
	htmlFlags := mhtml.CommonFlags | mhtml.HrefTargetBlank
	opts := mhtml.RendererOptions{Flags: htmlFlags}
	renderer := mhtml.NewRenderer(opts)

	for i, part := range p.parts {
		if part.IsFile {
			continue
		}

		p.parts[i].TextRaw = strings.ReplaceAll(p.parts[i].TextRaw, "\r", "")
		p.parts[i].Text = template.HTML(markdown.ToHTML([]byte(p.parts[i].TextRaw), nil, renderer))
	}
}

func (p *Page) HTML() template.HTML {
	for _, part := range p.parts {
		if part.IsFile {
			part.File.prepareParts()
		}
	}

	buf := bytes.NewBuffer(nil)
	err := p.tmpl.Execute(buf, p.repr.Page(p))
	if err != nil {
		return template.HTML(err.Error())
	}

	return template.HTML(buf.String())
}
