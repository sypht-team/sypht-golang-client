# Sypht HTTP Go Library

The Golang library to interact with Sypht's API

## Sypht

Sypht is an service which extracts key fields from documents. For example, you can upload an image or pdf of a bill or invoice and extract the amount due, due date, invoice number and biller information.

Pixels in, json out.

Checkout [sypht.com](https://sypht.com) for more details.

### API

Sypht provides a REST api for interaction with the service. Full documentation is available at: [docs.sypht.com](https://docs.sypht.com/).
This repository is an open-source Golang reference client implementation for working with the API.


### Getting started

To get started you'll need some API credentials, i.e. a `client_id` and `client_secret`.

Sypht is currently in closed-beta, if you'd like to try out the service contact: [support@sypht.com](mailto://support@sypht.com).

### Prerequisites

* Go - supports **Go 1.13 or greater**.

### Installation

```sh
$ go get github.com/sypht-team/sypht-golang-client
```


### Usage
Populate these system environment variables with the credentials generated above:

```Bash
CLIENT_ID="client_id"
CLIENT_SECRET="client_secret"
```

then invoke the client with a file of your choice:
```go
import "github.com/sypht-team/sypht-golang-client"

client, err := sypht.NewSyphtClientFromEnv(nil)
if err != nil {
    // handle error
}
resp, err := client.Upload("mytaxireceipt.pdf", []string{
    sypht.Invoice,
    sypht.Document,
})
if err != nil {
    // handle error
} else {
    sypht.PrettyPrintResponse(resp)
}
```


