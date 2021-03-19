package bye

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

type IndexPage struct {
	common *PageCommonInfo

	title       string
	description string
	pages       []*Page
	tmpl        *template.Template

	repr *Representator
}

func NewIndexPage(common *PageCommonInfo, title string, description string, tmpl *template.Template, pages []*Page, repr *Representator) *IndexPage {
	return &IndexPage{
		common:      common,
		title:       title,
		description: description,
		pages:       pages,
		tmpl:        tmpl,
		repr:        repr,
	}
}

func (p *IndexPage) WithMarkdown() {
	p.parseMarkdown()
}

func (p *IndexPage) parseMarkdown() {
	p.description = strings.ReplaceAll(p.description, "\r", "")
	p.description = string(markdown.ToHTML([]byte(p.description), nil, nil))
	p.description = strings.TrimSuffix(p.description, "\n")
}

func (p *IndexPage) Generate() error {
	buf := bytes.NewBuffer(nil)
	err := p.tmpl.Execute(buf, p.repr.IndexPage(p))
	if err != nil {
		return err
	}
	data := buf.String()

	file, err := os.Create(filepath.Join(p.common.Dst, "index.html"))
	if err != nil {
		return err
	}

	fmt.Fprint(file, data)
	return nil
}
