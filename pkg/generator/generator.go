package generator

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"bye/pkg/bye"
	"bye/pkg/bye/config"

	"github.com/gookit/color"
)

type Generator struct {
	// Configuration file.
	config *config.Config

	// A set of example pages to generate.
	pages *bye.Pages
	// Template for a page with examples.
	pageTmpl *template.Template

	// Index page.
	index *bye.IndexPage
	// Template for a index page.
	indexPageTmpl *template.Template

	// A set of pages for tags.
	tags []*bye.TagPage
	// Template for a page with tag examples.
	tagsTmpl *template.Template

	// Common information for all pages.
	pageCommonInfo *bye.PageCommonInfo

	representator *bye.Representator
}

func NewGenerator(conf *config.Config, pageTmpl *template.Template, indexPageTmpl *template.Template, tagsTmpl *template.Template) *Generator {
	pageCommonInfo := &bye.PageCommonInfo{
		HeadTitle:   conf.Title,
		HeaderTitle: conf.Title,
		SiteRoot:    conf.SiteRoot,
		StylePath:   conf.StylePath,
		GithubLink:  conf.GithubLink,
		Theme:       conf.Theme,
		Dst:         conf.Dst,
	}
	repr := bye.NewRepresentator()

	return &Generator{
		config:         conf,
		pages:          &bye.Pages{PageCommonInfo: pageCommonInfo, SrcStylePath: conf.StylePath, TemplatePath: conf.TemplatesFolder},
		pageTmpl:       pageTmpl,
		index:          nil,
		indexPageTmpl:  indexPageTmpl,
		tags:           nil, // Set after processing the example pages in Generator.handleTags() method.
		tagsTmpl:       tagsTmpl,
		pageCommonInfo: pageCommonInfo,
		representator:  repr,
	}
}

func (g *Generator) handlePage(pagePath string, dir string) error {
	example, found := g.config.PageByFolder(dir)
	if !found {
		return fmt.Errorf("for %s not found description in config", dir)
	}

	indexFilePath := filepath.Join(pagePath, "index.md")

	pageCommonInfo := *g.pageCommonInfo
	pageCommonInfo.HeadTitle += " | " + example.Title
	page, err := bye.NewPage(&pageCommonInfo, indexFilePath, example, g.pageTmpl, dir, g.representator)
	if err != nil {
		return err
	}

	if len(page.Warnings()) != 0 {
		fmt.Println("  Warnings encountered during page preparation:")
		for _, warning := range page.Warnings() {
			color.Yellow.Println("    " + warning)
		}
	}

	g.pages.Pages = append(g.pages.Pages, page)
	return nil
}

func (g *Generator) handleTags() error {
	tags := make(map[string][]*bye.Page, len(g.pages.Pages))

	for _, page := range g.pages.Pages {
		for _, tag := range page.Tags() {
			tags[tag] = append(tags[tag], page)
		}
	}

	for tag, pages := range tags {
		pageCommonInfo := *g.pageCommonInfo
		pageCommonInfo.HeadTitle += " | #" + tag
		tagPage := bye.NewTagPage(&pageCommonInfo, tag, pages, g.tagsTmpl, g.representator)
		g.tags = append(g.tags, tagPage)
	}

	return nil
}

func (g *Generator) HandlePages() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	srcFolder := path.Join(wd, g.config.Src)

	files, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		return fmt.Errorf("error read dir %s: %v", srcFolder, err)
	}

	var indexPageDescription string

	for _, pageDescription := range g.config.Pages {
		pagePath := filepath.Join(srcFolder, pageDescription.Name)
		err = g.handlePage(pagePath, pageDescription.Name)
		if err != nil {
			return err
		}
	}

	for _, dir := range files {
		if dir.Name() == "index.md" {
			data, err := ioutil.ReadFile(filepath.Join(srcFolder, dir.Name()))
			if err != nil {
				return err
			}

			indexPageDescription = string(data)
		}
	}

	err = g.handleTags()
	if err != nil {
		return err
	}

	err = g.handleIndexPage(indexPageDescription)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) handleIndexPage(description string) error {
	index := bye.NewIndexPage(
		g.pageCommonInfo,
		g.config.IndexTitle,
		description,
		g.indexPageTmpl,
		g.pages.Pages,
		g.representator,
	)

	index.WithMarkdown()
	g.index = index

	return nil
}

func (g *Generator) Generate() error {
	fmt.Println("  Prepare generation")

	err := g.HandlePages()
	if err != nil {
		return fmt.Errorf("error handle pages: %v", err)
	}

	g.pages.Prepare()

	fmt.Println("  Creating build folders")
	err = os.MkdirAll(filepath.Join(g.config.Dst, "static", "styles"), 0777)
	if err != nil {
		return err
	}

	err = g.pages.Generate()
	if err != nil {
		return fmt.Errorf("error generate pages: %v", err)
	}

	fmt.Println("  Start index page generation")
	err = g.index.Generate()
	if err != nil {
		return fmt.Errorf("error generate index page: %v", err)
	}
	color.Green.Println("  Index page generation completed successfully")

	for _, tag := range g.tags {
		fmt.Printf("  Start tag '%s' page generation\n", tag.Name())
		err = tag.Generate()
		if err != nil {
			return fmt.Errorf("error generate tag '%s' page: %v", tag.Name(), err)
		}
		color.Green.Printf("  Tag '%s' page generation completed successfully\n", tag.Name())
	}

	return nil
}

func (g *Generator) isAllowedFile(name string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}

	fileExt := filepath.Ext(name)
	fileExt = strings.TrimPrefix(fileExt, ".")

	for _, ext := range exts {
		if ext == fileExt {
			return true
		}
	}

	return false
}
