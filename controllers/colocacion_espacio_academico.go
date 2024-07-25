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
	c.Mapping("GetColocacionesSegunGrupoEstudio", c.GetColocacionesSegunGrupoEstudio)
}

// @Title GetColocacionesSegunGrupoEstudio
// @Description get colocaciones de espacios academicos segun id de grupo estudio
// @Param	grupo-estudio-id	query	string	false	"Se recibe parametro: id del grupo estudio"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *ColocacionEspacioAcademicoController) GetColocacionesSegunGrupoEstudio() {
	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudioId := c.GetString("grupo-estudio-id")

	respuesta := services.GetColocacionesSegunGrupoEstudio(grupoEstudioId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
