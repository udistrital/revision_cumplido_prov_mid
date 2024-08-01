package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
)

func GetJsonWSO2Test(urlp string, target interface{}) (status int, err error) {
	fmt.Println("URL: ", urlp)
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

func GetJsonTest(url string, target interface{}) (status int, err error) {
	fmt.Println("URL: ", url)
	r, err := http.Get(url)
	fmt.Println(err)
	if err != nil {
		fmt.Println("r", r)
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

func LimpiezaRespuestaRefactor(respuesta map[string]interface{}, v interface{}) {
	b, err := json.Marshal(respuesta["Data"])
	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, v)
}

func SendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	client := &http.Client{}
	fmt.Print("Json que se le va a apasar a la funcion: ")
	fmt.Println(b)
	req, err := http.NewRequest(trequest, url, b)
	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}
