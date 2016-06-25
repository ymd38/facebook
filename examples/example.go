package main

import (
	"log"
	"net/http"
	"os"

	. "github.com/ymd38/facebook/messenger"

	"github.com/gin-gonic/gin"
)

var fb_token string

func main() {

	fb_token = os.Getenv("FBTOKEN")

	router := gin.Default()

	router.GET("/webhook", GetWebHook)
	router.POST("/webhook", PostWebHook)

	router.Run(":9000")
}

// Setup Webhook
// see #2: https://developers.facebook.com/docs/messenger-platform/quickstart
func GetWebHook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	if mode == "subscribe" && token == fb_token {
		challenge := c.Query("hub.challenge")
		c.JSON(http.StatusOK, challenge)
	} else {
		c.String(http.StatusNotFound, "token error")
	}
}

// Recive message and send message
func PostWebHook(c *gin.Context) {
	receiver := &ReceivedMessage{}

	if err := c.BindJSON(&receiver); err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	for _, messaging := range receiver.Entry[0].Messaging {
		fb := NewFacebookMessenger(fb_token)
		m := NewTextMessage(messaging.Sender.ID, "Hello")
		//m = NewImageMessage(messaging.Sender.ID, "http://1093.up.n.seesaa.net/1093/image/takokora.jpg")
		if err := fb.SendMessage(m); err != nil {
			log.Println(err.Error())
		}
	}

	c.String(http.StatusOK, "OK")
}
