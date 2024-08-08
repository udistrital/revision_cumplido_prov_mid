package models

type ContratoSupervisor struct {
	NombreSupervisor        string                         `json:"nombre_supervisor"`
	Dependencias_supervisor []Dependencia                  `json:"dependencias_supervisor"`
	Contratos               []InformacionContratoProveedor `json:"contratos"`
}
