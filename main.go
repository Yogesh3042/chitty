package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron"
)

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID   int64  `json:"id"`
			NAME string `json:"first_name"`
		} `json:"chat"`
	} `json:"message"`
}

func Handler(res http.ResponseWriter, req *http.Request) {

	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	if strings.Contains(strings.ToLower(body.Message.Text), "hi") {
		if err := reply(body.Message.Chat.ID, body.Message.Chat.NAME, "Welcome to Solar"); err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}
	}
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		if err := reply(body.Message.Chat.ID, body.Message.Chat.NAME, "Welcome to Solar"); err != nil {
			fmt.Println("error in sending reply:", err)
			return
		}
	})
	c.Start()

	fmt.Println("reply sent")
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func main() {
	port := ":" + os.Getenv("PORT")
	fmt.Println(port)
	http.ListenAndServe(port, http.HandlerFunc(Handler))
}

func reply(chatID int64, name string, msg string) error {
	t := time.Now().Hour()
	var txt string
	if t < 12 && t > 5 {
		txt = "Good Morning"
	} else if t > 12 && t < 17 {
		txt = "Good Afternoon"
	} else if t > 17 && t < 22 {
		txt = "Good Evening"
	} else {
		txt = "Its Night "
	}
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   txt + " " + name + " " + msg,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post("https://api.telegram.org/bot5332294644:AAFSugXpoJXaZNHCZtS6JaXOEfkJiQx1tQE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
