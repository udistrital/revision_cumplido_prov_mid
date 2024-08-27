package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
)

func GetJsonWSO2Test(urlp string, target interface{}) (status int, err error) {
	//fmt.Println("URL: ", urlp)
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
	//fmt.Println("URL: ", url)
	r, err := http.Get(url)
	fmt.Println(err)
	if err != nil {
		//fmt.Println("r", r)
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

func SendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	client := &http.Client{}
	//fmt.Print("Json que se le va a apasar a la funcion: ")
	//fmt.Println(b)
	req, err := http.NewRequest(trequest, url, b)
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	//fmt.Println("Respuesta de la peticion: ", r.StatusCode)
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

func ValorLetras(n int) string {
	var unidades = []string{"", "un", "dos", "tres", "cuatro", "cinco", "seis", "siete", "ocho", "nueve", "diez", "once", "doce", "trece", "catorce", "quince", "dieciséis", "diecisiete", "dieciocho", "diecinueve", "veinte"}
	var decenas = []string{"", "", "veinti", "treinta", "cuarenta", "cincuenta", "sesenta", "setenta", "ochenta", "noventa"}
	var centenas = []string{"", "cien", "doscientos", "trescientos", "cuatrocientos", "quinientos", "seiscientos", "setecientos", "ochocientos", "novecientos"}

	if n == 0 {
		return "cero"
	}

	// Millones
	if n >= 1000000 {
		if n == 1000000 {
			return "un millón"
		}
		return fmt.Sprintf("%s millones %s", ValorLetras(n/1000000), ValorLetras(n%1000000))
	}

	// Miles
	if n >= 1000 {
		if n == 1000 {
			return "mil"
		}
		return fmt.Sprintf("%s mil %s", ValorLetras(n/1000), ValorLetras(n%1000))
	}

	// Centenas
	if n >= 100 {
		if n == 100 {
			return "cien"
		}
		return fmt.Sprintf("%s %s", centenas[n/100], ValorLetras(n%100))
	}

	// Decenas
	if n >= 20 {
		if n%10 == 0 {
			return decenas[n/10]
		}
		if n < 30 {
			return decenas[n/10] + unidades[n%10]
		}
		return fmt.Sprintf("%s y %s", decenas[n/10], unidades[n%10])
	}

	// Unidades
	return unidades[n]
}

func FormatNumber(value interface{}, precision int, thousand string, decimal string) string {
	v := reflect.ValueOf(value)
	var x string
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x = fmt.Sprintf("%d", v.Int())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x = fmt.Sprintf("%d", v.Uint())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Float32, reflect.Float64:
		x = fmt.Sprintf(fmt.Sprintf("%%.%df", precision), v.Float())
	case reflect.Ptr:
		switch v.Type().String() {
		case "*big.Rat":
			x = value.(*big.Rat).FloatString(precision)

		default:
			panic("Unsupported type - " + v.Type().String())
		}
	default:
		panic("Unsupported type - " + v.Kind().String())
	}

	return FormatNumberString(x, precision, thousand, decimal)
}

func FormatNumberString(x string, precision int, thousand string, decimal string) string {
	lastIndex := strings.Index(x, ".") - 1
	if lastIndex < 0 {
		lastIndex = len(x) - 1
	}

	var buffer []byte
	var strBuffer bytes.Buffer

	j := 0
	for i := lastIndex; i >= 0; i-- {
		j++
		buffer = append(buffer, x[i])

		if j == 3 && i > 0 && !(i == 1 && x[0] == '-') {
			buffer = append(buffer, ',')
			j = 0
		}
	}

	for i := len(buffer) - 1; i >= 0; i-- {
		strBuffer.WriteByte(buffer[i])
	}
	result := strBuffer.String()

	if thousand != "," {
		result = strings.Replace(result, ",", thousand, -1)
	}

	extra := x[lastIndex+1:]
	if decimal != "." {
		extra = strings.Replace(extra, ".", decimal, 1)
	}

	return result + extra
}
