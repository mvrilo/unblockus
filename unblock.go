package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const usage = `usage - cli interface for the unblock-us service

Usage:
  unblock <email> [us|ca|uk|ie|mx|br|se|dk|no|fi|nl|ar|co]

Examples:
  unblock email@example.com # will simply reactivate the ip with the email
  unblock he@example.com ie # will reactivate and change the country to Ireland
`

const reactivateURL = "https://check.unblock-us.com/get-status.js?reactivate=1"
const countryURL = "http://realcheck.unblock-us.com/set-country.php"

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

func reactivate() bool {
	body := newRequest(reactivateURL)
	if bytes.Contains(body, []byte(`"is_active":true`)) {
		return true
	}
	return false
}

func setCountry() bool {
	body := newRequest(countryURL + "?code=" + os.Args[2])
	if bytes.Contains(body, []byte(`"current":"`+strings.ToUpper(os.Args[2])+`"`)) {
		return true
	}
	return false
}

func main() {
	switch len(os.Args) {
	case 1:
		fmt.Println(usage)
	case 2:
		if !reactivate() {
			fmt.Println("Reactivation didn't work")
			os.Exit(1)
		}
	case 3:
		if !reactivate() {
			fmt.Println("Reactivation didn't work")
			os.Exit(1)
		}
		if !setCountry() {
			fmt.Println("Changing country didn't work")
			os.Exit(1)
		}
	}

	fmt.Println("Done!")
	os.Exit(0)
}
