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

//UploadResponse response struct
type UploadResponse struct {
	FileID     string `json:"fileId"`
	UploadedAt string `json:"uploadedAt"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

// Upload uploads files with fileName
// options are list of fieldSets constant
func (s *Client) Upload(fileName string, options []string, workflowID string) (resp UploadResponse, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(fileName)
	defer file.Close()

	cType, err := getFileContentType(file)
	if err != nil {
		log.Println(err)
		return
	}

	ok := validateFileFormat(cType, filepath.Ext(strings.TrimSpace(fileName)))
	if !ok {
		err = fmt.Errorf("unsupported file : %s", fileName)
		return
	}

	part, err := writer.CreateFormFile("fileToUpload", filepath.Base(fileName))
	_, err = io.Copy(part, file)

	fieldSets := parseOptions(options)
	_ = writer.WriteField("fieldSets", fieldSets)
	if workflowID != "" {
		_ = writer.WriteField("workflowId", workflowID)
	}
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
	err = json.Unmarshal(body, &resp)

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

func parseOptions(options []string) string {
	if len(options) == 0 {
		return "[]"
	}
	return "[" + strings.Join(options, ",") + "]"
}

func validateFileFormat(format, ext string) (ok bool) {
	ext = strings.ToLower(ext)
	supportedType := []string{"application/pdf", "image/jpeg", "image/png", "image/gif"}
	for _, t := range supportedType {
		if t == format {
			return true
		}
	}
	// have to check extension for tiff file since image/tiff not supported in go's mime contentType
	return ext == ".tiff"
}

func getFileContentType(out *os.File) (contentType string, err error) {
	buffer := make([]byte, 512)
	_, err = out.Read(buffer)
	if err != nil {
		return
	}
	out.Seek(0, 0)
	contentType = http.DetectContentType(buffer)
	return
}
