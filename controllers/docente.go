package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// Operations about GruposEstudio
type DocenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *DocenteController) URLMapping() {
	c.Mapping("GetDocente", c.GetDocente)
}

// @Title getDocenteYVincuaciones
// @Description get docente y sus vinculaciones por documento
// @Param	documento	query	string	false	"Se recibe parametro: documento del docente"
// @Success 200 {}
// @Failure 403 body is empty
// @router /vinculaciones [get]
func (c *DocenteController) GetDocente() {
	defer errorhandler.HandlePanic(&c.Controller)

	documento := c.GetString("documento")

	respuesta := services.GetDocenteYVinculacionesPorDocumento(documento)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
