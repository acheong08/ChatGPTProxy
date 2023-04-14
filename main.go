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
	http_proxy   = os.Getenv("http_proxy")
	auth_proxy   = os.Getenv("auth_proxy")
	openai_email = os.Getenv("OPENAI_EMAIL")
	openai_pass  = os.Getenv("OPENAI_PASS")
	admin_pass   = os.Getenv("ADMIN_PASS")
)

func main() {
	if access_token == "" && puid == "" && openai_email == "" && openai_pass == "" {
		println("Error: Authentication information not found.")
		return
	}

	if http_proxy != "" {
		client.SetProxy(http_proxy)
		println("Proxy set:" + http_proxy)
	}
	// Automatically refresh the puid cookie
	if openai_email != "" && openai_pass != "" {
		go refresh_access_token()
	} else if access_token != "" {
		go func() {
			for {
				refresh_puid()
				time.Sleep(6 * time.Hour)
			}
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

	handler.POST("/admin/update", func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != admin_pass {
			c.JSON(401, gin.H{"message": "unauthorized"})
			return
		}
		type Update struct {
			Value string `json:"value"`
			Field string `json:"field"`
		}
		var update Update
		c.BindJSON(&update)
		if update.Field == "puid" {
			puid = update.Value
			// export environment variable
			os.Setenv("PUID", puid)
		} else if update.Field == "access_token" {
			access_token = update.Value
			os.Setenv("ACCESS_TOKEN", access_token)
		} else if update.Field == "http_proxy" {
			http_proxy = update.Value
			client.SetProxy(http_proxy)
		} else if update.Field == "openai_email" {
			openai_email = update.Value
			os.Setenv("OPENAI_EMAIL", openai_email)
		} else if update.Field == "openai_pass" {
			openai_pass = update.Value
			os.Setenv("OPENAI_PASS", openai_pass)
		} else if update.Field == "admin_pass" {
			admin_pass = update.Value
			os.Setenv("ADMIN_PASS", admin_pass)
		} else if update.Field == "auth_proxy" {
			auth_proxy = update.Value
			os.Setenv("auth_proxy", auth_proxy)
		} else {
			c.JSON(400, gin.H{"message": "field not found"})
			return
		}
		c.JSON(200, gin.H{"message": "updated"})
	})
	gin.SetMode(gin.ReleaseMode)
	endless.ListenAndServe(os.Getenv("HOST")+":"+PORT, handler)
}

func proxy(c *gin.Context) {

	var url string
	var err error
	var request_method string
	var request *http.Request
	var response *http.Response

	if c.Request.URL.RawQuery != "" {
		url = "https://chat.openai.com/backend-api" + c.Param("path") + "?" + c.Request.URL.RawQuery
	} else {
		url = "https://chat.openai.com/backend-api" + c.Param("path")
	}
	request_method = c.Request.Method

	request, err = http.NewRequest(request_method, url, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	request.Header.Set("Host", "chat.openai.com")
	request.Header.Set("Origin", "https://chat.openai.com/chat")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json")
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
