package models

type SupervisorContratoProveedor struct {
	Contratos struct {
		Supervisor []struct {
			Documento string `json:"documento"`
			Cargo     string `json:"cargo"`
			Nombre    string `json:"nombre"`
		} `json:"supervisor"`
	} `json:"contratos"`
}
