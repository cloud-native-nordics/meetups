package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type EventData struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Duration int64      `json:"duration"`
	Date     string     `json:"local_date"`
	Time     string     `json:"local_time"`
	Venue    EventVenue `json:"venue"`
	RVSPs    uint64     `json:"yes_rsvp_count"`
}

func (ev *EventData) GetTime() (*Time, error) {
	dateTime := fmt.Sprintf("%sT%s:00Z", ev.Date, ev.Time)
	d, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return nil, err
	}
	return &Time{d}, nil
}

type EventVenue struct {
	Address string `json:"address_1"`
}

func setMeetupData(cfg *Config) error {
	for i := range cfg.MeetupGroups {
		mg := &cfg.MeetupGroups[i]
		if mg.MeetupID == "" {
			continue
		}
		if mg.AutogenMeetupGroup == nil {
			mg.AutogenMeetupGroup = &AutogenMeetupGroup{}
		}
		if mg.AutoMeetups == nil {
			mg.AutoMeetups = map[string]AutogenMeetup{}
		}
		events, err := GetMeetupEvents(mg.MeetupID)
		if err != nil {
			return err
		}
		for _, ev := range events {
			t, err := ev.GetTime()
			if err != nil {
				return err
			}
			dateStr := t.YYYYMMDD()
			for _, ignoreDate := range mg.IgnoreMeetupDates {
				if ignoreDate == dateStr {
					continue
				}
			}
			// Continue if this automatically generated meetup already exists
			if _, ok := mg.AutoMeetups[dateStr]; ok {
				continue
			}
			meetup := AutogenMeetup{}
			id, _ := strconv.Atoi(ev.ID)
			meetup.ID = uint64(id)
			meetup.Date = *t
			meetup.Name = ev.Name
			meetup.Address = ev.Venue.Address
			meetup.Duration = Duration{time.Duration(ev.Duration * 1000 * 1000)}

			if time.Now().UTC().After(meetup.Date.Time) {
				meetup.Attendees = ev.RVSPs
			} else {
				meetup.Attendees = 0
			}
			mg.AutoMeetups[t.YYYYMMDD()] = meetup
		}
		mg.ApplyGeneratedData()
	}
	return nil
}

type MeetupGroupAPI struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	City        string `json:"untranslated_city"`
	Country     string `json:"localized_country_name"`
	Members     uint64 `json:"members"`
	Photo       Photo  `json:"key_photo"`
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

type Attendance struct {
	Member AttendanceMember `json:"member"`
	RSVP   AttendanceRSVP   `json:"rsvp"`
}

type AttendanceMember struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type AttendanceRSVP struct {
	Response string `json:"response"`
	Guests   uint64 `json:"guests"`
}

func GetAttendanceList(meetupGroupID string, meetupID uint64) ([]Attendance, error) {
	url := fmt.Sprintf("https://api.meetup.com/%s/events/%d/attendance?&sign=true&photo-host=public&page=20", meetupGroupID, meetupID)
	att := []Attendance{}
	if err := GetJSON(url, &att); err != nil {
		return nil, err
	}
	return att, nil
}

func GetJSON(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("GetJSON failed for url %s with error %v", url, err)
	}
	return nil
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

func aggregateStats(cfg *Config) (*StatsFile, error) {
	s := &StatsFile{
		MeetupGroups: uint64(len(cfg.MeetupGroups)),
		PerMeetup:    map[string]MeetupStats{},
	}

	var wg sync.WaitGroup
	wg.Add(len(cfg.MeetupGroups))

	for _, mg := range cfg.MeetupGroups {
		go func(mg MeetupGroup) {
			defer wg.Done()

			mgStat := MeetupStats{}
			mgStat.Members = mg.members
			totalRSVPs := uint64(0)
			// allRSVPs maps an user ID to the amount of RSVPs for that user
			allRSVPs := map[uint64]uint64{}
			sponsors := map[string]bool{}
			speakers := map[string]bool{}
			priorMeetups := uint64(0)
			for _, m := range mg.Meetups {
				if m.Date.UTC().After(time.Now().UTC()) {
					continue
				}
				priorMeetups++
				totalRSVPs += m.Attendees
				for _, pres := range m.Presentations {
					for _, s := range pres.Speakers {
						if _, ok := speakers[string(s.ID)]; !ok {
							speakers[string(s.ID)] = true
						}
					}
				}
				for _, c := range append(m.Sponsors.Other, m.Sponsors.Venue) {
					if c == nil {
						continue
					}
					if _, ok := sponsors[string(c.ID)]; !ok {
						sponsors[string(c.ID)] = true
					}
				}

				attendance, err := GetAttendanceList(mg.MeetupID, m.ID)
				if err != nil {
					log.Fatalf("error: %v", err)
					return
				}
				for _, attendee := range attendance {
					if attendee.RSVP.Response != "yes" {
						continue
					}
					rsvps, ok := allRSVPs[attendee.Member.ID]
					if ok {
						allRSVPs[attendee.Member.ID] = rsvps + attendee.RSVP.Guests
					} else {
						allRSVPs[attendee.Member.ID] = 1 + attendee.RSVP.Guests
					}
				}
			}
			mgStat.Speakers = uint64(len(speakers))
			mgStat.Sponsors = uint64(len(sponsors))
			mgStat.TotalRSVPs = totalRSVPs
			if priorMeetups > 0 {
				mgStat.Meetups = priorMeetups
				mgStat.AverageRSVPs = uint64(math.Floor(float64(totalRSVPs / priorMeetups)))
			}
			for _, num := range allRSVPs {
				mgStat.UniqueRSVPs += num
			}

			s.PerMeetup[mg.CityLowercase()] = mgStat

			s.AllMeetups.Meetups += mgStat.Meetups
			s.AllMeetups.Members += mgStat.Members
			s.AllMeetups.TotalRSVPs += mgStat.TotalRSVPs
			s.AllMeetups.UniqueRSVPs += mgStat.UniqueRSVPs
		}(mg)
	}
	wg.Wait()
	s.AllMeetups.Sponsors = uint64(len(cfg.Sponsors))
	s.AllMeetups.Speakers = uint64(len(cfg.Speakers))
	s.AllMeetups.AverageRSVPs = uint64(math.Floor(float64(s.AllMeetups.TotalRSVPs / s.AllMeetups.Meetups)))
	return s, nil
}
