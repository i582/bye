package file

import (
	"bytes"
	"html/template"
)

type Part struct {
	Comment    template.HTML
	RawComment string
	Code       template.HTML
	RawCode    string

	BeforeCode    bool
	AfterCode     bool
	FirstCodePart bool
	LastCodePart  bool

	Ext string
}

type Parser struct {
	filename []byte
	lang     []byte
	contents []byte
	lines    [][]byte

	parts    []Part
	warnings []string
}

func NewParser(filename, lang []byte, content []byte) *Parser {
	lines := bytes.Split(content, []byte("\n"))

	return &Parser{
		filename: filename,
		contents: content,
		lang:     lang,
		lines:    lines,
	}
}

func (p *Parser) Parts() []Part {
	return p.parts
}

func (p *Parser) Parse() {
	startCommentIndex := -1
	endCommentIndex := -1

	startCodeIndex := -1
	endCodeIndex := -1

	var firstComment = true

	var comment [][]byte
	var code [][]byte

	var parts []Part

	for i := range p.lines {
		switch {
		case isEmptyLine(p.lines[i]) && startCommentIndex != -1:
			// get full comment part
			comment = p.lines[startCommentIndex : endCommentIndex+1]
			for i := range comment {
				comment[i] = trimComment(comment[i])
			}

			parts = append(parts, Part{
				RawComment: joinLines(comment),
				RawCode:    "",
			})

			startCommentIndex = -1
			comment = comment[:]

		case isAllowedComment(p.lines[i]):
			if startCodeIndex != -1 {
				// get all code part
				code = p.lines[startCodeIndex : endCodeIndex+1]
				startCodeIndex = -1
				for i := range code {
					code[i] = bytes.ReplaceAll(code[i], []byte("\t"), []byte("  "))
				}

				parts = append(parts, Part{
					RawComment: joinLines(comment),
					RawCode:    string(bytes.Join(code, []byte("\n"))),
				})
			}

			if startCommentIndex == -1 {
				startCommentIndex = i
				endCommentIndex = i
			} else {
				endCommentIndex = i
			}

		case startCommentIndex != -1:
			comment = p.lines[startCommentIndex : endCommentIndex+1]
			for i := range comment {
				comment[i] = trimComment(comment[i])
			}

			startCommentIndex = -1
			startCodeIndex = i
			endCodeIndex = i
			firstComment = false

		default:
			if firstComment {
				firstComment = false
			}

			if startCodeIndex == -1 {
				startCodeIndex = i
				endCodeIndex = i
			} else {
				endCodeIndex = i
			}
		}
	}

	if startCodeIndex != -1 {
		// get last code part
		code = p.lines[startCodeIndex : endCodeIndex+1]
		for i := range code {
			code[i] = bytes.ReplaceAll(code[i], []byte("\t"), []byte("  "))
		}

		parts = append(parts, Part{
			RawComment: joinLines(comment),
			RawCode:    string(bytes.Join(code, []byte("\n"))),
		})
	}

	if startCommentIndex != -1 {
		// get last comment
		comment = p.lines[startCommentIndex : endCommentIndex+1]
		for i := range comment {
			comment[i] = trimComment(comment[i])
		}
		parts = append(parts, Part{
			RawComment: joinLines(comment),
			RawCode:    "",
		})
	}

	if len(parts) > 0 {
		// setup BeforeCode and AfterCode flags
		partIndexBeforeCode := -1
		for i, part := range parts {
			if i > partIndexBeforeCode && part.RawCode == "" {
				parts[i].BeforeCode = true
				partIndexBeforeCode = i
			} else if part.RawCode != "" {
				partIndexBeforeCode = len(parts)
			}
		}
		partIndexAfterCode := len(parts)
		for i := len(parts) - 1; i >= 0; i-- {
			part := parts[i]
			if i < partIndexAfterCode && part.RawCode == "" {
				parts[i].AfterCode = true
				partIndexAfterCode = i
			} else if part.RawCode != "" {
				partIndexAfterCode = -1
			}
		}

		// setup FirstCodePart and LastCodePart flags
		firstIndexPartWithCode := len(parts) - 1
		lastIndexPartWithCode := 0
		for i, part := range parts {
			if i < firstIndexPartWithCode && part.RawCode != "" {
				firstIndexPartWithCode = i
			}
			if i > lastIndexPartWithCode && part.RawCode != "" {
				lastIndexPartWithCode = i
			}
		}
		parts[firstIndexPartWithCode].FirstCodePart = true
		parts[lastIndexPartWithCode].LastCodePart = true
	}

	for i, part := range parts {
		parts[i].Code = template.HTML(part.RawCode)
	}

	p.parts = parts
}
