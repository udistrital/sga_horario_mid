package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// Operations about GruposEstudio
type ColocacionEspacioAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ColocacionEspacioAcademicoController) URLMapping() {
	c.Mapping("GetColocacionesSegunGrupoEstudioYPeriodo", c.GetColocacionesSegunGrupoEstudioYPeriodo)
}

// @Title GetColocacionesSegunGrupoEstudioYPeriodo
// @Description get colocaciones de espacios academicos segun id de grupo estudio y id del periodo
// @Param	grupo-estudio-id	query	string	false	"Se recibe parametro: id del grupo estudio"
// @Param	periodo-id	query	string	false	"Se recibe parametro: id del periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *ColocacionEspacioAcademicoController) GetColocacionesSegunGrupoEstudioYPeriodo() {
	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudioId := c.GetString("grupo-estudio-id")
	periodoId := c.GetString("periodo-id")

	respuesta := services.GetColocacionesSegunGrupoEstudioYPeriodo(grupoEstudioId, periodoId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
