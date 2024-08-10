package models

import "time"

type CumplidoProveedor struct {
	Id                 int
	TipoContrato       string
	NumeroContrato     string
	VigenciaContrato   int
	Rp                 string
	NombreProveedor    string
	Dependencia        string
	Cdp                string
	NombreOrdenador    string
	DocumentoOrdenador string
	VigenciaRP         string
	FechaModificacion  time.Time
	FechaCreacion      time.Time
	Activo             bool
}
