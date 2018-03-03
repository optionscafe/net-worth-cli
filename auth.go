//
// Date: 3/1/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

//
// If we do not have an access token in our .net-worth-cli we call this to auth.
//
func DoAuth() {

	// Ask questions of the user to get login information.
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What is the URL to your Net Worth Server (https://example.com): ")
	url, _ := reader.ReadString('\n')

	fmt.Print("What is your client_id to your Net Worth Server: ")
	clientId, _ := reader.ReadString('\n')

	fmt.Print("What is your email address? (john@example.com): ")
	email, _ := reader.ReadString('\n')

	fmt.Print("What is your password?: ")
	password, _ := reader.ReadString('\n')

	// Get the current user.
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	// Send request to the server to get an access token.
	access_token := SendLoginGetAccessToken(strings.TrimSuffix(url, "\n"), strings.TrimSuffix(email, "\n"), strings.TrimSuffix(password, "\n"), strings.TrimSuffix(clientId, "\n"))

	// WRite the config file.
	env, err := godotenv.Unmarshal("TIMEZONE=America/Los_Angeles\nSERVER_URL=" + url + "\nACCESS_TOKEN=" + access_token)

	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Write(env, usr.HomeDir+"/.net-worth-cli")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Auth Successful.")
}

//
// Send the login request to the server. And get an access token
//
func SendLoginGetAccessToken(url string, email string, password string, clientId string) string {
	var postStr = []byte(`{"username": "` + email + `","password": "` + password + `", "client_id": "` + clientId + `", "grant_type": "password"}`)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", url+"/oauth/token", bytes.NewBuffer(postStr))

	// Headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", `application/json; charset=utf-8`)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("response Status : ", resp.Status)
		fmt.Println("response Headers : ", resp.Header)
		fmt.Println("response Body : ", string(respBody))
		return ""
	}

	// Return the access token
	return gjson.Get(string(respBody), "access_token").String()
}

/* End File */
