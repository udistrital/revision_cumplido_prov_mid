package models

import "time"

type CambioEstadoCumplido struct {
	Id                   int
	EstadoCumplidoId     *EstadoCumplido
	CumplidoProveedorId  *CumplidoProveedor
	DocumentoResponsable int
	CargoReponsable      string
	FechaCreacion        time.Time
	FechaModificacion    time.Time
	Activo               bool
}
