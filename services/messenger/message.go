package messenger

import (
	"bytes"
	"encoding/json"
	"github.com/s4kibs4mi/jally-commerce-bot/log"
	"io/ioutil"
	"net/http"
)

type MessageResponse struct {
	RecipientID  string `json:"recipient_id"`
	MessageID    string `json:"message_id"`
	AttachmentID string `json:"attachment_id,omitempty"`
}

type rawMessage struct {
	Recipient
	MessageQuery
}

func (m *Messenger) sendCustomMessage(i interface{}) ([]byte, error) {
	byt, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	log.Log().Warningln(string(byt))

	resp, err := m.doRequest("POST", GraphAPI+"/v15.0/me/messages", bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	log.Log().Warningln(string(read))

	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, er.Error
	}
	return read, err
}

func (m *Messenger) SendMessage(mq MessageQuery) (*MessageResponse, error) {
	b, err := m.sendCustomMessage(mq)
	if err != nil {
		return nil, err
	}
	response := &MessageResponse{}
	err = json.Unmarshal(b, response)
	return response, err
}

func (m *Messenger) SendSimpleMessage(recipient string, message string) (*MessageResponse, error) {
	return m.SendMessage(MessageQuery{
		Recipient: Recipient{
			ID: recipient,
		},
		Message: SendMessage{
			Text: message,
		},
		MessagingType: MessagingTypeRegular,
	})
}
