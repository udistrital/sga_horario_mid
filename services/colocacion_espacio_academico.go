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

func GetColocacionesSegunGrupoEstudioYPeriodo(grupoEstudioId, periodoId string) requestresponse.APIResponse {
	colocacionesTotales, err := helpers.GetColocacionesSegunGrupoEstudioYPeriodo(grupoEstudioId, periodoId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en metodo GetColocacionesSegunGrupoEstudioYPeriodo: %v", err), err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(colocacionesTotales))

	for _, colocacion := range colocacionesTotales {
		wg.Add(1)
		go func(colocacion map[string]interface{}) {
			defer wg.Done()
			_, err := helpers.AgregarInfoAdicionalColocacion(colocacion)
			if err != nil {
				errChan <- err
			}
		}(colocacion.(map[string]interface{}))
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en metodo AgregarInfoAdicionalColocacion: %v", err), err)
		}
	}

	return requestresponse.APIResponseDTO(true, 200, colocacionesTotales, "")
}
