package services

import (
	// "fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetGruposEstudio(proyectoAcademicoId, planEstudiosId, semestreId string) requestresponse.APIResponse {
	url := beego.AppConfig.String("HorariosService") + "grupo-estudio"
	query := "Activo:true,ProyectoAcademicoId:" + proyectoAcademicoId + ",PlanEstudiosId:" + planEstudiosId + ",SemestreId:" + semestreId + "&limit:0"

	var gruposEstudioResp map[string]interface{}
	if err := request.GetJson(url+"?query="+query, &gruposEstudioResp); err != nil || gruposEstudioResp["Success"] == false {
		return requestresponse.APIResponseDTO(false, 404, nil, "Error al listar horarios")
	}

	for _, grupoData := range gruposEstudioResp["Data"].([]interface{}) {
		grupo := grupoData.(map[string]interface{})
		espaciosCompletos := make([]map[string]interface{}, 0)

		for _, espacioId := range grupo["EspaciosAcademicos"].([]interface{}) {
			if espacio, errEspacio := helpers.ObtenerEspacioAcademicoSegunId(espacioId.(string)); errEspacio != nil || espacio["Success"] == false {
				return requestresponse.APIResponseDTO(false, 404, nil, "Error al obtener espacio académico")
			} else if espacioData := espacio["Data"].([]interface{}); len(espacioData) > 0 {
				espaciosCompletos = append(espaciosCompletos, espacioData[0].(map[string]interface{}))
			}
		}
		grupo["EspaciosAcademicosCompletos"] = espaciosCompletos
	}
	return requestresponse.APIResponseDTO(true, 200, gruposEstudioResp["Data"], nil)
}
