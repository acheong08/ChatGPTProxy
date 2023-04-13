package main

import (
	"io"
	"time"

	"github.com/acheong08/OpenAIAuth/auth"
	http "github.com/bogdanfinn/fhttp"
)

func refresh_puid() {
	url := "https://chat.openai.com/backend-api/models"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Host", "chat.openai.com")
	req.Header.Set("origin", "https://chat.openai.com/chat")
	req.Header.Set("referer", "https://chat.openai.com/chat")
	req.Header.Set("sec-ch-ua", `Chromium";v="110", "Not A(Brand";v="24", "Brave";v="110`)
	req.Header.Set("sec-ch-ua-platform", "Linux")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "text/event-stream")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+access_token)
	// Initial puid cookie
	req.AddCookie(
		&http.Cookie{
			Name:  "_puid",
			Value: puid,
		},
	)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	println("Got response: " + resp.Status)
	if resp.StatusCode != 200 {
		println("Error: " + resp.Status)
		// Print response body
		body, _ := io.ReadAll(resp.Body)
		println(string(body))
		return
	}
	// Get cookies from response
	cookies := resp.Cookies()
	// Find _puid cookie
	for _, cookie := range cookies {
		if cookie.Name == "_puid" {
			puid = cookie.Value
			println("puid: " + puid)
			break
		}
	}
}

func refresh_access_token() {
	authenticator := auth.NewAuthenticator(openai_email, openai_pass, puid, auth_proxy)
	for {
		err := authenticator.Begin()
		if err.Error != nil {
			println("Error: " + err.Details)
			// Sleep for 30 minutes
			time.Sleep(30 * time.Minute)
			continue
		}
		access_token, err = authenticator.GetAccessToken()
		if err.Error != nil {
			println("Error: " + err.Details)
			// Sleep for 20 minutes
			time.Sleep(20 * time.Minute)
			continue
		}
		println("Got access token: " + access_token)
		refresh_puid()
		// Sleep for 12 hour
		time.Sleep(24 * time.Hour)
	}
}
