package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

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
