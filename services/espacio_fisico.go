package services

import (
	"fmt"
	"sync"

	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetEspaciosFisicosOcupadosSegunPeriodo(espacioFisicoId, periodoId string) requestresponse.APIResponse {
	fmt.Println(espacioFisicoId)
	var espaciosOcupadosHorario, espaciosOcupadosPlanDocente []map[string]interface{}
	var errPlanDocente, errHorario error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		espaciosOcupadosHorario, errHorario = helpers.GetEspaciosFisicosOcupadosDeHorario(espacioFisicoId, periodoId)
	}()

	go func() {
		defer wg.Done()
		espaciosOcupadosPlanDocente, errPlanDocente = helpers.GetEspaciosFisicosOcupadosDePlanDocente(espacioFisicoId, periodoId)
	}()

	wg.Wait()

	if errPlanDocente != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errPlanDocente.Error())
	}

	if errHorario != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, errHorario.Error())
	}

	//en este se unen las colocaciones del modulo de plan docente y horario
	ocupadosMap := make(map[string]map[string]interface{})

	//Si una colocacion se repite se deja una
	for _, espacioOcupado := range append(espaciosOcupadosPlanDocente, espaciosOcupadosHorario...) {
		id, ok := espacioOcupado["_id"].(string)
		if ok && espacioOcupado != nil && ocupadosMap[id] == nil {
			ocupadosMap[id] = espacioOcupado
		}
	}

	espaciosOcupados := make([]map[string]interface{}, 0, len(ocupadosMap))
	for _, colocacion := range ocupadosMap {
		espaciosOcupados = append(espaciosOcupados, colocacion)
	}

	return requestresponse.APIResponseDTO(true, 200, espaciosOcupados, "")
}
