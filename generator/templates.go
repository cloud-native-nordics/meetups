package main

import "text/template"

var (
	readmeTmpl   = template.Must(template.New("").Parse(readmeTmplStr))
	toplevelTmpl = template.Must(template.New("").Parse(toplevelTmplStr))
)

const (
	readmeTmplStr = `## Meetups organized in {{ .City }}

#### Organizers

{{ range .Organizers }} - {{ . }}
{{end}}{{if .CFP}}
#### Submit a talk

if you're interested in speaking in this meetup, fill out this form: {{.CFP}}
{{end}}{{ range .Meetups }}
### {{ .Name }}

 - Date: {{ .DateTime }}
 - Meetup link: https://www.meetup.com/{{ $.MeetupID }}/events/{{ .ID }}{{ if .Recording }}
 - Recording: {{ .Recording }}{{end}}{{ if .Attendees }}
 - Attendees (according to meetup.com): {{ .Attendees }}{{end}}
{{ if .Sponsors.Venue }} - Venue sponsor: [{{ .Sponsors.Venue.Name }}]({{ .Sponsors.Venue.WebsiteURL }}){{end}}
{{ if .Sponsors.Other }} - Meetup sponsors:
{{ range .Sponsors.Other }}   - [{{ .Name }}]({{ .WebsiteURL }})
{{end}}{{end}}
#### Agenda

{{ range .Presentations }} - {{ .StartTime }} - {{ .EndTime }}: {{ .Title }} {{ range .Speakers }}
   - {{ . }}{{end}}{{ if .Slides }}
   - Slides: {{ .Slides }}{{end}}{{ if .Recording }}
   - Recording: {{ .Recording }}{{end}}
{{end}}{{end}}`

	toplevelTmplStr = `# Cloud Native Nordics Meetups

Repository to gather all meetup information and slides from Cloud Native Nordic meetups:

{{ range .MeetupGroups }}* [{{ .City }}]({{ .CityLowercase }}/README.md){{ range .Organizers }}
  * {{ . }}{{end}}
{{end}}
## Join our Community!

To facilitate and help each other in between meetups and different geographical locations, we have set up a joined Slack Community.

In order to sign-up, go to [www.cloudnativenordics.com](https://www.cloudnativenordics.com) and enter your e-mail. Shortly hereafter you will receive an email with instructions to join the community.
`
)
