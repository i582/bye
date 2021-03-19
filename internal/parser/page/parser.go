package page

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"regexp"
	"strings"

	fparser "bye/internal/parser/file"
)

var fileRegexp = regexp.MustCompile(`@\[\]\(.*?\)`)

type Part struct {
	Text    template.HTML
	TextRaw string

	Lang        []byte
	Filename    []byte
	FileContent []byte
	FileParts   []fparser.Part
	IsFile      bool
}

type Parser struct {
	filename string
	contents []byte
	lines    [][]byte
	parts    []Part

	fileRegexp *regexp.Regexp
	warnings   []string
}

func NewParser(filename string) (*Parser, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewParserContent(filename, contents), nil
}

func NewParserContent(filename string, content []byte) *Parser {
	lines := bytes.Split(content, []byte("\n"))

	return &Parser{
		filename:   filename,
		contents:   content,
		lines:      lines,
		fileRegexp: fileRegexp,
	}
}

func (p *Parser) Parts() []Part {
	return p.parts
}

func (p *Parser) Warnings() []string {
	return p.warnings
}

func (p *Parser) Parse() {
	last := 0
	startCode := false
	for i, line := range p.lines {
		line = bytes.TrimLeft(line, " \t")
		if !bytes.HasPrefix(line, []byte("```")) {
			continue
		}

		startCode = !startCode
		if startCode {
			text := string(bytes.Join(p.lines[last:i], []byte("\n")))
			text = strings.TrimLeft(text, "\r")
			if text != "" {
				p.parts = append(p.parts, Part{
					TextRaw: text,
				})
			}
			last = i
		} else {
			lines := p.lines[last : i+1]
			fileContent := bytes.Join(lines[1:len(lines)-1], []byte("\n"))

			lang, filename := p.parseCodeHeader(lines[0], last)
			fileParser := fparser.NewParser(filename, lang, fileContent)
			fileParser.Parse()

			p.parts = append(p.parts, Part{
				Lang:        lang,
				Filename:    filename,
				FileContent: fileContent,
				FileParts:   fileParser.Parts(),
				IsFile:      true,
			})
			last = i + 1
		}
	}

	if last < len(p.lines) {
		p.parts = append(p.parts, Part{
			TextRaw: string(bytes.Join(p.lines[last:len(p.lines)], []byte("\n"))),
		})
	}
}

func (p *Parser) warn(message string) {
	p.warnings = append(p.warnings, message)
}

func (p *Parser) parseCodeHeader(line []byte, i int) (lang, filename []byte) {
	parts := bytes.Fields(bytes.TrimPrefix(line, []byte("```")))
	if len(parts) == 0 {
		p.warn(fmt.Sprintf("%s:%d please specify the filename and language for correct work of syntax highlighting", p.filename, i+1))
		return lang, filename
	}

	if len(parts) == 1 {
		if parts[0][0] == '(' {
			p.warn(fmt.Sprintf("%s:%d please specify the code language first, and then the file name", p.filename, i+1))
			filename = parts[0][1 : len(parts[0])-1]
			return lang, filename
		}

		p.warn(fmt.Sprintf("%s:%d please specify the file name, after the code language in brackets ```php (filename.php)", p.filename, i+1))
		lang = parts[0]
		return lang, filename
	}

	if len(parts) == 2 {
		if parts[0][0] == '(' {
			p.warn(fmt.Sprintf("%s:%d please specify the code language first, and then the file name", p.filename, i+1))
			filename = parts[0][1 : len(parts[0])-1]
			lang = parts[1]
			return lang, filename
		}

		lang = parts[0]
		filename = parts[1][1 : len(parts[1])-1]
	}

	return lang, filename
}
