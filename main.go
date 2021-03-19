package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"bye/pkg/bye/config"
	"bye/pkg/bye/server"
	"bye/pkg/generator"
	"bye/pkg/wizard"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

func main() {
	var configPath string
	var local bool
	var port int64

	app := &cli.App{
		Name:        "bye",
		Version:     "v0.0.1",
		Usage:       "By example documentation generator",
		Description: "By example documentation generator",
		Commands: []*cli.Command{
			{
				Name:      "create",
				Usage:     "Create and initialize new project",
				ArgsUsage: "<project-name>",
				Action: func(c *cli.Context) error {
					wizard.Run(c.Args())
					return nil
				},
			},
			{
				Name:  "serve",
				Usage: "Generate and run local server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "config",
						Usage:       "path to config",
						Value:       "bye.yml",
						Destination: &configPath,
					},
					&cli.Int64Flag{
						Name:        "port",
						Aliases:     []string{"p"},
						Usage:       "use specific port",
						Value:       3005,
						Destination: &port,
					},
				},
				Action: func(c *cli.Context) error {
					conf, err := config.New(configPath)
					if err != nil {
						return fmt.Errorf("error open config: %v", err)
					}

					err = generate(conf, false, true)
					if err != nil {
						return err
					}

					fmt.Println()
					fmt.Printf("Started local server on port %d.\n", port)
					fmt.Printf("Go to http://localhost:%d to view the site.\n", port)
					server.Run(conf.Dst, port)
					return nil
				},
			},
			{
				Name:  "generate",
				Usage: "Start generation",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "local",
						Usage:       "use local site root",
						Value:       false,
						Destination: &local,
					},
					&cli.StringFlag{
						Name:        "config",
						Usage:       "path to config",
						Value:       "bye.yml",
						Destination: &configPath,
					},
				},
				Action: func(c *cli.Context) error {
					conf, err := config.New(configPath)
					if err != nil {
						return fmt.Errorf("error open config: %v", err)
					}

					return generate(conf, local, false)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(color.Red.Sprintf("Error: %v", err))
	}
}

func generate(conf *config.Config, local bool, server bool) error {
	pageTmpl, indexTmpl, tagsTmpl, err := Templates(conf)
	if err != nil {
		return err
	}

	fmt.Println("Bye generator v0.0.1")

	if server {
		fmt.Println("Server Mode (serve command provided)")
		conf.SiteRoot = ""
	}

	if local {
		fmt.Println("Local Mode (flag -local provided)")
		if conf.SiteRootLocal == "" {
			conf.SiteRootLocal = conf.Dst
		}

		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("some unexpected error, try again: %v", err)
		}

		conf.SiteRoot = filepath.ToSlash(filepath.Join(workingDir, conf.SiteRootLocal))
	}

	g := generator.NewGenerator(conf, pageTmpl, indexTmpl, tagsTmpl)

	fmt.Println("Start generation")
	err = g.Generate()
	if err != nil {
		return fmt.Errorf("error generate pages: %v", err)
	}
	fmt.Println("End generation")
	fmt.Println()

	color.Green.Println("Generation success")
	return nil
}

func Templates(conf *config.Config) (pageTmpl *template.Template, indexTmpl *template.Template, tagsTmpl *template.Template, err error) {
	pageTmpl, err = template.ParseFiles(conf.TemplatesFolder+"/page.html", conf.TemplatesFolder+"/page-part.html", conf.TemplatesFolder+"/file.html", conf.TemplatesFolder+"/page-preview.html")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error parse templates: %v", err)
	}

	indexTmpl, err = template.ParseFiles(conf.TemplatesFolder + "/index.html")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error parse index page template: %v", err)
	}

	tagsTmpl, err = template.ParseFiles(conf.TemplatesFolder+"/tags.html", conf.TemplatesFolder+"/tag-page-preview.html")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error parse index page template: %v", err)
	}

	return pageTmpl, indexTmpl, tagsTmpl, nil
}
