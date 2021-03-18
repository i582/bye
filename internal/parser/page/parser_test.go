package page

import (
	"log"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	content := `
# Header 1
~~~php (filename.php)
$a = 100;
echo $a;
~~~
Some text
`

	content = strings.ReplaceAll(content, "~", "`")

	p := NewParserContent("test_file", []byte(content))
	p.Parse()

	for _, warning := range p.Warnings() {
		log.Println(warning)
	}
}
