package models

type CumplidosFiltrados struct {
	NumeroContrato  string `json:"NumeroContrato"`
	Vigencia        string `json:"Vigencia"`
	Rp              string `json:"Rp"`
	Mes             int    `json:"Mes"`
	FechaAprobacion string `json:"FechaAprobacion"`
	NombreProveedor string `json:"NombreProveedor"`
	Dependencia     string `json:"Dependencia"`
	Estado          string `json:"Estado"`
}
