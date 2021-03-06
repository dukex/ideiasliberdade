package buffer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	URL = "https://api.bufferapp.com/1"
)

type Client struct {
	AccessToken string
	Url         string
}

type Update struct {
	Id        string
	Text      string
	ProfileId string
}

type Profile struct {
	Avatar            string
	CreatedAt         int64
	Default           bool
	FormattedUsername string
	Id                string
	Schedules         []map[string][]string
	Service           string
	ServiceId         string
	ServiceUsername   string
	Statistics        map[string]interface{}
	TeamMembers       []string
	Timezone          string
	UserId            string
}

type Profiles []Profile

func (c *Client) Profiles() Profiles {
	bufferResponse := c.sendGET("profiles")
	fmt.Println(string(bufferResponse[:]))
	var response Profiles
	err := json.Unmarshal(bufferResponse, &response)

	if err != nil {
		panic(err)
	}
	return response
}

func (c *Client) CreateUpdate(text string, profileIds []string, options map[string]interface{}) []Update {

	params := url.Values{}
	params.Set("text", text)
	for _, p := range profileIds {
		params.Add("profile_ids[]", p)
	}

	bufferResponse := c.sendPOST("updates/create", params)

	var response struct {
		Success          bool
		BufferCount      int
		BufferPercentage int
		Updates          []Update
	}

	err := json.Unmarshal(bufferResponse, &response)

	if err != nil {
		panic(err)
	}

	return response.Updates
}

func (c *Client) sendGET(resource string) []byte {
	urlEndpoint := c.Url + "/" + resource + ".json?access_token=" + c.AccessToken
	request, err := http.Get(urlEndpoint)
	if err != nil {
		panic(err)
	}

	defer request.Body.Close()
	requestBodyByte, _ := ioutil.ReadAll(request.Body)

	return requestBodyByte
}

func (c *Client) sendPOST(resource string, params url.Values) []byte {
	urlEndpoint := c.Url + "/" + resource + ".json?access_token=" + c.AccessToken
	request, err := http.PostForm(urlEndpoint, params)
	if err != nil {
		panic(err)
	}

	defer request.Body.Close()
	requestBodyByte, _ := ioutil.ReadAll(request.Body)

	return requestBodyByte
}

func NewClient(accessToken string) *Client {
	return &Client{Url: URL, AccessToken: accessToken}
}
