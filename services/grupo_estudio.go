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
		return requestresponse.APIResponseDTO(false, 404, nil, "Error al listar horarios", err.Error())
	}

	for _, grupoData := range gruposEstudioResp["Data"].([]interface{}) {
		grupo := grupoData.(map[string]interface{})
		espacios := make([]map[string]interface{}, 0)

		for _, espacioId := range grupo["EspaciosAcademicos"].([]interface{}) {
			if espacio, errEspacio := helpers.ObtenerEspacioAcademicoSegunId(espacioId.(string)); errEspacio == nil {
				espacios = append(espacios, espacio)
			} else {
				return requestresponse.APIResponseDTO(false, 404, nil, "Error al obtener espacio académico")
			}
		}
		grupo["EspaciosAcademicos"] = espacios
		grupo["Nombre"] = fmt.Sprintf("%s %s", grupo["CodigoProyecto"], grupo["IndicadorGrupo"])
	}
	return requestresponse.APIResponseDTO(true, 200, gruposEstudioResp["Data"], nil)
}

func DeleteGrupoEstudio(grupoEstudioId string) requestresponse.APIResponse {
	//eliminar grupo estudio
	var grupoEstudio map[string]interface{}
	urlGrupoEstudio := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	if err := request.GetJson(urlGrupoEstudio, &grupoEstudio); err != nil {
		return requestresponse.APIResponseDTO(false, 404, nil, "Error en el servicio horario"+err.Error())
	}

	grupoEstudioData := grupoEstudio["Data"].(map[string]interface{})
	grupoEstudioData["Activo"] = false

	urlGrupoEstudioPost := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioData["_id"].(string)
	var colocacionPost map[string]interface{}
	if err := request.SendJson(urlGrupoEstudioPost, "PUT", &colocacionPost, grupoEstudioData); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de horario", err.Error())
	}

	//elminar las colocaciones del grupo de estudio
	urlColocaciones := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico?query=Activo:true,GrupoEstudioId:" + grupoEstudioId + "&limit=0"
	var colocaciones map[string]interface{}
	if err := request.GetJson(urlColocaciones, &colocaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 404, nil, "Error en el servicio horario"+err.Error())
	}

	if len(colocaciones["Data"].([]interface{})) > 0 {
		for _, colocacionData := range colocaciones["Data"].([]interface{}) {
			colocacion := colocacionData.(map[string]interface{})

			urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacion["_id"].(string)
			if err := request.GetJson(urlColocacion, &colocacion); err != nil {
				return requestresponse.APIResponseDTO(false, 404, nil, "Error en el servicio horario"+err.Error())
			}

			colocacionData := colocacion["Data"].(map[string]interface{})
			colocacionData["Activo"] = false

			urlColocacionPost := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionData["_id"].(string)
			var colocacionPost map[string]interface{}
			if err := request.SendJson(urlColocacionPost, "PUT", &colocacionPost, colocacionData); err != nil {
				return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de horario", err.Error())
			}
		}
	}
	return requestresponse.APIResponseDTO(true, 200, nil, "delete success")
}
