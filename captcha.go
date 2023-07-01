package main

import (
	"encoding/json"
	"net/http"

	"github.com/acheong08/funcaptcha"
	"github.com/gin-gonic/gin"
)

func captchaStart(c *gin.Context) {
	token, hex, err := funcaptcha.GetOpenAIToken()
	if err == nil {
		c.JSON(200, gin.H{"token": token, "status": "success"})
		return
	}
	if err.Error() != "captcha required" {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session, err := funcaptcha.StartChallenge(token, hex)
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to log requests"})
		return
	}
	err = session.RequestChallenge(false)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to request challenge"})
		return
	}
	// Get form data (check if download_images is true)
	download_images := c.Query("download_images")
	var images []string
	if download_images == "true" {
		// Get Base64 encoded image
		images, err = funcaptcha.DownloadChallenge(session.ConciseChallenge.URLs, true)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to download images"})
			return
		}
	}
	c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"token": token, "session": session, "status": "captcha", "images": images})
}

func captchaVerify(c *gin.Context) {
	type submissionRequest struct {
		Index   int                `json:"index"`
		Session funcaptcha.Session `json:"session"`
	}
	var request submissionRequest
	// Map the request body to the submissionRequest struct
	if c.Request.Body != nil {
		err := json.NewDecoder(c.Request.Body).Decode(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(400, gin.H{"error": "request body not provided"})
		return
	}
	// Verify the captcha
	err := request.Session.SubmitAnswer(request.Index)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Success
	c.JSON(200, gin.H{"status": "success"})
}
