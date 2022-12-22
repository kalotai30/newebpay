package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"util_api/util/log"
)

type Connect interface {
	Call(mehod string, params []interface{}) string
}

type Request struct {
	Url      string
	UserId   string
	Password string
	IP       string
	Port     string
}

type Response struct {
	Id     int         `json:"id"`
	Error  interface{} `json:"error"`
	Result interface{} `json:"result"`
}

func NewRequest(ip string, port string, id string, password string) (req *Request) {

	req = &Request{
		"http://" + id + ":" + password + "@" + ip + ":" + port,
		id,
		password,
		ip,
		port,
	}

	return req
}

func (req *Request) Call(method string, params []interface{}) Response {

	request := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	body, err := json.Marshal(request)

	if err != nil {
		log.Error(err)
		return Response{}
	}

	res, err := http.NewRequest("POST", req.Url, bytes.NewBuffer(body))

	if err != nil {
		log.Error(err)
		return Response{}
	}

	res.Header.Set("Content-Type", "application/json")

	// Configure basic access authorization.
	res.SetBasicAuth(req.UserId, req.Password)

	client := &http.Client{}

	resp, err := client.Do(res)

	if err != nil {
		log.Error(err)
		return Response{}
	}

	defer closeRespBody(resp.Body)

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error(err)
		return Response{}
	}

	fmt.Println("response Body:", string(body))

	var u Response

	log.Error(json.Unmarshal(body, &u))

	return u
}

func closeRespBody(body io.ReadCloser) {
	log.Error(body.Close())
}

func ApiCall(url, method string, params map[string]interface{}) (string, error) {
	body, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
		return "", err
	}

	res, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Error(err)
		return "", err
	}

	res.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(res)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer closeRespBody(resp.Body)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	fmt.Println("response Body:", string(body))

	return string(body), nil
}

func AuthTokeApiCall(url, method string, authToken string, params map[string]interface{}) (string, error) {
	body, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
		return "", err
	}

	res, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Error(err)
		return "", err
	}

	if authToken != "" {
		res.Header.Add("Content-Type", "application/json")
		res.Header.Add("Auth-Token", authToken)
	} else {
		res.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}

	resp, err := client.Do(res)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer closeRespBody(resp.Body)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	fmt.Println("response Body:", string(body))
	return string(body), nil
}
