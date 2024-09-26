package services

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetEspaciosFisicosOcupadosSegunPeriodo(espacioFisicoId, periodoId string) requestresponse.APIResponse {
	urlColocaciones := beego.AppConfig.String("HorarioService") +
		"colocacion-espacio-academico?query=PeriodoId:" + periodoId + ",EspacioFisicoId:" + espacioFisicoId + ",Activo:true&limit=0"

	var colocaciones map[string]interface{}
	if err := request.GetJson(urlColocaciones, &colocaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "error en el servicio de horarios"+err.Error())
	}

	var ocupados []map[string]interface{}

	for _, colocacion := range colocaciones["Data"].([]interface{}) {
		colocacionData, _ := colocacion.(map[string]interface{})

		var colocacionEspacio map[string]interface{}
		_ = json.Unmarshal([]byte(colocacionData["ColocacionEspacioAcademico"].(string)), &colocacionEspacio)

		ocupados = append(ocupados, map[string]interface{}{
			"_id":           colocacionData["_id"],
			"horas":         int(colocacionEspacio["horas"].(float64)),
			"finalPosition": colocacionEspacio["finalPosition"],
			"horaFormato":   colocacionEspacio["horaFormato"],
		})
	}
	return requestresponse.APIResponseDTO(true, 200, ocupados, "")
}
