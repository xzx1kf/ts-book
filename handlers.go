package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

type Booking struct {
	Time	string
	Court	string
}

func BookCourt(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	days := q.Get("days")
	court := q.Get("court")
	hour := q.Get("hour")
	min := q.Get("min")
	timeslot := q.Get("timeSlot")

	// create a cookiejar this is required because the website uses cookies
	// and without it the booking of a court fails. 
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	c := &http.Client{
		Jar: jar,
	}

	// Create the get request to retrieve the court booking page.
	req, err := http.NewRequest("GET", "http://tynemouth-squash.herokuapp.com/bookings/new?" +
		"court=" + court +
		"&days=" + days +
		"&hour=" + hour +
		"&min=" + min +
		"&timeSlot=" + timeslot, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Perform the get request.
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Use goquery to parse the court booking page.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	token, time := ParseCourtBookingPage(doc)

	// Create the parameters for the POST request utilising the values 
	// parsed from the booking page.
	v := url.Values{}
	v.Set("utf8", "&#x2713;")
	v.Set("authenticity_token", token)
	v.Set("booking[full_name]", "Nick Hale")
	v.Set("booking[membership_number]", "s119")
	v.Set("booking[vs_player_name]", "")
	v.Set("booking[booking_number]", "1")
	v.Set("booking[start_time]", time)
	v.Set("booking[time_slot_id]", timeslot)
	v.Set("booking[court_time]", "40")
	v.Set("booking[court_id]", court)
	v.Set("booking[days]", days)
	v.Set("commit", "Book Court")

	// Create the POST request.
	req, err = http.NewRequest("POST", "http://tynemouth-squash.herokuapp.com/bookings", strings.NewReader(v.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Perform the POST request
	resp, err = c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}

	// Create a JSON response which indicates the time and court booked.
	b := Booking{time, court}
	js, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func ParseCourtBookingPage(doc *goquery.Document) (token string, time string) {
	s := doc.Find("form.booking")
	s.Find("input").Each(func(i int, sel *goquery.Selection) {
		input, _ := sel.Attr("name")
		if (input == "authenticity_token") {
			token, _ = sel.Attr("value")
		} else if (input == "booking[start_time]") {
			time, _ = sel.Attr("value")
		}
	})

	return token, time
}
