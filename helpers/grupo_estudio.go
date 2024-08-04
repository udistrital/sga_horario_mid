package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func ObtenerEspacioAcademicoSegunId(espacioId string) (map[string]interface{}, error) {
	var resultado map[string]interface{}
	urlGetEspacio := beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico?query=_id:" + espacioId + "&fields=nombre,grupo,espacio_academico_padre"

	fmt.Println(urlGetEspacio)
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
