package wizard

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/i582/cfmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"

	"bye/static"
)

func Run(args cli.Args) {
	if args.Len() == 0 {
		fmt.Println("Please specify folder for new project")
		return
	}

	dir := args.First()

	fullPath, err := createProjectDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	copyStatic(fullPath)
	copyConfig(fullPath)
	copyExample(fullPath)
	copyIndex(fullPath)

	cfmt.Println("{{v}}::green The project was {{successfully}}::green created.")
}

func copyExample(fullPath string) {
	examplesDir := filepath.Join(fullPath, "pages")
	err := os.MkdirAll(examplesDir, 0777)
	if err != nil {
		log.Fatalf("error create folder '%s': %v", examplesDir, err)
	}

	pagesDir, err := static.Content.ReadDir("pages")
	for _, entry := range pagesDir {
		info, err := entry.Info()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		if !info.IsDir() {
			continue
		}

		exampleDir := "pages/" + info.Name()
		exampleDirFs, err := static.Content.ReadDir(exampleDir)

		err = os.MkdirAll(filepath.Join(fullPath, exampleDir), 0777)
		if err != nil {
			log.Fatalf("error create folder '%s': %v", exampleDir, err)
		}

		for _, entry := range exampleDirFs {
			info, err := entry.Info()
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			name := info.Name()
			filename := filepath.ToSlash(filepath.Join(exampleDir, name))
			exampleData := readFile(filename)

			exampleFile, err := os.Create(filepath.Join(fullPath, exampleDir, name))
			if err != nil {
				log.Fatalf("error create file '%s': %v", filename, err)
			}

			fmt.Fprint(exampleFile, string(exampleData))
		}
	}
}

func copyConfig(fullPath string) {
	configData := readFile("config/config.yml")

	configFileName := filepath.Join(fullPath, "bye.yml")
	configFile, err := os.Create(configFileName)
	if err != nil {
		log.Fatalf("error create config file '%s': %v", configFileName, err)
	}

	fmt.Fprint(configFile, string(configData))
}

func createProjectDir(dir string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("some unexpected error, try again: %v", err)
	}

	fullPath := filepath.Join(workingDir, dir)

	exist, err := exists(fullPath)
	if err != nil {
		return "", fmt.Errorf("some unexpected error, try again: %v", err)
	}

	if exist {
		// return "", fmt.Errorf("Folder '%s' already exists\n", dir)
	}

	err = os.MkdirAll(fullPath, 0777)
	if err != nil {
		return "", fmt.Errorf("error create folder '%s': %v", dir, err)
	}

	return fullPath, nil
}

func copyStatic(fullPath string) {
	err := copyTemplates(fullPath)
	if err != nil {
		log.Fatal(err)
	}

	err = copyStyles(fullPath)
	if err != nil {
		log.Fatal(err)
	}
}

func copyStyles(fullPath string) error {
	styleFileDir := filepath.Join(fullPath, "static", "styles")
	err := os.MkdirAll(styleFileDir, 0777)
	if err != nil {
		return fmt.Errorf("error create folder '%s': %v", styleFileDir, err)
	}

	styleData := readFile("styles/style.css")

	styleFileName := filepath.Join(styleFileDir, "style.css")
	styleFile, err := os.Create(styleFileName)
	if err != nil {
		return fmt.Errorf("error create file '%s': %v", styleFileName, err)
	}

	fmt.Fprint(styleFile, string(styleData))
	return nil
}

func copyIndex(fullPath string) {
	indexData := readFile("index.md")

	indexFileName := filepath.Join(fullPath, "pages", "index.md")
	indexFile, err := os.Create(indexFileName)
	if err != nil {
		log.Fatalf("error create file '%s': %v", indexFileName, err)
	}

	fmt.Fprint(indexFile, string(indexData))
}

func copyTemplates(fullPath string) error {
	templatesFileDir := filepath.Join(fullPath, "static", "templates")
	err := os.MkdirAll(templatesFileDir, 0777)
	if err != nil {
		return fmt.Errorf("error create folder 'static/templates': %v", err)
	}

	templatesDir, err := static.Content.ReadDir("templates")
	for _, entry := range templatesDir {
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}

		name := info.Name()
		filename := filepath.ToSlash(filepath.Join("templates", name))
		templateData := readFile(filename)

		templateFile, err := os.Create(filepath.Join(templatesFileDir, name))
		if err != nil {
			return fmt.Errorf("error create file '%s': %v", filename, err)
		}

		fmt.Fprint(templateFile, string(templateData))
	}

	return err
}

func readFile(name string) []byte {
	file, err := static.Content.Open(name)
	if err != nil {
		log.Fatalf("error open file: %v", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error open style file: %v", err)
	}
	return data
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func projectNamePrompt() string {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("project name must not be empty")
		}
		return nil
	}

	promptProjectName := promptui.Prompt{
		Label:       "Project name",
		Validate:    validate,
		HideEntered: true,
	}

	projectName, err := promptProjectName.Run()
	if err != nil {
		log.Fatalf("Some unexpected error, try again: %v", err)
	}

	return projectName
}
