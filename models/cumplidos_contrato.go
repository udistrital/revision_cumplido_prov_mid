package models

type CumplidosContrato struct {
	NumeroContrato      string             `json:"NumeroContrato"`
	FechaCreacion       string             `json:"FechaCreacion"`
	Periodo             string             `json:"Periodo"`
	EstadoCumplido      string             `json:"EstadoCumplido"`
	CumplidoProveedorId *CumplidoProveedor `json:"CumplidoProveedorId"`
}
