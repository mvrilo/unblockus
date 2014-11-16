/*
unblockus
Copyright (C) 2014, Murilo Santana <mvrilo@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

const USAGE = `unblockus - command line interface for the unblock-us service

Usage:
  unblockus <email> [%s]

Examples:
  unblockus email@example.com # will simply reactivate the ip with the email
  unblockus he@example.com ie # will reactivate and change the country to Ireland
`

const reactivateURL = "https://check.unblock-us.com/get-status.js?reactivate=1"
const countryURL = "http://realcheck.unblock-us.com/set-country.php"

var countries = []string{"us", "ar", "at", "be", "br", "ca", "co", "dk",
	"fi", "fr", "de", "ie", "lu", "mx", "nl", "no", "se", "ch", "uk"}

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
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		reactivate()
		wg.Done()
	}()

	if args > 2 {
		wg.Add(1)
		go func() {
			setCountry()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Done!")
	os.Exit(0)
}
