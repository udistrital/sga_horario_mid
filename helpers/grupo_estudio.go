package helpers

import (
	// "fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func obtenerEspacioAcademicoSegunId(espacioId string) (map[string]interface{}, error) {
	var espacio map[string]interface{}
	urlGetEspacio := beego.AppConfig.String("EspaciosAcademicosServices") + "espacio-academico/" + espacioId
	err := request.GetJson(urlGetEspacio, &espacio)
	if err != nil {
		return nil, err
	}
	return espacio, nil
}
