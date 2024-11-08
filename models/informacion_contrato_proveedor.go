package models

import "time"

type InformacionContratoProveedor struct {
	TipoContrato           string
	NumeroContratoSuscrito string
	Vigencia               string
	NumeroRp               string
	VigenciaRp             string
	RPFechaRegistro        time.Time
	NombreProveedor        string
	NombreDependencia      string
	NumeroCdp              string
	VigenciaCdp            string
	CDPFechaExpedicion     time.Time
	Rubro                  string
	IdProveedor            int
}
