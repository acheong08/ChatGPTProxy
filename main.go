package main

import (
	"io"
	"log"
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
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _  = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	user_agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	http_proxy = os.Getenv("http_proxy")
)

func main() {

	if http_proxy != "" {
		client.SetProxy(http_proxy)
		println("Proxy set:" + http_proxy)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "9090"
	}
	handler := gin.Default()
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	handler.Any("/api/*path", proxy)

	gin.SetMode(gin.ReleaseMode)
	endless.ListenAndServe(os.Getenv("HOST")+":"+PORT, handler)
}

func proxy(c *gin.Context) {
	// Remove _cfuvid cookie from session
	jar.SetCookies(c.Request.URL, []*http.Cookie{})

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
	request.Header.Set("Authorization", c.Request.Header.Get("Authorization"))
	request.Header.Set("sec-ch-ua", "\"Chromium\";v=\"112\", \"Brave\";v=\"112\", \"Not:A-Brand\";v=\"99\"")
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("sec-gpc", "1")
	request.Header.Set("user-agent", user_agent)
	if os.Getenv("PUID") != "" {
		request.AddCookie(&http.Cookie{Name: "_puid", Value: os.Getenv("PUID")})
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

	buf := make([]byte, 4096)
	for {
		n, err := response.Body.Read(buf)
		if n > 0 {
			_, writeErr := c.Writer.Write(buf[:n])
			if writeErr != nil {
				log.Printf("Error writing to client: %v", writeErr)
				break
			}
			c.Writer.Flush() // flush buffer to make sure the data is sent to client in time.
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading from response body: %v", err)
			break
		}
	}
}
