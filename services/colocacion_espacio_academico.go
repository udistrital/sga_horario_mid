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

// GetSobreposicionEspacioFisico verifica si una colocación sobrepone a cualquier otra existente
// con respecto al espacio fisico durante un período determinado.
func GetSobreposicionEspacioFisico(colocacionId, periodoId string) requestresponse.APIResponse {
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
			_ = helpers.GetSedeEdificioSalon(espacioOcupado)

			colocacionSobrepuesta = map[string]interface{}{
				"sobrepuesta":         true,
				"colocacionConflicto": espacioOcupado,
			}
			break
		}
	}

	return requestresponse.APIResponseDTO(true, 200, colocacionSobrepuesta, "")
}

func GetColocacionesDeGrupoEstudio(grupoEstudioId, periodoId string) requestresponse.APIResponse {
	colocacionesTotales, err := helpers.GetColocacionesDeGrupoEstudio(grupoEstudioId, periodoId)
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
			return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en metodo AgregarInfoAdicionalColocacion: %v", err))
		}
	}

	return requestresponse.APIResponseDTO(true, 200, colocacionesTotales, "")
}

func GetColocacionInfoAdicional(colocacionId string) requestresponse.APIResponse {
	urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionId
	var colocacionEspacioAcademico map[string]interface{}
	if err := request.GetJson(urlColocacion, &colocacionEspacioAcademico); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio horario"+err.Error())
	}

	colocacionInfoAdicional, err := helpers.AgregarInfoAdicionalColocacion(colocacionEspacioAcademico["Data"].(map[string]interface{}))
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, err.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, colocacionInfoAdicional, "")
}

func DeleteColocacionEspacioAcademico(colocacionId string) requestresponse.APIResponse {
	_, errPlan := helpers.DesactivarCargaPlanSegunColocacion(colocacionId)
	if errPlan != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en metodo DesactivarColocacion: %v", errPlan))
	}

	_, err := helpers.DesactivarColocacion(colocacionId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en metodo DesactivarColocacion: %v", err))
	}

	return requestresponse.APIResponseDTO(true, 200, nil, "delete success")
}

func GetColocacionesGrupoSinDetalles(grupoEstudioId, periodoId string) requestresponse.APIResponse {
	colocaciones, err := helpers.GetColocacionesDeGrupoEstudio(grupoEstudioId, periodoId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en metodo GetColocacionesDeGrupoEstudio: %v", err), err)
	}

	var colocacionesSinDetalles []map[string]interface{}
	for _, colocacion := range colocaciones {
		colocacionData, _ := colocacion.(map[string]interface{})

		var colocacionEspacio map[string]interface{}
		_ = json.Unmarshal([]byte(colocacionData["ColocacionEspacioAcademico"].(string)), &colocacionEspacio)

		colocacionesSinDetalles = append(colocacionesSinDetalles, map[string]interface{}{
			"_id":           colocacionData["_id"],
			"horas":         int(colocacionEspacio["horas"].(float64)),
			"finalPosition": colocacionEspacio["finalPosition"],
			"horaFormato":   colocacionEspacio["horaFormato"],
		})
	}
	return requestresponse.APIResponseDTO(true, 200, colocacionesSinDetalles, "")
}

// Obtiene si una colocación se sobrepone a alguna colocación del grupo de estudio
func GetSobreposicionEnGrupoEstudio(grupoEstudioId, periodoId, colocacionId string) requestresponse.APIResponse {
	colocacionesGrupoEstudio, err := helpers.GetColocacionesDeGrupoEstudio(grupoEstudioId, periodoId)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en metodo GetColocacionesDeGrupoEstudio: %v", err), err)
	}

	urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionId
	var colocacionEspacioAcademico map[string]interface{}
	if err := request.GetJson(urlColocacion, &colocacionEspacioAcademico); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en el servicio horario: %v", err), err)
	}

	colocacionSobrepuesta, err := helpers.VerificarSobreposicionEnColocaciones(colocacionEspacioAcademico["Data"].(map[string]interface{}), colocacionesGrupoEstudio)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, fmt.Sprintf("error en VerificarSobreposicion: %v", err), err)
	}

	return requestresponse.APIResponseDTO(true, 200, colocacionSobrepuesta, "")
}

func CopiarColocacionesAGrupoEstudio(infoParaCopiado []byte) requestresponse.APIResponse {
	var infoParaCopiadoMap map[string]interface{}
	_ = json.Unmarshal(infoParaCopiado, &infoParaCopiadoMap)
	grupoEstudioId := infoParaCopiadoMap["grupoEstudioId"]
	periodoId := infoParaCopiadoMap["periodoId"]
	colocaciones := infoParaCopiadoMap["colocaciones"].([]interface{})

	var colocacionesPost []map[string]interface{}
	for _, colocacion := range colocaciones {
		colocacionMap := colocacion.(map[string]interface{})
		urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionMap["colocacionId"].(string)
		var colocacion map[string]interface{}
		if err := request.GetJson(urlColocacion, &colocacion); err != nil {
			return requestresponse.APIResponseDTO(true, 200, nil, "Error en el servicio horario"+err.Error())
		}
		colocacion = colocacion["Data"].(map[string]interface{})
		delete(colocacion, "_id")
		colocacion["GrupoEstudioId"] = grupoEstudioId
		colocacion["PeriodoId"] = periodoId
		colocacion["EspacioAcademicoId"] = colocacionMap["espacioAcademicoId"]

		urlColocacionPost := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico"
		var colocacionPost map[string]interface{}
		if err := request.SendJson(urlColocacionPost, "POST", &colocacionPost, colocacion); err != nil {
			return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de horario", err.Error())
		}
		colocacionesPost = append(colocacionesPost, colocacionPost["Data"].(map[string]interface{}))
	}

	horarioCopiado := map[string]interface{}{
		"grupoEstudioId": grupoEstudioId,
		"colocaciones":   colocacionesPost,
	}

	return requestresponse.APIResponseDTO(true, 200, horarioCopiado, "")
}
