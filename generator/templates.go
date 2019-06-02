package main

import "text/template"

var (
	readmeTmpl   = template.Must(template.New("").Parse(readmeTmplStr))
	toplevelTmpl = template.Must(template.New("").Parse(toplevelTmplStr))
)

const (
	readmeTmplStr = `## Meetups organized in {{ .City }}

#### Organizers

{{ range .Organizers }} - {{ .Name }}{{if .Github }} ([@{{ .Github }}](https://github.com/{{ .Github }})){{end}}{{if .Title }}, {{ .Title }}{{end}}{{if .Company }}, [{{ .Company.Name }}]({{ .Company.WebsiteURL }}){{end}}{{if .Email }}, {{ .Email }}{{end}}
{{end}}{{ range .Meetups }}
### {{ .Name }}

 - Date: {{ .DateInternal }}
 - Meetup link: {{ $.MeetupURL }}/events/{{ .ID }}/
 - Recording: {{ .Recording }}
 - Attendees (according to meetup.com): {{ .Attendees }}
{{ if .Sponsors.Venue }} - Venue sponsor: [{{ .Sponsors.Venue.Name }}]({{ .Sponsors.Venue.WebsiteURL }}){{end}}
{{ if .Sponsors.Other }} - Meetup sponsors:
{{ range .Sponsors.Other }}   - [{{ .Name }}]({{ .WebsiteURL }})
{{end}}{{end}}
#### Agenda

{{ range .Presentations }} - {{ .StartTime }} - {{ .EndTime }}: {{ .Title }} {{ range .Speakers }}
   - {{ .Name }}{{ if .Title }}, {{ .Title }}{{end}}{{ if .Company }}, {{ .Company.Name }}{{end}}{{end}}{{ if .Slides }}
   - Slides: {{ .Slides }}{{end}}{{ if .Recording }}
   - Recording: {{ .Recording }}{{end}}
{{end}}{{end}}`

	toplevelTmplStr = `# Cloud Native Nordics Meetups

Repository to gather all meetup information and slides from Cloud Native Nordic meetups:

{{ range .MeetupGroups }}* [{{ .City }}]({{ .CityLowercase }}/README.md)
{{end}}
## Join our Community!

To facilitate and help each other in between meetups and different geographical locations, we have set up a joined Slack Community.

In order to sign-up, go to [www.cloudnativenordics.com](https://www.cloudnativenordics.com) and enter your e-mail. Shortly hereafter you will receive an email with instructions to join the community.
`
)
