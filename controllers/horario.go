package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// Operations about GruposEstudio
type HorarioController struct {
	beego.Controller
}

// URLMapping ...
func (c *HorarioController) URLMapping() {
	c.Mapping("GetActividadesParaHorarioYPlanDocente", c.GetActividadesParaHorarioYPlanDocente)
}

// @Title GetActividadesParaHorarioYPlanDocente
// @Description Obtener si hay eventos de calendario para hacer acciones con el modulo de horario
// @Param	periodo-id		query	string	false	"Se recibe parametro: id del periodo"
// @Param	nivel-id		query	string	false	"Se recibe parametro: id del nivel"
// @Param	dependencia-id	query	string	false	"Se recibe parametro: id del dependecia"
// @Success 200 {}
// @Failure 403 body is empty
// @router /calendario [get]
func (c *HorarioController) GetActividadesParaHorarioYPlanDocente() {
	defer errorhandler.HandlePanic(&c.Controller)

	periodoId := c.GetString("periodo-id")
	nivelId := c.GetString("nivel-id")
	dependeciaId := c.GetString("dependencia-id")

	respuesta := services.GetActividadesParaHorarioYPlanDocente(periodoId, nivelId, dependeciaId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
