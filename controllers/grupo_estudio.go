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
	c.Mapping("GetGruposEstudioSegunHorarioYSemestre", c.GetGruposEstudioSegunHorarioYSemestre)
	c.Mapping("DeleteGrupoEstudio", c.DeleteGrupoEstudio)
	c.Mapping("PostGrupoEstudio", c.PostGrupoEstudio)
	c.Mapping("PutGrupoEstudio", c.PutGrupoEstudio)
}

// @Title getGrupoEstudioSegunHorarioYSemestre
// @Description get grupos de estudio segun horario y semestre
// @Param	horario-id	query	string	true	"Se recibe parametro: id del horario"
// @Param	semestre-id	query	string	true	"Se recibe parametro: id del semestre"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *GrupoEstudioController) GetGruposEstudioSegunHorarioYSemestre() {
	defer errorhandler.HandlePanic(&c.Controller)

	horarioId := c.GetString("horario-id")
	semestreId := c.GetString("semestre-id")

	respuesta := services.GetGruposEstudioSegunHorarioYSemestre(horarioId, semestreId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// @Title deleteGrupoEstudio
// @Description delete grupo estudio
// @Param   id      path    string  true        "grupo de estudio id"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *GrupoEstudioController) DeleteGrupoEstudio() {
	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudioId := c.Ctx.Input.Param(":id")

	respuesta := services.DeleteGrupoEstudio(grupoEstudioId)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// @Title PostGrupoEstudio
// @Description create grupo de estudio
// @Param   body        body    {}  true		"body"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *GrupoEstudioController) PostGrupoEstudio() {

	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudio := c.Ctx.Input.RequestBody

	respuesta := services.CreateGrupoEstudio(grupoEstudio)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// @Title PutGrupoEstudio
// @Description Update grupo de estudio
// @Param   id      path    string  true        "grupo de estudio id"
// @Param   body        body    {}  true		"body"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *GrupoEstudioController) PutGrupoEstudio() {

	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudio := c.Ctx.Input.RequestBody
	grupoEstudioId := c.Ctx.Input.Param(":id")

	respuesta := services.UpdateGrupoEstudio(grupoEstudioId, grupoEstudio)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}
