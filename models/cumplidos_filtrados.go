package models

type CumplidosFiltrados struct {
	NumeroContrato    string `json:"NumeroContrato"`
	Vigencia          string `json:"Vigencia"`
	Rp                string `json:"Rp"`
	Mes               int    `json:"Mes"`
	FechaCambioEstado string `json:"FechaCambioEstado"`
	NombreProveedor   string `json:"NombreProveedor"`
	Dependencia       string `json:"Dependencia"`
	Estado            string `json:"Estado"`
	TipoContrato      string `json:"TipoContrato"`
	IdCumplido        int    `json:"IdCumplido"`
}
