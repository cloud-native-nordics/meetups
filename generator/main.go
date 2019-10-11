package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/pflag"
	"sigs.k8s.io/yaml"
)

var speakersFile = pflag.String("speakers-file", "speakers.yaml", "Point to the speakers.yaml file")
var companiesFile = pflag.String("companies-file", "companies.yaml", "Point to the companies.yaml file")
var rootDir = pflag.String("meetups-dir", ".", "Point to the directory that has all meetup groups as subfolders, each with a meetup.yaml file")
var dryRun = pflag.Bool("dry-run", true, "Whether to actually apply the changes or not")
var statsFlag = pflag.Bool("stats", false, "With this flag, the generator generates only the stats.json file")
var validateFlag = pflag.Bool("validate", false, "Whether to validate the current state of the repo content with the spec")
var isTesting = false
var unmarshal = yaml.UnmarshalStrict

// this maps the locations returned from meetup.com to what we want to use here.
// TODO: Maybe skip this and just use "Århus" directly in our
var cityNameExceptions = map[string]string{
	"Århus": "Aarhus",
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	pflag.Parse()
	cfg, err := load(*companiesFile, *speakersFile, *rootDir)
	if err != nil {
		return err
	}
	if err := update(cfg); err != nil {
		return err
	}
	if *statsFlag {
		return writeStats(cfg)
	}
	out, err := exec(cfg)
	if err != nil {
		return err
	}
	if *validateFlag {
		return validate(out, *rootDir)
	}
	return apply(out, *rootDir)
}

func load(companiesPath, speakersPath, meetupsDir string) (*Config, error) {
	companiesObj := &CompaniesFile{}
	companiesContent, err := ioutil.ReadFile(companiesPath)
	if err != nil {
		return nil, err
	}
	if err := unmarshal(companiesContent, companiesObj); err != nil {
		return nil, err
	}
	speakersObj := &SpeakersFile{}
	speakersContent, err := ioutil.ReadFile(speakersPath)
	if err != nil {
		return nil, err
	}
	if err := unmarshal(speakersContent, speakersObj); err != nil {
		return nil, err
	}
	meetupGroups := []MeetupGroup{}

	err = filepath.Walk(meetupsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		// Consider only subdirectories of the root path
		if !isTesting && filepath.Dir(path) != "." {
			return nil
		}
		meetupsFile := filepath.Join(path, "meetup.yaml")
		if _, err := os.Stat(meetupsFile); os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return err
		}
		mg := MeetupGroup{}
		mgContent, err := ioutil.ReadFile(meetupsFile)
		if err != nil {
			return err
		}
		if err := unmarshal(mgContent, &mg); err != nil {
			return err
		}
		city := filepath.Base(filepath.Dir(meetupsFile))
		autogenFile := filepath.Join(meetupsDir, autogenPath(city))
		autoMGContent, err := ioutil.ReadFile(autogenFile)
		if err != nil {
			return err
		}
		mg.AutogenMeetupGroup = &AutogenMeetupGroup{}
		if err := unmarshal(autoMGContent, mg.AutogenMeetupGroup); err != nil {
			return err
		}
		mg.ApplyGeneratedData()
		meetupGroups = append(meetupGroups, mg)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Config{
		Speakers:     speakersObj.Speakers,
		Sponsors:     companiesObj.Sponsors,
		Members:      companiesObj.Members,
		MeetupGroups: meetupGroups,
	}, nil
}

func apply(files map[string][]byte, rootDir string) error {
	for path, fileContent := range files {
		fullPath := filepath.Join(rootDir, path)
		if err := writeFile(fullPath, fileContent); err != nil {
			return err
		}
	}
	return nil
}

func validate(files map[string][]byte, rootDir string) error {
	for path, fileContent := range files {
		fullPath := filepath.Join(rootDir, path)
		actual, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return err
		}
		if !bytes.Equal(actual, fileContent) {
			return fmt.Errorf("%s differs from expected state. expected: \"%s\", actual: \"%s\"", fullPath, fileContent, actual)
		}
	}
	fmt.Println("Validation succeeded!")
	return nil
}

func tmpl(t *template.Template, obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func exec(cfg *Config) (map[string][]byte, error) {
	result := map[string][]byte{}
	shouldMarshalSpeakerID = true
	shouldMarshalCompanyID = true
	for _, mg := range cfg.MeetupGroups {
		mg.SetMeetupList()
		b, err := tmpl(readmeTmpl, mg)
		if err != nil {
			return nil, err
		}
		city := mg.CityLowercase()
		path := filepath.Join(city, "README.md")
		result[path] = b

		autoMGYAML, err := yaml.Marshal(*mg.AutogenMeetupGroup)
		if err != nil {
			return nil, err
		}
		path = autogenPath(city)
		result[path] = autoMGYAML
	}
	shouldMarshalSpeakerID = false
	shouldMarshalCompanyID = false
	companiesYAML, err := yaml.Marshal(CompaniesFile{
		Sponsors: cfg.Sponsors,
		Members:  cfg.Members,
	})
	if err != nil {
		return nil, err
	}
	result["companies.yaml"] = companiesYAML
	shouldMarshalCompanyID = true
	speakersYAML, err := yaml.Marshal(SpeakersFile{Speakers: cfg.Speakers})
	if err != nil {
		return nil, err
	}
	result["speakers.yaml"] = speakersYAML
	readmeBytes, err := tmpl(toplevelTmpl, cfg)
	if err != nil {
		return nil, err
	}
	result["README.md"] = readmeBytes
	shouldMarshalCompanyID = false
	// Don't output the autoMeetups thing in config.json
	for i, mg := range cfg.MeetupGroups {
		mg.AutoMeetups = nil
		cfg.MeetupGroups[i] = mg
	}
	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}
	result["config.json"] = configJSON
	return result, nil
}

func writeStats(cfg *Config) error {
	result := map[string][]byte{}
	stats, err := aggregateStats(cfg)
	if err != nil {
		return err
	}
	statsJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	result["stats.json"] = statsJSON
	*dryRun = false
	return apply(result, *rootDir)
}

func update(cfg *Config) error {
	if !isTesting {
		if err := setMeetupData(cfg); err != nil {
			return err
		}
	}
	for i := range cfg.MeetupGroups {
		mg := &cfg.MeetupGroups[i]
		if !isTesting {
			data, err := GetMeetupInfo(mg.MeetupID)
			if err != nil {
				return err
			}
			mg.members = data.Members
			mg.Photo = data.Photo.Link
			mg.Country = strings.ToLower(data.Country)
			mg.City = data.City
			if newName, ok := cityNameExceptions[data.City]; ok {
				mg.City = newName
			}
			mg.Name = data.Name
		}
		for _, s := range mg.Organizers {
			cfg.SetSpeakerCountry(s, mg.Country)
		}
		for j, m := range mg.Meetups {
			if err := setPresentationTimestamps(&m); err != nil {
				return err
			}
			for _, pres := range m.Presentations {
				for _, s := range pres.Speakers {
					cfg.SetSpeakerCountry(s, mg.Country)
				}
			}
			cfg.SetCompanyCountry(m.Sponsors.Venue, mg.Country)
			for _, s := range m.Sponsors.Other {
				cfg.SetCompanyCountry(s, mg.Country)
			}
			mg.Meetups[j] = m
		}
	}
	return nil
}

func writeFile(path string, b []byte) error {
	if *dryRun {
		fmt.Printf("Would write file %q with contents \"%s\"\n", path, string(b))
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

func autogenPath(city string) string {
	return fmt.Sprintf("autogenerated/%s.yaml", city)
}
