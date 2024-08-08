package models

import "time"

type CDPRP struct {
	CDPNumeroDisponibilidad string    `json:"CDP_NUMERO_DISPONIBILIDAD"`
	RPVigencia              string    `json:"RP_VIGENCIA"`
	CDPVigencia             string    `json:"CDP_VIGENCIA"`
	RPFechaRegistro         time.Time `json:"RP_FECHA_REGISTRO"`
	RPNumeroRegistro        string    `json:"RP_NUMERO_REGISTRO"`
	CDPFechaExpedicion      time.Time `json:"CDP_FECHA_EXPEDICION"`
}
