package models

import "time"

type ContratoProveedor struct {
	TipoContrato       string
	NumeroContrato     string
	Vigencia           int
	Rp                 string
	NombreProveedor    string
	Dependencia        string
	Cdp                string
	NombreOrdenador    string
	DocumentoOrdenador string
	VigenciaRP         string
	FechaCreacion      time.Time
}

type FormaPago struct {
	Id                int    `json:"Id"`
	Descripcion       string `json:"Descripcion"`
	CodigoContraloria string `json:"CodigoContraloria"`
	EstadoRegistro    bool   `json:"EstadoRegistro"`
	FechaRegistro     string `json:"FechaRegistro"`
}

//Contrato disponibilidad
type ContratoDisponibilidad struct {
	Id             int       `json:"Id"`
	NumeroCdp      int       `json:"NumeroCdp"`
	NumeroContrato string    `json:"NumeroContrato"`
	Vigencia       int       `json:"Vigencia"`
	Estado         bool      `json:"Estado"`
	FechaRegistro  time.Time `json:"FechaRegistro"`
	VigenciaCdp    int       `json:"VigenciaCdp"`
}

///////ContratoCompleto
type Cargo struct {
	Id    int    `json:"Id"`
	Cargo string `json:"Cargo"`
}

type Supervisor struct {
	Id                    int    `json:"Id"`
	Nombre                string `json:"Nombre"`
	Documento             int    `json:"Documento"`
	Cargo                 string `json:"Cargo"`
	SedeSupervisor        string `json:"SedeSupervisor"`
	DependenciaSupervisor string `json:"DependenciaSupervisor"`
	Tipo                  int    `json:"Tipo"`
	Estado                bool   `json:"Estado"`
	DigitoVerificacion    int    `json:"DigitoVerificacion"`
	FechaInicio           string `json:"FechaInicio"`
	FechaFin              string `json:"FechaFin"`
	CargoId               Cargo  `json:"CargoId"`
}

type UnidadEjecucion struct {
	Id                int    `json:"Id"`
	Descripcion       string `json:"Descripcion"`
	CodigoContraloria string `json:"CodigoContraloria"`
	EstadoRegistro    bool   `json:"EstadoRegistro"`
	FechaRegistro     string `json:"FechaRegistro"`
}
