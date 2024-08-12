package models

type Banco struct {
	Id                     int    `json:"Id"`
	NombreBanco            string `json:"NombreBanco"`
	DenominacionBanco      string `json:"DenominacionBanco"`
	Descripcion            string `json:"Descripcion"`
	Nit                    string `json:"Nit"`
	CodigoSuperintendencia int    `json:"CodigoSuperintendencia"`
	CodigoAch              int    `json:"CodigoAch"`
	EstadoActivo           bool   `json:"EstadoActivo"`
}
