package main

import (
	"log"
	"os"

	"github.com/sypht-team/sypht-golang-client/pkg/sypht"
)

const (
	fileName = "pkg/sypht/assets/sample_invoice.pdf"
)

func main() {
	client, err := sypht.NewSyphtClient(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), nil)
	if err != nil {
		log.Fatalf("Fail to initialize Sypht client %v", err)
	}
	token, err := client.RefreshToken()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("api token is %s", token)
	}

	resp, err := client.Upload(fileName, []string{
		sypht.Invoice,
		sypht.Document,
	})
	if err != nil {
		log.Println(err)
	} else {
		sypht.PrettyPrintResponse(resp)
	}

	// result, err := client.Results(resp["fileId"].(string))
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	sypht.PrettyPrintResponse(result)
	// }

	// _, err = client.Image(resp["fileId"].(string), 1)
	// if err != nil {
	// 	log.Println(err)
	// }

}
