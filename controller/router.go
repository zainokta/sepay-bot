package router

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func errCheck(err error) {
	if err != nil {
		log.Print(err.Error())
	}
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func Callback(c *gin.Context) {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(c.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.AbortWithError(http.StatusBadGateway, err)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				reply := message.Text
				_, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage(strings.ToUpper(reply))).Do()
				errCheck(err)
			}
		}
	}
}
