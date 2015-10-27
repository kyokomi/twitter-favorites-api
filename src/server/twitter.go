package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	autoTokenURL      = "https://api.twitter.com/oauth2/token"
	favoritesImageURL = "https://api.twitter.com/1.1/favorites/list.json"
)

type TwitterAPI interface {
	GetFavoritesImages(count int, screenName string) ([]string, error)
}

type twitterAPI struct {
	client      *http.Client
	accessToken string
}

var _ TwitterAPI = (*twitterAPI)(nil)

func NewTwitterAPI(client *http.Client, apiKey, apiSecret string) (TwitterAPI, error) {
	t := &twitterAPI{
		client: client,
	}
	token, err := t.newAccessToken(apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("generate token error %s", err)
	}
	t.accessToken = token

	return t, nil
}

func (t *twitterAPI) newAccessToken(key, secret string) (string, error) {
	if key == "" {
		key = os.Getenv("TWITTER_API_KEY")
	}
	if secret == "" {
		secret = os.Getenv("TWITTER_API_SECRET")
	}
	token := base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(key) + ":" + url.QueryEscape(secret)))

	data := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequest("POST", autoTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+token)

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})

	return m["access_token"].(string), nil
}

type TwitterFavoritesResponse struct {
	Entities struct {
		Media []struct {
			MediaURL      string `json:"media_url"`
			MediaURLHTTPs string `json:"media_url_https"`
		} `json:"media"`
	} `json:"extended_entities"`
}

func (r TwitterFavoritesResponse) GetMediaURLs() []string {
	urls := map[string]bool{}
	for _, m := range r.Entities.Media {
		if urls[m.MediaURL] {
			continue
		}
		urls[m.MediaURL] = true
	}

	results := make([]string, len(urls))
	i := 0
	for url, _ := range urls {
		results[i] = url
		i++
	}
	return results
}

func (t *twitterAPI) GetFavoritesImages(count int, screenName string) ([]string, error) {
	req, err := t.NewRequest("GET", fmt.Sprintf("%s?count=%d&screen_name=%s", favoritesImageURL, count, screenName))
	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	searchBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response []TwitterFavoritesResponse
	if err := json.Unmarshal(searchBody, &response); err != nil {
		return nil, err
	}

	mediaURLs := []string{}
	for _, re := range response {
		mediaURLs = append(mediaURLs, re.GetMediaURLs()...)
	}
	return mediaURLs, nil
}

func (t *twitterAPI) NewRequest(method string, apiURL string) (*http.Request, error) {
	req, err := http.NewRequest(method, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+t.accessToken)
	return req, err
}
