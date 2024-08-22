package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/helpers"
	"github.com/udistrital/utils_oas/request"
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

	//en este se unen las colocaciones del modulo de plan docente y horario
	colocacionesMap := make(map[string]map[string]interface{})

	//Si una colocacion se repite se deja una, priorizando la del modulo de plan docente
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

// GetSobreposicionColocacion verifica si una colocación sobrepone a cualquier otra existente
// con respecto al espacio fisico durante un período determinado.
func GetSobreposicionColocacion(colocacionId, periodoId string) requestresponse.APIResponse {
	urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionId
	var colocacionEspacioAcademico map[string]interface{}
	if err := request.GetJson(urlColocacion, &colocacionEspacioAcademico); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio horario"+err.Error())
	}

	colocacionEspacioAcademico = colocacionEspacioAcademico["Data"].(map[string]interface{})

	var colocacion map[string]interface{}
	_ = json.Unmarshal([]byte(colocacionEspacioAcademico["ColocacionEspacioAcademico"].(string)), &colocacion)

	espacioFisicoId := fmt.Sprintf("%v", colocacionEspacioAcademico["EspacioFisicoId"])
	espaciosOcupados, _ := GetEspaciosFisicosOcupadosSegunPeriodo(espacioFisicoId, periodoId).Data.([]map[string]interface{})

	colocacionSobrepuesta := map[string]interface{}{"sobrepuesta": false}
	for _, espacioOcupado := range espaciosOcupados {
		if helpers.HaySobreposicion(colocacion, espacioOcupado) {
			urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + espacioOcupado["_id"].(string)
			if err := request.GetJson(urlColocacion, &espacioOcupado); err != nil {
				return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio horario"+err.Error())
			}

			espacioOcupado = espacioOcupado["Data"].(map[string]interface{})

			colocacionSobrepuesta = map[string]interface{}{
				"sobrepuesta":         true,
				"colocacionConflicto": espacioOcupado,
			}
			break
		}
	}

	return requestresponse.APIResponseDTO(true, 200, colocacionSobrepuesta, "")
}
