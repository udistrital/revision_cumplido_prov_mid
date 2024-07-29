package models

import "time"

type CumplidoProveedor struct {
	Id                int
	NumeroContrato    string
	VigenciaContrato  int
	FechaModificacion time.Time
	FechaCreacion     time.Time
	Activo            bool
}
