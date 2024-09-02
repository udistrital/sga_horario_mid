package services

import (
	// "fmt"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetGruposEstudioSegunHorarioYSemestre(horarioId, semestreId string) requestresponse.APIResponse {
	url := beego.AppConfig.String("HorarioService") + "grupo-estudio"
	query := "HorarioId:" + horarioId + ",SemestreId:" + semestreId + ",Activo:true&limit:0"

	var gruposEstudioResp map[string]interface{}
	if err := request.GetJson(url+"?query="+query, &gruposEstudioResp); err != nil || gruposEstudioResp["Success"] == false {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al listar horarios", err.Error())
	}

	for _, grupoData := range gruposEstudioResp["Data"].([]interface{}) {
		grupo := grupoData.(map[string]interface{})
		espacios := make([]map[string]interface{}, 0)

		for _, espacioId := range grupo["EspaciosAcademicos"].([]interface{}) {
			if espacio, errEspacio := helpers.ObtenerEspacioAcademicoSegunId(espacioId.(string)); errEspacio == nil {
				espacios = append(espacios, espacio)
			} else {
				return requestresponse.APIResponseDTO(false, 500, nil, "Error al obtener espacio académico"+errEspacio.Error())
			}
		}
		grupo["EspaciosAcademicos"] = espacios
		grupo["Nombre"] = fmt.Sprintf("%s %s", grupo["CodigoProyecto"], grupo["IndicadorGrupo"])
	}
	return requestresponse.APIResponseDTO(true, 200, gruposEstudioResp["Data"], nil)
}

func DeleteGrupoEstudio(grupoEstudioId string) requestresponse.APIResponse {
	//eliminar grupo de estudio
	_, err := helpers.DesactivarGrupoEstudio(grupoEstudioId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	//eliminar colocaciones del grupo de estudio
	urlColocaciones := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico?query=Activo:true,GrupoEstudioId:" + grupoEstudioId + "&limit=0"
	var colocaciones map[string]interface{}
	if err := request.GetJson(urlColocaciones, &colocaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 404, nil, "Error en el servicio horario"+err.Error())
	}

	if len(colocaciones["Data"].([]interface{})) > 0 {
		for _, colocacionData := range colocaciones["Data"].([]interface{}) {
			colocacion := colocacionData.(map[string]interface{})
			_, err := helpers.DesactivarColocacion(colocacion["_id"].(string))
			if err != nil {
				return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
			}
		}
	}
	return requestresponse.APIResponseDTO(true, 200, nil, "delete success")
}
