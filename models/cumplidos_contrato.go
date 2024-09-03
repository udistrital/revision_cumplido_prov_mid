package models

import "time"

type CumplidosContrato struct {
	ConsecutivoCumplido int                `json:"ConsecutivoCumplido"`
	NumeroContrato      string             `json:"NumeroContrato"`
	FechaCreacion       time.Time          `json:"FechaCreacion"`
	Periodo             string             `json:"Periodo"`
	EstadoCumplido      string             `json:"EstadoCumplido"`
	CumplidoProveedorId *CumplidoProveedor `json:"CumplidoProveedorId"`
}
