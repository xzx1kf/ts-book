package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func BookCourt(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	court := q.Get("court")
	hour := q.Get("hour")
	min := q.Get("min")
	timeslot := q.Get("timeSlot")

	doc, _ := GetCourtBookingPage(court, hour, min, timeslot)
	
	ParseCourtBookingPage(doc)

	v := url.Values{}
	v.Set("authenticity_token", token)
	v.Add("booking[booking_number]", "1")
	v.Add("booking[full_name]", "Nick Hale")
	v.Add("booking[membership_number]", "s119")
	v.Add("booking[start_time]", "yyyy-MM-dd HH:mm:ss zzz") // TODO
	v.Add("booking[time_slot_id]", timeslot)
	v.Add("booking[court_time]", "40")
	v.Add("booking[court_id]", court)
	v.Add("booking[days]", "21") // TODO
	v.Add("commit", "Book Court")

	//resp, err := http.PostForm("http://tynemouth-squash.herokuapp.com/bookings", v)

	/*
	if err != nill {
		log.Fatal(err)
	}
	*/
	fmt.Println("Success")
}

func GetCourtBookingPage(court string, hour string, min string, timeSlot string) (*goquery.Document, error) {
	res, err := http.Get("http://tynemouth-squash.herokuapp.com/bookings/new?" +
		"court=" + court +
		"&hour=" + hour +
		"&min=" + min +
		"&timeSlot=" + timeSlot)
	if err != nil {
		fmt.Println(err)

	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	return doc, err
}

func ParseCourtBookingPage(doc goquery.Document) {
	s := doc.Find("form.booking")

	s.Find("input").Each(func(i int, sel *goquery.Selection) {
		input, exists := sel.Attr("authenticity_token")
		if exists {
			fmt.Println("AT=" + input.Text())
		}
		
		input, exists := sel.Attr("booking_start_time")
		if exists {
			fmt.Println("start time=" + input.Text())
		}
		
		input, exists := sel.Attr("booking_days")
		if exists {
			fmt.Println("days=" + input.Text())
		}
	})
}
	

