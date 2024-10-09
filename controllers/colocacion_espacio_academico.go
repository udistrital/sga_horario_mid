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
	c.Mapping("GetColocacionesDeGrupoEstudio", c.GetColocacionesDeGrupoEstudio)
	c.Mapping("GetColocacionInfoAdicional", c.GetColocacionInfoAdicional)
	c.Mapping("DeleteColocacionEspacioAcademico", c.DeleteColocacionEspacioAcademico)
	c.Mapping("GetColocacionesGrupoSinDetalles", c.GetColocacionesGrupoSinDetalles)
	c.Mapping("GetSobreposicionEnGrupoEstudio", c.GetSobreposicionEnGrupoEstudio)
	c.Mapping("GetSobreposicionEspacioFisico", c.GetSobreposicionEspacioFisico)
	c.Mapping("PostCopiarColocacionesAGrupoEstudio", c.PostCopiarColocacionesAGrupoEstudio)
}

// @Title GetColocacionesDeGrupoEstudio
// @Description get colocaciones de espacios academicos segun id de grupo estudio y id del periodo
// @Param	grupo-estudio-id	query	string	false	"Se recibe parametro: id del grupo estudio"
// @Param	periodo-id	query	string	false	"Se recibe parametro: id del periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [get]
func (c *ColocacionEspacioAcademicoController) GetColocacionesDeGrupoEstudio() {
	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudioId := c.GetString("grupo-estudio-id")
	periodoId := c.GetString("periodo-id")

	respuesta := services.GetColocacionesDeGrupoEstudio(grupoEstudioId, periodoId)

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

// @Title GetColocacionesSinInfoAdicionalDeGrupoEstudio
// @Description get colocaciones de espacios academicos segun id de grupo estudio y id del periodo
// @Param	grupo-estudio-id	query	string	false	"Se recibe parametro: id del grupo estudio"
// @Param	periodo-id	query	string	false	"Se recibe parametro: id del periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /sin-detalles [get]
func (c *ColocacionEspacioAcademicoController) GetColocacionesGrupoSinDetalles() {
	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudioId := c.GetString("grupo-estudio-id")
	periodoId := c.GetString("periodo-id")

	respuesta := services.GetColocacionesGrupoSinDetalles(grupoEstudioId, periodoId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// @Title GetSobreposicionEnGrupoEstudio
// @Description Obtiene si una colocación se sobrepone a alguna colocación del grupo de estudio
// @Param	grupo-estudio-id	query	string	false	"Se recibe parametro: id del grupo de estudio"
// @Param	periodo-id			query	string	false	"Se recibe parametro: id del periodo del grupo de estudio"
// @Param	colocacion-id		query	string	false	"Se recibe parametro: id de la colocacion"
// @Success 200 {}
// @Failure 403 body is empty
// @router /grupo-estudio/sobreposicion [get]
func (c *ColocacionEspacioAcademicoController) GetSobreposicionEnGrupoEstudio() {
	defer errorhandler.HandlePanic(&c.Controller)

	grupoEstudioId := c.GetString("grupo-estudio-id")
	periodoId := c.GetString("periodo-id")
	colocacionId := c.GetString("colocacion-id")

	respuesta := services.GetSobreposicionEnGrupoEstudio(grupoEstudioId, periodoId, colocacionId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// @Title GetSobreposicionEspacioFisico
// @Description get si hay una colocacion puesta en donde se quiere poner otra segun el espacio fisico
// @Param	colocacion-id	query	string	false	"Se recibe parametro: id de la colocacion"
// @Param	periodo-id	query	string	false	"Se recibe parametro: id del periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /espacio-fisico/sobreposicion [get]
func (c *ColocacionEspacioAcademicoController) GetSobreposicionEspacioFisico() {
	defer errorhandler.HandlePanic(&c.Controller)

	colocacionId := c.GetString("colocacion-id")
	periodoId := c.GetString("periodo-id")

	respuesta := services.GetSobreposicionEspacioFisico(colocacionId, periodoId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// @Title PostCopiarColocacionesAGrupoEstudio
// @Description copia las colocaciones de un un grupo de estudio a otro
// @Param   body        body    {}  true		"body"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /copiar [post]
func (c *ColocacionEspacioAcademicoController) PostCopiarColocacionesAGrupoEstudio() {

	defer errorhandler.HandlePanic(&c.Controller)

	infoParaCopiado := c.Ctx.Input.RequestBody

	respuesta := services.CopiarColocacionesAGrupoEstudio(infoParaCopiado)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}
