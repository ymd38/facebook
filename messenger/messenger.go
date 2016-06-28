package messenger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const endpoint = "https://graph.facebook.com/v2.6/me/messages?access_token="

// FacebookMessenger ...
type FacebookMessenger struct {
	Token string
}

// ReceivedMessage ...
type ReceivedMessage struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

// Entry ...
type Entry struct {
	ID        string      `json:"id"`
	Time      int         `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

// Messaging ...
type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int       `json:"timestamp"`
	Message   Message   `json:"message"`
}

// Sender ...
type Sender struct {
	ID string `json:"id"`
}

// Recipient ...
type Recipient struct {
	ID string `json:"id"`
}

// Message ...
type Message struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

// TextMessage ...
type TextMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Text string `json:"text"`
	} `json:"message"`
}

// ImageMessage ...
type ImageMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Attachment struct {
			Type    string `json:"type"`
			Payload struct {
				Url string `json:"url"`
			} `json:"payload"`
		} `json:"attachment"`
	} `json:"message"`
}

// GenericTemplate ...
type GenericTemplate struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Attachment struct {
			Type     string `json:"type"`
			Elements struct {
				Title    string `json:"title"`
				ItemUrl  string `json:"item_url,omitempty"`
				ImageUrl string `json:"image_url,omitempty"`
			} `json:"elements"`
		} `json:"attachment"`
	} `json:"message"`
}

// NewFacebookMessenger ...
func NewFacebookMessenger(token string) *FacebookMessenger {
	return &FacebookMessenger{Token: token}
}

// NewTextMessage ...
func NewTextMessage(senderid string, text string) *TextMessage {
	t := &TextMessage{}
	t.Recipient.ID = senderid
	t.Message.Text = text

	log.Println(t)

	return t
}

// NewImageMessage ...
func NewImageMessage(senderid string, image_url string) *ImageMessage {
	i := &ImageMessage{}
	i.Recipient.ID = senderid
	i.Message.Attachment.Type = "image"
	i.Message.Attachment.Payload.Url = image_url

	log.Println(i)

	return i
}

func (fb *FacebookMessenger) SendMessage(m interface{}) error {
	log.Println(m)
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint+fb.Token, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var result map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	log.Println(result)

	return nil
}
