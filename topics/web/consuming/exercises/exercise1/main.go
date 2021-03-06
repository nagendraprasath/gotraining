// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Call the GitHub API to get a list of repository contributors.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Create a type where we can decode contributor json values
// It needs the fields "login" and "contributions"

type contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

func main() {

	// Get an access token from the environment
	tkn := os.Getenv("GITHUB_TOKEN")
	if tkn == "" {
		log.Print("Token not found. You must set it in your environment like")
		log.Print("export GITHUB_TOKEN=000a0aaaa0000a00000000aaa00000000a000000")
		log.Print("You can generate a token at https://github.com/settings/tokens")
		os.Exit(1)
	}

	// Create a request for the contributors api endpoint
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/ardanlabs/gotraining/contributors", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add the access token in the "Authorization" header. The value should be like "token 000aa0a0..."
	req.Header.Set("Authorization", "token "+tkn)

	// Create an http.Client and make the request
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Defer closing the response body
	defer res.Body.Close()

	// Ensure we get a 200 OK status back
	if res.StatusCode != http.StatusOK {
		log.Println("API responded with", res.Status)
		io.Copy(os.Stderr, res.Body)
		os.Exit(1)
	}

	// Decode the results into a []contributor
	var c []contributor
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		log.Fatal(err)
	}

	// Loop through the []contributor and print values
	for i, c := range c {
		fmt.Println(i, c.Login, c.Contributions)
	}
}
