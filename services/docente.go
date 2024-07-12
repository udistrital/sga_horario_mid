package services

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
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
	urlVinculacion := beego.AppConfig.String("TercerosService") + "vinculacion?query=TerceroPrincipalId.Id:" + docenteId + "&fields=Id,TipoVinculacionId"
	var vinculacionesId []map[string]interface{}
	if err := request.GetJson(urlVinculacion, &vinculacionesId); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de terceros")
	}

	if len(vinculacionesId) == 0 || vinculacionesId[0]["Id"] == nil {
		return requestresponse.APIResponseDTO(true, 200, nil, "El docente no tiene vinculaciones")
	}

	urlParametroBase := beego.AppConfig.String("ParametroService") + "parametro/"
	var vinculaciones []map[string]interface{}
	for _, vinculacion := range vinculacionesId {
		var resVinculacion map[string]interface{}
		if err := request.GetJson(urlParametroBase+fmt.Sprintf("%v", vinculacion["TipoVinculacionId"]), &resVinculacion); err == nil {
			res := resVinculacion["Data"].(map[string]interface{})
			vinculaciones = append(vinculaciones, map[string]interface{}{
				"Id":     res["Id"],
				"Nombre": strings.ToUpper(string(res["Nombre"].(string)[0])) + strings.ToLower(res["Nombre"].(string)[1:]),
			})
		}
	}

	docente := map[string]interface{}{
		"Id":     resTercero[0]["TerceroId"].(map[string]interface{})["Id"],
		"Nombre": strings.Title(strings.ToLower(resTercero[0]["TerceroId"].(map[string]interface{})["NombreCompleto"].(string))),
	}
	return requestresponse.APIResponseDTO(true, 200, map[string]interface{}{"Docente": docente, "Vinculaciones": vinculaciones}, "")
}
