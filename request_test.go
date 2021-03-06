package sypht

import (
	"log"
	"os"
	"testing"
	"time"
)

func initTest() (client *Client, err error) {
	client, err = NewSyphtClient(os.Getenv("SYPHT_API_KEY"), nil)
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
	resp, err := client.Upload("example/assets/sample_invoice.pdf", []string{
		Invoice,
		Document,
	}, "")
	if err != nil {
		t.Error(err)
	}
	if resp.FileID == "" {
		t.Error("Empty fileID")
	}
}

func TestSypht_PredictionWithCustomFieldSet(t *testing.T) {
	client, err := initTest()
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Upload("example/assets/sample_invoice.pdf", []string{
		Invoice,
		Document,
	}, "")
	if err != nil {
		t.Error(err)
	}
	if resp.FileID == "" {
		t.Error("Empty fileID")
	}
	time.Sleep(time.Second * 10)
	_, err = client.Results(resp.FileID)
	if err != nil {
		t.Error(err)
	}
}
