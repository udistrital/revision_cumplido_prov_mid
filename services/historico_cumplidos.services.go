package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerHistoricoCumplidosFiltro(anios []int, meses []int, vigencias []string, nombres_proveedores []string, estados []string, dependencias []string, contratos []string) (cumplidos_filtrados []models.CumplidosFiltrados, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerHistoricoCumplidosFiltro", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var cumplidos_filtrados_final []models.CambioEstadoCumplido

	if len(dependencias) == 0 {
		outputError = map[string]interface{}{"funcion": "ObtenerHistoricoCumplidosFiltro", "err": "Debe seleccionar al menos una dependencia", "status": "404"}
		return cumplidos_filtrados, outputError
	}

	var contratos_dependencias models.ContratoDependencia
	for _, dependencia := range dependencias {
		contratos, err := helpers.ObtenerContratosDependencia(dependencia)
		if err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerHistoricoCumplidosFiltro", "err": err, "status": "404"}
			return cumplidos_filtrados, outputError
		}
		contratos_dependencias.Contratos.Contrato = append(contratos_dependencias.Contratos.Contrato, contratos.Contratos.Contrato...)
	}

	cumplidos_filtro, err := ObtenerCambiosCumplidosFiltro(contratos, vigencias, estados)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "ObtenerHistoricoCumplidosFiltro", "err": err, "status": "404"}
		return cumplidos_filtrados, outputError
	}

	var primer_filtro []models.CambioEstadoCumplido
	for _, contrato_dependencia := range contratos_dependencias.Contratos.Contrato {
		for _, cumplido_filtro := range cumplidos_filtro {
			if contrato_dependencia.NumeroContrato == cumplido_filtro.CumplidoProveedorId.NumeroContrato && contrato_dependencia.Vigencia == strconv.Itoa(cumplido_filtro.CumplidoProveedorId.VigenciaContrato) {
				primer_filtro = append(primer_filtro, cumplido_filtro)
			}
		}
	}

	if len(primer_filtro) == 0 {
		outputError = map[string]interface{}{"funcion": "ObtenerHistoricoCumplidosFiltro", "err": "No hay ningun Cumplido Proveedor que coincida con los filtros ingresados", "status": "404"}
		return cumplidos_filtrados, outputError
	}

	// En caso de que no se aplique filtro ni de meses, años o proveedores se retornan las coincidencias ya obtenidas
	if len(meses) == 0 && len(anios) == 0 && len(nombres_proveedores) == 0 {
		for _, cumplido := range primer_filtro {
			if err != nil {
				logs.Error(err)
				continue
			}
			cumplidos_filtrados_final = append(cumplidos_filtrados_final, cumplido)
		}
	}

	// Aplicar filtro por meses
	var filtro_meses []models.CambioEstadoCumplido
	if len(meses) > 0 {
		for _, mes := range meses {
			for _, cumplido := range primer_filtro {
				if int(cumplido.FechaCreacion.Month()) == mes {
					filtro_meses = append(filtro_meses, cumplido)
				}
			}
		}
	}

	// Aplicar filtro por años
	var filtro_anios []models.CambioEstadoCumplido
	if len(anios) > 0 {
		if len(filtro_meses) > 0 {
			for _, anio := range anios {
				for _, cumplido := range filtro_meses {
					if strconv.Itoa(cumplido.FechaCreacion.Year()) == strconv.Itoa(anio) {
						filtro_anios = append(filtro_anios, cumplido)
					}
				}
			}
		} else {
			for _, anio := range anios {
				for _, cumplido := range primer_filtro {
					if strconv.Itoa(cumplido.FechaCreacion.Year()) == strconv.Itoa(anio) {
						filtro_anios = append(filtro_anios, cumplido)
					}
				}
			}
		}
	}

	// Aplicar filtro por proveedores
	if len(nombres_proveedores) > 0 {
		if len(filtro_anios) > 0 {
			for _, nombre_proveedor := range nombres_proveedores {
				for _, cumplido := range filtro_anios {
					informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
					if err != nil {
						logs.Error(err)
						continue
					}
					if informacion_contrato[0].NombreProveedor == nombre_proveedor {
						cumplidos_filtrados_final = append(cumplidos_filtrados_final, cumplido)
					}
				}
			}
		} else if len(filtro_meses) > 0 {
			for _, nombre_proveedor := range nombres_proveedores {
				for _, cumplido := range filtro_meses {
					informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
					if err != nil {
						logs.Error(err)
						continue
					}
					if informacion_contrato[0].NombreProveedor == nombre_proveedor {
						cumplidos_filtrados_final = append(cumplidos_filtrados_final, cumplido)
					}
				}
			}
		} else {
			for _, nombre_proveedor := range nombres_proveedores {
				for _, cumplido := range primer_filtro {
					informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
					if err != nil {
						logs.Error(err)
						continue
					}
					if informacion_contrato[0].NombreProveedor == nombre_proveedor {
						cumplidos_filtrados_final = append(cumplidos_filtrados_final, cumplido)
					}
				}
			}
		}
	} else {
		if len(filtro_anios) > 0 {
			cumplidos_filtrados_final = filtro_anios
		} else if len(filtro_meses) > 0 {
			cumplidos_filtrados_final = filtro_meses
		} else {
			cumplidos_filtrados_final = primer_filtro
		}
	}

	for _, cumplido := range cumplidos_filtrados_final {
		informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
		if err != nil {
			logs.Error(err)
			continue
		}
		cumplidos_filtrados = append(cumplidos_filtrados, models.CumplidosFiltrados{
			NumeroContrato:  cumplido.CumplidoProveedorId.NumeroContrato,
			Vigencia:        strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato),
			Rp:              informacion_contrato[0].NumeroRp,
			Mes:             int(cumplido.FechaCreacion.Month()),
			FechaAprobacion: cumplido.FechaCreacion.Format("2006/01/02"),
			NombreProveedor: informacion_contrato[0].NombreProveedor,
			Dependencia:     informacion_contrato[0].NombreDependencia,
		})
	}

	return cumplidos_filtrados, nil
}

func ObtenerCambiosCumplidosFiltro(contratos []string, vigencias []string, estados []string) (cambios_estados_cumplidos []models.CambioEstadoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerCambiosCumplidosFiltro", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}

	//Se contruye dinamicamente el query

	query := strings.TrimSuffix(("?query=" + buildQuery(contratos, "CumplidoProveedorId.NumeroContrato") + buildQuery(vigencias, "CumplidoProveedorId.Vigencia") + buildQuery(estados, "EstadoCumplidoId.CodigoAbreviacion")), ",")
	order := "&order=desc"
	sortby := "&sortby=FechaCreacion"
	limit := "&limit=0"

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/"+query+sortby+order+limit, &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_estados_cumplidos)
			return cambios_estados_cumplidos, nil
		} else {
			outputError = map[string]interface{}{"funcion": "ObtenerCambiosCumplidosFiltro", "err": "No hay ningun Cumplido Proveedor que coincida con los filtros ingresados", "status": "404"}
			return cambios_estados_cumplidos, outputError
		}

	}
	return cambios_estados_cumplidos, outputError
}

func buildQuery(slices []string, columna string) string {

	query := ""

	if len(slices) == 1 {
		query += fmt.Sprintf("%s.in:%v,", columna, slices[0])
	}
	if len(slices) > 1 {
		for i, dato := range slices {
			if i == 0 {
				query += fmt.Sprintf("%s.in:%v|", columna, dato)
			} else if i < len(slices)-1 {
				query += fmt.Sprintf("%s|", dato)
			} else {
				query += fmt.Sprintf("%s,", dato)
			}
		}
		return query
	}
	return query
}
