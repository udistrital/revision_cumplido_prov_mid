package models

import "time"

type CumplidoProveedor struct {
	Id                int       "json:id"
	NumeroContrato    string    "json:numero_contrato"
	VigenciaContrato  int       "json:vigencia_contrato"
	FechaModificacion time.Time "json:fecha_modificacion"
	FechaCreacion     time.Time "json:fecha_creacion"
	Activo            bool      "json:activo"
}

type SolicituRevisionCumplidoProveedor struct {
	CumplidoId        int
	TipoContrato      string
	NumeroContrato    string
	VigenciaContrato  int
	Rp                string
	NombreProveedor   string
	Dependencia       string
	Cdp               string
	VigenciaRP        string
	FechaModificacion time.Time
	FechaCreacion     time.Time
	Activo            bool
}
