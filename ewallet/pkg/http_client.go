package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"simple-go/application/config"
	"simple-go/application/constants"
	"strconv"
	"time"
)

type HttpClient struct {
	url    string
	client *http.Client
}

type RequestHeader struct {
	IpAddress     string
	UserAgent     string
	Cookie        string
	Authorization string
	Apikey        string
	CookieKey     string
	ContentType   string `json:"Content-Type"`
}

func NewHttpClient() HttpClient {
	gateway := config.GetString(config.CFG_GATEWAY, "")
	stage := config.GetString(config.CFG_STAGE, "")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return HttpClient{
		url:    gateway + stage + "/",
		client: client,
	}
}

func (h HttpClient) GetThirdParty(headers RequestHeader, path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	client, req := h.SetHeaders(headers, h.client, req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return body, errors.New(fmt.Sprintf("error from rpi with code %s _ %s, response: %s", strconv.Itoa(resp.StatusCode), path, string(body)))
	}

	return body, nil
}

func (h HttpClient) PostThirdParty(headers RequestHeader, path string, data []byte) ([]byte, error) {
	log.Println(path, " :path")
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	client, req := h.SetHeaders(headers, h.client, req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return body, errors.New(fmt.Sprintf("error from rpi with code %s _ %s, response: %s", strconv.Itoa(resp.StatusCode), path, string(body)))
	}

	return body, nil
}

func (h HttpClient) Patch(headers RequestHeader, path string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPatch, h.url+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	client, req := h.SetHeaders(headers, h.client, req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return body, errors.New(fmt.Sprintf("error from rpi with code %s _ %s, response: %s", strconv.Itoa(resp.StatusCode), path, string(body)))
	}

	return body, nil
}

func (h HttpClient) SetHeaders(headers RequestHeader, client *http.Client, req *http.Request) (*http.Client, *http.Request) {
	if headers.Authorization != "" {
		req.Header.Set(constants.HeaderKeyAuthorization, headers.Authorization)
	}

	if headers.Cookie != "" {
		jar, _ := cookiejar.New(nil)
		client.Jar = jar
		req.AddCookie(&http.Cookie{
			Name:  headers.CookieKey,
			Value: headers.Cookie,
		})
	}

	if headers.ContentType != "" {
		req.Header.Set("Content-Type", headers.ContentType)
	} else {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	return client, req
}
