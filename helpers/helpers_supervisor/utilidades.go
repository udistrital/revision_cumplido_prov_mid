package helpers_supervisor

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
)

func getJsonWSO2Test(urlp string, target interface{}) (status int, err error) {
	b := new(bytes.Buffer)
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlp, b)
	req.Header.Set("Accept", "application/json")
	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(nil, err)
		}
	}()

	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

func getJsonTest(url string, target interface{}) (status int, err error) {
	r, err := http.Get(url)
	if err != nil {
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}
