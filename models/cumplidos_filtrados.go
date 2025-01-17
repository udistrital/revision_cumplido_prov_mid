package models

type CumplidosFiltrados struct {
	IdCumplido      int    `json:"IdCumplido"`
	NumeroContrato  string `json:"NumeroContrato"`
	Vigencia        string `json:"Vigencia"`
	Rp              string `json:"Rp"`
	NombreProveedor string `json:"NombreProveedor"`
	Dependencia     string `json:"Dependencia"`
	Estado          string `json:"Estado"`
	TipoContrato    string `json:"TipoContrato"`
	InformacionPago string `json:"InformacionPago"`
}
