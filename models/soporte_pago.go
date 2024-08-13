package models

import "time"

type SoportePago struct {
	Id                  int
	DocumentoId         int
	CumplidoProveedorId *CumplidoProveedor
	FechaCreacion       time.Time
	FechaModificacion   time.Time
	Activo              bool
}
