package models

import "time"

type ComentarioSoporte struct {
	Id                     int
	Comentario             string
	FechaCreacion          time.Time
	FechaModificacion      time.Time
	CambioEstadoCumplidoId *CambioEstadoCumplido
	Activo                 bool
	SoportePagoId          *SoportePago
}

type RespuestaComentarioSoporte struct {
	SoportePagoId          int    `json:"soporte_pago_id"`
	CambioEstadoCumplidoId int    `json:"cambio_estado_cumplido_id"`
	Comentario             string `json:"comentario"`
}
