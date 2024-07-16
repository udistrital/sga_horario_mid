package helpers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func GetVinculacionesDeDocente(docenteId string) ([]map[string]interface{}, error) {
	urlVinculacion := beego.AppConfig.String("TercerosService") + "vinculacion?query=TerceroPrincipalId.Id:" + docenteId + "&fields=Id,TipoVinculacionId"
	var vinculacionesId []map[string]interface{}
	if err := request.GetJson(urlVinculacion, &vinculacionesId); err != nil {
		return nil, fmt.Errorf("Error en el servicio de terceros")
	}

	if len(vinculacionesId) == 0 {
		return nil, nil // No hay vinculaciones
	}

	idsVinculacionTipoDocente := map[int64]struct{}{293: {}, 294: {}, 296: {}, 297: {}, 298: {}, 299: {}}

	var vinculaciones []map[string]interface{}
	urlParametroBase := beego.AppConfig.String("ParametroService") + "parametro/"
	for _, vinculacion := range vinculacionesId {
		tipoVinculacionId := int64(vinculacion["TipoVinculacionId"].(float64))
		if _, ok := idsVinculacionTipoDocente[tipoVinculacionId]; ok {
			var resVinculacion map[string]interface{}
			if err := request.GetJson(urlParametroBase+fmt.Sprintf("%v", tipoVinculacionId), &resVinculacion); err == nil {
				res := resVinculacion["Data"].(map[string]interface{})
				vinculaciones = append(vinculaciones, map[string]interface{}{
					"Id":     res["Id"],
					"Nombre": strings.ToUpper(string(res["Nombre"].(string)[0])) + strings.ToLower(res["Nombre"].(string)[1:]),
				})
			}
		}
	}

	if len(vinculaciones) == 0 {
		return nil, nil // No hay vinculaciones de tipo docente
	}

	return vinculaciones, nil
}
