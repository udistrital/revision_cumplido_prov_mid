package models

import "time"

type HistoricoCumplido struct {
	NombreResponsable string    `json:"nombreResponsable"`
	Estado            string    `json:"estado"`
	Fecha             time.Time `json:"fecha"`
	CargoResponsable  string    `json:"cargo"`
}
