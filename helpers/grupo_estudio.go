package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func DesactivarGrupoEstudio(grupoEstudioId string) (map[string]interface{}, error) {
	grupoEstudio := map[string]interface{}{
		"Activo": false,
	}

	urlGrupoEstudioPut := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	var grupoEstudioPut map[string]interface{}
	if err := request.SendJson(urlGrupoEstudioPut, "PUT", &grupoEstudioPut, grupoEstudio); err != nil {
		return nil, fmt.Errorf("error en el servicio de horario: %v", err)
	}

	return grupoEstudioPut["Data"].(map[string]interface{}), nil
}

func ObtenerEspacioAcademicoSegunId(espacioId string) (map[string]interface{}, error) {
	var resultado map[string]interface{}
	urlGetEspacio := beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico?query=_id:" + espacioId + "&fields=nombre,grupo,espacio_academico_padre,activo"

	if err := request.GetJson(urlGetEspacio, &resultado); err != nil {
		return nil, err
	}

	if data, ok := resultado["Data"].([]interface{}); ok && len(data) > 0 {
		if espacio, ok := data[0].(map[string]interface{}); ok {
			return espacio, nil
		}
	}
	return nil, fmt.Errorf("no se encontró el espacio académico con el ID %s", espacioId)
}

func AsignarGrupoEstudioAEspaciosAcademicos(grupoEstudioId string, espaciosAcademicos []interface{}) (map[string]interface{}, error) {
	for _, espacioId := range espaciosAcademicos {
		espacioAsignarGrupoEstudio := map[string]interface{}{
			"grupo_estudio_id": grupoEstudioId,
		}

		urlGrupoEstudioPut := beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico/" + espacioId.(string)
		var colocacionPut map[string]interface{}
		if err := request.SendJson(urlGrupoEstudioPut, "PUT", &colocacionPut, espacioAsignarGrupoEstudio); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func DesasignarEspaciosAcademicosDeGrupoEstudio(grupoEstudioId string) (map[string]interface{}, error) {
	urlGrupoEstudio := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	var grupoEstudio map[string]interface{}
	if err := request.GetJson(urlGrupoEstudio, &grupoEstudio); err != nil {
		return nil, fmt.Errorf("error en servicio de horarios" + err.Error())
	}

	espaciosParaDesasignar := grupoEstudio["Data"].(map[string]interface{})["EspaciosAcademicos"].([]interface{})
	for _, espacioId := range espaciosParaDesasignar {
		espacioAsignarGrupoEstudio := map[string]interface{}{
			"grupo_estudio_id": "0",
		}

		urlGrupoEstudioPut := beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico/" + espacioId.(string)
		var colocacionPut map[string]interface{}
		if err := request.SendJson(urlGrupoEstudioPut, "PUT", &colocacionPut, espacioAsignarGrupoEstudio); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func AsignarEspacioAcademicoAGrupoEstudio(espacioAcademicoId, grupoEstudioId string) (map[string]interface{}, error) {
	urlGrupoEstudio := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	var grupoEstudio map[string]interface{}
	if err := request.GetJson(urlGrupoEstudio, &grupoEstudio); err != nil {
		return nil, fmt.Errorf("error en el servicio horario" + err.Error())
	}

	grupoEstudio = grupoEstudio["Data"].(map[string]interface{})
	grupoEstudio["EspaciosAcademicos"] = append(grupoEstudio["EspaciosAcademicos"].([]interface{}), espacioAcademicoId)

	urlGrupoEstudioPut := beego.AppConfig.String("HorarioService") + "grupo-estudio/" + grupoEstudioId
	var grupoEstudioPut map[string]interface{}
	if err := request.SendJson(urlGrupoEstudioPut, "PUT", &grupoEstudioPut, grupoEstudio); err != nil {
		return nil, fmt.Errorf("error en el servicio de horario" + err.Error())
	}

	return grupoEstudioPut["Data"].(map[string]interface{}), nil

}
