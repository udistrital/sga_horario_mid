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
// @Description get grupos de estudio
// @Param	proyecto-academico	query	string	false	"Se recibe parametro: id del proyecto academico"
// @Param	plan-estudios		query	string	false	"Se recibe parametro: id del plan de estudios"
// @Param	semestre			query	string	false	"Se recibe parametro: id del semestre"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *GrupoEstudioController) GetGruposEstudio() {
	defer errorhandler.HandlePanic(&c.Controller)

	proyectoAcademicoId := c.GetString("proyecto-academico")
	planEstudiosId := c.GetString("plan-estudios")
	semestreId := c.GetString("semestre")

	respuesta := services.GetGruposEstudio(proyectoAcademicoId, planEstudiosId, semestreId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
