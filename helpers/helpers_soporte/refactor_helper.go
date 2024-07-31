package helpers_soporte

import (
	"encoding/json"
)

func LimpiezaRespuestaRefactor(respuesta map[string]interface{}, v interface{}) {
	b, err := json.Marshal(respuesta["Data"])
	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, v)
}
