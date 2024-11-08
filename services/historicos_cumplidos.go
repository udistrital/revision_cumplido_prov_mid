package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/helpers"
	"github.com/udistrital/revision_cumplidos_proveedores_mid/models"
)

func ObtberHistoricoEstado(cumplido_proveedor_id string) (historicos []models.HistoricoCumplido, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var peticion_info_persona []models.InformacionPersonaNatural
	var respuesta_histroricos []models.CambioEstadoCumplido
	var peticion_historicos map[string]interface{}

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.Id:"+cumplido_proveedor_id+"&sortby=FechaCreacion&order=desc&limit=-1", &peticion_historicos); err == nil && response == 200 {

		data := peticion_historicos["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {

			return historicos, nil
		}
		helpers.LimpiezaRespuestaRefactor(peticion_historicos, &respuesta_histroricos)

	} else {
		outputError = fmt.Errorf("Error al consultar el historico del cumplido")
		return nil, outputError
	}

	urlRequet := "/informacion_persona_natural?fields=PrimerNombre,SegundoNombre,PrimerApellido,SegundoApellido&limit=0&query=Id:" + strconv.Itoa(respuesta_histroricos[0].DocumentoResponsable)

	for _, historico_estado := range respuesta_histroricos {

		if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlcrudAgora")+urlRequet, &peticion_info_persona); response == 200 && err == nil && peticion_info_persona != nil {

			historico := models.HistoricoCumplido{
				NombreResponsable: peticion_info_persona[0].PrimerNombre + " " + peticion_info_persona[0].SegundoNombre + " " + peticion_info_persona[0].PrimerApellido + " " + peticion_info_persona[0].SegundoApellido,
				Estado:            historico_estado.EstadoCumplidoId.Nombre,
				Fecha:             historico_estado.FechaCreacion,
				CargoResponsable:  historico_estado.CargoResponsable,
			}
			historicos = append(historicos, historico)
		} else {
			outputError = fmt.Errorf("Error al consultar la informacion de la persona")
			return nil, outputError
		}

	}

	return historicos, nil
}

func ObtenerHistoricoCumplidosFiltro(anios []int, meses []int, vigencias []int, proveedores []int, estados []string, dependencias []string, contratos []string, tipos_contratos []int) (cumplidos_filtrados []models.CumplidosFiltrados, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	if len(dependencias) == 0 {
		outputError = fmt.Errorf("Debe seleccionar al menos una dependencia")
		return cumplidos_filtrados, outputError
	}

	var contratos_dependencias models.ContratoDependencia
	for _, dependencia := range dependencias {
		contratos, err := helpers.ObtenerContratosDependencia(dependencia)
		if err != nil {
			outputError = fmt.Errorf("Error al obtener contratos de la dependencia %v", dependencia)
			return cumplidos_filtrados, outputError
		}
		contratos_dependencias.Contratos.Contrato = append(contratos_dependencias.Contratos.Contrato, contratos.Contratos.Contrato...)
	}

	var vigencias_string []string
	if len(vigencias) > 0 {
		for _, vigencia := range vigencias {
			vigencias_string = append(vigencias_string, strconv.Itoa(vigencia))
		}
	}

	cumplidos_filtro, err := ObtenerCambiosCumplidosFiltro(contratos, vigencias_string, estados)
	if err != nil {
		outputError = fmt.Errorf("Error al obtener los cumplidos filtrados")
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
		outputError = fmt.Errorf("No hay ningun Cumplido Proveedor que coincida con los filtros ingresados")
		return cumplidos_filtrados, outputError
	}

	// En caso de que no se aplique filtro ni de meses, años, tipos de contrato o proveedores se retornan las coincidencias ya obtenidas
	if len(meses) == 0 && len(anios) == 0 && len(proveedores) == 0 && len(tipos_contratos) == 0 {
		for _, cumplido := range primer_filtro {
			informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
			if err != nil {
				logs.Error(err)
				continue
			}
			informacion_pago := ObtenerPeriodoInformacionPago(cumplido.CumplidoProveedorId.Id)
			cumplido_filtrado := models.CumplidosFiltrados{
				IdCumplido:      cumplido.CumplidoProveedorId.Id,
				NumeroContrato:  cumplido.CumplidoProveedorId.NumeroContrato,
				Vigencia:        strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato),
				Rp:              informacion_contrato[0].NumeroRp,
				NombreProveedor: informacion_contrato[0].NombreProveedor,
				Dependencia:     informacion_contrato[0].NombreDependencia,
				Estado:          cumplido.EstadoCumplidoId.Nombre,
				TipoContrato:    informacion_contrato[0].TipoContrato,
				InformacionPago: informacion_pago,
			}
			cumplidos_filtrados = append(cumplidos_filtrados, cumplido_filtrado)
		}
	}

	// Aplicar filtros de meses, anios, tipos de contrato y nombres proveedores
	for _, cumplido := range primer_filtro {

		informacion_contrato, err := helpers.ObtenerInformacionContratoProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
		if err != nil {
			logs.Error(err)
			continue
		}
		contrato_general, err := helpers.ObtenerContratoGeneralProveedor(cumplido.CumplidoProveedorId.NumeroContrato, strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato))
		if err != nil {
			logs.Error(err)
			continue
		}

		fecha_inicio, fecha_fin := ObtenerInformacionPagoProveedor(cumplido.CumplidoProveedorId.Id)
		if (fecha_inicio == "" || fecha_fin == "") && (len(meses) > 0 || len(anios) > 0) {
			continue
		}
		layout := "2006-01-02"
		FechaInicio, _ := time.Parse(layout, fecha_inicio)
		FechaFin, _ := time.Parse(layout, fecha_fin)

		// Usar MesesEntreFechas para obtener todos los meses entre las fechas de inicio y fin
		periodoMeses := MesesEntreFechas(FechaInicio, FechaFin)

		// Validar filtros de meses y años
		cumplimiento_mes := len(meses) == 0
		cumplimiento_anio := len(anios) == 0

		cumplimiento_anio_mes := false
		if !cumplimiento_mes && !cumplimiento_anio {
			for _, periodo := range periodoMeses {
				if contieneInt(meses, periodo.Mes) {
					if contieneInt(anios, periodo.Anio) {
						cumplimiento_anio_mes = true
						break
					}
				}
			}
		}

		if !cumplimiento_mes && cumplimiento_anio {
			for _, periodo := range periodoMeses {
				if contieneInt(meses, periodo.Mes) {
					cumplimiento_anio_mes = true
					break
				}
			}
		}

		if cumplimiento_mes && !cumplimiento_anio {
			for _, periodo := range periodoMeses {
				if contieneInt(anios, periodo.Anio) {
					cumplimiento_anio_mes = true
					break
				}
			}
		}

		if cumplimiento_mes && cumplimiento_anio {
			cumplimiento_anio_mes = true
		}

		cumplimiento_proveedor := len(proveedores) == 0 || contieneInt(proveedores, contrato_general.Contratista)
		cumplimiento_tipo_contrato := len(tipos_contratos) == 0 || contieneInt(tipos_contratos, contrato_general.TipoContrato.Id)

		if cumplimiento_anio_mes && cumplimiento_proveedor && cumplimiento_tipo_contrato {
			informacion_pago := ObtenerPeriodoInformacionPago(cumplido.CumplidoProveedorId.Id)
			cumplido_filtrado := models.CumplidosFiltrados{
				NumeroContrato:  cumplido.CumplidoProveedorId.NumeroContrato,
				Vigencia:        strconv.Itoa(cumplido.CumplidoProveedorId.VigenciaContrato),
				Rp:              informacion_contrato[0].NumeroRp,
				NombreProveedor: informacion_contrato[0].NombreProveedor,
				Dependencia:     informacion_contrato[0].NombreDependencia,
				Estado:          cumplido.EstadoCumplidoId.Nombre,
				TipoContrato:    informacion_contrato[0].TipoContrato,
				IdCumplido:      cumplido.CumplidoProveedorId.Id,
				InformacionPago: informacion_pago,
			}
			cumplidos_filtrados = append(cumplidos_filtrados, cumplido_filtrado)
		}

	}
	return cumplidos_filtrados, nil
}

func MesesEntreFechas(fechaInicio, fechaFin time.Time) []models.MesAnio {
	// Ajuste las fechas al primer día del mes para incluir todos los meses completos
	fechaInicio = time.Date(fechaInicio.Year(), fechaInicio.Month(), 1, 0, 0, 0, 0, fechaInicio.Location())
	fechaFin = time.Date(fechaFin.Year(), fechaFin.Month(), 1, 0, 0, 0, 0, fechaFin.Location())

	var meses []models.MesAnio

	// Recorremos los meses y años entre las fechas dadas
	for fecha := fechaInicio; !fecha.After(fechaFin); fecha = fecha.AddDate(0, 1, 0) {
		meses = append(meses, models.MesAnio{
			Mes:  int(fecha.Month()),
			Anio: fecha.Year(),
		})
	}

	return meses
}

func ObtenerInformacionPagoProveedor(cumplido_proveedor_id int) (fecha_inicio, fecha_fin string) {

	var respuesta_peticion map[string]interface{}
	var informacion_pago_proveedor []models.InformacionPago
	//fmt.Println("URL", beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/informacion_pago/?query=Activo:true,CumplidoProveedorId.Id:"+strconv.Itoa(cumplido_proveedor_id))
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/informacion_pago/?query=Activo:true,CumplidoProveedorId.Id:"+strconv.Itoa(cumplido_proveedor_id), &respuesta_peticion); err == nil && response == 200 {
		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			return "", ""
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &informacion_pago_proveedor)
		fecha_inicio := informacion_pago_proveedor[0].FechaInicial.Format("2006-01-02")
		fecha_fin := informacion_pago_proveedor[0].FechaFinal.Format("2006-01-02")

		return fecha_inicio, fecha_fin

	} else {
		return "", ""
	}
}

/**
func contiene(lista []string, elemento string) bool {
	for _, v := range lista {
		if strings.ToLower(v) == strings.ToLower(elemento) {
			return true
		}
	}
	return false
}
**/

func contieneInt(lista []int, elemento int) bool {
	for _, v := range lista {
		if v == elemento {
			return true
		}
	}
	return false
}

func ObtenerCambiosCumplidosFiltro(contratos []string, vigencias []string, estados []string) (cambios_estados_cumplidos []models.CambioEstadoCumplido, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerCambiosCumplidosFiltro", "err": err, "status": "404"}
			panic(outputError)
		}
	}()

	var respuesta_peticion map[string]interface{}

	//Se contruye dinamicamente el query //	estadoFiltrado := buildQuery([]string{"Aprobado ordenador"}, "EstadoCumplidoId.Nombre")
	query := strings.TrimSuffix(("?query=" + buildQuery(contratos, "CumplidoProveedorId.NumeroContrato") + buildQuery(vigencias, "CumplidoProveedorId.VigenciaContrato") + buildQuery(estados, "EstadoCumplidoId.CodigoAbreviacion")), ",")
	order := "&order=desc"
	sortby := ",Activo:true&sortby=FechaCreacion,CumplidoProveedorId__Id"
	limit := "&limit=0"

	fmt.Println("URL Filtros", beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/"+query+sortby+order+limit)
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
