package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"html/template"
	"gopkg.in/yaml.v3"
)

type eventConfig struct {
	EventName	string `yaml:"Event-name"`
	Date		string `yaml:"Date"`
	Time		string `yaml:"Time"`
	Location	string `yaml:"Location"`
	Description	string `yaml:"Description"`
	HostName	string `yaml:"Host-name"`
	Contact		string `yaml:"Contact"`
}

func rsvpHandler(w http.ResponseWriter, r *http.Request, conf eventConfig) {
    t, _ := template.ParseFiles("./web/templates/rsvp.html")
	t.Execute(w, conf)  // should later swap nil for parametrized event details
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func makeRsvpHandler(fn func(http.ResponseWriter, *http.Request, eventConfig)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		configFile, err := ioutil.ReadFile("./configs/event.yml")

		if err != nil {
			log.Fatal(err)
		}
	
		var configData eventConfig
		err2 := yaml.Unmarshal(configFile, &configData)
	
		if err2 != nil {
			log.Fatal(err2)
		}

		fn(w, r, configData)
	}
}

func main() {
    http.HandleFunc("/", makeRsvpHandler(rsvpHandler))
	http.HandleFunc("/submit", submitHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}