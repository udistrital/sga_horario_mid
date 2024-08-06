package helpers

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func GetEspaciosOcupadoSegunPeriodo(espacioFisicoId, periodoId string) ([]map[string]interface{}, error) {
	var ocupados []map[string]interface{}
	var mu sync.Mutex

	var horarios map[string]interface{}
	urlHorario := beego.AppConfig.String("HorarioService") + "horario?query=PeriodoId:" + periodoId + ",Activo:true&limit=0"
	if err := request.GetJson(urlHorario, &horarios); err != nil {
		return nil, fmt.Errorf("error en el servicio de horario: %w", err)
	}

	horariosData, _ := horarios["Data"].([]interface{})

	// Channels para comunicación entre goroutines
	gruposEstudioChan := make(chan map[string]interface{}, len(horariosData))
	colocacionesChan := make(chan map[string]interface{}, len(horariosData)*100)

	var wg sync.WaitGroup

	// Fetch grupos de estudio en paralelo
	for _, horario := range horariosData {
		horarioMap, _ := horario.(map[string]interface{})
		wg.Add(1)
		go func(horarioID string) {
			defer wg.Done()
			urlGrupoEstudio := beego.AppConfig.String("HorarioService") + "grupo-estudio?query=HorarioId:" + horarioID + ",Activo:true&limit=0"
			var gruposEstudio map[string]interface{}
			if err := request.GetJson(urlGrupoEstudio, &gruposEstudio); err != nil {
				fmt.Printf("error en el servicio de grupo estudio: %v\n", err)
				return
			}

			gruposEstudioData, _ := gruposEstudio["Data"].([]interface{})
			for _, grupoEstudio := range gruposEstudioData {
				gruposEstudioChan <- grupoEstudio.(map[string]interface{})
			}
		}(horarioMap["_id"].(string))
	}

	// Cierra el canal de gruposEstudioChan una vez que todas las goroutines hayan terminado
	go func() {
		wg.Wait()
		close(gruposEstudioChan)
	}()

	// Fetch colocaciones en paralelo
	go func() {
		for grupoEstudioMap := range gruposEstudioChan {
			urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico?query=GrupoEstudioId:" + grupoEstudioMap["_id"].(string) + ",EspacioFisicoId:" + espacioFisicoId + ",Activo:true&limit=0"
			var colocaciones map[string]interface{}
			if err := request.GetJson(urlColocacion, &colocaciones); err != nil {
				fmt.Printf("error en el servicio de colocacion: %v\n", err)
				continue
			}

			colocacionesData, _ := colocaciones["Data"].([]interface{})
			for _, colocacion := range colocacionesData {
				colocacionesChan <- colocacion.(map[string]interface{})
			}
		}
		close(colocacionesChan)
	}()

	// Recolecta resultados de colocaciones
	for colocacion := range colocacionesChan {
		var colocacionEspacio map[string]interface{}
		_ = json.Unmarshal([]byte(colocacion["ColocacionEspacioAcademico"].(string)), &colocacionEspacio)
		mu.Lock()
		ocupados = append(ocupados, map[string]interface{}{
			"_id":           colocacion["_id"].(string),
			"horas":         int(colocacionEspacio["horas"].(float64)),
			"finalPosition": colocacionEspacio["finalPosition"],
		})
		mu.Unlock()
	}

	return ocupados, nil
}

func GetEspaciosOcupadoSegunPlanDocente(espacioFisicoId, periodoId string) ([]map[string]interface{}, error) {
	var ocupados []map[string]interface{}
	var mu sync.Mutex

	// Fetch the plan docente data
	urlPlanDocente := beego.AppConfig.String("PlanDocenteService") + "plan_docente?query=periodo_id:" + periodoId + ",activo:true&limit=0"
	var planesDocente map[string]interface{}
	if err := request.GetJson(urlPlanDocente, &planesDocente); err != nil {
		return nil, fmt.Errorf("error en el servicio de plan docente: %w", err)
	}

	planesDocenteData, _ := planesDocente["Data"].([]interface{})

	// Channels para comunicación entre goroutines
	cargaPlanChan := make(chan map[string]interface{}, len(planesDocenteData))
	colocacionesChan := make(chan map[string]interface{}, len(planesDocenteData)*100) // Adjust buffer size as needed

	var wg sync.WaitGroup

	// Fetch grupos de estudio en paralelo
	for _, planDocente := range planesDocenteData {
		planDocenteMap, _ := planDocente.(map[string]interface{})
		wg.Add(1)
		go func(planDocenteID string) {
			defer wg.Done()
			urlCargaPlan := beego.AppConfig.String("PlanDocenteService") + "carga_plan?query=plan_docente_id:" + planDocenteID + ",salon_id:" + espacioFisicoId + ",activo:true&limit=0"
			var cargaPlanes map[string]interface{}
			if err := request.GetJson(urlCargaPlan, &cargaPlanes); err != nil {
				fmt.Printf("error en el servicio de carga plan: %v\n", err)
				return
			}

			cargaPlanesData, _ := cargaPlanes["Data"].([]interface{})
			for _, cargaPlan := range cargaPlanesData {
				cargaPlanMap, _ := cargaPlan.(map[string]interface{})
				if cargaPlanMap["colocacion_espacio_academico_id"] == nil {
					continue
				}
				cargaPlanChan <- cargaPlanMap
			}
		}(planDocenteMap["_id"].(string))
	}

	// Cierra el canal de gruposEstudioChan una vez que todas las goroutines hayan terminado
	go func() {
		wg.Wait()
		close(cargaPlanChan)
	}()

	// Fetch colocaciones en paralelo
	go func() {
		for cargaPlanMap := range cargaPlanChan {
			urlColocacion := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + cargaPlanMap["colocacion_espacio_academico_id"].(string)
			var colocacion map[string]interface{}
			if err := request.GetJson(urlColocacion, &colocacion); err != nil {
				fmt.Printf("error en el servicio de colocacion: %v\n", err)
				continue
			}

			colocacionData, _ := colocacion["Data"].(map[string]interface{})

			var colocacionEspacio map[string]interface{}
			_ = json.Unmarshal([]byte(colocacionData["ColocacionEspacioAcademico"].(string)), &colocacionEspacio)

			colocacionesChan <- map[string]interface{}{
				"_id":           colocacionData["_id"],
				"horas":         int(colocacionEspacio["horas"].(float64)),
				"finalPosition": colocacionEspacio["finalPosition"],
			}
		}
		close(colocacionesChan)
	}()

	// Recolecta resultados de colocaciones
	for colocacion := range colocacionesChan {
		mu.Lock()
		ocupados = append(ocupados, colocacion)
		mu.Unlock()
	}

	return ocupados, nil
}
