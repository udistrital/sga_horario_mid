package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// Operations about GruposEstudio
type EspacioFisicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *EspacioFisicoController) URLMapping() {
	c.Mapping("GetOcupadosDeEspacioFisicoSegunPeriodo", c.GetEspaciosOCupadoSegunPeriodo)
}

// @Title GetEspaciosOCupadoSegunPeriodo
// @Description get espacios ocupacdos de un espacio fisico segun el periodo
// @Param	espacio-fisico-id	query	string	false	"Se recibe parametro: id del espacio fisico"
// @Param	periodo-id			query	string	false	"Se recibe parametro: id del periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /ocupados [get]
func (c *EspacioFisicoController) GetEspaciosOCupadoSegunPeriodo() {
	defer errorhandler.HandlePanic(&c.Controller)

	espacioFisicoId := c.GetString("espacio-fisico-id")
	periodoId := c.GetString("periodo-id")

	respuesta := services.GetEspaciosFisicosOcupadosSegunPeriodo(espacioFisicoId, periodoId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
