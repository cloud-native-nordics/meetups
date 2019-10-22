package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

var (
	globalSpeakerMap        = map[SpeakerID]*Speaker{}
	globalCompanyMap        = map[CompanyID]*Company{}
	shouldMarshalCompanyID  = false
	shouldMarshalSpeakerID  = false
	shouldMarshalAutoMeetup = false
)

type CompanyID string
type SpeakerID string

type CompaniesFile struct {
	Sponsors []Company `json:"sponsors"`
	Members  []Company `json:"members"`
}

type SpeakersFile struct {
	Speakers []Speaker `json:"speakers"`
}

type StatsFile struct {
	MeetupGroups uint64                 `json:"meetupGroups"`
	AllMeetups   MeetupStats            `json:"allMeetups"`
	PerMeetup    map[string]MeetupStats `json:"perMeetup"`
}

type MeetupStats struct {
	Sponsors     uint64 `json:"sponsors"`
	Speakers     uint64 `json:"speakers"`
	Meetups      uint64 `json:"meetups"`
	Members      uint64 `json:"members"`
	TotalRSVPs   uint64 `json:"totalRSVPs"`
	AverageRSVPs uint64 `json:"averageRSVPs"`
	UniqueRSVPs  uint64 `json:"uniqueRSVPs"`
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

type SponsorRole string

var (
	SponsorRoleVenue    SponsorRole = "Venue"
	SponsorRoleLongterm SponsorRole = "Longterm"
	SponsorRoleCloud    SponsorRole = "Cloud"
	SponsorRoleFood     SponsorRole = "Food"
	SponsorRoleOther    SponsorRole = "Other"

	ValidSponsorRoles = map[SponsorRole]struct{}{
		SponsorRoleVenue:    {},
		SponsorRoleLongterm: {},
		SponsorRoleCloud:    {},
		SponsorRoleFood:     {},
		SponsorRoleOther:    {},
	}
)

func (c *SponsorRole) UnmarshalJSON(b []byte) error {
	str := ""
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	if _, ok := ValidSponsorRoles[SponsorRole(str)]; !ok {
		return fmt.Errorf("not a valid sponsor role: %q", str)
	}
	*c = SponsorRole(str)
	return nil
}

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
		if _, ok := globalCompanyMap[c.ID]; ok {
			log.Printf("duplicate company found: %q", c.ID)
		}
		globalCompanyMap[c.ID] = c
		return nil
	}
	cid := CompanyID("")
	err := json.Unmarshal(b, &cid)
	if err == nil {
		company, ok := globalCompanyMap[cid]
		if !ok {
			log.Fatalf("company reference not found: %s", cid)
		}
		*c = *company
		return nil
	}
	return fmt.Errorf("couldn't marshal company %q: %v", string(b), err)
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

func (s Speaker) String() string {
	str := s.Name
	if len(s.Github) != 0 {
		str += fmt.Sprintf(" [@%s](https://github.com/%s)", s.Github, s.Github)
	}
	if len(s.Title) != 0 {
		str += fmt.Sprintf(", %s", s.Title)
	}
	if s.Company != nil {
		str += fmt.Sprintf(", [%s](%s)", s.Company.Name, s.Company.WebsiteURL)
	}
	if len(s.SpeakersBureau) != 0 {
		str += fmt.Sprintf(", [Contact](https://www.cncf.io/speaker/%s)", s.SpeakersBureau)
	}
	return str
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
		if _, ok := globalSpeakerMap[s.ID]; ok {
			// TODO: Make this Fatal, after figuring out the combination of member/sponsor companies
			log.Printf("duplicate speaker found: %q", s.ID)
		}
		globalSpeakerMap[s.ID] = s
		return nil
	}
	sid := SpeakerID("")
	err := json.Unmarshal(b, &sid)
	if err == nil {
		speaker, ok := globalSpeakerMap[sid]
		if !ok {
			log.Fatalf("speaker reference not found: %s", sid)
		}
		*s = *speaker
		return nil
	}
	return fmt.Errorf("couldn't marshal speaker %q: %v", string(b), err)
}

type AutogenMeetupGroup struct {
	Photo       string                   `json:"photo,omitempty"`
	Name        string                   `json:"name"`
	City        string                   `json:"city"`
	Country     string                   `json:"country"`
	Description string                   `json:"description"`
	AutoMeetups map[string]AutogenMeetup `json:"-"`

	members uint64
}

type MeetupGroup struct {
	*AutogenMeetupGroup `json:",inline,omitempty"`

	MeetupID          string            `json:"meetupID"`
	Organizers        []*Speaker        `json:"organizers"`
	IgnoreMeetupDates []string          `json:"ignoreMeetupDates,omitempty"`
	CFP               string            `json:"cfpLink"`
	Latitude          float64           `json:"latitude"`
	Longitude         float64           `json:"longitude"`
	Meetups           map[string]Meetup `json:"meetups"`
	MeetupList        MeetupList        `json:"-"`
}

func (mg *MeetupGroup) ApplyGeneratedData() {
	for key := range mg.AutoMeetups {
		m, ok := mg.Meetups[key]
		if !ok {
			found := false
			for _, date := range mg.IgnoreMeetupDates {
				if key == date {
					found = true
				}
			}
			if !found {
				fmt.Printf("%s: didn't find meetup with date %q\n", mg.Name, key)
			}
			continue
		}
		autoMeetup := mg.AutoMeetups[key]
		m.AutogenMeetup = &autoMeetup
		mg.Meetups[key] = m
	}
}

// CityLowercase gets the lowercase variant of the city
func (mg *MeetupGroup) CityLowercase() string {
	return strings.ToLower(mg.City)
}

func (mg *MeetupGroup) SetMeetupList() {
	marr := []Meetup{}
	for _, m := range mg.Meetups {
		marr = append(marr, m)
	}
	mg.MeetupList = MeetupList(marr)
	sort.Sort(mg.MeetupList)
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

type AutogenMeetup struct {
	ID        uint64   `json:"id"`
	Name      string   `json:"name"`
	Date      Time     `json:"date,omitempty"`
	Duration  Duration `json:"duration,omitempty"`
	Attendees uint64   `json:"attendees,omitempty"`
	Address   string   `json:"address"`

	// rsvps map the user ID to how many rsvp's they used at this event (themselves + guests)
	rsvps map[uint64]uint64
}

type HumanMeetup struct {
	Recording     string          `json:"recording"`
	Sponsors      []MeetupSponsor `json:"sponsors"`
	Presentations []Presentation  `json:"presentations"`
}

type Meetup struct {
	*AutogenMeetup `json:",inline,omitempty"`
	HumanMeetup    `json:",inline"`
}

type fullMeetup struct {
	*AutogenMeetup `json:",inline,omitempty"`
	HumanMeetup    `json:",inline"`
}

func (m Meetup) MarshalJSON() ([]byte, error) {
	if shouldMarshalAutoMeetup {
		return json.Marshal(fullMeetup{
			AutogenMeetup: m.AutogenMeetup,
			HumanMeetup:   m.HumanMeetup,
		})
	}
	return json.Marshal(m.HumanMeetup)
}

type MeetupSponsor struct {
	Role    SponsorRole `json:"role"`
	Company *Company    `json:"company"`
}

func (m *Meetup) DateTime() string {
	d := m.Date.UTC()
	year, month, day := d.Date()
	hour, min, _ := d.Clock()
	hour2, min2, _ := d.Add(m.Duration.Duration).Clock()
	return fmt.Sprintf("%d %s, %d at %d:%02d - %d:%02d", day, month, year, hour, min, hour2, min2)
}

type Presentation struct {
	Duration  Duration   `json:"duration"`
	Delay     *Duration  `json:"delay,omitempty"`
	Title     string     `json:"title"`
	Slides    string     `json:"slides"`
	Recording string     `json:"recording,omitempty"`
	Speakers  []*Speaker `json:"speakers"`

	start time.Time
	end   time.Time
}

func (p *Presentation) StartTime() string {
	return fmt.Sprintf("%d:%02d", p.start.UTC().Hour(), p.start.UTC().Minute())
}

func (p *Presentation) EndTime() string {
	return fmt.Sprintf("%d:%02d", p.end.UTC().Hour(), p.end.UTC().Minute())
}
