package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type WebhookData struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Content   string `json:"content"`
}

const webhookURL string = "https://discord.com/api/webhooks/1148013439818674246/6ZFnpo4aCn1aopCV_OPQP0cST9_7oOf7DmaCoucwtorSkjT0AIFvX8T2dQk17jYWuMGN"

func main() {

	log.Println("Webhook caller with file uploads sample")

	webhookData := WebhookData{
		Username:  "Webhook Friend",
		AvatarURL: "",
		Content:   "Hello, check these attached files!",
	}

	payloadJson, _ := json.MarshalIndent(webhookData, "  ", "    ")
	log.Println("Sending payload:\n", string(payloadJson))

	multipartStuff := map[string]string{
		"file1":        "@files/birds.mp4",
		"file2":        "@files/sample_img.png",
		"payload_json": string(payloadJson),
	}

	contentType, formBody, err := createForm(multipartStuff)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Using content type:", contentType)
	resp, err := http.Post(webhookURL, contentType, formBody)
	if err != nil {
		log.Panicln(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	bodyString := string(bodyBytes)

	log.Println(resp.Status)
	log.Println(bodyString)
}

// Working with multipart/formdata in Go
// https://stackoverflow.com/a/67865385
func createForm(form map[string]string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil {
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}
