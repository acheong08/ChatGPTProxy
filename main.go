package main

import (
	"io"
	"os"

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
	client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
)

func main() {
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

	url = "https://api.openai.com" + c.Param("path")
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
