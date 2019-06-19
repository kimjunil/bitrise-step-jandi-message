package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	Debug bool `env:"is_debug_mode,opt[yes,no]"`

	WebhookURL stepconf.Secret `env:"jandi_webhook_url"`

	// Message
	MessageBody        string `env:"message_body"`
	MessageBodyOnError string `env:"message_body_on_error"`
	Color              string `env:"color,required"`
	ColorOnError       string `env:"color_on_error"`

	// Attachment
	Title              string `env:"attachment_area_title"`
	TitleOnError       string `env:"attachment_area_title_on_error"`
	Descriotion        string `env:"attachment_area_description"`
	DescriotionOnError string `env:"attachment_area_description_on_error"`
	ImageURL           string `env:"attachment_area_image_url"`
	ImageURLOnError    string `env:"attachment_area_image_url_on_error"`
}

// success is true if the build is successful, false otherwise.
var success = os.Getenv("BITRISE_BUILD_STATUS") == "0"

// selectValue chooses the right value based on the result of the build.
func selectValue(ifSuccess, ifFailed string) string {
	if success || ifFailed == "" {
		return ifSuccess
	}
	return ifFailed
}

// ensureNewlines replaces all \n substrings with newline characters.
func ensureNewlines(s string) string {
	return strings.Replace(s, "\\n", "\n", -1)
}

func newMessage(c Config) Message {
	msg := Message{
		MessageBody: selectValue(c.MessageBody, c.MessageBodyOnError),
		Color:       selectValue(c.Color, c.ColorOnError),
		Attachments: []Attachment{{
			Title:       selectValue(c.Title, c.TitleOnError),
			Description: selectValue(c.Descriotion, c.DescriotionOnError),
			ImageURL:    selectValue(c.ImageURL, c.ImageURLOnError),
		}},
	}
	return msg
}

func postMessage(conf Config, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Debugf("Request to Jandi: %s\n", b)

	url := string(conf.WebhookURL)

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Add("Accept", headerAccept)
	req.Header.Add("Content-Type", headerContentType)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("server error: %s, failed to read response: %s", resp.Status, err)
		}
		return fmt.Errorf("server error: %s, response: %s", resp.Status, body)
	}

	return nil
}

func main() {
	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(conf)
	log.SetEnableDebugLog(conf.Debug)

	msg := newMessage(conf)
	if err := postMessage(conf, msg); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nJandi message successfully sent! 🚀\n")
}
