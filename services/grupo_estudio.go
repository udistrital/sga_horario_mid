package services

import (
	// "fmt"

	"encoding/json"
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
		espaciosActivos := make([]map[string]interface{}, 0)
		espaciosDesactivos := make([]map[string]interface{}, 0)

		for _, espacioId := range grupo["EspaciosAcademicos"].([]interface{}) {
			if espacio, errEspacio := helpers.ObtenerEspacioAcademicoSegunId(espacioId.(string)); errEspacio == nil {
				if espacio["activo"] == true {
					espaciosActivos = append(espaciosActivos, espacio)
				} else {
					espaciosDesactivos = append(espaciosDesactivos, espacio)
				}
			} else {
				return requestresponse.APIResponseDTO(false, 500, nil, "error al obtener espacio académico "+errEspacio.Error())
			}
		}

		grupo["EspaciosAcademicos"] = map[string]interface{}{
			"activos":    espaciosActivos,
			"desactivos": espaciosDesactivos,
		}
		grupo["Nombre"] = fmt.Sprintf("%s %s", grupo["CodigoProyecto"], grupo["IndicadorGrupo"])
	}
	return requestresponse.APIResponseDTO(true, 200, gruposEstudioResp["Data"], nil)
}

func DeleteGrupoEstudio(grupoEstudioId string) requestresponse.APIResponse {
	urlColocaciones := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico?query=Activo:true,GrupoEstudioId:" + grupoEstudioId + "&limit=0"
	var colocaciones map[string]interface{}
	if err := request.GetJson(urlColocaciones, &colocaciones); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio horario"+err.Error())
	}

	if len(colocaciones["Data"].([]interface{})) > 0 {
		return requestresponse.APIResponseDTO(true, 200, nil, "tiene colocaciones")
	}

	_, err := helpers.DesactivarGrupoEstudio(grupoEstudioId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	_, errDesasignar := helpers.DesasignarEspaciosAcademicosDeGrupoEstudio(grupoEstudioId)
	if errDesasignar != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errDesasignar.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, nil, "delete success")
}

func CreateGrupoEstudio(grupoEstudio []byte) requestresponse.APIResponse {
	var grupoEstudioMap map[string]interface{}
	_ = json.Unmarshal(grupoEstudio, &grupoEstudioMap)

	urlGrupoEstudioPost := beego.AppConfig.String("HorarioService") + "grupo-estudio"
	var grupoEstudioPost map[string]interface{}
	if err := request.SendJson(urlGrupoEstudioPost, "POST", &grupoEstudioPost, grupoEstudioMap); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de horario", err.Error())
	}

	espaciosAcademicos := grupoEstudioMap["EspaciosAcademicos"].([]interface{})

	_, err := helpers.AsignarGrupoEstudioAEspaciosAcademicos(grupoEstudioPost["Data"].(map[string]interface{})["_id"].(string), espaciosAcademicos)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, grupoEstudioPost["Data"], "")
}

func UpdateGrupoEstudio(grupoEstudioId string, grupoEstudioEditar []byte) requestresponse.APIResponse {
	_, errDesasignar := helpers.DesasignarEspaciosAcademicosDeGrupoEstudio(grupoEstudioId)
	if errDesasignar != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errDesasignar.Error())
	}

	var grupoEstudioMap map[string]interface{}
	_ = json.Unmarshal(grupoEstudioEditar, &grupoEstudioMap)

	urlGrupoEstudioPut := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	var grupoEstudioPut map[string]interface{}
	if err := request.SendJson(urlGrupoEstudioPut, "PUT", &grupoEstudioPut, grupoEstudioMap); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de horario", err.Error())
	}

	espaciosAcademicosAsigmar := grupoEstudioMap["EspaciosAcademicos"].([]interface{})

	_, errAsignar := helpers.AsignarGrupoEstudioAEspaciosAcademicos(grupoEstudioId, espaciosAcademicosAsigmar)
	if errAsignar != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errAsignar.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, grupoEstudioPut["Data"], "")
}

func CreateEspacioAcademico(espacioAcademico []byte) requestresponse.APIResponse {
	var espacioAcademicoMap map[string]interface{}
	_ = json.Unmarshal(espacioAcademico, &espacioAcademicoMap)

	urlEspacioAcademicoPost := beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico"
	var espacioAcademicoPost map[string]interface{}
	if err := request.SendJson(urlEspacioAcademicoPost, "POST", &espacioAcademicoPost, espacioAcademicoMap); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de espacios academios", err.Error())
	}

	espacioCreadoId := espacioAcademicoPost["Data"].(map[string]interface{})["_id"]
	grupoEstudioId := espacioAcademicoMap["grupo_estudio_id"]

	grupoEditado, errAsignar := helpers.AsignarEspacioAcademicoAGrupoEstudio(espacioCreadoId.(string), grupoEstudioId.(string))
	if errAsignar != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errAsignar.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, grupoEditado, "")
}
