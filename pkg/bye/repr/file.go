package repr

import "bye/internal/parser/file"

type File struct {
	// File name.
	Name string

	// File extension.
	//
	// Used to set a class in html, for syntax highlighting.
	Ext string

	// Full file Code.
	Code string

	// Separate parts of the file obtained after parsing the code.
	Parts []file.Part

	// The parent page where this file resides.
	Parent *Page
}
