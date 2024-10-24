package models

type Historico struct {
	NombreResponsable string `json:"nombreResponsable"`
	Estado            string `json:"estado"`
	Fecha             string `json:"fecha"`
	Cargo             string `json:"cargo"`
}
