package models

import "time"

type InformacionInformeSatisfaccion struct {
	Dependencia            string
	NombreProveedor        string
	DocumentoProveedor     string
	CumplimientoTotal      bool
	TipoContrato           string
	FechaInicio            time.Time
	NumeroContratoSuscrito string
	Cdp                    string
	VigenciaCdp            time.Time
	Rp                     string
	VigenciaRp             time.Time
	CargoSupervisor        string
	ValorTotalContrato     int
	SaldoContrato          int
	FechaFin               time.Time
	Supervisor             string
	DocumentoSupervisor    string
}
