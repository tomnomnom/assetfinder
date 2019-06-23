package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func fetchFacebook(domain string) ([]string, error) {

	appId := os.Getenv("FB_APP_ID")
	appSecret := os.Getenv("FB_APP_SECRET")
	if appId == "" || appSecret == "" {
		// fail silently because it's reasonable not to have
		// the Facebook API creds
		return []string{}, nil
	}

	accessToken, err := facebookAuth(appId, appSecret)
	if err != nil {
		return []string{}, err
	}

	domains, err := getFacebookCerts(accessToken, domain)
	if err != nil {
		return []string{}, err
	}

	return domains, nil
}

func getFacebookCerts(accessToken, query string) ([]string, error) {
	out := make([]string, 0)
	fetchURL := fmt.Sprintf(
		"https://graph.facebook.com/certificates?fields=domains&access_token=%s&query=*.%s",
		accessToken, query,
	)

	for {

		wrapper := struct {
			Data []struct {
				Domains []string `json:"domains"`
			} `json:"data"`

			Paging struct {
				Next string `json:"next"`
			} `json:"paging"`
		}{}

		err := fetchJSON(fetchURL, &wrapper)
		if err != nil {
			return out, err
		}

		for _, data := range wrapper.Data {
			for _, d := range data.Domains {
				out = append(out, d)
			}
		}

		fetchURL = wrapper.Paging.Next
		if fetchURL == "" {
			break
		}
	}
	return out, nil
}

func facebookAuth(appId, appSecret string) (string, error) {
	authUrl := fmt.Sprintf(
		"https://graph.facebook.com/oauth/access_token?client_id=%s&client_secret=%s&grant_type=client_credentials",
		appId, appSecret,
	)

	resp, err := http.Get(authUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	auth := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = dec.Decode(&auth)
	if err != nil {
		return "", err
	}

	if auth.AccessToken == "" {
		return "", errors.New("no access token in Facebook API response")
	}

	return auth.AccessToken, nil
}
