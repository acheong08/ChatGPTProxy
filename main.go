package main

import (
	"io"
	"net/http"
	"os"
	"time"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

var (
	jar          = tls_client.NewCookieJar()
	options      = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Chrome_110),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}
	client, _    = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	access_token = os.Getenv("ACCESS_TOKEN")
	puid         = os.Getenv("PUID")
	http_proxy   = os.Getenv("http_proxy")
)

func main() {
	if access_token == "" && puid == "" {
		println("Error: ACCESS_TOKEN and PUID are not set")
		return
	}

	if http_proxy != "" {
		client.SetProxy(http_proxy)
		println("Proxy set:" + http_proxy)
	}

	// Automatically refresh the puid cookie
	if access_token != "" {
		go refreshPuidCookie()
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

func refreshPuidCookie() {
	url := "https://chat.openai.com/backend-api/models"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		println("Error creating request:", err.Error())
		return
	}

	setRequestHeaders(req)

	for {
		resp, err := client.Do(req)
		if err != nil {
			println("Error:", err.Error())
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			println("Error: " + resp.Status)
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
		// Sleep for 6 hours
		time.Sleep(6 * time.Hour)
	}
	println("Error: Failed to refresh puid cookie")
}

func proxy(c *gin.Context) {
	url := "https://chat.openai.com/backend-api" + c.Param("path")
	request_method := c.Request.Method

	request, err := http.NewRequest(request_method, url, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	setRequestHeaders(request)
	request.Header.Set("Authorization", c.Request.Header.Get("Authorization"))

	puidValue := c.Request.Header.Get("Puid")
	if puidValue == "" {
		puidValue = puid
	}

	request.AddCookie(&http.Cookie{
		Name:  "_puid",
		Value: puidValue,
	})

	response, err := client.Do(request)
	if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer response.Body.Close()

    // Check if the puid cookie is invalid
    if response.StatusCode == 401 {
        refreshPuidCookie()
        response, err = client.Do(request)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        defer response.Body.Close()
    }

    c.Header("Content-Type", response.Header.Get("Content-Type"))
    c.Status(response.StatusCode)
    c.Stream(func(w io.Writer) bool {
        io.Copy(w, response.Body)
        return false
    })
}

func setRequestHeaders(req *http.Request) {
    req.Header.Set("Host", "chat.openai.com")
    req.Header.Set("origin", "https://chat.openai.com/chat")
    req.Header.Set("referer", "https://chat.openai.com/chat")
    req.Header.Set("sec-ch-ua", Chromium";v="110", "Not A(Brand";v="24", "Brave";v="110)
    req.Header.Set("sec-ch-ua-platform", "Linux")
    req.Header.Set("content-type", "application/json")
    req.Header.Set("accept", "text/event-stream")
    req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
    // Set authorization header
    req.Header.Set("Authorization", "Bearer "+access_token)
}
