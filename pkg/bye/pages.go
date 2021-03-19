package bye

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gookit/color"
)

type Pages struct {
	*PageCommonInfo

	SrcStylePath string
	TemplatePath string
	Pages        []*Page
}

func (p *Pages) Prepare() {
	if len(p.Pages) == 0 || len(p.Pages) == 1 {
		return
	}

	for i, page := range p.Pages {
		if i == 0 {
			page.nextPage = p.Pages[i+1]
			continue
		} else if i == len(p.Pages)-1 {
			page.prevPage = p.Pages[i-1]
			continue
		}

		page.nextPage = p.Pages[i+1]
		page.prevPage = p.Pages[i-1]
	}
}

func (p *Pages) Generate() error {
	fmt.Println("  Preparation of static files")

	err := p.copyStyleFile()
	if err != nil {
		return err
	}

	err = p.copyTemplateFiles()
	if err != nil {
		return err
	}

	fmt.Println("  Preparation of static files completed successfully")

	for _, page := range p.Pages {
		fmt.Printf("  Start page '%s' generation\n", page.title)

		err = p.handlePage(page)
		if err != nil {
			return err
		}

		color.Green.Printf("  Page '%s' generation completed successfully\n", page.title)
	}

	return nil
}

func (p *Pages) handlePage(page *Page) error {
	path := filepath.Join(p.Dst, page.dir)

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, "index.html"))
	if err != nil {
		return err
	}

	data := page.HTML()

	fmt.Fprint(file, data)
	return nil
}

func (p *Pages) copyStyleFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	stylePath := filepath.Join(wd, p.SrcStylePath)
	newStylePath := filepath.Join(p.Dst, "static", "styles", "style.css")

	err = copyFileContents(stylePath, newStylePath)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pages) copyTemplateFiles() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	templatePath := filepath.Join(wd, p.TemplatePath)
	newTemplatePath := filepath.Join(p.Dst, "static")

	files, err := ioutil.ReadDir(templatePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		singleTemplatePath := filepath.Join(templatePath, file.Name())
		singleNewTemplatePath := filepath.Join(newTemplatePath, file.Name())
		err = copyFileContents(singleTemplatePath, singleNewTemplatePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFileContents(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
