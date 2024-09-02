package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtenerCumplidosPendientesOrdenador(documento_ordenador string) (cambios_estado []models.CambioEstadoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=DocumentoResponsable:" + documento_ordenador + ",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=DocumentoResponsable:"+documento_ordenador+",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true", &respuesta_peticion); (err == nil) && (response == 200) {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "message": "No hay cumplidos pendientes de revision por el ordenador", "status": "404"}
			return nil, outputError
		} else {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_estado)
			return cambios_estado, nil
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "message": "Error al consultar los cumplidos pendeinetes de revision por el ordenador", "err": err, "status": "404"}
		return nil, outputError
	}
}

func ObtenerSolicitudesCumplidos(documento_ordenador string) (cumplidosInfo []models.SolicituRevisionCumplidoProveedor, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ObtenerSolicitudesCumplidos", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	cumplidos, err := ObtenerCumplidosPendientesOrdenador(documento_ordenador)
	if err == nil {
		for _, cumplido := range cumplidos {
			informacion_contrato_proveedor, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
			if err == nil {
				vigencia, _ := strconv.Atoi(informacion_contrato_proveedor[0].Vigencia)
				solicitudes_cumplido := models.SolicituRevisionCumplidoProveedor{
					TipoContrato:      informacion_contrato_proveedor[0].TipoContrato,
					NumeroContrato:    informacion_contrato_proveedor[0].NumeroContratoSuscrito,
					VigenciaContrato:  vigencia,
					Dependencia:       informacion_contrato_proveedor[0].NombreDependencia,
					NombreProveedor:   informacion_contrato_proveedor[0].NombreProveedor,
					Cdp:               informacion_contrato_proveedor[0].NumeroCdp,
					Rp:                informacion_contrato_proveedor[0].NumeroRp,
					VigenciaRP:        informacion_contrato_proveedor[0].VigenciaRp,
					CumplidoId:        cumplido.CumplidoProveedorId.Id,
					Activo:            cumplido.Activo,
					FechaCreacion:     cumplido.CumplidoProveedorId.FechaCreacion,
					FechaModificacion: cumplido.CumplidoProveedorId.FechaModificacion,
				}
				cumplidosInfo = append(cumplidosInfo, solicitudes_cumplido)
			} else {
				outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "message": "Error al consultar los cumplidos pendeinetes de revision por el ordenador", "err": err, "status": "404"}
				return nil, outputError
			}
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/ObtenerCumplidosPendientesOrdenador", "message": "No hay cumplidos pendientes de aprobacion para el ordenador " + documento_ordenador, "status": "404"}
		return nil, outputError
	}

	return cumplidosInfo, nil
}

func ListaCumplidosReversibles(documento_ordenador string) (soliciudes_revertibles []models.SolicituRevisionCumplidoProveedor, outputError interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	fechaActual := time.Now()
	fechaMenosQuinceDias := fechaActual.AddDate(0, 0, -15)
	fechaFormateada := fechaMenosQuinceDias.Format("01/02/2006")

	var respuesta_peticion map[string]interface{}
	var cumplidos []models.CambioEstadoCumplido
	//fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=DocumentoResponsable:" + documento_ordenador + ",EstadoCumplidoId.CodigoAbreviacion:AO,Activo:true,FechaModificacion__gte:" + fechaFormateada)
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=DocumentoResponsable:"+documento_ordenador+",EstadoCumplidoId.CodigoAbreviacion:AO,Activo:true,FechaModificacion__gte:"+fechaFormateada, &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cumplidos)
			if len(cumplidos) > 0 {
				for _, cumplido := range cumplidos {
					informacion_contrato_proveedor, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
					if err == nil {
						vigencia, _ := strconv.Atoi(informacion_contrato_proveedor[0].Vigencia)
						solicitudes_cumplido := models.SolicituRevisionCumplidoProveedor{
							TipoContrato:      informacion_contrato_proveedor[0].TipoContrato,
							NumeroContrato:    informacion_contrato_proveedor[0].NumeroContratoSuscrito,
							VigenciaContrato:  vigencia,
							Dependencia:       informacion_contrato_proveedor[0].NombreDependencia,
							NombreProveedor:   informacion_contrato_proveedor[0].NombreProveedor,
							Cdp:               informacion_contrato_proveedor[0].NumeroCdp,
							Rp:                informacion_contrato_proveedor[0].NumeroRp,
							VigenciaRP:        informacion_contrato_proveedor[0].VigenciaRp,
							CumplidoId:        cumplido.CumplidoProveedorId.Id,
							Activo:            cumplido.Activo,
							FechaCreacion:     cumplido.CumplidoProveedorId.FechaCreacion,
							FechaModificacion: cumplido.CumplidoProveedorId.FechaModificacion,
						}
						soliciudes_revertibles = append(soliciudes_revertibles, solicitudes_cumplido)
					} else {
						outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "message": "Error al consultar los cumplidos pendientes de revision por el ordenador", "err": err, "status": "404"}
						return nil, outputError
					}
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "message": "No hay cumplidos que se puedan revertir", "status": "404"}
				return nil, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ListaCumplidosReversibles", "message": "El ordenador " + documento_ordenador + " no tiene cumplidos que se puedan revertir", "status": "404"}
			return nil, outputError
		}

	}

	return soliciudes_revertibles, nil
}

func GenerarAutorizacionPago(id_solicitud_pago string) (autorizacion_pago models.DocumentoAutorizacionPago, outputError interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var respuesta_cambioEstado map[string]interface{}
	var cambio_estado []models.CambioEstadoCumplido
	fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:" + id_solicitud_pago + ",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:"+id_solicitud_pago+",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true", &respuesta_cambioEstado); err == nil && response == 200 {
		data := respuesta_cambioEstado["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) > 0 {
			helpers.LimpiezaRespuestaRefactor(respuesta_cambioEstado, &cambio_estado)
			informacion_contrato_proveedor, err := helpers.ObtenerInformacionContratoProveedor(cambio_estado[0].CumplidoProveedorId.NumeroContrato, strconv.Itoa(cambio_estado[0].CumplidoProveedorId.VigenciaContrato))
			if err == nil {
				contrato_general, err := helpers.ObtenerContratoGeneralProveedor(cambio_estado[0].CumplidoProveedorId.NumeroContrato, strconv.Itoa(cambio_estado[0].CumplidoProveedorId.VigenciaContrato))
				if err == nil {
					proveedor, err := ObtenerInformacionProveedor(strconv.Itoa(contrato_general.Contratista))
					if err == nil {
						ordenador_contrato, err := helpers.ObtenerOrdenadorContrato(cambio_estado[0].CumplidoProveedorId.NumeroContrato, strconv.Itoa(cambio_estado[0].CumplidoProveedorId.VigenciaContrato))
						if err == nil {
							var respuesta_soportes_cumplido map[string]interface{}
							var soportes_cumplido []models.SoporteCumplido
							var id_documentos []string
							if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/?query=CumplidoProveedorId.id:"+id_solicitud_pago+",Activo:true", &respuesta_soportes_cumplido); err == nil && response == 200 {
								data := respuesta_soportes_cumplido["Data"].([]interface{})
								if len(data[0].(map[string]interface{})) > 0 {
									helpers.LimpiezaRespuestaRefactor(respuesta_soportes_cumplido, &soportes_cumplido)
									if len(soportes_cumplido) > 0 {
										for _, soporte := range soportes_cumplido {
											id_documentos = append(id_documentos, strconv.Itoa(soporte.DocumentoId))
										}
										documentos_busqueda := strings.Join(id_documentos, "|")
										var documentos []models.Documento
										if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?query=Id.in:"+documentos_busqueda+",Activo:true&limit=0", &documentos); err == nil && response == 200 {
											if len(documentos) == 0 {
												outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "No se encontraron documentos cargados", "status": "404"}
												return autorizacion_pago, outputError
											}
											var lista_documentos_cargados []string
											for _, documento := range documentos {
												lista_documentos_cargados = append(lista_documentos_cargados, documento.TipoDocumento.CodigoAbreviacion)
											}
											var respuesta_soporte map[string]interface{}
											var informacion_pago []models.InformacionPago
											if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/informacion_pago/?query=CumplidoProveedorId.Id:"+id_solicitud_pago, &respuesta_soporte); err == nil && response == 200 {
												data := respuesta_soporte["Data"].([]interface{})
												if len(data[0].(map[string]interface{})) > 0 {
													helpers.LimpiezaRespuestaRefactor(respuesta_soporte, &informacion_pago)
												}
											}
											valor_pago := int(informacion_pago[0].ValorCumplido)
											datos_documento := models.DatosAutorizacionPago{
												NombreOrdenador:    ordenador_contrato.Contratos.Ordenador[0].NombreOrdenador,
												DocumentoOrdenador: ordenador_contrato.Contratos.Ordenador[0].Documento,
												Rubro:              informacion_contrato_proveedor[0].Rubro,
												NombreProveedor:    proveedor.NomProveedor,
												DocumentoProveedor: proveedor.NumDocumento,
												DocumentosCargados: lista_documentos_cargados,
												ValorPago:          valor_pago,
											}

											autorizacion := helpers.GenerarPdfAutorizacionPago(datos_documento)
											if autorizacion != "" {
												nombre := "AutorizacionPago_" + strings.Join(strings.Fields(proveedor.NomProveedor), "")
												autorizacion_pago = models.DocumentoAutorizacionPago{
													File:    nombre,
													Archivo: autorizacion,
												}
												return autorizacion_pago, nil
											} else {
												outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "Error al generar el archivo de autorización de pago", "status": "404"}
												return autorizacion_pago, outputError
											}
										} else {
											outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "Error al consultar los documentos cargados", "err": err, "status": "404"}
											return autorizacion_pago, outputError
										}
									} else {
										outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "No se encontraron soportes del cumplido", "status": "404"}
										return autorizacion_pago, outputError
									}
								} else {
									outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "No se encontro este soporte", "status": "404"}
									return autorizacion_pago, outputError
								}

							} else {
								outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "Error al obtener los soportes del cumplido", "err": err, "status": "404"}
								return autorizacion_pago, outputError
							}
						} else {
							outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "Error al consultar el ordenador", "err": err, "status": "404"}
							return autorizacion_pago, outputError
						}
					} else {
						outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "Error al consultar el proveedor", "err": err, "status": "404"}
						return autorizacion_pago, outputError
					}

				} else {
					outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "Error al consultar el contrato", "err": err, "status": "404"}
					return autorizacion_pago, outputError
				}
			}

		} else {
			outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "No se encontró información del cumplido o el cumplido no esta pendiente de aprobación por parte del ordenador", "status": "404"}
			return autorizacion_pago, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/GenerarAutorizacionPago", "message": "Error al obtener el cambio cumplido", "err": err, "status": "404"}
		return autorizacion_pago, outputError
	}

	return autorizacion_pago, nil
}

func ObtenerInformacionProveedor(IdProveedor string) (provedor models.InformacionProveedor, outputError interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"Success": false,
				"Status":  400,
				"Message": "Error al consultar las dependencias: " + IdProveedor,
				"Error":   err,
			}
		}
	}()

	var informacion_proveedor []models.InformacionProveedor
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=id:"+IdProveedor, &informacion_proveedor); err == nil && response == 200 {
		if len(informacion_proveedor) > 0 {
			return informacion_proveedor[0], nil
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerInformacionProveedor", "message": "No se encontró información del proveedor", "status": "404"}
			return provedor, outputError
		}
	}
	return
}
