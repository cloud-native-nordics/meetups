Cloud Native Nordics is an organisation for cloud native (CNCF) Meetup Groups in the nordic countries to collaborate . Participating Meetup Groups can apply for CNCF membership and become official Cloud Native Computing Foundation Meetup Groups. The details for applying can be found at [https://github.com/cncf/meetups#how-to-apply](https://github.com/cncf/meetups#how-to-apply) once the group is created and has proven to be an active group.

Below you can find a “how to” guide for getting started with the creation of a cloud native meetup group. Please join the cloud-native-nordics slack via [www.cloudnativenordics.com](www.cloudnativenordics.com).

# Prepare the Meetup Group creation
The premise is that you want to create a Meetup Group in a nordic city, where no cloud-native meetups connected to Cloud Native Nordics already exist. The easiest way to see what meetups already exist under the Cloud Native Nordics brand, is probably to look at the [list](https://github.com/cloud-native-nordics/meetups/blob/master/README.md). 
If you do not find your city in the list please suggest the creation of a new meetup in your city in our Slack channel and there will be people that responds to that suggestion and help you getting started.

Once the decision for the creation of the meetup group has been decided, there are a couple of things that needs to be addressed.

# Create artwork for your meetup group
Initially you go to the [artwork repository](https://github.com/cloud-native-nordics/artwork) and create the background artwork for your future meetup group. Follow the artwork guide and remember to set the permissions for your API key. The end result from following the [guide](https://github.com/cloud-native-nordics/artwork/blob/master/slide-background/generator/README.md) is a picture of your &lt;city&gt; with streets overlayed onto the background and a Cloud Native Computing Foundation logo. 

The easiest way to get started is to fork the repo and create a branch with your artwork and create a pull-request for adding your city in the form of a &lt;city&gt;.png, you may want to consider creating a couple of frontpage-images for your Meetup and a couple of images for events.

# Create the Meetup on Meetup.com
Go to [https://meetup.com](https://meetup.com) and create an account there. Once you have that you can continue with the remaining tasks. 

* Create the meetup group on meetup.com
* The naming used is Cloud-Native-&lt;City&gt; or Kubernetes-&lt;City&gt;
* Find some local people that will help you with the creation of the meetup.com group
* Follow the guides for creating a new [meetup group](https://help.meetup.com/hc/en-us/articles/360002882111-Starting-a-Meetup-group).
* Follow the guide for creating the [first event](https://help.meetup.com/hc/en-us/articles/360002881251).
* Remember to use the artwork for the Meetup group and the event.

#Create the Meetup in Cloud Native Nordics
To finish the setup for the Cloud Native Nordics and to be able to have your meetup listed as a meetup under the Cloud Native Nordics umbrella the meetup is created in the [meetup repo](https://github.com/cloud-native-nordics/meetups). 

The easiest way to get started is to fork the [meetup repo](https://github.com/cloud-native-nordics/meetups) and create a branch with your meetup, create the &lt;city folder&gt;, please note that the location you chose for the for the meetup group on meetup.com and the folder name must be the same.  
Once you have created the folder you can create the _meetup.yaml_ file inside the &lt;city folder&gt;, fill in the potential missing speakers and organisation in the speakers.yaml and companies.yaml. 

* Create the folder named &lt;city&gt; 
* Creating a minimal _meetup.yaml_ file in the &lt;city&gt; folder.

Fill in the _meetup.yaml_ file:

``` yaml

city: <your city>
country: <your country>
meetupID: Cloud-Native-<city> - the name of the group from url
name: Cloud-Native-<City> / Kubernetes


```
The central file is the file mentioned above _meetup.yaml_ which is placed in the &lt;city&gt; folder. 

Now you should be ready for the initial population of information from meetup.com into the _meetup.yaml_ file in &lt;your city&gt; folder. To try that please run: 

```
$ make generate
```

Then you should see stuff being populated into the _meetup.yaml_ file, this is derived from the meetup, using the API offered by meetup.com, e.g.

``` yaml

city: <your city>
country: <your country>
meetupID: Cloud-Native-<city> - the name of the group from url
meetups: null  
name: <Cloud-Native-<City> / Kubernetes
organizers: null


```
If you have e.g. a single meetup event defined on meetup.com, it will pick parts of that as well and you will see e.g.:

``` yaml

city: <your city>
country: <your country>
meetupID: Cloud-Native-<city> - the name of the group from url 
meetups:
- address: <address for the event>
  date: "<date and time for event>"
  duration: <duration of the event>
  id: <meetup event id> - same as the url on meetup.com/<meetup>/events/<id>
  name: '#1 - The title of the initial Meetup event'
  presentations: null
  sponsors:
    other: null
    venue: null
name: Cloud-Native-<City> / Kubernetes
organizers: null


``` 
You can continue adding organizers etc. as seen below by hand.

``` yaml

city: <your city>
country: <your country>
meetupID: Cloud-Native-<City> - the name of the group from url 
meetups:
- address: <address for the event>
  date: "<date and time for event>"
  duration: <duration of the event>
  id: <meetup event id> - same as the url on meetup.com/<meetup>/events/<id>
  name: '#1 - The title of the initial Meetup event'
  presentations: null
		....
 sponsors:
   other: <sponsor for the event -please ensure included in companies.yaml>

   venue: <venue for the event - please ensure that its in companies.yaml>
name: Cloud Native <City>
organizers:
- <organizerid>  - remember to check/create the “organizer” <id> in speakers.yaml
- <co-organizerid>  - remember to check/create the “co-organizer” <id> in speakers.yaml


```
Take a look in the repo in github and see examples from the other cities and their _meetup.yaml_ files.

In order for it to be able to generate the remaining files for your meetup and add the meetup to the cloud-native-nordics group,  it is necessary to make sure that speakers, organizers, co-organizers etc are registered in the _speakers.yaml_ file and the venue and organizer’s/speaker’s company is registered in the _companies.yaml_ file. 

Generation of the _README.md_ file inside the &lt;city&gt; folder, the _config.json_ file etc. is done by invoking the *Makefile*, the build requires currently docker or go 1.12:
The command:

```
$ make generate   note: Other targets can be seen in the Makefile.
```

Once you can run the make generate without errors, you have a _README.md_ file in the &lt;city&gt; folder and a _config.json_

Commit and push your changes into your local fork and create  a pull-request for adding your city.

# Adding a new event to the _meetup.yaml_

You create a new event on meetup.com and add an event host etc. after that you run:

```
$ make generate
```

A new meetup is added to the meetup.yaml file in the &lt;city&gt; folder, here the example is taken from Tampere Finland.

``` yaml

.....
meetups:
.....
- address: Kelloportinkatu 1 D
 date: "2019-06-13T17:00:00Z"
 duration: 3h0m0s
 id: 261344385
 name: 'Summer Kubernetes Tampere Meetup'
 presentations: null
 sponsors:
   other: null
   venue: null
.....


```
Then you can add the meetup contents etc. in the _meetup.yaml_ file and complete the entries by hand, e.g. adding the following under 

``` yaml


meetups:
  .....
 name: Summer Kubernetes Tampere Meetup
 presentations:
 - duration: 10m0s
   slides: “”
   speakers: null
   title: Arrive to the venue, sit down and network with others ;)
 - duration: 5m0s
   slides: ""
   speakers: null
   title: Introductionary words from the venue sponsor for this time futurice
 - duration: 15m0s
   slides: “”
   speakers:
   - luxas
   title: Updates from the Cloud Native Nordics Community, Lucas Käldström
 - duration: 25m0s
   slides: “”
   speakers:
   - sergeysedelnikov
   title: Kubernetes as a service in Azure (Azure Kubernetes Service) for Real-Time
     API with NATS and HEMERA
 - delay: 5m0s
   duration: 30m0s
   slides: ""
   speakers: null
   title: Networking, food, drinks
 - duration: 25m0s
   slides: “”
   speakers:
   - carolchen
   title: Reflections from my first KubeCon - communities, operators, and more
 - delay: 5m0s
   duration: 25m0s
   slides: “”
   speakers:
   - cihanbebek
   title: Serverless - A natural step in DevOps thinking
 - delay: 5m0s
   duration: 30m0s
   slides: ""
   speakers: null
   title: Networking
 sponsors:
   other:
   - luxaslabs
   - cncf
   venue: futurice
   


```

Where the name of the speaker is equal to the speaker's name in the _speaker.yaml_ file

Examples:

``` yaml


- company: futurice
  countries:
  - finland
  email: cihan.m.bebek@gmail.com
  github: Keksike
  id: cihanbebek
  name: Cihan Bebek
  speakersBureau: ""
- company: redhat
  countries:
  - finland
  email: carol.chen@redhat.com
  github: cybette
  id: carolchen
  name: Carol Chen
  speakersBureau: ""


```
 
and the name of the speakers company is equal to the company name in the _companies.yaml_ file. 

Examples:

```yaml

- countries:
  - finland
  id: futurice
  logoURL: data:image/svg+xml;base64,PD94bWwgdm...c+CiAgICA8L2c+Cjwvc3ZnPg==
  name: Futurice
  websiteURL: https://www.futurice.com/
-  countries:
  - finland
  id: redhat
  logoURL: http://videos.cdn.redhat.com/...-RGB.svg
  name: Red Hat
  websiteURL: https://www.redhat.com


```

Remember to run:

```
$ make generate
```

That will generate the agenda into the _README.md_ file which includes the new meetup including agena. Using that makes it is easy to create the agenda on meetup.com as e.g:

```
Agenda:


 HH.mm - HH.mm Title - Speaker (Company)
 HH.mm - HH.mm Title - Speaker (Company)
 HH.mm - HH.mm Title

Example:
 Agenda:
18.15 - 18.30 Welcome from the host
18:30 - 18:55 Reflections from my first KubeCon - communities, operators, and more - Carol Chen (Red Hat)
19:00 - 19:25 Serverless: A natural step in DevOps thinking - Cihan Bebek (Futurice)
19:30 - 20:00 Networking


```

You will probably want to add some welcoming text, to the agenda, and a speaker bio in the agenda as well. Having a image with the highligths from the agenda as the featured photo on meetup.com for the event is nice for the invites and for sharing on SoMe’s. 

# After the Meetup Event
When the meetup event has finished you can add slides etc. to the _meetup.yaml_ file and then they will be listed in the _README.md_ in &lt;your-city&gt; folder.

```yaml

address: Kelloportinkatu 1 D
 attendees: 50
 date: "2019-06-13T17:00:00Z"
 duration: 3h0m0s
 id: 261344385
 name: Summer Kubernetes Tampere Meetup
 presentations:
 - duration: 10m0s
   slides: https://speakerdeck.com/luxas/kubernetes-and-cncf-meetup-tampere-june-2019
   speakers: null
   title: Arrive to the venue, sit down and network with others ;)
 - duration: 5m0s
   slides: ""
   speakers: null
   title: Introductionary words from the venue sponsor for this time futurice
 - duration: 15m0s
   slides: https://speakerdeck.com/luxas/kubernetes-and-cncf-meetup-tampere-june-2019
   speakers:
   - luxas
   title: Updates from the Cloud Native Nordics Community, Lucas Käldström
 - duration: 25m0s
   slides: https://sway.office.com/E8e1CyZyk9LCrhxY
   speakers:
   - sergeysedelnikov
   title: Kubernetes as a service in Azure (Azure Kubernetes Service) for Real-Time
     API with NATS and HEMERA
 - delay: 5m0s
   duration: 30m0s
   slides: ""
   speakers: null
   title: Networking, food, drinks
 - duration: 25m0s
   slides: https://www.slideshare.net/cybette/kubernetes-tampere-meetup-june-2019-community-operators-and-more
   speakers:
   - carolchen
   title: Reflections from my first KubeCon - communities, operators, and more
 - delay: 5m0s
   duration: 25m0s
   slides: https://drive.google.com/file/d/1mp2NvqBcRhDFvbf-z1OkJd4U7fh8waKE/view?usp=sharing
   speakers:
   - cihanbebek
   title: Serverless - A natural step in DevOps thinking
 - delay: 5m0s
   duration: 30m0s
   slides: ""
   speakers: null
   title: Networking
 recording: https://youtu.be/rnpMBRe0CY0
 sponsors:
   other:
   - luxaslabs
   - cncf
   venue: futurice


```

You complete these last  entries in _meetup.yaml_ by hand and run 

```
 $ make generate
```

after that you can check that everything is okay by using some of the other targets in the *Makefile*, such as e.g. _validate_ and after that push this to your forked repo and create a pull request for the event.

# Running an active group
Usually this means having an event or more per month, after running your meetup group for a while and having an active group, you can apply as mentioned above. The [guidelines are here](https://github.com/cncf/meetups#how-to-apply).

Your contributions and your experiences creating new groups will make this guide a better tool for the new groups to join, thus please help us making this as good as possible and accurate as necessary.

If you are interested in statistics you can populate information into the stats file by running:

```
 $ make stats
```

# Contributing to this guide
You are very welcome to contribute to this guide in order to make it as easy as possible to create a new meetup group and spread the good news about the Cloud Native technologies.

Thank you for reading. 