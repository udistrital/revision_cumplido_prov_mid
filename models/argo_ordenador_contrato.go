package models

import "time"

type ArgoOrdenadorContrato struct {
	Id              int       `json:"Id"`
	IdOrdenador     int       `json:"IdOrdenador"`
	InfoResolucion  string    `json:"InfoResolucion"`
	IdCiudad        int       `json:"IdCiudad"`
	FechaInicio     time.Time `json:"FechaInicio"`
	FechaFin        time.Time `json:"FechaFin"`
	Estado          bool      `json:"Estado"`
	Documento       int       `json:"Documento"`
	NombreOrdenador string    `json:"NombreOrdenador"`
	RolOrdenador    string    `json:"RolOrdenador"`
	RolId           int       `json:"RolId"`
}
