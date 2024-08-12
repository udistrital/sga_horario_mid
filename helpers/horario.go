package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func GetCalendarioDeDependencia(periodoId, nivelId, dependenciaId string) ([]map[string]interface{}, error) {
	dependeciaIdInt, _ := strconv.Atoi(dependenciaId)
	urlCalendario := beego.AppConfig.String("EventoService") + "calendario?query=Activo:true,Nivel:" + nivelId + ",PeriodoId:" + periodoId +
		"&limit:0&fields=Id,DependenciaId,DependenciaParticularId,AplicaExtension,Nombre"

	var calendarios []map[string]interface{}
	if err := request.GetJson(urlCalendario, &calendarios); err != nil {
		return nil, fmt.Errorf("error en el servicio de evento")
	}

	var calendariosDondeEstaDependecia []map[string]interface{}
	for _, calendario := range calendarios {
		aplicaExtension := calendario["AplicaExtension"].(bool)
		var listaProyectos map[string][]int

		if aplicaExtension {
			dependenciaParticularId := calendario["DependenciaParticularId"].(string)
			json.Unmarshal([]byte(dependenciaParticularId), &listaProyectos)
		} else {
			dependenciaId := calendario["DependenciaId"].(string)
			json.Unmarshal([]byte(dependenciaId), &listaProyectos)
		}

		for _, idProyecto := range listaProyectos["proyectos"] {
			if idProyecto == dependeciaIdInt {
				calendariosDondeEstaDependecia = append(calendariosDondeEstaDependecia, calendario)
				break
			}
		}
	}

	return calendariosDondeEstaDependecia, nil
}

func GetTipoEventosSegunEspecificacion(calendarios []map[string]interface{}, especificacion string) ([]map[string]interface{}, error) {
	// 'especificacion' es el criterio de filtrado, como por ejemplo "horario".
	// Si el nombre del tipo de evento es "Inscripción de horarios", entonces se agrega el tipoEvento al resultado.
	var tipoEventosConEspecificacion []map[string]interface{}

	for _, calendario := range calendarios {
		urlTipoEvento := beego.AppConfig.String("EventoService") + "tipo_evento?query=CalendarioID:" + strconv.Itoa(int(calendario["Id"].(float64))) + ",Activo:true&fields=Id,Nombre"
		var tipoEventos []map[string]interface{}
		if err := request.GetJson(urlTipoEvento, &tipoEventos); err != nil {
			return nil, fmt.Errorf("error en el servicio de evento")
		}

		if len(tipoEventos) > 0 && len(tipoEventos[0]) == 0 {
			continue
		}

		for _, tipoEvento := range tipoEventos {
			nombreEvento := strings.ToLower(strings.TrimSpace(tipoEvento["Nombre"].(string)))
			if strings.Contains(nombreEvento, especificacion) {
				tipoEventosConEspecificacion = append(tipoEventosConEspecificacion, tipoEvento)
			}
		}
	}

	return tipoEventosConEspecificacion, nil
}

func GetActividadesParaHorarioYPlanDocente(tipoEventosConEspecificacion []map[string]interface{}, dependenciaIdInt int) (map[string]interface{}, error) {
	var actividadesInscripcionHorario, actividadesInscripcionPlanDocente []map[string]interface{}

	for _, tipoEvento := range tipoEventosConEspecificacion {
		urlCalendarioEvento := beego.AppConfig.String("EventoService") + "calendario_evento?query=tipo_evento_id:" + strconv.Itoa(int(tipoEvento["Id"].(float64))) +
			",Activo:true&fields=Id,Nombre,FechaInicio,FechaFin,DependenciaId"
		var calendarioEventos []map[string]interface{}
		if err := request.GetJson(urlCalendarioEvento, &calendarioEventos); err != nil {
			return nil, fmt.Errorf("error en el servicio de evento")
		}

		for _, calendarioEvento := range calendarioEventos {
			var dependencia map[string]interface{}
			json.Unmarshal([]byte(calendarioEvento["DependenciaId"].(string)), &dependencia)

			nombreEvento := strings.ToLower(strings.TrimSpace(calendarioEvento["Nombre"].(string)))
			esHorario := strings.Contains(nombreEvento, "horario")
			esPlanDocente := strings.Contains(nombreEvento, "docente")

			if esHorario || esPlanDocente {
				for _, fecha := range dependencia["fechas"].([]interface{}) {
					if int(fecha.(map[string]interface{})["Id"].(float64)) == dependenciaIdInt {
						fechaInicio, _ := time.Parse(time.RFC3339, fecha.(map[string]interface{})["Inicio"].(string))
						fechaFin, _ := time.Parse(time.RFC3339, fecha.(map[string]interface{})["Fin"].(string))

						calendarioEvento["FechaInicio"] = fechaInicio
						calendarioEvento["FechaFin"] = fechaFin
						delete(calendarioEvento, "DependenciaId")

						fechaHoy := time.Now()
						calendarioEvento["DentroFechas"] = fechaHoy.After(fechaInicio) && fechaHoy.Before(fechaFin)

						if esHorario {
							actividadesInscripcionHorario = append(actividadesInscripcionHorario, calendarioEvento)
						}
						if esPlanDocente {
							actividadesInscripcionPlanDocente = append(actividadesInscripcionPlanDocente, calendarioEvento)
						}
					}
				}
			}
		}
	}

	actividades := map[string]interface{}{
		"actividadesInscripcionHorario":     actividadesInscripcionHorario,
		"actividadesInscripcionPlanDocente": actividadesInscripcionPlanDocente,
	}

	return actividades, nil
}
