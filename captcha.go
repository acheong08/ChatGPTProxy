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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = session.RequestChallenge(false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Get session as JSON
	session_json, _ := json.Marshal(session)
	c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"token": token, "session": string(session_json), "status": "captcha"})
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
