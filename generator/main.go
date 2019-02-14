package main // import "github.com/cloud-native-nordics/meetups/generator"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/pflag"
	"sigs.k8s.io/yaml"
)

var configFile = pflag.String("config", "meetups.yaml", "Point to the meetups")
var dryRun = pflag.Bool("dry-run", true, "Whether to actually apply the changes or not")
var validateFlag = pflag.Bool("validate", false, "Whether to validate the current state of the repo content with the spec")

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	pflag.Parse()
	if len(*configFile) == 0 {
		return fmt.Errorf("--config is a required argument")
	}
	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}
	cfg, err := load(b)
	if err != nil {
		return err
	}
	out, err := exec(cfg)
	if err != nil {
		return err
	}
	if *validateFlag {
		return validate(out)
	}
	return apply(out)
}

func load(b []byte) (*Config, error) {
	companiesObj := &Companies{}
	if err := yaml.UnmarshalStrict(b, companiesObj); err != nil {
		return nil, err
	}
	companiesObj.SetGlobalMap()
	speakersObj := &Speakers{}
	if err := yaml.UnmarshalStrict(b, speakersObj); err != nil {
		return nil, err
	}
	speakersObj.SetGlobalMap()
	meetupGroupsObj := &MeetupGroups{}
	if err := yaml.UnmarshalStrict(b, meetupGroupsObj); err != nil {
		return nil, err
	}
	return &Config{
		Speakers:     speakersObj.Speakers,
		Companies:    companiesObj.Companies,
		MeetupGroups: meetupGroupsObj.MeetupGroups,
	}, nil
}

func apply(files map[string][]byte) error {
	for city, fileContent := range files {
		readmePath := filepath.Join(strings.ToLower(city), "README.md")
		if err := writeFile(readmePath, fileContent); err != nil {
			return err
		}
	}
	return nil
}

func validate(files map[string][]byte) error {
	for city, fileContent := range files {
		readmePath := filepath.Join(strings.ToLower(city), "README.md")
		actual, err := ioutil.ReadFile(readmePath)
		if err != nil {
			return err
		}
		if !bytes.Equal(actual, fileContent) {
			return fmt.Errorf("%s differs from expected state. expected: \"%s\", actual: \"%s\"", readmePath, fileContent, actual)
		}
	}
	return nil
}

func exec(cfg *Config) (map[string][]byte, error) {
	tmpl, err := template.New("").Parse(readmeTmpl)
	if err != nil {
		return nil, err
	}
	result := map[string][]byte{}
	for _, mg := range cfg.MeetupGroups {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, mg); err != nil {
			return nil, err
		}
		result[mg.City] = buf.Bytes()
	}
	return result, nil
}

func writeFile(path string, b []byte) error {
	if *dryRun {
		fmt.Printf("Would write file %q with contents \"%s\"\n", path, string(b))
		return nil
	}
	return ioutil.WriteFile(path, b, 0644)
}
