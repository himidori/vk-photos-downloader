package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/himidori/vkapi"
)

func main() {
	var login string
	var pass string
	var UID string

	fmt.Print("Login: ")
	fmt.Scanln(&login)
	fmt.Print("Password: ")
	fmt.Scanln(&pass)
	fmt.Print("UID (empty to download from every dialog): ")
	fmt.Scanln(&UID)

	client, err := vkapi.NewVKClient(vkapi.DeviceIPhone, login, pass)
	if err != nil {
		log.Printf("error occured: %s\n", err)
		return
	}

	log.Printf("authorized as %s %s\n", client.Self.FirstName, client.Self.LastName)

	if !folderExists("photos") {
		mkdir("photos")
	}

	if UID != "" {
		id, err := strconv.ParseInt(UID, 10, 32)
		if err != nil {
			log.Printf("incorrect ID string: %s\n", err)
			return
		}

		photos, err := getAttachments(client, int(id))
		if err != nil {
			log.Printf("failed to get attachments: %s\n", err)
			return
		}

		if len(photos) == 0 {
			log.Printf("there's no attachments for UID %d\n", id)
			return
		}

		downloadPath := "photos/" + UID
		mkdir(downloadPath)

		downloaded := 0
		gocounter := 0
		limit := 10
		queue := len(photos)
		total := len(photos)
		var wg sync.WaitGroup

		for _, p := range photos {
			link := getBestLink(p)
			path := downloadPath + "/" + getFileName(link)
			wg.Add(1)

			go func() {
				err := downloadFile(link, path)
				if err != nil {
					log.Printf("failed to download file: %s\n", err)
				}
				wg.Done()
				queue--
				downloaded++
				log.Printf("downloaded %d/%d photos\n", downloaded, total)
			}()

			if queue < limit {
				limit = queue
			}

			gocounter++
			if gocounter == limit {
				wg.Wait()
				gocounter = 0
			}
		}
	} else {
		dialogs, err := getDialogs(client)
		if err != nil {
			log.Printf("failed to get dialogs: %s\n", err)
			return
		}

		curUser := 0
		for _, d := range dialogs {
			photos, err := getAttachments(client, d.UID)
			if err != nil {
				log.Printf("failed to get attachments for UID %d\n", d.UID)
				continue
			}

			if len(photos) == 0 {
				continue
			}

			downloadPath := "photos/" + strconv.Itoa(d.UID)
			mkdir(downloadPath)

			curUser++
			downloaded := 0
			gocounter := 0
			limit := 10
			queue := len(photos)
			total := len(photos)
			totalUsers := len(dialogs)
			var wg sync.WaitGroup

			for _, p := range photos {
				link := getBestLink(p)
				path := downloadPath + "/" + getFileName(link)
				wg.Add(1)

				go func() {
					err := downloadFile(link, path)
					if err != nil {
						log.Printf("failed to download file: %s\n", err)
					}
					wg.Done()
					queue--
					downloaded++
					log.Printf("downloaded %d/%d photos. user %d/%d\n",
						downloaded, total, curUser, totalUsers)
				}()

				if queue < limit {
					limit = queue
				}

				gocounter++
				if gocounter == limit {
					wg.Wait()
					gocounter = 0
				}
			}
		}
	}
}
