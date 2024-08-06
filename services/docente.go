package services

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetDocenteYVinculacionesPorDocumento(documento string) requestresponse.APIResponse {
	urlTercero := beego.AppConfig.String("TercerosService") + "datos_identificacion?query=numero:" + documento + "&fields=Id,TerceroId"
	var resTercero []map[string]interface{}
	if err := request.GetJson(urlTercero, &resTercero); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de terceros")
	}

	if len(resTercero) == 0 || resTercero[0]["Id"] == nil {
		return requestresponse.APIResponseDTO(true, 200, nil, "No hay docente con el documento dado")
	}

	docenteId := fmt.Sprintf("%v", resTercero[0]["TerceroId"].(map[string]interface{})["Id"])
	vinculaciones, err := helpers.GetVinculacionesDeDocente(docenteId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	if len(vinculaciones) == 0 {
		return requestresponse.APIResponseDTO(true, 200, nil, "El docente no tiene vinculaciones")
	}

	docente := map[string]interface{}{
		"Id":     resTercero[0]["TerceroId"].(map[string]interface{})["Id"],
		"Nombre": strings.Title(strings.ToLower(resTercero[0]["TerceroId"].(map[string]interface{})["NombreCompleto"].(string))),
	}
	return requestresponse.APIResponseDTO(true, 200, map[string]interface{}{"Docente": docente, "Vinculaciones": vinculaciones}, "")
}

func GetPreasignacionesSegunDocenteYPeriodo(docenteId, periodoId string) requestresponse.APIResponse {
	urlPreasignaciones := beego.AppConfig.String("PlanDocenteService") +
		"pre_asignacion?query=docente_id:" + docenteId + ",periodo_id:" + periodoId + ",aprobacion_docente:true,aprobacion_proyecto:true,activo:true&limit=0"
	var resPreAsignaciones map[string]interface{}
	if err := request.GetJson(urlPreasignaciones, &resPreAsignaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de plan docente")
	}

	return requestresponse.APIResponseDTO(true, 200, resPreAsignaciones["Data"], "")
}
