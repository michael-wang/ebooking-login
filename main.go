package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

const lineLoginURL string = "https://access.line.me/oauth2/v2.1/authorize"

var lineChannelID string
var lineAuthCallback string

func init() {
	flag.StringVar(&lineChannelID, "channelID", "", "LINE channel ID.")
	flag.StringVar(&lineAuthCallback, "authCallback", "", "LINE OAuth Callback URL.")
}

func main() {
	flag.Parse()
	if len(lineChannelID) == 0 || len(lineAuthCallback) == 0 {
		log.Fatal("missing argument: channelID or authCallback")
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", loginPage)
	http.HandleFunc("/login/line", requestLineLogin)
	http.HandleFunc("/auth/line", showRequestHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

	http.ListenAndServe(":8000", nil)
}

func showRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "Remote Addr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("loginPage template err:", err)
	}
}

func requestLineLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

	req, err := http.NewRequest("GET", lineLoginURL, nil)
	if err != nil {
		log.Fatalf("http.NewRequest failed: %q", err)
	}

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("redirect_uri", lineAuthCallback)
	q.Add("client_id", lineChannelID)
	q.Add("state", base64.StdEncoding.EncodeToString([]byte(randStringRunes(8))))
	q.Add("scope", "profile")
	req.URL.RawQuery = q.Encode()
	url := req.URL.String()
	log.Println("request URL: ", url)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

/*
func lineAuthCallback(w http.ResponseWriter, r *http.Request) {
}
*/

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
