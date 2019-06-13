package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

var (
	globalSpeakerMap       = map[SpeakerID]*Speaker{}
	globalCompanyMap       = map[CompanyID]*Company{}
	shouldMarshalCompanyID = false
	shouldMarshalSpeakerID = false
)

type CompanyID string
type SpeakerID string

type CompaniesFile struct {
	Sponsors []Company `json:"sponsors"`
	Members  []Company `json:"members"`
}

func (c *CompaniesFile) SetGlobalMap() {
	for i, co := range c.Sponsors {
		globalCompanyMap[co.ID] = &c.Sponsors[i]
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
	Sponsors     []Company     `json:"sponsors"`
	Members      []Company     `json:"members"`
	Speakers     []Speaker     `json:"speakers"`
	MeetupGroups []MeetupGroup `json:"meetupGroups"`
}

func (cfg *Config) SetSpeakerCountry(speaker *Speaker, country string) {
	if speaker == nil || country == "" {
		return
	}
	for i, s := range cfg.Speakers {
		if s.ID != speaker.ID {
			continue
		}
		found := false
		for _, c := range cfg.Speakers[i].Countries {
			if c == country {
				found = true
				break
			}
		}
		if !found {
			cfg.Speakers[i].Countries = append(cfg.Speakers[i].Countries, country)
		}
	}
	cfg.SetCompanyCountry(speaker.Company, country)
}

func (cfg *Config) SetCompanyCountry(company *Company, country string) {
	if company == nil || country == "" {
		return
	}
	for i, c := range cfg.Sponsors {
		if c.ID != company.ID {
			continue
		}
		found := false
		for _, c := range cfg.Sponsors[i].Countries {
			if c == country {
				found = true
				break
			}
		}
		if !found {
			cfg.Sponsors[i].Countries = append(cfg.Sponsors[i].Countries, country)
		}
	}
}

var _ json.Marshaler = &Company{}
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
	WhiteLogo  bool      `json:"whiteLogo,omitempty"`
	Countries  []string  `json:"countries"`
}

func (c Company) MarshalJSON() ([]byte, error) {
	if shouldMarshalCompanyID {
		return []byte(`"` + c.ID + `"`), nil
	}
	return json.Marshal(c.companyInternal)
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
	Title          string    `json:"title,omitempty"`
	Email          string    `json:"email"`
	Company        *Company  `json:"company"`
	Countries      []string  `json:"countries"`
	Github         string    `json:"github"`
	Twitter        string    `json:"twitter,omitempty"`
	SpeakersBureau string    `json:"speakersBureau"`
}

func (s Speaker) MarshalJSON() ([]byte, error) {
	if shouldMarshalSpeakerID {
		return []byte(`"` + s.ID + `"`), nil
	}
	return json.Marshal(s.speakerInternal)
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
	MeetupID        string     `json:"meetupID"`
	Name            string     `json:"name"`
	Members         uint64     `json:"members"`
	Photo           string     `json:"photo,omitempty"`
	City            string     `json:"city"`
	Country         string     `json:"country"`
	Organizers      []*Speaker `json:"organizers"`
	Meetups         MeetupList `json:"meetups"`
	IgnoreMeetupIDs []uint64   `json:"ignoreMeetupIDs,omitempty"`
}

// CityLowercase gets the lowercase variant of the city
func (mg *MeetupGroup) CityLowercase() string {
	return strings.ToLower(mg.City)
}

// MeetupList is a slice of meetups implementing sort.Interface
type MeetupList []Meetup

var _ sort.Interface = MeetupList{}

func (ml MeetupList) Len() int {
	return len(ml)
}

func (ml MeetupList) Less(i, j int) bool {
	return ml[i].Date.Time.After(ml[j].Date.Time)
}

func (ml MeetupList) Swap(i, j int) {
	ml[i], ml[j] = ml[j], ml[i]
}

type Meetup struct {
	ID            uint64         `json:"id"`
	Name          string         `json:"name"`
	Date          Time           `json:"date,omitempty"`
	Duration      Duration       `json:"duration,omitempty"`
	Recording     string         `json:"recording,omitempty"`
	Attendees     uint64         `json:"attendees"`
	Address       string         `json:"address"`
	Sponsors      Sponsors       `json:"sponsors"`
	Presentations []Presentation `json:"presentations"`
}

func (m *Meetup) DateTime() string {
	year, month, day := m.Date.Date()
	hour, min, _ := m.Date.Clock()
	hour2, min2, _ := m.Date.Add(m.Duration.Duration).Clock()
	return fmt.Sprintf("%d %s, %d at %d:%02d - %d:%02d", day, month, year, hour, min, hour2, min2)
}

type Presentation struct {
	StartTime string     `json:"startTime"`
	EndTime   string     `json:"endTime"`
	Title     string     `json:"title"`
	Slides    string     `json:"slides"`
	Recording string     `json:"recording,omitempty"`
	Speakers  []*Speaker `json:"speakers"`
}

type Sponsors struct {
	Venue *Company   `json:"venue"`
	Other []*Company `json:"other"`
}
