package main

import (
	"net/url"
	"strconv"

	"github.com/himidori/golang-vk-api"
)

func getDialogs(client *vkapi.VKClient) ([]*vkapi.DialogMessage, error) {
	offset := 0
	params := url.Values{}
	messages := []*vkapi.DialogMessage{}

	for {
		params.Set("offset", strconv.Itoa(offset))
		dialogs, err := client.DialogsGet(200, params)
		if err != nil {
			return nil, err
		}

		if len(dialogs.Messages) > 0 {
			for _, msg := range dialogs.Messages {
				messages = append(messages, msg.Message)
			}
		}

		offset += 200

		if len(messages) >= dialogs.Count {
			break
		}
	}

	return messages, nil
}
