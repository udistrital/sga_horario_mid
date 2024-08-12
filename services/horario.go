package services

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetActividadesParaHorario(periodoId, nivelId, dependenciaId string) requestresponse.APIResponse {
	dependeciaIdInt, _ := strconv.Atoi(dependenciaId)
	urlCalendario := beego.AppConfig.String("EventoService") + "calendario?query=Activo:true,Nivel:" + nivelId + ",PeriodoId:" + periodoId +
		"&limit:0&fields=Id,DependenciaId,DependenciaParticularId,AplicaExtension,Nombre"
	var calendarios []map[string]interface{}
	if err := request.GetJson(urlCalendario, &calendarios); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de evento")
	}

	var calendariosDondeEstaDependecia []map[string]interface{}
	for _, calendario := range calendarios {
		aplicaExtension := calendario["AplicaExtension"].(bool)
		if aplicaExtension {
			var listaProyectos map[string][]int
			dependenciaParticularId := calendario["DependenciaParticularId"].(string)
			json.Unmarshal([]byte(dependenciaParticularId), &listaProyectos)
			for _, idProyecto := range listaProyectos["proyectos"] {
				if idProyecto == dependeciaIdInt {
					calendariosDondeEstaDependecia = append(calendariosDondeEstaDependecia, calendario)
					break
				}
			}
		} else {
			var listaProyectos map[string][]int
			dependenciaId := calendario["DependenciaId"].(string)
			json.Unmarshal([]byte(dependenciaId), &listaProyectos)
			for _, idProyecto := range listaProyectos["proyectos"] {
				if idProyecto == dependeciaIdInt {
					calendariosDondeEstaDependecia = append(calendariosDondeEstaDependecia, calendario)
					break
				}
			}

		}
	}

	if len(calendariosDondeEstaDependecia) > 0 && len(calendariosDondeEstaDependecia[0]) == 0 {
		return requestresponse.APIResponseDTO(true, 200, nil, "No hay calendario con los parámetros dados")
	}

	var tipoEventosConEspecificacion []map[string]interface{}
	for _, calendario := range calendariosDondeEstaDependecia {
		urlTipoEvento := beego.AppConfig.String("EventoService") + "tipo_evento?query=CalendarioID:" + strconv.Itoa(int(calendario["Id"].(float64))) + ",Activo:true&fields=Id,Nombre"
		var tipoEventos []map[string]interface{}
		if err := request.GetJson(urlTipoEvento, &tipoEventos); err != nil {
			return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de evento")
		}

		if len(tipoEventos) > 0 && len(tipoEventos[0]) == 0 {
			continue
		}

		for _, tipoEvento := range tipoEventos {
			nombreEvento := strings.ToLower(tipoEvento["Nombre"].(string))
			if strings.Contains(strings.TrimSpace(nombreEvento), "periodo") {
				tipoEventosConEspecificacion = append(tipoEventosConEspecificacion, tipoEvento)
			}
		}
	}

	if len(tipoEventosConEspecificacion) > 0 && len(tipoEventosConEspecificacion[0]) == 0 {
		return requestresponse.APIResponseDTO(true, 200, nil, "No hay proceso de planeación de los períodos académicos")
	}

	var actividadesInscripcionHorario, actividadesInscripcionPlanDocente []map[string]interface{}

	for _, tipoEvento := range tipoEventosConEspecificacion {
		urlCalendarioEvento := beego.AppConfig.String("EventoService") + "calendario_evento?query=tipo_evento_id:" + strconv.Itoa(int(tipoEvento["Id"].(float64))) +
			",Activo:true&fields=Id,Nombre,FechaInicio,FechaFin,DependenciaId"
		var calendarioEventos []map[string]interface{}
		if err := request.GetJson(urlCalendarioEvento, &calendarioEventos); err != nil {
			return requestresponse.APIResponseDTO(false, 500, nil, "Error en el servicio de evento")
		}

		for _, calendarioEvento := range calendarioEventos {
			var dependencia map[string]interface{}
			json.Unmarshal([]byte(calendarioEvento["DependenciaId"].(string)), &dependencia)

			nombreEvento := strings.ToLower(calendarioEvento["Nombre"].(string))
			esHorario, esPlanDocente := strings.Contains(strings.TrimSpace(nombreEvento), "horario"), strings.Contains(strings.TrimSpace(nombreEvento), "docente")

			if esHorario || esPlanDocente {
				for _, fecha := range dependencia["fechas"].([]interface{}) {
					if int(fecha.(map[string]interface{})["Id"].(float64)) == dependeciaIdInt {
						fechaInicio, _ := time.Parse(time.RFC3339, fecha.(map[string]interface{})["Inicio"].(string))
						fechaFin, _ := time.Parse(time.RFC3339, fecha.(map[string]interface{})["Fin"].(string))

						calendarioEvento["FechaInicio"] = fechaInicio
						calendarioEvento["FechaFin"] = fechaFin
						delete(calendarioEvento, "DependenciaId") //quita este atributo que ya no nos sirve

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

	return requestresponse.APIResponseDTO(true, 200, actividades, "")
}
