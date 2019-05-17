package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	globalSpeakerMap = map[SpeakerID]*Speaker{}
	globalCompanyMap = map[CompanyID]*Company{}
)

type CompanyID string
type SpeakerID string

type CompaniesFile struct {
	Companies []Company `json:"companies"`
}

func (c *CompaniesFile) SetGlobalMap() {
	for i, co := range c.Companies {
		globalCompanyMap[co.ID] = &c.Companies[i]
	}
}

type SpeakersFile struct {
	Speakers []Speaker `json:"speakers"`
}

func (s *SpeakersFile) SetGlobalMap() {
	for i, sp := range s.Speakers {
		globalSpeakerMap[sp.ID] = &s.Speakers[i]
	}
}

type MeetupGroupsFile struct {
	MeetupGroups []MeetupGroup `json:"meetupGroups"`
}

type Config struct {
	Companies    []Company     `json:"companies"`
	Speakers     []Speaker     `json:"speakers"`
	MeetupGroups []MeetupGroup `json:"meetupGroups"`
}

var _ json.Unmarshaler = &Company{}
var _ json.Unmarshaler = &Speaker{}

type Company struct {
	companyInternal
}

type companyInternal struct {
	ID         CompanyID `json:"id"`
	Name       string    `json:"name"`
	WebsiteURL string    `json:"websiteURL"`
	LogoURL    string    `json:"logoURL"`
	Countries  []string  `json:"countries"`
}

func (c *Company) UnmarshalJSON(b []byte) error {
	ctest := companyInternal{}
	if err := json.Unmarshal(b, &ctest); err == nil {
		c.companyInternal = ctest
		return nil
	}
	cid := CompanyID("")
	if err := json.Unmarshal(b, &cid); err == nil {
		*c = *globalCompanyMap[cid]
		return nil
	}
	return fmt.Errorf("couldn't marshal company")
}

type Speaker struct {
	speakerInternal
}

type speakerInternal struct {
	ID             SpeakerID `json:"id"`
	Name           string    `json:"name"`
	Title          string    `json:"title"`
	Email          string    `json:"email"`
	Company        *Company  `json:"company"`
	Github         string    `json:"github"`
	Twitter        string    `json:"twitter"`
	SpeakersBureau string    `json:"speakersBureau"`
}

func (s *Speaker) UnmarshalJSON(b []byte) error {
	stest := speakerInternal{}
	if err := json.Unmarshal(b, &stest); err == nil {
		s.speakerInternal = stest
		return nil
	}
	sid := SpeakerID("")
	if err := json.Unmarshal(b, &sid); err == nil {
		speaker, ok := globalSpeakerMap[sid]
		if !ok {
			return fmt.Errorf("speaker reference not found: %s", sid)
		}
		*s = *speaker
		return nil
	}
	return fmt.Errorf("couldn't marshal speaker")
}

type MeetupGroup struct {
	Name       string    `json:"name"`
	MeetupID   string    `json:"meetupID"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	Organizers []Speaker `json:"organizers"`
	Meetups    []Meetup  `json:"meetups"`
}

// CityLowercase gets the lowercase variant of the city
func (mg *MeetupGroup) CityLowercase() string {
	return strings.ToLower(mg.City)
}

type Meetup struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Date          time.Time      `json:"date"`
	Duration      time.Duration  `json:"duration"`
	DateInternal  string         `json:"dateInternal"`
	Recording     string         `json:"recording"`
	Attendees     uint32         `json:"attendees"`
	Address       string         `json:"address"`
	Sponsors      Sponsors       `json:"sponsors"`
	Presentations []Presentation `json:"presentations"`
}

type Presentation struct {
	StartTime string     `json:"startTime"`
	EndTime   string     `json:"endTime"`
	Title     string     `json:"title"`
	Slides    string     `json:"slides"`
	Recording string     `json:"recording"`
	Speakers  []*Speaker `json:"speakers"`
}

type Sponsors struct {
	Venue *Company  `json:"venue"`
	Other []Company `json:"other"`
}
