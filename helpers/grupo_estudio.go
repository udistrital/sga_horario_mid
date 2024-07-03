package helpers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func ObtenerEspacioAcademicoSegunId(espacioId string) (map[string]interface{}, error) {
	var espacio map[string]interface{}
	urlGetEspacio := beego.AppConfig.String("EspaciosAcademicosServices") + "espacio-academico?query=_id:" + espacioId + "&fields=nombre,grupo,espacio_academico_padre"
	err := request.GetJson(urlGetEspacio, &espacio)
	if err != nil {
		return nil, err
	}
	return espacio, nil
}
