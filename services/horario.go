package services

import (
	"strconv"

	"github.com/udistrital/sga_horario_mid/helpers"
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
