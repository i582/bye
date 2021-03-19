package bye

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/gomarkdown/markdown"
	mhtml "github.com/gomarkdown/markdown/html"

	fparser "bye/internal/parser/file"
	"bye/internal/parser/page"
)

type File struct {
	// File name.
	name string

	// File extension.
	//
	// Used to set a class in html, for syntax highlighting.
	ext string

	// Separate parts of the file obtained after parsing the code.
	parts []fparser.Part

	// Full file code.
	code string

	// The parent page where this file resides.
	parent *Page

	// Template for the file.
	tmpl *template.Template

	// Flag indicating that the file is a special header file.
	header bool

	repr *Representator
}

func NewFile(part page.Part, parent *Page, repr *Representator) *File {
	f := &File{
		name:   string(part.Filename),
		ext:    string(part.Lang),
		parts:  part.FileParts,
		code:   string(part.FileContent),
		parent: parent,
		repr:   repr,
	}

	for i := range f.parts {
		f.parts[i].Ext = f.ext
	}

	return f
}

func (f *File) WithMarkdown() *File {
	f.parseMarkdown()
	return f
}

func (f *File) WithHighlights() *File {
	f.addHighlights()
	return f
}

func (f *File) HTML() template.HTML {
	f.prepareParts()

	buf := bytes.NewBuffer(nil)
	err := f.tmpl.Execute(buf, f.repr.File(f))
	if err != nil {
		return template.HTML(err.Error())
	}

	return template.HTML(buf.String())
}

func (f *File) prepareParts() {
	for _, part := range f.parts {
		part.Ext = f.ext
	}
}

func (f *File) parseMarkdown() {
	htmlFlags := mhtml.CommonFlags | mhtml.HrefTargetBlank
	opts := mhtml.RendererOptions{Flags: htmlFlags}
	renderer := mhtml.NewRenderer(opts)

	for i := range f.parts {
		f.parts[i].RawComment = strings.ReplaceAll(f.parts[i].RawComment, "\r", "")
		f.parts[i].Comment = template.HTML(markdown.ToHTML([]byte(f.parts[i].RawComment), nil, renderer))
	}
}

func (f *File) addHighlights() {
	lexer := lexers.Get(f.ext)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get("base16-snazzy")
	if style == nil {
		style = styles.Fallback
	}

	// w, _ := os.Create("file.css")
	// formatter1 := html.New(html.WithClasses(true))
	// formatter1.WriteCSS(w, style)
	// w.Close()

	formatter := html.New(html.WithClasses(true), html.PreventSurroundingPre(true))
	// if formatter == nil {
	// 	// formatter = formatters.Fallback
	// }

	for i, part := range f.parts {
		iterator, err := lexer.Tokenise(nil, part.RawCode)
		if err != nil {
			continue
		}

		buf := bytes.NewBufferString("")
		err = formatter.Format(buf, style, iterator)
		if err != nil {
			continue
		}

		// htmlCode := strings.TrimSuffix(buf.String(), "\n")
		htmlCode := buf.String()

		f.parts[i].Code = template.HTML(htmlCode)
	}
}
