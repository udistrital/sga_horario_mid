package services

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetActividadesParaHorarioYPlanDocente(periodoId, nivelId, dependenciaId string) requestresponse.APIResponse {
	dependenciaIdInt, _ := strconv.Atoi(dependenciaId)
	calendariosDondeEstaDependecia, err := helpers.GetCalendarioDeDependencia(periodoId, nivelId, dependenciaId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	if len(calendariosDondeEstaDependecia) > 0 && len(calendariosDondeEstaDependecia[0]) == 0 {
		return requestresponse.APIResponseDTO(true, 200, nil, "No hay calendario con los parámetros dados")
	}

	tipoEventosConEspecificacion, err := helpers.GetTipoEventosSegunEspecificacion(calendariosDondeEstaDependecia, "periodo")
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	if len(tipoEventosConEspecificacion) > 0 && len(tipoEventosConEspecificacion[0]) == 0 {
		return requestresponse.APIResponseDTO(true, 200, nil, "No hay proceso de planeación de los períodos académicos")
	}

	actividades, err := helpers.GetActividadesParaHorarioYPlanDocente(tipoEventosConEspecificacion, dependenciaIdInt)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, actividades, "Actividades obtenidas exitosamente")
}

func CreateHorarioCopia(infoParaCopiado []byte) requestresponse.APIResponse {
	var infoParaCopiadoMap map[string]interface{}
	_ = json.Unmarshal(infoParaCopiado, &infoParaCopiadoMap)
	grupoEstudio := infoParaCopiadoMap["grupoEstudio"]
	colocacionesIds := infoParaCopiadoMap["colocacionesIds"]

	urlGrupoEstudioPost := beego.AppConfig.String("HorarioService") + "grupo-estudio"
	var grupoEstudioPost map[string]interface{}
	if err := request.SendJson(urlGrupoEstudioPost, "POST", &grupoEstudioPost, grupoEstudio); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de horario", err.Error())
	}

	idGrupoEstudio := grupoEstudioPost["Data"].(map[string]interface{})["_id"].(string)

	var colocacionesPost []map[string]interface{}
	for _, colocacionId := range colocacionesIds.([]interface{}) {
		urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionId.(string)
		var colocacion map[string]interface{}
		if err := request.GetJson(urlColocacion, &colocacion); err != nil {
			return requestresponse.APIResponseDTO(true, 200, nil, "Error en el servicio horario"+err.Error())
		}
		colocacion = colocacion["Data"].(map[string]interface{})
		delete(colocacion, "_id")
		colocacion["GrupoEstudioId"] = idGrupoEstudio

		urlColocacionPost := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico"
		var colocacionPost map[string]interface{}
		if err := request.SendJson(urlColocacionPost, "POST", &colocacionPost, colocacion); err != nil {
			return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de horario", err.Error())
		}
		colocacionesPost = append(colocacionesPost, colocacionPost["Data"].(map[string]interface{}))
	}

	horarioCopiado := map[string]interface{}{
		"grupoEstudio": grupoEstudioPost["Data"].(map[string]interface{}),
		"colocaciones": colocacionesPost,
	}

	return requestresponse.APIResponseDTO(true, 200, horarioCopiado, "")
}
