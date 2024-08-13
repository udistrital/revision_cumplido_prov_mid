package models

type ContratoDependencia struct {
	Contratos struct {
		Contrato []ContratoDep `json:"contrato"`
	} `json:"contratos"`
}

type ContratoDep struct {
	Vigencia       string `json:"vigencia"`
	NumeroContrato string `json:"numero_contrato"`
}
