package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const USAGE = `unblock - cli interface for the unblock-us service

Usage:
  unblock <email> [%s]

Examples:
  unblock email@example.com # will simply reactivate the ip with the email
  unblock he@example.com ie # will reactivate and change the country to Ireland
`

const reactivateURL = "https://check.unblock-us.com/get-status.js?reactivate=1"
const countryURL = "http://realcheck.unblock-us.com/set-country.php"

var countries = []string{"us", "ca", "uk", "ie", "mx", "br", "se", "dk", "no", "fi", "nl", "ar", "co"}

func newRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	// :(
	req.AddCookie(&http.Cookie{
		Name:  "_stored_email_",
		Value: os.Args[1],
	})

	cli := new(http.Client)
	res, err := cli.Do(req)
	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	return body
}

func reactivate() {
	body := newRequest(reactivateURL)
	if bytes.Contains(body, []byte(`"is_active":true`)) {
		return
	}

	fmt.Println("Reactivation didn't work")
	os.Exit(1)
}

func setCountry() {
	body := newRequest(countryURL + "?code=" + os.Args[2])
	if bytes.Contains(body, []byte(`"current":"`+strings.ToUpper(os.Args[2])+`"`)) {
		return
	}

	fmt.Println("Changing country didn't work")
	os.Exit(1)
}

func validate() {
	if !strings.Contains(os.Args[1], "@") {
		usage()
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		return
	}

	var code string
	for _, c := range countries {
		if c == os.Args[2] {
			code = c
		}
	}

	if code == "" {
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(USAGE, strings.Join(countries, "|"))
}

func main() {
	args := len(os.Args)
	if args < 2 {
		usage()
		os.Exit(0)
	}

	validate()
	reactivate()

	if args > 2 {
		setCountry()
	}
	fmt.Println("Done!")
}
