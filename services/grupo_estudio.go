package services

import (
	// "fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/sga_horario_mid/helpers"
)

func GetGruposEstudio(proyectoAcademicoId, planEstudiosId, semestreId string) requestresponse.APIResponse {
	var gruposEstudio map[string]interface{}

	urlGetHorarios := beego.AppConfig.String("HorariosService") + "grupo-estudio?query=Activo:true,ProyectoAcademicoId:" + proyectoAcademicoId +
		",PlanEstudiosId:" + planEstudiosId + ",SemestreId:" + semestreId + "&limit:0"

	errSolicitudHorarios := request.GetJson(urlGetHorarios, &gruposEstudio)

	if errSolicitudHorarios != nil || gruposEstudio["Success"] == false {
		return requestresponse.APIResponseDTO(false, 404, nil, "Error al listar horarios")
	}

	gruposEstudioData := gruposEstudio["Data"].([]interface{})
	for _, grupoData := range gruposEstudioData {
		grupo := grupoData.(map[string]interface{})
		espaciosAcademicos := grupo["EspaciosAcademicos"].([]interface{})
		espaciosCompletos := make([]map[string]interface{}, len(espaciosAcademicos))

		for i, espacioId := range espaciosAcademicos {
			espacio, errEspacio := helpers.obtenerEspacioAcademicoSegunId(espacioId.(string))
			if errEspacio != nil || espacio["Success"] == false {
				return requestresponse.APIResponse{Success: false, Status: 404, Message: "Error al obtener espacio académico", Data: nil}
			}
			espaciosCompletos[i] = espacio
		}
		grupo["EspaciosAcademicosCompletos"] = espaciosCompletos
	}

	return requestresponse.APIResponse{Success: true, Status: 200, Message: nil, Data: gruposEstudio["Data"]}
}
