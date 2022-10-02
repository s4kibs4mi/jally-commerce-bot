package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/s4kibs4mi/jally-commerce-bot/log"
	"github.com/s4kibs4mi/jally-commerce-bot/services/messenger/template"
	"io/ioutil"
	"net/http"
)

type settingType string

const (
	settingTypeGreeting      settingType = "greeting"
	settingTypeCallToActions settingType = "call_to_actions"
)

type threadState string

const (
	threadStateNew      threadState = "new_thread"
	threadStateExisting threadState = "existing_thread"
)

type greeting struct {
	Text string `json:"text,omitempty"`
}

type result struct {
	Result string `json:"result"`
}

type persistentMenuSettings struct {
	Locale                string      `json:"locale"`
	ComposerInputDisabled bool        `json:"composer_input_disabled"`
	CallToActions         interface{} `json:"call_to_actions,omitempty"`
}

type persistentMenu struct {
	PersistentMenu []persistentMenuSettings `json:"persistent_menu"`
}

func (m *Messenger) changePersistentMenuSettings(httpMethod string, set *persistentMenu) (*result, error) {
	body, err := json.Marshal(set)
	if err != nil {
		return nil, err
	}

	log.Log().Infoln(string(body))

	resp, err := m.doRequest(httpMethod, GraphAPI+"/v15.0/me/messenger_profile", bytes.NewReader(body))
	if err != nil {
		log.Log().Errorln(err)
		return nil, err
	}
	defer resp.Body.Close()
	pld, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Log().Infoln(string(pld))

	result := &result{}
	err = json.Unmarshal(pld, result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}
	return result, nil
}

type getStartedPayload struct {
	Payload string `json:"payload"`
}

type getStarted struct {
	GetStarted getStartedPayload `json:"get_started"`
}

func (m *Messenger) changeGetStartedSettings(httpMethod string, set *getStarted) (*result, error) {
	body, err := json.Marshal(set)
	if err != nil {
		return nil, err
	}

	log.Log().Infoln(string(body))

	resp, err := m.doRequest(httpMethod, GraphAPI+"/v15.0/me/messenger_profile", bytes.NewReader(body))
	if err != nil {
		log.Log().Errorln(err)
		return nil, err
	}
	defer resp.Body.Close()
	pld, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Log().Infoln(string(pld))

	result := &result{}
	err = json.Unmarshal(pld, result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}
	return result, nil
}

// SetGreetingText sets a greeting text which is only rendered the first time user interacts with the Page on Messenger
// https://developers.facebook.com/docs/messenger-platform/thread-settings/greeting-text
// text must be UTF-8 and have a 160 character limit
func (m *Messenger) SetGreetingText(text string) error {
	if len(text) > 160 {
		return fmt.Errorf("Greeting text is too long. It has to be at most 160 characters long.")
	}

	result, err := m.changePersistentMenuSettings(http.MethodPost, &persistentMenu{})
	if err != nil {
		return err
	}
	if result.Result != "Successfully updated greeting" {
		return fmt.Errorf("Error occured while setting greeting, invalid result: %s", result.Result)
	}
	return nil
}

// Seems to not be supported yet
// func (m *Messenger) DeleteGreetingText() error {
// 	result, err := m.changePersistentMenuSettings(http.MethodDelete, &threadSettings{
// 		Type: settingTypeGreeting,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(result)
// 	return nil
// }

type ctaPayload struct {
	Payload string `json:"payload"`
}

// SetGetStartedButton sets a button which is shown at the bottom of the window ans is only rendered the first time the user interacts with the Page on Messenger
// When this button is tapped, we will trigger the postback received callback and deliver the person's page-scoped ID (PSID).
// You can then present a personalized message to greet the user or present buttons to prompt him or her to take an action.
// https://developers.facebook.com/docs/messenger-platform/thread-settings/get-started-button
func (m *Messenger) SetGetStartedButton(payload string) error {
	result, err := m.changeGetStartedSettings(http.MethodPost, &getStarted{
		GetStarted: getStartedPayload{
			Payload: payload,
		},
	})
	if err != nil {
		return err
	}
	if result.Result != "success" {
		return fmt.Errorf("Error occured while setting get started button, invalid result: %s", result.Result)
	}
	return nil
}

// DeleteGetStartedButton delets a button set by SetGetStartedButton
func (m *Messenger) DeleteGetStartedButton() error {
	result, err := m.changePersistentMenuSettings(http.MethodDelete, &persistentMenu{})
	if err != nil {
		return err
	}
	if result.Result != "Successfully deleted all new_thread's CTAs" {
		return fmt.Errorf("Error occured while deleting get started button, invalid result: %s", result.Result)
	}
	return nil
}

// SetPersistentMenu sets a Persistent Menu is a persistent menu that is always available to the user.
// This menu should contain top-level actions that users can enact at any point.
// Having a persistent menu easily communicates the basic capabilities of your bot for first-time and returning users.
// The menu will automatically appear in a thread if the person has been away for a certain period of time and return..
// https://developers.facebook.com/docs/messenger-platform/thread-settings/persistent-menu
func (m *Messenger) SetPersistentMenu(buttons []template.Button) error {
	if len(buttons) > 5 {
		return fmt.Errorf("Too many menu buttons, number of elements in the buttons array is limited to 5.")
	}
	result, err := m.changePersistentMenuSettings(http.MethodPost, &persistentMenu{
		PersistentMenu: []persistentMenuSettings{
			{
				Locale:                "default",
				ComposerInputDisabled: false,
				CallToActions:         buttons,
			},
		},
	})
	if err != nil {
		log.Log().Errorln(err)
		return err
	}
	if result.Result != "success" {
		return fmt.Errorf("Error occured while setting persistent menu, invalid result: %s", result.Result)
	}
	return nil
}

// DeletePersistentMenu deletes a menu set by SetPersistentMenu
func (m *Messenger) DeletePersistentMenu() error {
	result, err := m.changePersistentMenuSettings(http.MethodDelete, &persistentMenu{})
	if err != nil {
		return err
	}
	if result.Result != "Successfully deleted structured menu CTAs" {
		return fmt.Errorf("Error occured while deleting persistent menu, invalid result: %s", result.Result)
	}
	return nil
}
