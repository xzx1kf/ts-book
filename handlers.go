package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

func BookCourt (w http.ResponseWriter, r *http.Request) {

}

/* I think I should limit the scope of this service to simply scape the booking 
   web page and return the appropriate content in a JSON format. 
*/
func ScrapeNewBooking (bl string) {

	res, err := http.Get("http://tynemouth-squash.herokuapp.com/bookings/new?court=4&days=0&hour=14&min=20&timeSlot=69")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s",
			res.StatusCode,
			res.Status)
	}

}
