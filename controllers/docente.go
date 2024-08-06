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
	c.Mapping("GetDocenteYVincuaciones", c.GetDocenteYVincuaciones)
	c.Mapping("GetPreasignacionesSegunDocenteYPeriodo", c.GetPreasignacionesSegunDocenteYPeriodo)
}

// @Title GetDocenteYVincuaciones
// @Description get docente y sus vinculaciones por documento
// @Param	documento	query	string	false	"Se recibe parametro: documento del docente"
// @Success 200 {}
// @Failure 403 body is empty
// @router /vinculaciones [get]
func (c *DocenteController) GetDocenteYVincuaciones() {
	defer errorhandler.HandlePanic(&c.Controller)

	documento := c.GetString("documento")

	respuesta := services.GetDocenteYVinculacionesPorDocumento(documento)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// @Title getPreasignacionesSegunDocenteYPeriodo
// @Description get preasignaciones de docente segun el periodo
// @Param	docente-id	query	string	false	"Se recibe parametro: id del docente"
// @Param	periodo-id	query	string	false	"Se recibe parametro: ide del periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /pre-asignacion [get]
func (c *DocenteController) GetPreasignacionesSegunDocenteYPeriodo() {
	defer errorhandler.HandlePanic(&c.Controller)

	docenteId := c.GetString("docente-id")
	periodoId := c.GetString("periodo-id")

	respuesta := services.GetPreasignacionesSegunDocenteYPeriodo(docenteId, periodoId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
