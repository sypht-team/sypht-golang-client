package sypht

import (
	"log"
	"os"
	"testing"
	"time"
)

func initTest() (client *Client, err error) {
	client, err = NewSyphtClient(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), nil)
	if err != nil {
		log.Fatalf("Fail to initialize Sypht client %v", err)
	}

	_, err = client.RefreshToken()
	if err != nil {
		log.Println(err)
	}
	return
}

func TestSypht_UploadFile(t *testing.T) {
	client, err := initTest()
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Upload("assets/sample_invoice.pdf", []string{})
	if err != nil {
		t.Error(err)
	}
	if resp["fileId"].(string) == "" {
		t.Error("Empty fileID")
	}
}

func TestSypht_PredictionWithCustomFieldSet(t *testing.T) {
	client, err := initTest()
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Upload("assets/sample_invoice.pdf", []string{
		Invoice,
		Document,
	})
	if err != nil {
		t.Error(err)
	}
	if resp["fileId"].(string) == "" {
		t.Error("Empty fileID")
	}
	time.Sleep(time.Second * 10)
	_, err = client.Results(resp["fileId"].(string))
	if err != nil {
		t.Error(err)
	}
}
