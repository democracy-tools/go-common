package whatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/democracy-tools/go-common/env"
	"github.com/sirupsen/logrus"
)

type Client interface {
	SendReminderTemplate(template string, phone string, userId string) error
}

type ClientWrapper struct {
	auth string
	from string
}

func NewClientWrapper() Client {
	return &ClientWrapper{
		auth: fmt.Sprintf("Bearer %s", env.GetWhatAppToken()),
		from: env.GetWhatsAppFromPhone(),
	}
}

func (c *ClientWrapper) SendReminderTemplate(template string, to string, userId string) error {

	return send(c.from, to, c.auth,
		newTemplate(template, to, fmt.Sprintf("?user-id=%s", userId), nil))
}

func send(from string, to string, auth string, message TemplateMessageRequest) error {

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(message)
	if err != nil {
		err = fmt.Errorf("failed to encode whatsapp template message '%s' with '%v' phone '%s'",
			message.Template.Name, err, to)
		logrus.Error(err.Error())
		return err
	}

	r, err := http.NewRequest(http.MethodPost, getMessageUrl(from), &body)
	if err != nil {
		logrus.Errorf("failed to create HTTP request for sending a whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", auth)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		logrus.Errorf("failed to send whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("failed to send whatsapp message to '%s' with '%s'", to, response.Status)
		logrus.Info(msg)
		return errors.New(msg)
	}

	return nil
}

func newTemplate(name string, to string, buttonUrlParam string, bodyTextParams []string) TemplateMessageRequest {

	res := TemplateMessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: Template{
			Name: name,
			Language: Language{
				Policy: "deterministic",
				Code:   "he",
			},
		},
	}

	if buttonUrlParam != "" {
		res.Template.Components = []Component{{
			Type:    "button",
			SubType: "url",
			Index:   "0",
			Parameters: []Parameter{{
				Type: "text",
				Text: buttonUrlParam,
			}},
		}}
	}

	if len(bodyTextParams) > 0 {
		var params []Parameter
		for _, currParam := range bodyTextParams {
			params = append(params, Parameter{
				Type: "text",
				Text: currParam,
			})
		}
		res.Template.Components = append(res.Template.Components, Component{
			Type:       "body",
			Parameters: params})
	}

	return res
}

func getMessageUrl(from string) string {

	return fmt.Sprintf("https://graph.facebook.com/v16.0/%s/messages", from)
}
