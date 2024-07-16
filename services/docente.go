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
	fmt.Println(err)
	fmt.Println(vinculaciones)
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
