# Sypht Golang Client
This repository is a Golang reference client implementation for working with the Sypht API at https://api.sypht.com.

## About Sypht
[Sypht](https://sypht.com) is a SAAS [API]((https://docs.sypht.com/)) which extracts key fields from documents. For 
example, you can upload an image or pdf of a bill or invoice and extract the amount due, due date, invoice number 
and biller information. 

### Getting started
To get started you'll need API credentials, i.e. a `client_id` and `client_secret`, which can be obtained by registering
for an [account](https://www.sypht.com/signup/developer)

### Prerequisites
* Go - supports **Go 1.13 or greater**.

### Installation
```sh
$ go get github.com/sypht-team/sypht-golang-client
```

### Usage
Populate system environment variable with the credentials generated above:

```Bash
SYPHT_API_KEY="$client_id:$client_secret"
```

then invoke the client with a file of your choice:
```go
client, _ := sypht.NewSyphtClient(os.Getenv("SYPHT_API_KEY"), nil)

	uploaded, _ := client.Upload(fileName, []string{
		sypht.Invoice,
		sypht.Document,
	})

	result, _ := client.Results(uploaded["fileId"].(string))
	sypht.PrettyPrintResponse(result)
```


