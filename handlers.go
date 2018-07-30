package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func BookCourt(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	days := q.Get("days")
	court := q.Get("court")
	hour := q.Get("hour")
	min := q.Get("min")
	timeslot := q.Get("timeSlot")

	doc, _ := GetCourtBookingPage(days, court, hour, min, timeslot)

	token, time := ParseCourtBookingPage(doc)

	v := url.Values{}
	v.Set("authenticity_token", token)
	v.Add("booking[booking_number]", "1")
	v.Add("booking[full_name]", "Nick Hale")
	v.Add("booking[membership_number]", "s119")
	v.Add("booking[start_time]", time)
	v.Add("booking[time_slot_id]", timeslot)
	v.Add("booking[court_time]", "40")
	v.Add("booking[court_id]", court)
	v.Add("booking[days]", days)
	v.Add("utf8", "&#x2713;")
	v.Add("commit", "Book Court")

	fmt.Println(v)
	resp, err := http.PostForm("http://tynemouth-squash.herokuapp.com/bookings", v)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)
	fmt.Println("Success")
}

func GetCourtBookingPage(days string, court string, hour string, min string, timeSlot string) (*goquery.Document, error) {
	res, err := http.Get("http://tynemouth-squash.herokuapp.com/bookings/new?" +
		"court=" + court +
		"&days=" + days +
		"&hour=" + hour +
		"&min=" + min +
		"&timeSlot=" + timeSlot)
	if err != nil {
		fmt.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	return doc, err
}

func ParseCourtBookingPage(doc *goquery.Document) (token string, time string) {
	s1 := doc.Find("input#booking_start_time")
	s2, exist := s1.Attr("value")
	if exist {
		fmt.Println("Nick " +  s2)
	} else {
		fmt.Println("boo")
	}

	s := doc.Find("form.booking")
	s.Find("input").Each(func(i int, sel *goquery.Selection) {
		input, exists := sel.Attr("name")
		if (input == "authenticity_token") {
			token, exists = sel.Attr("value")
			fmt.Println(token)
		} else if (input == "booking[start_time]") {
			time, exists = sel.Attr("value")
			if exists {
				fmt.Println("days: " + time)
			} else {
				fmt.Println("BOOK")
			}
		}
	})


	return token, time
}


