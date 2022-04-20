# tiny-rsvp

A web application to collect RSVPs for an event. Populate your event details in *configs/* and update the *configPath* setting at the top of *main.go*! - currently set to an example event. Timestamps for RSVP window start and end are expected in RFC3339 format.

The data entry form, as defined by *web/templates/rsvp.html* will ask for the attendees' family name, number of attendees, contact info and any comments (ex. allergies).

Those details are then saved to a timestamped .csv saved in databases/ which you will need to fetch (ex. via scp) to see your attendee list. At the top of the database file is a header, defined by *totalHeader* which keeps a running total of attendees for you. Input validation is done client side through form types in the *web/templates/rsvp.html*

Feel free to fork, add/change features as desired. The intent was the features minimal but useful.

# Deployment

*compose.sh* is provided as a convenience script to build and run tiny-rsvp as a distroless container. Simply clone this repo and run:

```
./compose.sh build
./compose.sh up
```

If you wish to change the listening port (currently set to 8080) ensure to update *serverPort* in main.go as well as the **PORT** build-arg in *compose.sh* 