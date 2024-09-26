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

func ObtenerCumplidosPendientesOrdenador(documento_ordenador string) (cambios_estado []models.CambioEstadoCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}
	//fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=DocumentoResponsable:" + documento_ordenador + ",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=DocumentoResponsable:"+documento_ordenador+",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true", &respuesta_peticion); (err == nil) && (response == 200) {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			outputError = fmt.Errorf("No hay cumplidos pendientes de revision por el ordenador")
			return nil, outputError
		} else {
			helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_estado)
			return cambios_estado, nil
		}
	} else {
		outputError = fmt.Errorf("Error al consultar los cumplidos pendientes de revision por el ordenador")
		return nil, outputError
	}
}

func ObtenerSolicitudesCumplidos(documento_ordenador string) (cumplidosInfo []models.SolicituRevisionCumplidoProveedor, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
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
				outputError = fmt.Errorf("Error al consultar los cumplidos pendientes de revision por el ordenador")
				return nil, outputError
			}
		}
	} else {
		outputError = fmt.Errorf("No hay cumplidos pendientes de aprobacion para el ordenador " + documento_ordenador)
		return nil, outputError
	}

	return cumplidosInfo, nil
}

func ListaCumplidosReversibles(documento_ordenador string) (soliciudes_revertibles []models.SolicituRevisionCumplidoProveedor, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	fechaActual := time.Now()
	fechaMenosQuinceDias := fechaActual.AddDate(0, 0, -15)
	fechaFormateada := fechaMenosQuinceDias.Format("01/02/2006")
	var listaDocumentos []string
	var ordenadores []models.ArgoOrdenadorContrato

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/ordenadores/?query=Documento:"+documento_ordenador, &ordenadores); err == nil && response == 200 {
		if len(ordenadores) == 0 {
			outputError = fmt.Errorf("No se encontro el ordenador con el documento " + documento_ordenador)
			return nil, outputError
		}
	}

	//Verificar si hace 15 días el ordenador estaba activo
	if fechaMenosQuinceDias.After(ordenadores[0].FechaInicio) && fechaMenosQuinceDias.Before(ordenadores[0].FechaFin) {
		listaDocumentos = append(listaDocumentos, strconv.Itoa(ordenadores[0].Documento))
	} else {
		// Si el oredenador no estuvo activo hace 15 días buscar el ordenador que si lo estaba
		var ordenadores_contrato []models.ArgoOrdenadorContrato
		if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+"/ordenadores/?query=RolId:"+strconv.Itoa(ordenadores[0].RolId)+"FechaFin__gte:"+fechaFormateada+"&limit=-1&sortby=FechaFin&order=desc", &ordenadores_contrato); err == nil && response == 200 {
			if len(ordenadores_contrato) > 0 {
				for _, ordenador_contrato := range ordenadores_contrato {
					listaDocumentos = append(listaDocumentos, strconv.Itoa(ordenador_contrato.Documento))
				}
			} else {
				outputError = fmt.Errorf("No se encontro ningun ordenador activo")
				return nil, outputError
			}
		} else {
			outputError = fmt.Errorf("Error al consultar los ordenadores anteriores")
			return nil, outputError
		}
	}

	documentos := strings.Join(listaDocumentos, ",")
	var respuesta_peticion map[string]interface{}
	var cumplidos []models.CambioEstadoCumplido
	//fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=DocumentoResponsable:" + documento_ordenador + ",EstadoCumplidoId.CodigoAbreviacion:AO,Activo:true,FechaModificacion__gte:" + fechaFormateada)
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=DocumentoResponsable:"+documentos+",EstadoCumplidoId.CodigoAbreviacion:AO,Activo:true,FechaModificacion__gte:"+fechaFormateada, &respuesta_peticion); err == nil && response == 200 {
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
						outputError = fmt.Errorf("Error al consultar los cumplidos pendientes de revision por el ordenador")
						return nil, outputError
					}
				}
			} else {
				outputError = fmt.Errorf("No hay cumplidos que se puedan revertir")
				return nil, outputError
			}
		} else {
			outputError = fmt.Errorf("El ordenador " + documento_ordenador + " no tiene cumplidos que se puedan revertir")
			return nil, outputError
		}

	}

	return soliciudes_revertibles, nil
}

func GenerarAutorizacionGiro(id_solicitud_pago string) (autorizacion_pago models.DocumentoAutorizacionPago, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta_cambioEstado map[string]interface{}
	var cambio_estado []models.CambioEstadoCumplido
	//fmt.Println(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores") + "/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:" + id_solicitud_pago + ",EstadoCumplidoId.CodigoAbreviacion:PRO,Activo:true")
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
							//fmt.Println("URL soportes: ", beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/?query=CumplidoProveedorId.id:"+id_solicitud_pago+",Activo:true")
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
										//fmt.Println("URL documentos: ", beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?query=Id.in:"+documentos_busqueda+",Activo:true&limit=0")
										if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+"/documento/?query=Id.in:"+documentos_busqueda+",Activo:true&limit=0", &documentos); err == nil && response == 200 {
											if len(documentos) == 0 {
												outputError = fmt.Errorf("No se encontraron documentos cargados")
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
											if len(informacion_pago) == 0 {
												outputError = fmt.Errorf("No se encontró información de pago")
												return autorizacion_pago, outputError
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

											autorizacion := helpers.GenerarPdfAutorizacionGiro(datos_documento)
											if autorizacion != "" {
												nombre := "AutorizacionPago_" + strings.Join(strings.Fields(proveedor.NomProveedor), "") + "_" + cambio_estado[0].CumplidoProveedorId.NumeroContrato + "_" + strconv.Itoa(cambio_estado[0].CumplidoProveedorId.VigenciaContrato)
												autorizacion_pago = models.DocumentoAutorizacionPago{
													NombreArchivo:        nombre,
													Archivo:              autorizacion,
													NombreResponsable:    ordenador_contrato.Contratos.Ordenador[0].NombreOrdenador,
													CargoResponsable:     ordenador_contrato.Contratos.Ordenador[0].RolOrdenador,
													DescripcionDocumento: "Autorización de pago para el cumplido " + cambio_estado[0].CumplidoProveedorId.NumeroContrato + " de " + strconv.Itoa(cambio_estado[0].CumplidoProveedorId.VigenciaContrato) + " - " + strconv.Itoa(cambio_estado[0].CumplidoProveedorId.Id),
												}
												return autorizacion_pago, nil
											} else {
												outputError = fmt.Errorf("Error al generar el archivo de autorización de pago")
												return autorizacion_pago, outputError
											}
										} else {
											outputError = fmt.Errorf("Error al consultar los documentos cargados")
											return autorizacion_pago, outputError
										}
									} else {
										outputError = fmt.Errorf("No se encontraron soportes del cumplido")
										return autorizacion_pago, outputError
									}
								} else {
									outputError = fmt.Errorf("No se encontro este soporte")
									return autorizacion_pago, outputError
								}

							} else {
								outputError = fmt.Errorf("Error al obtener los soportes del cumplido")
								return autorizacion_pago, outputError
							}
						} else {
							outputError = fmt.Errorf("Error al consultar el ordenador")
							return autorizacion_pago, outputError
						}
					} else {
						outputError = fmt.Errorf("Error al consultar el proveedor")
						return autorizacion_pago, outputError
					}

				} else {
					outputError = fmt.Errorf("Error al consultar el contrato")
					return autorizacion_pago, outputError
				}
			}

		} else {
			outputError = fmt.Errorf("No se encontró información del cumplido o el cumplido no esta pendiente de aprobación por parte del ordenador")
			return autorizacion_pago, outputError
		}
	} else {
		outputError = fmt.Errorf("Error al obtener el cambio cumplido")
		return autorizacion_pago, outputError
	}

	return autorizacion_pago, nil
}

func ObtenerInformacionProveedor(IdProveedor string) (provedor models.InformacionProveedor, outputError interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
		}
	}()

	var informacion_proveedor []models.InformacionProveedor
	fmt.Println("URL proveedor: ", beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=id:"+IdProveedor)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlcrudAgora")+"/informacion_proveedor/?query=id:"+IdProveedor, &informacion_proveedor); err == nil && response == 200 {
		if len(informacion_proveedor) > 0 {
			return informacion_proveedor[0], nil
		} else {
			outputError = fmt.Errorf("No se encontró información del proveedor")
			return provedor, outputError
		}
	}
	return
}
