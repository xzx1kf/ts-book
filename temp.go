type Booking struct {
    Time string
    Court string
}

booking := Booking{"5", "21/08/2018 19:40"}

js, err := json.Marshal(booking)
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerErrror)
    return
}

w.Header().Set("Content-Type", "application/json; charset=utf-8")
w.WriteHeader(http.StatusOK)
w.Write(js)
