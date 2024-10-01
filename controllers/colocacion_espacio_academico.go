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
	c.Mapping("GetSobreposicionColocacion", c.GetSobreposicionColocacion)
	c.Mapping("GetSobreposicionColocacion", c.GetColocacionInfoAdicional)
	c.Mapping("DeleteColocacionEspacioAcademico", c.DeleteColocacionEspacioAcademico)
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

// @Title GetSobreposicionColocacion
// @Description get si hay una colocacion puesta en donde se quiere poner otra
// @Param	colocacion-id	query	string	false	"Se recibe parametro: id de la colocacion"
// @Param	periodo-id	query	string	false	"Se recibe parametro: id del periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /sobreposicion [get]
func (c *ColocacionEspacioAcademicoController) GetSobreposicionColocacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	colocacionId := c.GetString("colocacion-id")
	periodoId := c.GetString("periodo-id")

	respuesta := services.GetSobreposicionColocacion(colocacionId, periodoId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// @Title GetColocacionInfoAdicional
// @Description obtiene la colocacion con mas datos sobre la misma
// @Param   id      path    string  true "id de la colocacion"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info-adicional/:id [get]
func (c *ColocacionEspacioAcademicoController) GetColocacionInfoAdicional() {
	defer errorhandler.HandlePanic(&c.Controller)

	colocacionId := c.Ctx.Input.Param(":id")

	respuesta := services.GetColocacionInfoAdicional(colocacionId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// @Title deleteGrupoEstudio
// @Description delete colocacion espacio academico (delete carga plan si tiene)
// @Param   id      path    string  true        "colocacion espacio academico id"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *ColocacionEspacioAcademicoController) DeleteColocacionEspacioAcademico() {
	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudioId := c.Ctx.Input.Param(":id")

	respuesta := services.DeleteColocacionEspacioAcademico(grupoEstudioId)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// // @Title PostCopiarColocaciones
// // @Description copia las colocaciones de un un grupo de estudio a otro
// // @Param   body        body    {}  true		"body"
// // @Success 200 {}
// // @Failure 400 the request contains incorrect syntax
// // @router /copiar [post]
// func (c *ColocacionEspacioAcademicoController) PostCopiarColocaciones() {

// 	defer errorhandler.HandlePanic(&c.Controller)

// 	infoParaCopiado := c.Ctx.Input.RequestBody

// 	respuesta := services.CreateHorarioCopia(infoParaCopiado)

// 	c.Ctx.Output.SetStatus(respuesta.Status)

// 	c.Data["json"] = respuesta
// 	c.ServeJSON()
// }
