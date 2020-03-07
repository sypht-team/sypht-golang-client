package main

import (
	"fmt"
	"os"

	"github.com/sypht-team/sypht-golang-client"
)

const (
	fileName = "example/assets/sample_invoice.pdf"
)

func main() {
	client, _ := sypht.NewSyphtClient(os.Getenv("SYPHT_API_KEY"), nil)
	uploaded, _ := client.Upload(fileName, []string{
		sypht.Invoice,
		sypht.Document,
	}, "")
	result, _ := client.Results(uploaded.FileID)
	fmt.Println(result)
}
