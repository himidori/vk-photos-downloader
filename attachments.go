package main

import (
	"net/url"

	"github.com/himidori/vkapi"
)

func getAttachments(client *vkapi.VKClient, UID int) ([]*vkapi.PhotoAttachment, error) {
	attachments := []*vkapi.PhotoAttachment{}
	params := url.Values{}

	for {
		att, err := client.GetHistoryAttachments(UID, "photo", 200, params)
		if err != nil {
			return nil, err
		}

		if len(att.Attachments) == 0 {
			break
		}

		for _, photo := range att.Attachments {
			attachments = append(attachments, photo.Attachment.Photo)
		}

		params.Set("start_from", att.NextFrom)
	}

	return attachments, nil
}
