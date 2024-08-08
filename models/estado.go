package models

type EstadoCumplidoId struct {
	Id                int    `json:"Id"`
	Nombre            string `json:"Nombre"`
	Descripcion       string `json:"Descripcion"`
	Activo            bool   `json:"Activo"`
	CodigoAbreviación string `json:"CodigoAbreviación"`
}

type CumplidoProveedorId struct {
	Id                int    `json:"Id"`
	NumeroContrato    string `json:"NumeroContrato"`
	VigenciaContrato  int    `json:"VigenciaContrato"`
	FechaModificacion string `json:"FechaModificacion"`
	FechaCreacion     string `json:"FechaCreacion"`
	Activo            bool   `json:"Activo"`
}

type ContratoEstado struct {
	Id                   int                 `json:"Id"`
	EstadoCumplidoId     EstadoCumplidoId    `json:"EstadoCumplidoId"`
	CumplidoProveedorId  CumplidoProveedorId `json:"CumplidoProveedorId"`
	DocumentoResponsable int                 `json:"DocumentoResponsable"`
	CargoReponsable      string              `json:"CargoReponsable"`
	FechaCreacion        string              `json:"FechaCreacion"`
	FechaModificacion    string              `json:"FechaModificacion"`
	Activo               bool                `json:"Activo"`
}
