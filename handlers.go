package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

func BookCourt(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	days := q.Get("days")
	court := q.Get("court")
	hour := q.Get("hour")
	min := q.Get("min")
	timeslot := q.Get("timeSlot")

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	
	c := &http.Client{
		Jar: jar,
	}
	
	req, err := http.NewRequest("GET", "http://tynemouth-squash.herokuapp.com/bookings/new?" +
		"court=" + court +
		"&days=" + days +
		"&hour=" + hour +
		"&min=" + min +
		"&timeSlot=" + timeslot, nil)
	
	
	doc, _ := GetCourtBookingPage(days, court, hour, min, timeslot)

	token, time := ParseCourtBookingPage(doc)

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
		
	req, err = http.NewRequest("POST", "http://tynemouth-squash.herokuap.com/bookings", strings.NewReader(v.Encode()))
	//resp, err := http.PostForm("http://tynemouth-squash.herokuapp.com/bookings", strings.NewReader(values))

	if err != nil {
		log.Fatal(err)
	}
	
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*.*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, defalte")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.8,en;q=0.6")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "tynemouth-squash.herokuapp.com")
	req.Header.Set("Origin", "http://tynemouth-squash.herokuapp.com")
	referer := "http://tynemouth-squash.herokuapp.com/bookings/new?court=" + court + "&days=" + days + "&hour=" + hour + "&min=" + min + "&timeSlot=" + timeslot
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux armv7l) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.84 Safari/537.36")
	req.Header.Add("Content-Length", strconv.Itoa(len(v.Encode())))
	req.Header.Set("Cookie", "_squash_session=OFdlSHcwSkVKZFg2dG94aGR3YWdxVlVuY3pkMFc5Mm5wWFJ5NEdPYWRlTWpVM2VzRHQ3Q05LRDA4NXpQSDJFOTMrZXByYTd5V0ZDWUhhNXBLZ2hCQkVzYVlRRS8vcEtqRTZMMHVPZW9pUlFycmpKQU1HaU5nWDFYRTVBeXh2Z2lTM2JQOE5UbUk2Z3FDZktweE9waklVS3BsMDZpR2dxeEFrU1d3Wi9BQ2Q3WEIxVmZGanU2RjFJRlVFQ2pER3hqLS1xYUxFYVhLSUc5Y3lPYUY2NHhPb3FBPT0%3D--eb7c6608b60eb4d36429c828cee126540147f39c")
	
	
	
	fmt.Println(req.Header)
	
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

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


