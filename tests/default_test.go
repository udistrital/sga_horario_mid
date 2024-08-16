package acceptance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/udistrital/sga_horario_mid/controllers"
)

var (
	resStatus   string
	resBody     []byte
	savepostres map[string]interface{}
	response    *httptest.ResponseRecorder
	debug       = true
)

func TestMain(m *testing.M) {
	opts := godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",      // Puedes cambiar el formato a json, cucumber, etc.
		Paths:  []string{"."}, // Directorio donde están los archivos .feature

	}

	status := godog.TestSuite{
		Name:                 "acceptance",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// Inicialización global antes de ejecutar la suite de pruebas
		log.Println("Starting test suite...")
		// Por ejemplo: Inicializar la base de datos, lanzar un contenedor, etc.
	})

	ctx.AfterSuite(func() {
		// Limpieza global después de ejecutar la suite de pruebas
		log.Println("Test suite finished. Cleaning up...")
		// Por ejemplo: Detener contenedores, limpiar base de datos, etc.
	})
}

// @getPages convierte en un tipo el json
func getPages(ruta string) []byte {
	raw, err := os.ReadFile(ruta)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return raw
}

// @iSendRequestToWhereBodyIsJson realiza la solicitud a la API
func iSendRequestToWhereBodyIsJson(method, endpoint, bodyreq string) error {

	if debug {
		fmt.Println("Step: iSendRequestToWhereBodyIsJson")
	}

	var url string
	baseURL := "http://:localhost:8080" + endpoint

	switch method {
	case "GET", "POST":
		url = baseURL

	case "PUT", "DELETE", "GETID":
		//str := strconv.FormatFloat(Id, 'f', 0, 64)
		//url = baseURL + "/" + str

		if method == "GETID" {
			method = "GET"
		}
	}

	if debug {
		fmt.Println("Test: " + method + " to " + url)
	}

	beego.BeeApp.Handlers.Add("/v1/docente", &controllers.DocenteController{}, "get:GetAll")

	pages := getPages(bodyreq)

	// Crear la solicitud usando httptest y la ruta en Beego
	req, err := http.NewRequest(method, url, bytes.NewBuffer(pages))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Usa el ResponseRecorder para capturar la respuesta
	response = httptest.NewRecorder()

	// Llama al handler correspondiente
	beego.BeeApp.Handlers.ServeHTTP(response, req)

	resStatus = response.Result().Status
	resBody = response.Body.Bytes()

	if method == "POST" && resStatus == "201 Created" {
		json.Unmarshal(resBody, &savepostres)
		//Id = savepostres["Id"].(float64)
	}

	return nil

}

// @theResponseCodeShouldBe valida el codigo de respuesta
func theResponseCodeShouldBe(arg1 string) error {
	if debug {
		fmt.Println("Step: theResponseCodeShouldBe")
	}

	if resStatus != arg1 {
		return fmt.Errorf("se esperaba el codigo de respuesta .. %s .. y se obtuvo el codigo de respuesta .. %s .. ", arg1, resStatus)
	}
	return nil
}

// @theResponseShouldMatchJson valida el JSON de respuesta
func theResponseShouldMatchJson(arg1 string) error {
	if debug {
		fmt.Println("Step: theResponseShouldMatchJson")
	}

	div := strings.Split(arg1, "/")
	pages := getPages(arg1)

	pages_s := string(pages)
	body_s := string(resBody)

	var data1, data2 interface{}

	if err := json.Unmarshal([]byte(pages_s), &data1); err != nil {
		fmt.Println("Error unmarshalling JSON1:", err)
		return err
	}

	if err := json.Unmarshal([]byte(body_s), &data2); err != nil {
		fmt.Println("Error unmarshalling JSON2:", err)
		return err
	}

	prefix := div[3]

	switch {
	case strings.HasPrefix(prefix, "V"):
		if sameStructure(data1, data2) {
			return nil
		} else {
			return fmt.Errorf("Errores: La estructura del objeto recibido no es la que se esperaba %s != %s", pages_s, body_s)
		}

	case strings.HasPrefix(prefix, "I"):
		areEqual, _ := AreEqualJSON(pages_s, body_s)
		if areEqual {
			return nil
		} else {
			return fmt.Errorf("Se esperaba el body de respuesta %s y se obtuvo %s", pages_s, resBody)
		}
	}

	return fmt.Errorf("Respuesta no validada")
}

// @AreEqualJSON comparar dos JSON si son iguales retorna true de lo contrario false
func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

// @extractKeysTypes Extraer las llaves de un json
func extractKeysTypes(data interface{}) map[string]reflect.Type {
	keysTypes := make(map[string]reflect.Type)
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Map {
		for _, key := range value.MapKeys() {
			val := value.MapIndex(key).Interface()
			if val == nil {
				keysTypes[key.String()] = nil
			} else if reflect.TypeOf(val).Kind() == reflect.Map {
				// Recursively check nested objects
				keysTypes[key.String()] = reflect.TypeOf(extractKeysTypes(val))
			} else {
				keysTypes[key.String()] = reflect.TypeOf(val)
			}
		}
	}
	return keysTypes
}

// @sameStructure comparar dos JSON si su estructura es igual retorna true de lo contrario false
func sameStructure(data1, data2 interface{}) bool {
	if data1 == nil || data2 == nil {
		return false
	}

	type1 := reflect.TypeOf(data1)
	type2 := reflect.TypeOf(data2)

	if type1.Kind() != type2.Kind() {
		return false
	}

	if type1.Kind() == reflect.Slice {
		v1 := reflect.ValueOf(data1)
		v2 := reflect.ValueOf(data2)
		if v1.Len() == 0 || v2.Len() == 0 {
			return false
		}
		return sameStructure(v1.Index(0).Interface(), v2.Index(0).Interface())
	} else if type1.Kind() == reflect.Map {
		keysTypes1 := extractKeysTypes(data1)
		keysTypes2 := extractKeysTypes(data2)
		return reflect.DeepEqual(keysTypes1, keysTypes2)
	}

	return false
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Registra los steps correspondientes a cada feature
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" where body is json "([^"]*)"$`, iSendRequestToWhereBodyIsJson)
	ctx.Step(`^the response code should be "([^"]*)"$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should match json "([^"]*)"$`, theResponseShouldMatchJson)

	// Uso de Before y After con el formato actualizado
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		log.Printf("Starting scenario: %s\n", sc.Name)
		// Realiza cualquier inicialización necesaria antes de cada escenario
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			log.Printf("Scenario failed: %s\n", sc.Name)
		}
		log.Printf("Finished scenario: %s\n", sc.Name)
		// Realiza cualquier limpieza necesaria después de cada escenario
		return ctx, nil
	})
}
