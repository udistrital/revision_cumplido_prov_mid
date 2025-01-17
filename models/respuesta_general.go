package models

type RespuestaGeneral struct {
	Data    []map[string]interface{}
	Message string
	Status  string
	Success bool
}
