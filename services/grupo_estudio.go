package services

import (
	// "encoding/json"
	// "fmt"
	// "strconv"

	// "github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetGruposEstudio() (APIResponseDTO requestresponse.APIResponse) {
	return requestresponse.APIResponseDTO(true, 200, nil, nil)
}
