package models

type ContratoDependencia struct {
	Contratos struct {
		Contrato []Contrato `json:"contrato"`
	} `json:"contratos"`
}

type Contrato struct {
	Vigencia       string `json:"vigencia"`
	NumeroContrato string `json:"numero_contrato"`
}
