package helpers_estado_pago

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func sendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	client := &http.Client{}
	fmt.Print("Json que se le va a apasar a la funcion: ")
	//fmt.Println(b)
	req, err := http.NewRequest(trequest, url, b)
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	fmt.Println("Respuesta de la peticion: ", r.StatusCode)
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

func sendJson3(url string, trequest string, target interface{}, datajson interface{}) error {
	// Convertir datajson en un Formato Json para poderlo enviar como parametro
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			fmt.Println(err)
			beego.Error(err)
		}
	}

	// Crear una nueva solicitud POST con el cuerpo del JSON
	req, err := http.NewRequest(trequest, url, b)
	if err != nil {
		fmt.Println("Error al crear la solicitud POST:", err)
		return err
	}

	//Configurar el encabezado Accept
	req.Header.Set("Accept", "application/json")
	// Configurar el encabezado Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Configurar el cliente HTTP con tiempo de espera y tamaño de búfer
	client := &http.Client{
		Timeout: time.Second * 10, // Tiempo de espera máximo de 10 segundos
		Transport: &http.Transport{
			MaxIdleConns:        100, // Número máximo de conexiones inactivas permitidas
			MaxIdleConnsPerHost: 100, // Número máximo de conexiones inactivas permitidas por host
		},
	}

	// Realizar la solicitud POST
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al enviar la solicitud POST:", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(resp.Body).Decode(target)
}
