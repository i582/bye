package bye

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type TagPage struct {
	common *PageCommonInfo

	title string
	pages []*Page
	tmpl  *template.Template
	repr  *Representator
}

func NewTagPage(common *PageCommonInfo, title string, pages []*Page, tmpl *template.Template, repr *Representator) *TagPage {
	return &TagPage{
		common: common,
		title:  title,
		pages:  pages,
		tmpl:   tmpl,
		repr:   repr,
	}
}

func (p *TagPage) Name() string {
	return p.title
}

func (p *TagPage) Generate() error {
	buf := bytes.NewBuffer(nil)
	p1 := p.repr.TagPage(p)
	err := p.tmpl.Execute(buf, p1)
	if err != nil {
		return err
	}
	data := buf.String()

	path := filepath.Join(p.common.Dst, "tags", p.title)
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, "index.html"))
	if err != nil {
		return err
	}

	fmt.Fprint(file, data)
	return nil
}
