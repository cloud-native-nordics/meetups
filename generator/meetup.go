package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type IntOrString uint64

func (ios *IntOrString) UnmarshalJSON(b []byte) error {
	b = []byte(strings.ReplaceAll(string(b), `"`, ""))
	var i uint64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*ios = IntOrString(i)
	return nil
}

type EventData struct {
	ID       IntOrString `json:"id"`
	Name     string      `json:"name"`
	Duration int64       `json:"duration"`
	Date     string      `json:"local_date"`
	Time     string      `json:"local_time"`
	Venue    EventVenue  `json:"venue"`
	RVSPs    uint64      `json:"yes_rsvp_count"`
}

type EventVenue struct {
	Address string `json:"address_1"`
}

func setMeetupData(cfg *Config) error {
	for i, mg := range cfg.MeetupGroups {
		if mg.MeetupID == "" {
			continue
		}
		events, err := GetMeetupEvents(mg.MeetupID)
		if err != nil {
			return err
		}
		for _, ev := range events {
			found := false
			for _, meetup := range mg.Meetups {
				if meetup.ID == uint64(ev.ID) {
					found = true
				}
			}
			for _, ign := range mg.IgnoreMeetupIDs {
				if ign == uint64(ev.ID) {
					found = true
				}
			}
			if !found {
				cfg.MeetupGroups[i].Meetups = append(cfg.MeetupGroups[i].Meetups, Meetup{
					ID: uint64(ev.ID),
				})
			}
		}
		for _, ev := range events {
			for j, meetup := range cfg.MeetupGroups[i].Meetups {
				if meetup.ID != uint64(ev.ID) {
					continue
				}
				meetup.Name = ev.Name
				meetup.Address = ev.Venue.Address
				meetup.Attendees = ev.RVSPs
				meetup.Duration = Duration{time.Duration(ev.Duration * 1000 * 1000)}
				dateTime := fmt.Sprintf("%sT%s:00Z", ev.Date, ev.Time)
				d, err := time.Parse(time.RFC3339, dateTime)
				if err != nil {
					return err
				}
				meetup.Date = Time{d}
				cfg.MeetupGroups[i].Meetups[j] = meetup
			}
		}
	}
	return nil
}

type MeetupGroupAPI struct {
	ID      uint64 `json:"id"`
	Members uint64 `json:"members"`
	Photo   Photo  `json:"group_photo"`
}

type Photo struct {
	Link string `json:"highres_link"`
}

func GetMeetupInfo(meetupID string) (*MeetupGroupAPI, error) {
	url := fmt.Sprintf("https://api.meetup.com/%s", meetupID)
	mg := &MeetupGroupAPI{}
	if err := GetJSON(url, mg); err != nil {
		return nil, err
	}
	return mg, nil
}

func GetMeetupEvents(meetupID string) ([]EventData, error) {
	url := fmt.Sprintf("https://api.meetup.com/%s/events?sign=true&photo-host=public&page=100&status=past,upcoming", meetupID)
	ed := []EventData{}
	if err := GetJSON(url, &ed); err != nil {
		return nil, err
	}
	return ed, nil
}

func GetJSON(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func setPresentationTimestamps(m *Meetup) error {
	for i := range m.Presentations {
		p := &m.Presentations[i]
		var t time.Time
		if i == 0 {
			t = m.Date.Time
		} else {
			p2 := m.Presentations[i-1]
			t = p2.end
		}
		if p.Delay != nil {
			t = t.Add((*p.Delay).Duration)
		}
		p.start = t
		p.end = p.start.Add(p.Duration.Duration)
	}
	return nil
}
