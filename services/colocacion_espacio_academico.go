package services

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// Obtiene colocaciones según el horario de semestre
func GetColocacionesSegunHorarioSemestre(horarioSemestreId string) requestresponse.APIResponse {
	// Trae las colocaciones de espacio segun el id del semestre horario
	urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico?query=HorarioSemestreId:" + horarioSemestreId + ",Activo:true&limit=0"
	var resColocaciones map[string]interface{}
	if err := request.GetJson(urlColocacion, &resColocaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de terceros")
	}

	data := resColocaciones["Data"].([]interface{})

	for i, colocacion := range data {
		colocacionMap := colocacion.(map[string]interface{})
		//agrega objeto completo para sede edificio y salon
		if err := helpers.GetSedeEdificioSalon(colocacionMap); err != nil {
			return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
		}

		//agrega objeto completo de espacio academico
		if id, ok := colocacionMap["EspacioAcademicoId"].(string); ok {
			if espacioAcademico, err := helpers.ObtenerEspacioAcademicoSegunId(id); err == nil {
				colocacionMap["EspacioAcademico"] = espacioAcademico
			}
		}

		data[i] = colocacionMap
	}

	return requestresponse.APIResponseDTO(true, 200, data, "")
}
