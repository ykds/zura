package main

import (
	"encoding/json"
	"github.com/ykds/zura/internal/logic/services/user"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	base := "http://localhost:8000/api"
	login := &user.LoginResponse{}
	request(http.MethodPost, base+"/users/login", map[string]string{"login_type": "username", "username": "test", "password": "123456"}, login, "")

}

func request(method string, api string, data map[string]string, result interface{}, token string) {
	vs := url.Values{}
	if len(data) > 0 {
		for k, v := range data {
			vs.Set(k, v)
		}
	}
	req, err := http.NewRequest(method, api, strings.NewReader(vs.Encode()))
	if err != nil {
		panic(err)
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("token", token)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if result != nil {
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bytes, result)
		if err != nil {
			panic(err)
		}
	}
}
