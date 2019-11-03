package rc

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var DefaultPath = fmt.Sprintf("%s/.chirc", os.Getenv("HOME"))

type Config struct {
	savePath string
	current  string
	targets  []Target
}

type fileConfig struct {
	Current string   `yaml:"current"`
	Targets []Target `yaml:"targets"`
}

func (c Config) Current() *Target {
	return c.Find(c.current)
}

func (c *Config) SetCurrent(name string) error {
	if c.Find(name) == nil {
		return fmt.Errorf("No target with name `%s' exists", name)
	}

	c.current = name
	return nil
}

func (c *Config) Add(target Target) error {
	for c.Find(target.Name) != nil {
		c.Delete(target.Name)
	}

	c.targets = append(c.targets, target)
	return nil
}

func (c Config) Find(name string) *Target {
	for _, target := range c.targets {
		if target.Name == name {
			return &target
		}
	}

	return nil
}

func (c *Config) Delete(name string) {
	for i, target := range c.targets {
		if target.Name == name {
			c.targets[i], c.targets[len(c.targets)-1] = c.targets[len(c.targets)-1], c.targets[i]
			c.targets = c.targets[:len(c.targets)-1]
			break
		}
	}
}

type Target struct {
	Name       string `yaml:"name"`
	Address    string `yaml:"address"`
	Token      string `yaml:"token"`
	SkipVerify bool   `yaml:"skip_verify,omitempty"`
}

func Load(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("Could not open config at `%s': %s", path, err)
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Could not read from config (%s): %s", path, err)
	}

	unmarshalInto := fileConfig{}
	err = yaml.Unmarshal(fileContents, &unmarshalInto)
	if err != nil {
		return nil, fmt.Errorf("Could not parse config (%s) as YAML: %s", path, err)
	}

	conf := Config{
		savePath: path,
		current:  unmarshalInto.Current,
		targets:  unmarshalInto.Targets,
	}

	return &conf, nil
}

func (c *Config) Save() error {
	file, err := os.OpenFile(c.savePath, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Could not open config file at `%s' for writing: %s", c.savePath, err)
	}

	jEncoder := yaml.NewEncoder(file)
	toWrite := fileConfig{
		Current: c.current,
		Targets: c.targets,
	}
	err = jEncoder.Encode(&toWrite)
	if err != nil {
		return fmt.Errorf("Could not write YAML to file at `%s': %s", c.savePath, err)
	}

	return nil
}
