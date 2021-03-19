package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type PageDescription struct {
	Name  string   `yaml:"name"`
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
}

type Config struct {
	Title           string            `yaml:"title"`
	Src             string            `yaml:"source"`
	Dst             string            `yaml:"destination"`
	SiteRoot        string            `yaml:"site_root"`
	SiteRootLocal   string            `yaml:"site_root_local"`
	StylePath       string            `yaml:"styles"`
	TemplatesFolder string            `yaml:"templates"`
	Theme           string            `yaml:"theme"`
	Pages           []PageDescription `yaml:"pages"`

	IndexTitle string `yaml:"index_title"`
	GithubLink string `yaml:"github_link"`
}

func New(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config *Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) PageByFolder(folder string) (*PageDescription, bool) {
	for _, example := range c.Pages {
		if example.Name == folder {
			return &example, true
		}
	}

	return nil, false
}
