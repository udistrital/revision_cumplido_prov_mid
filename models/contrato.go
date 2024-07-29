package models

import "time"

type Contrato struct {
	TipoContrato       string
	NumeroContrato     string
	Vigencia           int
	Rp                 string
	NombreProveedor    string
	Dependencia        string
	Cdp                string
	NombreOrdenador    string
	DocumentoOrdenador string
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

type TipoContrato struct {
	Id           int    `json:"Id"`
	TipoContrato string `json:"TipoContrato"`
	Estado       bool   `json:"Estado"`
}

type LugarEjecucion struct {
	Id          int    `json:"Id"`
	Direccion   string `json:"Direccion"`
	Sede        string `json:"Sede"`
	Dependencia string `json:"Dependencia"`
	Ciudad      int    `json:"Ciudad"`
}

type UnidadEjecucion struct {
	Id                int    `json:"Id"`
	Descripcion       string `json:"Descripcion"`
	CodigoContraloria string `json:"CodigoContraloria"`
	EstadoRegistro    bool   `json:"EstadoRegistro"`
	FechaRegistro     string `json:"FechaRegistro"`
}

type ContratoSuscrito struct {
	Id                     int            `json:"Id"`
	NumeroContrato         NumeroContrato `json:"NumeroContrato"`
	Vigencia               int            `json:"Vigencia"`
	FechaRegistro          string         `json:"FechaRegistro"`
	Usuario                string         `json:"Usuario"`
	FechaSuscripcion       string         `json:"FechaSuscripcion"`
	NumeroContratoSuscrito string         `json:"NumeroContratoSuscrito"`
}

type NumeroContrato struct {
	Id                           string             `json:"Id"`
	VigenciaContrato             int                `json:"VigenciaContrato"`
	ObjetoContrato               string             `json:"ObjetoContrato"`
	PlazoEjecucion               int                `json:"PlazoEjecucion"`
	FormaPago                    FormaPago          `json:"FormaPago"`
	OrdenadorGasto               int                `json:"OrdenadorGasto"`
	ClausulaRegistroPresupuestal bool               `json:"ClausulaRegistroPresupuestal"`
	SedeSolicitante              string             `json:"SedeSolicitante"`
	DependenciaSolicitante       string             `json:"DependenciaSolicitante"`
	Contratista                  int                `json:"Contratista"`
	ValorContrato                int                `json:"ValorContrato"`
	Justificacion                string             `json:"Justificacion"`
	DescripcionFormaPago         string             `json:"DescripcionFormaPago"`
	Condiciones                  string             `json:"Condiciones"`
	FechaRegistro                string             `json:"FechaRegistro"`
	TipologiaContrato            int                `json:"TipologiaContrato"`
	TipoCompromiso               int                `json:"TipoCompromiso"`
	ModalidadSeleccion           int                `json:"ModalidadSeleccion"`
	Procedimiento                int                `json:"Procedimiento"`
	RegimenContratacion          int                `json:"RegimenContratacion"`
	TipoGasto                    int                `json:"TipoGasto"`
	TemaGastoInversion           int                `json:"TemaGastoInversion"`
	OrigenPresupueso             int                `json:"OrigenPresupueso"`
	OrigenRecursos               int                `json:"OrigenRecursos"`
	TipoMoneda                   int                `json:"TipoMoneda"`
	ValorContratoMe              int                `json:"ValorContratoMe"`
	ValorTasaCambio              int                `json:"ValorTasaCambio"`
	TipoControl                  int                `json:"TipoControl"`
	Observaciones                string             `json:"Observaciones"`
	Supervisor                   Supervisor         `json:"Supervisor"`
	ClaseContratista             int                `json:"ClaseContratista"`
	Convenio                     string             `json:"Convenio"`
	NumeroConstancia             int                `json:"NumeroConstancia"`
	Estado                       bool               `json:"Estado"`
	TipoContrato                 TipoContrato       `json:"TipoContrato"`
	LugarEjecucion               LugarEjecucion     `json:"LugarEjecucion"`
	UnidadEjecucion              UnidadEjecucion    `json:"UnidadEjecucion"`
	UnidadEjecutora              int                `json:"UnidadEjecutora"`
	ContratoSuscrito             []ContratoSuscrito `json:"ContratoSuscrito"`
}

type ContratoProveedor struct {
	Id                           string             `json:"Id"`
	VigenciaContrato             int                `json:"VigenciaContrato"`
	ObjetoContrato               string             `json:"ObjetoContrato"`
	PlazoEjecucion               int                `json:"PlazoEjecucion"`
	FormaPago                    FormaPago          `json:"FormaPago"`
	OrdenadorGasto               int                `json:"OrdenadorGasto"`
	ClausulaRegistroPresupuestal bool               `json:"ClausulaRegistroPresupuestal"`
	SedeSolicitante              string             `json:"SedeSolicitante"`
	DependenciaSolicitante       string             `json:"DependenciaSolicitante"`
	Contratista                  int                `json:"Contratista"`
	ValorContrato                int                `json:"ValorContrato"`
	Justificacion                string             `json:"Justificacion"`
	DescripcionFormaPago         string             `json:"DescripcionFormaPago"`
	Condiciones                  string             `json:"Condiciones"`
	FechaRegistro                string             `json:"FechaRegistro"`
	TipologiaContrato            int                `json:"TipologiaContrato"`
	TipoCompromiso               int                `json:"TipoCompromiso"`
	ModalidadSeleccion           int                `json:"ModalidadSeleccion"`
	Procedimiento                int                `json:"Procedimiento"`
	RegimenContratacion          int                `json:"RegimenContratacion"`
	TipoGasto                    int                `json:"TipoGasto"`
	TemaGastoInversion           int                `json:"TemaGastoInversion"`
	OrigenPresupueso             int                `json:"OrigenPresupueso"`
	OrigenRecursos               int                `json:"OrigenRecursos"`
	TipoMoneda                   int                `json:"TipoMoneda"`
	ValorContratoMe              int                `json:"ValorContratoMe"`
	ValorTasaCambio              int                `json:"ValorTasaCambio"`
	TipoControl                  int                `json:"TipoControl"`
	Observaciones                string             `json:"Observaciones"`
	Supervisor                   Supervisor         `json:"Supervisor"`
	ClaseContratista             int                `json:"ClaseContratista"`
	Convenio                     string             `json:"Convenio"`
	NumeroConstancia             int                `json:"NumeroConstancia"`
	Estado                       bool               `json:"Estado"`
	TipoContrato                 TipoContrato       `json:"TipoContrato"`
	LugarEjecucion               LugarEjecucion     `json:"LugarEjecucion"`
	UnidadEjecucion              UnidadEjecucion    `json:"UnidadEjecucion"`
	UnidadEjecutora              int                `json:"UnidadEjecutora"`
	ContratoSuscrito             []ContratoSuscrito `json:"ContratoSuscrito"`
}
