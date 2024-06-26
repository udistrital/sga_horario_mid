package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// Operations about GruposEstudio
type GrupoEstudioController struct {
	beego.Controller
}

// URLMapping ...
func (c *GrupoEstudioController) URLMapping() {
	c.Mapping("GetGruposEstudio", c.GetGruposEstudio)
}

// @Title getGrupoEstudio
// @Description get events
// @Success 200 {}
// @Failure 403 body is empty
// @router /getAll [get]
func (c *GrupoEstudioController) GetGruposEstudio() {

	defer errorhandler.HandlePanic(&c.Controller)

	respuesta := services.GetGruposEstudio()

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
