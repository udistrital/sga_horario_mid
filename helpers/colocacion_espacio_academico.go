package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func GetColocacionesDeModuloHorario(grupoEstudioId string) ([]map[string]interface{}, error) {
	urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico?query=GrupoEstudioId:" + grupoEstudioId + ",Activo:true&limit=0"
	var resColocaciones map[string]interface{}
	if err := request.GetJson(urlColocacion, &resColocaciones); err != nil {
		return nil, fmt.Errorf("error en el servicio de terceros: %w", err)
	}

	data, ok := resColocaciones["Data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de datos inesperado")
	}

	colocaciones := make([]map[string]interface{}, len(data))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(data))

	for i, colocacion := range data {
		wg.Add(1)
		go func(i int, colocacion interface{}) {
			defer wg.Done()
			colocacionMap, ok := colocacion.(map[string]interface{})
			if !ok {
				errCh <- fmt.Errorf("error al convertir colocacion a mapa")
				return
			}

			// Agrega objeto completo para sede, edificio y salón
			if err := GetSedeEdificioSalon(colocacionMap); err != nil {
				errCh <- fmt.Errorf("error al obtener sede, edificio y salón: %w", err)
				return
			}

			// Agrega objeto completo de espacio académico
			if id, ok := colocacionMap["EspacioAcademicoId"].(string); ok {
				if espacioAcademico, err := ObtenerEspacioAcademicoSegunId(id); err == nil {
					colocacionMap["EspacioAcademico"] = espacioAcademico
				} else {
					errCh <- fmt.Errorf("error al obtener espacio académico: %w", err)
					return
				}
			}

			mu.Lock()
			colocaciones[i] = colocacionMap
			mu.Unlock()
		}(i, colocacion)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return colocaciones, nil
}

func GetColocacionesDeModuloPlanDocente(grupoEstudioId, periodoId string) ([]map[string]interface{}, error) {
	//Obtengo los planes docente que pertenecen al periodo dado
	urlPlanDocente := beego.AppConfig.String("PlanDocenteService") + "plan_docente?query=periodo_id:" + periodoId + ",activo:true&limit=0"

	var planesDocente map[string]interface{}
	if err := request.GetJson(urlPlanDocente, &planesDocente); err != nil {
		return nil, fmt.Errorf("error en el servicio de plan docente: %w", err)
	}

	//Obtengo el grupo de estudio segun el id
	urlColocacion := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	var resColocaciones map[string]interface{}
	if err := request.GetJson(urlColocacion, &resColocaciones); err != nil {
		return nil, fmt.Errorf("error en el servicio de horario: %w", err)
	}

	//Se extrae los espacios academicos del grupo de estudio
	espaciosAcademicos, ok := resColocaciones["Data"].(map[string]interface{})["EspaciosAcademicos"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error al procesar EspaciosAcademicos")
	}

	var colocaciones []map[string]interface{}
	//go routines
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(planesDocente["Data"].([]interface{})))

	for _, planDocente := range planesDocente["Data"].([]interface{}) {
		planDocenteMap, ok := planDocente.(map[string]interface{})
		if !ok {
			continue
		}

		wg.Add(1)
		go func(planDocenteMap map[string]interface{}) {
			defer wg.Done()
			//accedo a las cargas plan del plan docente
			urlCargaPlan := beego.AppConfig.String("PlanDocenteService") + "carga_plan?query=plan_docente_id:" + planDocenteMap["_id"].(string) + ",activo:true&limit=0"
			var cargaPlanes map[string]interface{}
			if err := request.GetJson(urlCargaPlan, &cargaPlanes); err != nil {
				errCh <- fmt.Errorf("error en el servicio de plan docente: %w", err)
				return
			}

			for _, cargaPlan := range cargaPlanes["Data"].([]interface{}) {
				cargaPlanMap, ok := cargaPlan.(map[string]interface{})
				if !ok || cargaPlanMap["colocacion_espacio_academico_id"] == nil {
					continue
				}
				//accedo a la colocacion espacio de la carga plan
				urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + cargaPlanMap["colocacion_espacio_academico_id"].(string)

				var colocacion map[string]interface{}
				if err := request.GetJson(urlColocacion, &colocacion); err != nil {
					errCh <- fmt.Errorf("error en el servicio horario: %w", err)
					return
				}

				colocacionData, ok := colocacion["Data"].(map[string]interface{})
				if !ok {
					continue
				}

				if colocacionData["Activo"] == false {
					continue
				}

				espacioAcademicoId, exists := colocacionData["EspacioAcademicoId"]
				if !exists {
					continue
				}

				// reviso si ese espacio academico esta en el arreglo de espacios academicos del grupo de estudio
				if !Contains(espaciosAcademicos, espacioAcademicoId.(string)) {
					continue
				}

				//si existe, se le colocan los atributos para responder al cliente
				if err := GetSedeEdificioSalon(colocacionData); err != nil {
					errCh <- fmt.Errorf("error al obtener sede, edificio y salón: %w", err)
					return
				}

				if id, ok := colocacionData["EspacioAcademicoId"].(string); ok {
					if espacioAcademico, err := ObtenerEspacioAcademicoSegunId(id); err == nil {
						colocacionData["EspacioAcademico"] = espacioAcademico
					} else {
						errCh <- fmt.Errorf("error al obtener espacio académico: %w", err)
						return
					}
				}

				urlDocente := beego.AppConfig.String("TercerosService") + "tercero/" + planDocenteMap["docente_id"].(string)
				var docente map[string]interface{}
				if err := request.GetJson(urlDocente, &docente); err != nil {
					errCh <- fmt.Errorf("error en el servicio terceros: %w", err)
					return
				}

				colocacionData["CargaPlanId"] = cargaPlanMap["_id"]
				colocacionData["Docente"] = docente

				mu.Lock()
				colocaciones = append(colocaciones, colocacionData)
				mu.Unlock()
			}
		}(planDocenteMap)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return colocaciones, nil
}

// Con el objeto de colocación trae los datos completos de sede, edificio y salón
func GetSedeEdificioSalon(colocacion map[string]interface{}) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error

	resumen := colocacion["ResumenColocacionEspacioFisico"].(string)
	var resumenMap map[string]interface{}
	if json.Unmarshal([]byte(resumen), &resumenMap) != nil {
		return fmt.Errorf("Error al deserializar resumen")
	}
	espacioFisico := resumenMap["espacio_fisico"].(map[string]interface{})

	for _, idField := range []string{"edificio_id", "salon_id", "sede_id"} {
		wg.Add(1)
		go func(idField string, id interface{}) {
			defer wg.Done()
			url := beego.AppConfig.String("OikosService") + "espacio_fisico/" + fmt.Sprintf("%v", id)
			var res map[string]interface{}
			if fetchErr := request.GetJson(url, &res); fetchErr != nil {
				mu.Lock()
				if err == nil {
					err = fmt.Errorf("Error en el servicio de oikos: %s", idField)
				}
				mu.Unlock()
				return
			}
			mu.Lock()
			espacioFisico[strings.TrimSuffix(idField, "_id")] = res
			mu.Unlock()
		}(idField, espacioFisico[idField])
	}

	wg.Wait()
	resumenMap["espacio_fisico"] = espacioFisico
	colocacion["ResumenColocacionEspacioFisico"] = resumenMap

	return err
}

func Contains(slice []interface{}, item string) bool {
	for _, v := range slice {
		if str, ok := v.(string); ok && str == item {
			return true
		}
	}
	return false
}

func DesactivarColocacion(colocacionId string) (map[string]interface{}, error) {
	urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionId
	var colocacion map[string]interface{}
	if err := request.GetJson(urlColocacion, &colocacion); err != nil {
		return nil, fmt.Errorf("Error en el servicio horario: %v", err)
	}

	colocacionData := colocacion["Data"].(map[string]interface{})
	colocacionData["Activo"] = false

	urlColocacionPost := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionData["_id"].(string)
	var colocacionPost map[string]interface{}
	if err := request.SendJson(urlColocacionPost, "PUT", &colocacionPost, colocacionData); err != nil {
		return nil, fmt.Errorf("Error en el servicio de horario: %v", err)
	}

	return colocacion, nil
}
