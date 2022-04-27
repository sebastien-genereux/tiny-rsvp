package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"strconv"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// eventConfig is a struct used to unpack the event configuration yaml
type eventConfig struct {
	EventName	string `yaml:"Event-name"`
	Date		string `yaml:"Date"`
	Time		string `yaml:"Time"`
	Location	string `yaml:"Location"`
	Description	string `yaml:"Description"`
	HostName	string `yaml:"Host-name"`
	Contact		string `yaml:"Contact"`
	StartStr	string `yaml:"RSVP-start"`
	EndStr		string `yaml:"RSVP-end"`
	StartTime	time.Time
	EndTime		time.Time
}

/* SETTINGS */

// totalHeader defines the label used to keep track of the total attendees at the top of the database file
const totalHeader string = "Total number of attendees: "

// configPath contains the relative or absolute path to the event configuration being used to parse event details
const configPath string = "./configs/event_example.yml"

// serverPort indicates which tcp port to listen on and server http
const serverPort string = "8080"

/* GLOBALS */ 

// configData saves the event configuration as global state to share with the http handlers. It is initialized in main, before any http is served.
var configData eventConfig

// databaseFile stores a pointer to the database file as global state to share with the http handlers. It is initialized in main, before any http is served.
var databaseFile *os.File

// headerReg stores a regular expression used to match on the total attendee header and submatch on the exact current count stored. It is saved as part of the global scope to share with the http handlers. It is compiled in main, before any http is served.
var headerReg *regexp.Regexp

func rsvpHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the appropriate page based on the time relative to the RSVP window
	var t *template.Template
	now := time.Now()
	if now.Before(configData.StartTime) {
		t, _ = template.ParseFiles("./web/templates/early.html")
	}else if now.After(configData.EndTime){
		t, _ = template.ParseFiles("./web/templates/late.html")
	}else {
		t, _ = template.ParseFiles("./web/templates/rsvp.html")
	}

	t.Execute(w, configData)
}

func thanksHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("./web/templates/thanks.html")
	t.Execute(w, nil) 
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Record submission details in the database
	fmt.Fprintf(databaseFile, "%s,%s,%s,%s\n", r.FormValue("name"), r.FormValue("count"), r.FormValue("contact"), r.FormValue("comments"))

	// Update the total attendee header, by reading the first line, grabbing the first submatch (i.e the current count) and adding the new attendee count. Lastly update the header. Note this stops counting past '9'*20 attendees due to the byte array size
	currentHeader := make([]byte, len(totalHeader) + 20)
	_, err := databaseFile.ReadAt(currentHeader, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	total, err2 := strconv.Atoi(headerReg.FindStringSubmatch(string(currentHeader))[1])
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
        return
	}

	oldTotal, err3 := strconv.Atoi(r.FormValue("count"))
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
        return
	}

	newTotal := strconv.Itoa(total + oldTotal)
	databaseFile.WriteAt([]byte(totalHeader + newTotal + "\n"), 0)

	http.Redirect(w, r, "/thanks", http.StatusFound)
}

func main() {
	// Get current time a unix timestamp for file naming
	timeStr := strconv.Itoa(int(time.Now().Unix()))

	// Read the config file
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Parse the configuration file and unmarshal into global config variable
	err = yaml.Unmarshal(configFile, &configData)
	if err != nil {
		log.Fatalln(err)
	}

	// Convert the rfc3339 time strings into time structs
	startTime, err := time.Parse(time.RFC3339, configData.StartStr)
	endTime, err2 := time.Parse(time.RFC3339, configData.EndStr)
	if err != nil || err2 != nil{
		log.Fatalln(err, err2)
	}
	configData.StartTime, configData.EndTime = startTime, endTime
	
	// Create a timestamped database file and save a pointer to the global database variable
	databaseName := "./databases/" + timeStr + "_rsvp_" + configData.EventName + ".csv"
	databaseLocal, err2 := os.Create(databaseName)
	if err2 != nil {
		log.Fatalln(err2)
	}
	databaseFile = databaseLocal
	databaseFile.Chmod(0664)
	defer databaseFile.Close()

	// Write the header info in the database file
	databaseFile.WriteString(totalHeader + "0\n")
	fmt.Fprintf(databaseFile, "%s,%s,%s,%s\n", "Family Name", "Count of Attendees", "Contact Info", "Comments")

	// Compile the header regex for later use by the submit handler
	headerReg = regexp.MustCompile(totalHeader + "([0-9]+)")

	// Serve the web application
    http.HandleFunc("/", rsvpHandler)
	http.HandleFunc("/submit", submitHandler)
    http.HandleFunc("/thanks", thanksHandler)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	log.Printf("Starting server on 127.0.0.1:%s\n", serverPort)

	log.Fatalln(http.ListenAndServe(":" + serverPort, nil))
}