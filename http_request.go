package PushSDK

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func postReqUrlencoded(requestPath string, forms map[string]string, header map[string]string) ([]byte, error) {
	fmt.Println(requestPath, forms, header)
	form := url.Values{}
	for k, v := range forms {
		form.Add(k, v)
	}
	param := form.Encode()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", requestPath), strings.NewReader(param))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, fmt.Sprintf(v))
		}
	}

	var client = &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return body, err
	}

	return body, nil
}

func postReqJson(requestPath string, forms string, header map[string]string) ([]byte, error) {
	fmt.Println(requestPath, forms, header)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", requestPath), strings.NewReader(forms))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, fmt.Sprintf(v))
		}
	}

	var client = &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return body, err
	}

	return body, nil
}
