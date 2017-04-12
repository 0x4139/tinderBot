package main

import (
	"net/http"
	"io"
	"compress/gzip"
	"io/ioutil"
	"log"
	"time"
	"encoding/json"
	"fmt"
	"strings"
	"os"
)

const hunterToken = `Token token="c462d46b-ca1e-4368-8ae9-f33117bbaabd"`

const requestURL = "https://api.gotinder.com/recs/core?locale=en-RO"

type GetProfileResponse struct {
	Status int `json:"status"`
	Results []struct {
		Type         string `json:"type"`
		GroupMatched bool `json:"group_matched"`
		User struct {
			DistanceMi int `json:"distance_mi"`
			CommonConnections []struct {
				Degree int `json:"degree"`
				Photo struct {
					Large  string `json:"large"`
					Medium string `json:"medium"`
					Small  string `json:"small"`
				} `json:"photo"`
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"common_connections"`
			ConnectionCount int `json:"connection_count"`
			CommonLikes     []interface{} `json:"common_likes"`
			CommonInterests []interface{} `json:"common_interests"`
			CommonFriends   []interface{} `json:"common_friends"`
			ContentHash     string `json:"content_hash"`
			ID              string `json:"_id"`
			Bio             string `json:"bio"`
			BirthDate       time.Time `json:"birth_date"`
			Name            string `json:"name"`
			PingTime        time.Time `json:"ping_time"`
			Photos []struct {
				ID  string `json:"id"`
				URL string `json:"url"`
				ProcessedFiles []struct {
					URL    string `json:"url"`
					Height int `json:"height"`
					Width  int `json:"width"`
				} `json:"processedFiles"`
			} `json:"photos"`
			Jobs []struct {
				Company struct {
					Name string `json:"name"`
				} `json:"company"`
			} `json:"jobs"`
			Schools []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"schools"`
			Teaser struct {
				Type   string `json:"type"`
				String string `json:"string"`
			} `json:"teaser"`
			Teasers []struct {
				Type   string `json:"type"`
				String string `json:"string"`
			} `json:"teasers"`
			SNumber       int `json:"s_number"`
			Gender        int `json:"gender"`
			BirthDateInfo string `json:"birth_date_info"`
			GroupMatched  bool `json:"group_matched"`
		} `json:"user"`
	} `json:"results"`
}

func generateTinderRequest(url string) *http.Request {
	getProfilesRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	getProfilesRequest.Header.Set("Host", "api.gotinder.com")
	getProfilesRequest.Header.Set("Authorization", fmt.Sprintf(`Token token="%s"`,os.Getenv("TOKEN")))
	getProfilesRequest.Header.Set("x-client-version", "71104")
	getProfilesRequest.Header.Set("app-version", "1885")
	getProfilesRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	getProfilesRequest.Header.Set("If-None-Match", `W/"582464522"`)
	getProfilesRequest.Header.Set("platform", "ios")
	getProfilesRequest.Header.Set("Accept-Language", "en-RO;q=1, ro-RO;q=0.9")
	getProfilesRequest.Header.Set("Accept", "*/*")
	getProfilesRequest.Header.Set("User-Agent", "Tinder/7.1.1 (iPhone; iOS 10.2.1; Scale/2.00)")
	getProfilesRequest.Header.Set("Connection", "keep-alive")
	getProfilesRequest.Header.Set("X-Auth-Token", os.Getenv("TOKEN"))
	getProfilesRequest.Header.Set("os_version", "c100000200001")
	return getProfilesRequest
}

func performRequest(request *http.Request) []byte {
	client := new(http.Client)
	response, err := client.Do(request)
	defer response.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err.Error())
	}

	return body

}
func main() {
	if os.Getenv("TOKEN")==""{
		panic("you need a token!")
	}
	leLoop()

}

func leLoop() {
begin:
	getProfilesRequest := generateTinderRequest(requestURL)
	profileRequestResponse := performRequest(getProfilesRequest)

	var profileResponse *GetProfileResponse
	if strings.Contains(string(profileRequestResponse), "recs limited") {
		log.Printf("you've been limited we got to wait %s", string(profileRequestResponse))
		time.Sleep(time.Minute * 15)
		goto begin
	}

	if err := json.Unmarshal(profileRequestResponse, &profileResponse); err != nil {
		log.Printf("you've been limited %s", string(profileRequestResponse))
		panic("done")
	}

	for _, result := range profileResponse.Results {
		matchRequest := generateTinderRequest(fmt.Sprintf("https://api.gotinder.com/like/%s?content_hash=%s&s_number=%s", result.User.ID, result.User.ContentHash, result.User.SNumber))
		matchRequestResponse := performRequest(matchRequest)
		if strings.Contains(string(matchRequestResponse), `"match":false`) {
			log.Printf("No Luck with %s", result.User.Name)
		} else if strings.Contains(string(matchRequestResponse), `"match":true`) {
			log.Printf("Lucky STRIKE!!!!! %s", result.User.Name)
		} else {
			log.Printf("tinder might banned us :( %s", matchRequestResponse)
		}

		time.Sleep(time.Second * 2)
	}
	leLoop()
}
