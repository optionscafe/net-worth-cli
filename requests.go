//
// Date: 10/28/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
  "os"
  "errors"
  "net/http"
  "io/ioutil"
)

//
// Make a get request to server.
//
func MakeGetRequest(url string) (string, error) {

  // Setup http client
  client := &http.Client{}

  // Setup api request
  req, _ := http.NewRequest("GET", os.Getenv("SERVER_URL") + url, nil) 
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Bearer " + os.Getenv("ACCESS_TOKEN"))

  res, err := client.Do(req)

  if err != nil {
    return "", err
  } 

  // Close Body
  defer res.Body.Close()

  // Make sure the api responded with a 200
  if res.StatusCode == 404 {
    return "", errors.New("Server returned 404.")
  }

  // Read the data we got.
  body, _ := ioutil.ReadAll(res.Body)

  if err != nil {
    return "", err
  }  

  // Return happy
  return string(body), nil
}

/* End File */