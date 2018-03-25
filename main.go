package main

import (
	"flag"
	"log"
	"strconv"
	"sync"

	"github.com/himidori/golang-vk-api"
)

var (
	user     string
	pass     string
	token    string
	uid      int
	routines int
	device   int
)

func init() {
	flag.StringVar(&user, "u", "", "username or phone number")
	flag.StringVar(&pass, "p", "", "password")
	flag.StringVar(&token, "t", "", "access token (can be provided for authorization instead of user/pass)")
	flag.IntVar(&uid, "uid", 0, "user id (default: 0. pictures from every dialog will be downloaded instead of an individual user)")
	flag.IntVar(&routines, "r", 0, "number of goroutines for concurrent photo download (default: 0)")
	flag.IntVar(&device, "d", 0, "device to use for authorization. (default: 0. 0 - iPhone, 1 - Android, 2 - WPhone)")
}

func auth() (client *vkapi.VKClient, err error) {
	switch {
	case token != "":
		client, err = vkapi.NewVKClientWithToken(token)
	default:
		client, err = vkapi.NewVKClient(device, user, pass)
	}

	return client, err
}

func download(photos []*vkapi.PhotoAttachment, uid int) {
	if !folderExists("photos") {
		mkdir("photos")
	}

	downloadPath := "photos/" + strconv.Itoa(uid)
	mkdir(downloadPath)

	counter := 0
	total := len(photos)
	queue := len(photos)
	gocounter := 0
	limit := routines
	var wg sync.WaitGroup
	for _, p := range photos {
		link := getBestLink(p)
		path := downloadPath + "/" + getFileName(link)

		if routines > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := downloadFile(link, path)
				if err != nil {
					log.Printf("failed to download file: %s\n", err)
				}
				queue--
				counter++
				log.Printf("downloaded: %d/%d\n", counter, total)
			}()

			if queue < limit {
				limit = queue
			}

			gocounter++
			if gocounter == limit {
				wg.Wait()
				gocounter = 0
			}

		} else {
			err := downloadFile(link, path)
			if err != nil {
				log.Printf("failed to download file: %s\n", err)
			}
			log.Printf("downloaded: %d/%d\n", counter, total)
			counter++
		}
	}
}

func download_mass(client *vkapi.VKClient) {
	dialogs, err := getDialogs(client)
	if err != nil {
		log.Fatal("failed to get dialogs: %s\n", err)
	}

	log.Printf("received %d dialogs\n", len(dialogs))

	for _, dlg := range dialogs {
		photos, err := getAttachments(client, dlg.UID)
		if err != nil {
			log.Printf("failed to get attachments for id%d\n", dlg.UID)
			continue
		}

		if len(photos) == 0 {
			continue
		}

		download(photos, dlg.UID)
	}
}

func main() {
	flag.Parse()

	if user == "" || pass == "" {
		if token == "" {
			log.Fatal("no credentials were provided")
		}
	}

	client, err := auth()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("authorized as %s %s (id%d)\n", client.Self.FirstName,
		client.Self.LastName, client.Self.UID)

	if uid != 0 {
		photos, err := getAttachments(client, uid)
		if err != nil {
			log.Fatal("failed to get attachments: %s\n", err)
		}
		if len(photos) == 0 {
			log.Fatal("no attachments found for given uid")
		}
		log.Printf("received %d attachments for id%d\n", len(photos), uid)
		download(photos, uid)
	} else {
		download_mass(client)
	}
}
