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
	c.Mapping("GetColocacionesSegunHorarioSemestre", c.GetColocacionesSegunHorarioSemestre)
}

// @Title getColocacionesSegunHorarioSemestre
// @Description get colocaciones de espacios academicos segun id de horario semestre
// @Param	horario-semestre-id	query	string	false	"Se recibe parametro: id del horario semestre"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *ColocacionEspacioAcademicoController) GetColocacionesSegunHorarioSemestre() {
	defer errorhandler.HandlePanic(&c.Controller)

	horarioSemestreId := c.GetString("horario-semestre-id")

	respuesta := services.GetColocacionesSegunHorarioSemestre(horarioSemestreId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
