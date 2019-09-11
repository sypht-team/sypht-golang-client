package sypht

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Upload uploads files with fileName
// options are list of fieldSets constant
func (s *Client) Upload(fileName string, options []string) (resp map[string]interface{}, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	ok := checkFileExt(filepath.Ext(strings.TrimSpace(fileName)))
	if !ok {
		err = fmt.Errorf("unsupported file : %s", fileName)
		return
	}

	file, err := os.Open(fileName)
	defer file.Close()
	part, err := writer.CreateFormFile("fileToUpload", filepath.Base(fileName))
	_, err = io.Copy(part, file)

	fieldSets := parseOptions(options)
	_ = writer.WriteField("fieldSets", fieldSets)
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", s.config.apiBaseURL+"/fileupload", payload)

	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", strings.Join([]string{"Bearer ", s.getToken()}, ""))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, &resp)
	return
}

// Results fetches results of uploaded file
func (s *Client) Results(fileID string) (out map[string]interface{}, err error) {
	url := strings.Join([]string{s.config.apiBaseURL, "/result/final/", fileID}, "")
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", strings.Join([]string{"Bearer ", s.getToken()}, ""))
	res, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, &out)
	return
}

// Image retrieves an image copy of the uploaded document.
func (s *Client) Image(fileID string, page int) (file []byte, err error) {
	if page <= 0 {
		page = 1
	}
	queryParam := fmt.Sprintf("?page=%d", page)
	url := strings.Join([]string{s.config.apiBaseURL, "/result/image/", fileID, queryParam}, "")

	log.Printf("get image url %s", url)
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", strings.Join([]string{"Bearer ", s.getToken()}, ""))

	res, err := s.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	file, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}

//PrettyPrintResponse pretty printing
func PrettyPrintResponse(mp map[string]interface{}) {
	b, err := json.MarshalIndent(mp, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func parseOptions(options []string) (out string) {
	if len(options) == 0 {
		return
	}
	return "[" + strings.Join(options, ",") + "]"
}

func checkFileExt(ext string) (ok bool) {
	ext = strings.ToLower(ext)
	supportedExt := []string{".pdf", ".jpeg", ".png", ".gif"}
	for _, e := range supportedExt {
		if e == ext {
			return true
		}
	}
	return
}
