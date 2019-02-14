package main

const (
	readmeTmpl = `## Meetups organized in {{ .City }}

#### Organizers

{{ range .Organizers }} - {{ .Name }}{{if .Github }} ([@{{ .Github }}])(https://github.com/{{ .Github }}){{end}}{{if .Title }}, {{ .Title }}{{end}}{{if .Company }}, {{ .Company.Name }}{{end}}{{if .Email }}, {{ .Email }}{{end}}
{{end}}
{{ range .Meetups }}### {{ .Name }}

 - Date: {{ .DateInternal }}
 - Meetup link: {{ $.MeetupURL }}/events/{{ .ID }}/
 - Recording: {{ .Recording }}
 - Attendees (according to meetup.com): {{ .Attendees }}
{{ if .Sponsors.Venue }} - Venue sponsor: [{{ .Sponsors.Venue.Name }}]({{ .Sponsors.Venue.WebsiteURL }}){{end}}
{{ if .Sponsors.Other }} - Meetup sponsors:
{{ range .Sponsors.Other }}   - [{{ .Name }}]({{ .WebsiteURL }})
{{end}}{{end}}
#### Agenda

{{ range .Presentations }} - {{ .StartTime }} - {{ .EndTime }}: {{ .Title }} {{ if .Speaker }}
   - {{ .Speaker.Name }}{{ if .Speaker.Title }}, {{ .Speaker.Title }}{{end}}{{ if .Speaker.Company }}, {{ .Speaker.Company.Name }}{{end}}{{end}}{{ if .Slides }}
   - Slides: {{ .Slides }}{{end}}
{{end}}{{end}}`
)
