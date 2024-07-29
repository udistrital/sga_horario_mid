package services

import (
	"sync"

	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetColocacionesSegunGrupoEstudioYPeriodo(grupoEstudioId, periodoId string) requestresponse.APIResponse {
	var colocacionesDeModuloPlanDocente []map[string]interface{}
	var colocacionesDeModuloHorario []map[string]interface{}
	var errPlanDocente, errHorario error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		colocacionesDeModuloPlanDocente, errPlanDocente = helpers.GetColocacionesDeModuloPlanDocente(grupoEstudioId, periodoId)
	}()

	go func() {
		defer wg.Done()
		colocacionesDeModuloHorario, errHorario = helpers.GetColocacionesDeModuloHorario(grupoEstudioId)
	}()

	wg.Wait()

	if errPlanDocente != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errPlanDocente.Error())
	}

	if errHorario != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errHorario.Error())
	}

	colocacionesMap := make(map[string]map[string]interface{})

	for _, colocacion := range append(colocacionesDeModuloPlanDocente, colocacionesDeModuloHorario...) {
		id, ok := colocacion["_id"].(string)
		if ok && colocacion != nil && colocacionesMap[id] == nil {
			colocacionesMap[id] = colocacion
		}
	}

	colocaciones := make([]map[string]interface{}, 0, len(colocacionesMap))
	for _, colocacion := range colocacionesMap {
		colocaciones = append(colocaciones, colocacion)
	}

	return requestresponse.APIResponseDTO(true, 200, colocaciones, "")
}
