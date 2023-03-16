package main

import (
	"io"
	"os"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Chrome_110),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _    = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	access_token = os.Getenv("ACCESS_TOKEN")
	puid         = os.Getenv("PUID")
)

func main() {
	if access_token == "" && puid == "" {
		println("Error: ACCESS_TOKEN and PUID are not set")
		return
	}
	// Automatically refresh the puid cookie
	if access_token != "" {
		go func() {
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
			for {
				resp, err := client.Do(req)
				if err != nil {
					break
				}
				defer resp.Body.Close()
				println("Got response: " + resp.Status)
				if resp.StatusCode != 200 {
					println("Error: " + resp.Status)
					// Print response body
					body, _ := io.ReadAll(resp.Body)
					println(string(body))
					break
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
				// Sleep for 6 hour
				time.Sleep(6 * time.Hour)
			}
			println("Error: Failed to refresh puid cookie")
		}()
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	handler := gin.Default()
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	handler.Any("/api/*path", proxy)

	endless.ListenAndServe(os.Getenv("HOST")+":"+PORT, handler)
}

func proxy(c *gin.Context) {

	var url string
	var err error
	var request_method string
	var request *http.Request
	var response *http.Response

	url = "https://chat.openai.com/backend-api" + c.Param("path")
	request_method = c.Request.Method

	request, err = http.NewRequest(request_method, url, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	request.Header.Set("Host", "chat.openai.com")
	request.Header.Set("Origin", "https://chat.openai.com/chat")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Keep-Alive", "timeout=360")
	request.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	request.Header.Set("Authorization", c.Request.Header.Get("Authorization"))
	if c.Request.Header.Get("Puid") == "" {
		request.AddCookie(
			&http.Cookie{
				Name:  "_puid",
				Value: puid,
			},
		)
	} else {
		request.AddCookie(
			&http.Cookie{
				Name:  "_puid",
				Value: c.Request.Header.Get("Puid"),
			},
		)
	}

	response, err = client.Do(request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()
	c.Header("Content-Type", response.Header.Get("Content-Type"))
	// Get status code
	c.Status(response.StatusCode)
	c.Stream(func(w io.Writer) bool {
		// Write data to client
		io.Copy(w, response.Body)
		return false
	})

}
