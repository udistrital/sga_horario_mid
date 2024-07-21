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
	c.Mapping("GetGruposEstudioSegunHorarioSemestre", c.GetGruposEstudioSegunHorarioSemestre)
}

// @Title getGrupoEstudioSegunHorarioSemestre
// @Description get grupos de estudio segun horario semestre
// @Param	horario-semestre-id	query	string	true	"Se recibe parametro: id del horario semestre"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *GrupoEstudioController) GetGruposEstudioSegunHorarioSemestre() {
	defer errorhandler.HandlePanic(&c.Controller)

	horarioSemestreId := c.GetString("horario-semestre-id")

	respuesta := services.GetGruposEstudioSegunHorarioSemestre(horarioSemestreId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
