package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func GetColocacionesSegunGrupoEstudioYPeriodo(grupoEstudioId, periodoId string) ([]interface{}, error) {
	urlGrupoEstudio := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	var resGrupoEstudio map[string]interface{}

	if err := request.GetJson(urlGrupoEstudio, &resGrupoEstudio); err != nil {
		return nil, fmt.Errorf("error en el servicio de horario: %w", err)
	}

	espaciosAcademicosIds := resGrupoEstudio["Data"].(map[string]interface{})["EspaciosAcademicos"].([]interface{})
	var colocacionesTotales []interface{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	errs := make(chan error, len(espaciosAcademicosIds))

	for _, espacioAcademicoId := range espaciosAcademicosIds {
		wg.Add(1)
		go func(espacioAcademicoId string) {
			defer wg.Done()

			colocaciones, err := GetColocacionesDeEspacioAcademicoPorPeriodo(espacioAcademicoId, periodoId)
			if err != nil {
				errs <- err
				return
			}

			if colocaciones != nil {
				mu.Lock()
				colocacionesTotales = append(colocacionesTotales, colocaciones...)
				mu.Unlock()
			}
		}(espacioAcademicoId.(string))
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return nil, <-errs
	}

	return colocacionesTotales, nil
}

func GetColocacionesDeEspacioAcademicoPorPeriodo(espacioAcademicoId, periodoId string) ([]interface{}, error) {
	urlColocacion := beego.AppConfig.String("HorarioService") +
		"colocacion-espacio-academico?query=PeriodoId:" + periodoId + ",EspacioAcademicoId:" + espacioAcademicoId + ",Activo:true&limit=0"

	var resColocaciones map[string]interface{}
	if err := request.GetJson(urlColocacion, &resColocaciones); err != nil {
		return nil, fmt.Errorf("error en el servicio de espacios académicos: %w", err)
	}

	if data, ok := resColocaciones["Data"].([]interface{}); ok && len(data) > 0 {
		return data, nil
	}

	return nil, nil
}

func AgregarInfoAdicionalColocacion(colocacion map[string]interface{}) (map[string]interface{}, error) {
	// Obtener Sede, Edificio y Salón
	if err := GetSedeEdificioSalon(colocacion); err != nil {
		return nil, fmt.Errorf("error al obtener sede, edificio y salón: %w", err)
	}

	// Agregar objeto completo de Espacio Académico
	if id, ok := colocacion["EspacioAcademicoId"].(string); ok {
		if espacioAcademico, err := ObtenerEspacioAcademicoSegunId(id); err == nil {
			colocacion["EspacioAcademico"] = espacioAcademico
		} else {
			return nil, fmt.Errorf("error al obtener espacio académico: %w", err)
		}
	}

	// Obtener el Plan Docente
	urlCargaPlan := beego.AppConfig.String("PlanDocenteService") + "carga_plan?query=colocacion_espacio_academico_id:" + colocacion["_id"].(string) + ",activo:true"
	var cargaPlanes map[string]interface{}
	if err := request.GetJson(urlCargaPlan, &cargaPlanes); err != nil {
		return nil, fmt.Errorf("error en el servicio de plan docente: %w", err)
	}

	if data, ok := cargaPlanes["Data"].([]interface{}); ok && len(data) > 0 {
		planDocenteId := data[0].(map[string]interface{})["plan_docente_id"].(string)

		// Obtener detalle del Plan Docente
		urlPlanDocente := beego.AppConfig.String("PlanDocenteService") + "plan_docente/" + planDocenteId
		var planDocente map[string]interface{}
		if err := request.GetJson(urlPlanDocente, &planDocente); err != nil {
			return nil, fmt.Errorf("error en el servicio de plan docente: %w", err)
		}

		// Obtener información del docente
		docenteId := planDocente["Data"].(map[string]interface{})["docente_id"].(string)
		urlDocente := beego.AppConfig.String("TercerosService") + "tercero/" + docenteId
		var docente map[string]interface{}
		if err := request.GetJson(urlDocente, &docente); err != nil {
			return nil, fmt.Errorf("error en el servicio de terceros: %w", err)
		}

		colocacion["Docente"] = docente
	}
	return colocacion, nil
}

// Con el objeto de colocación trae los datos completos de sede, edificio y salón
func GetSedeEdificioSalon(colocacion map[string]interface{}) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error

	resumen := colocacion["ResumenColocacionEspacioFisico"].(string)
	var resumenMap map[string]interface{}
	if json.Unmarshal([]byte(resumen), &resumenMap) != nil {
		return fmt.Errorf("error al deserializar resumen")
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
					err = fmt.Errorf("error en el servicio de oikos: %s", idField)
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
		return nil, fmt.Errorf("error en el servicio horario: %v", err)
	}

	colocacionData := colocacion["Data"].(map[string]interface{})
	colocacionData["Activo"] = false

	urlColocacionPost := beego.AppConfig.String("HorarioService") + "colocacion-espacio-academico/" + colocacionData["_id"].(string)
	var colocacionPost map[string]interface{}
	if err := request.SendJson(urlColocacionPost, "PUT", &colocacionPost, colocacionData); err != nil {
		return nil, fmt.Errorf("error en el servicio de horario: %v", err)
	}

	return colocacion, nil
}

// haySobreposicion verifica si hay una superposición entre la colocación a poner
// con la colocacion ya puesta en base de datos
//
// Retorna:
//
// - true si la colocacion a poner y la colocacion ya establecida se solapan
//
// - false en caso contrario.
func HaySobreposicion(colocacionAPoner, colocacionPuesta map[string]interface{}) bool {
	finalX := colocacionAPoner["finalPosition"].(map[string]interface{})["x"].(float64)
	horaFormatoColocacion := colocacionAPoner["horaFormato"].(string)

	inicioColocacion, finColocacion, err := ObtenerMinutosDeRangoHora(horaFormatoColocacion)
	if err != nil {
		fmt.Println("Error al parsear horaFormato de colocacion:", err)
		return false
	}

	if colocacionPuesta["finalPosition"].(map[string]interface{})["x"].(float64) == finalX {
		horaFormatoEspacio := colocacionPuesta["horaFormato"].(string)
		inicioEspacio, finEspacio, err := ObtenerMinutosDeRangoHora(horaFormatoEspacio)
		if err != nil {
			fmt.Println("Error al parsear horaFormato de espacio ocupado:", err)
			return false
		}

		if (inicioColocacion < finEspacio && finColocacion > inicioEspacio) ||
			(inicioEspacio < finColocacion && finEspacio > inicioColocacion) {
			return true
		}
	}

	return false
}

// ObtenerMinutosDeRangoHora convierte un rango de horas en formato "HH:MM - HH:MM"
// a minutos desde medianoche.
//
// Parámetros:
// - rangoHora: Cadena en formato "HH:MM - HH:MM".
//
// Retorna:
// - inicio: Minutos desde medianoche del inicio del rango.
// - fin: Minutos desde medianoche del fin del rango.
// - err: Error si el formato es inválido.
//
// Ejemplo:
// Si rangoHora es "08:45 - 17:30", retorna:
//
// - inicio = 525 -> (8*60 + 45)
//
// - fin = 1050 -> (17*60 + 30)
func ObtenerMinutosDeRangoHora(rangoHora string) (inicio, fin int, err error) {
	horas := strings.Split(rangoHora, " - ")
	if len(horas) != 2 {
		return 0, 0, fmt.Errorf("formato de hora inválido")
	}

	convertir := func(hora string) (int, error) {
		partes := strings.Split(hora, ":")
		if len(partes) != 2 {
			return 0, fmt.Errorf("formato de hora inválido")
		}
		h, _ := strconv.Atoi(partes[0])
		m, _ := strconv.Atoi(partes[1])
		return h*60 + m, nil
	}

	if inicio, err = convertir(horas[0]); err != nil {
		return 0, 0, err
	}
	if fin, err = convertir(horas[1]); err != nil {
		return 0, 0, err
	}

	return inicio, fin, nil
}
