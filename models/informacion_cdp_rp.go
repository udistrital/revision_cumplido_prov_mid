package models

import (
	"time"
)

type InformacionCdpRp struct {
	CdpXRp struct {
		CdpRp []struct {
			RpFechaRegistro         time.Time `json:"RP_FECHA_REGISTRO"`
			RpNumeroRegistro        string    `json:"RP_NUMERO_REGISTRO"`
			RpVigencia              string    `json:"RP_VIGENCIA"`
			CdpFechaExpedicion      time.Time `json:"CDP_FECHA_EXPEDICION"`
			CdpVigencia             string    `json:"CDP_VIGENCIA"`
			CdpNumeroDisponibilidad string    `json:"CDP_NUMERO_DISPONIBILIDAD"`
		} `json:"cdprp"`
	} `json:"cdpxrp"`
}
