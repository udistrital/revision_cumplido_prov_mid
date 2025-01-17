package models

type ContratoSupervisor struct {
	Dependencias_supervisor []Dependencia                  `json:"dependencias_supervisor"`
	Contratos               []InformacionContratoProveedor `json:"contratos"`
}
