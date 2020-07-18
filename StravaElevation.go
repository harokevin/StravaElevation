package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"bytes"
	"net/url"
	"time"
	"strings"
	"os"
)

func main() {
	stravaAccessToken := getStravaAccessToken()
	log.Println(stravaAccessToken)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, reqErr := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://www.strava.com/api/v3/athlete/activities?after=%s&before=%s&per_page=%s",
			"1593586799", //after Tuesday, June 30, 2020 11:59:59 PM GMT-07:00
			"1596265199", //before Friday, July 31, 2020 11:59:59 PM GMT-07:00
			"200", //per_page
		),
		nil,
	)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", stravaAccessToken))
	if reqErr != nil {
		log.Fatalln(reqErr)
	}

	response, respErr := client.Do(request)
	if respErr != nil {
		log.Fatalln(respErr)
	}

	defer response.Body.Close()

	var result []map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	var sum_total_elevation_gain float64 = 0
	newest_activity := make(map[string]interface{})
	newest_activity["start_date"] = ""
	for k, v := range result { 
		fmt.Printf("key[%d] value's id[%.0f] total_elevation_gain[%f]\n", k, v["id"], v["total_elevation_gain"])
		sum_total_elevation_gain += v["total_elevation_gain"].(float64)
		if v["start_date"].(string) > newest_activity["start_date"].(string) {
			newest_activity = v
		}
	}
	fmt.Printf("sum_total_elevation_gain(meters): %f\n", sum_total_elevation_gain)
	fmt.Printf("sum_total_elevation_gain(feet): %f\n", sum_total_elevation_gain*3.28084) // 1 Meter =  3.28084 Feet
	fmt.Printf("Newest Activity Id: %.0f\n", newest_activity["id"])

	if newest_activity["id"] != nil {
		existingDescription := newest_activity["description"]
		descriptionLabel := "Syd And Macky Everesting Challenge"
		descriptionText := fmt.Sprintf("%.0f/14514 ft %.0f%s complete", sum_total_elevation_gain*3.28084, sum_total_elevation_gain*3.28084/14514*100, "%")
		// TODO
		// Add expected completion percentage to description
		// Impact: Feature

		var description string

		if newest_activity["description"] != nil {
			description = fmt.Sprintf(
				"%s %s %s %s %s", 
				existingDescription, 
				"\n\n",
				descriptionLabel,
				"\n", 
				descriptionText,
			)
			// TODO
			// Update not woking because the activites endpoint does not return description
			// The activty being updated needs to be requested to get the existing description
			// Impact: This will override existing activty descriptions
			if strings.Contains(newest_activity["description"].(string), descriptionLabel) {
				updateActivityDescription(newest_activity["id"].(float64), description, stravaAccessToken)
			} else {
				log.Fatalln("Activity already has an elevation challenge description:", newest_activity["id"])
			}
		} else {
			description = fmt.Sprintf(
				"%s %s %s", 
				descriptionLabel,
				"\n",
				descriptionText,
			)
			updateActivityDescription(newest_activity["id"].(float64), description, stravaAccessToken)
		}
		fmt.Printf("Description: %s\n", description);
	} else {
		log.Fatalln("Unable to find the newest activity")
	}
}

func updateActivityDescription(activityId float64, description string, stravaAccessToken string) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	requestBody, reqBodyErr := json.Marshal(map[string]string{
		"description": description,
	})
	if reqBodyErr != nil {
		log.Fatalln(reqBodyErr)
	}
	fmt.Printf("requestBody: %s\n", requestBody)

	request, reqErr := http.NewRequest(
		"PUT",
		fmt.Sprintf(
			"https://www.strava.com/api/v3/activities/%.0f",
			activityId,
		),
		bytes.NewReader(requestBody),
	)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", stravaAccessToken))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	if reqErr != nil {
		log.Fatalln(reqErr)
	}

	response, respErr := client.Do(request)
	if respErr != nil {
		log.Fatalln(respErr)
	}

	defer response.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	log.Println("updatedDescription:\n", result["description"])
}

func getStravaAccessToken() string {
	if os.Getenv("client_id") == "" ||
		os.Getenv("client_secret") == "" ||
		os.Getenv("refresh_token") == "" {
		log.Fatalln("Secrets required. Provide client_id, client_secret, refresh_token")
	}

	formData := url.Values{
		"client_id": {os.Getenv("client_id")},
		"client_secret": {os.Getenv("client_secret")},
		"grant_type": {"refresh_token"},
		"refresh_token": {os.Getenv("refresh_token")},
	}

	resp, respErr := http.PostForm(
		"https://www.strava.com/api/v3/oauth/token", 
		formData,
	)
	if respErr != nil {
		log.Fatalln(respErr)
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["access_token"].(string)
}

// TODO
// Move the GET request to activites into a function
// Impact: Makes code easier to read
// func getStravaActivites() {}