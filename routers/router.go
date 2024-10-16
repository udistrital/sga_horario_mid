// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_horario_mid/controllers"
	"github.com/udistrital/utils_oas/errorhandler"
)

func init() {
	beego.ErrorController(&errorhandler.ErrorHandlerController{})

	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/colocacion-espacio-academico",
			beego.NSInclude(
				&controllers.ColocacionEspacioAcademicoController{},
			),
		),
		beego.NSNamespace("/espacio-fisico",
			beego.NSInclude(
				&controllers.EspacioFisicoController{},
			),
		),
		beego.NSNamespace("/grupo-estudio",
			beego.NSInclude(
				&controllers.GrupoEstudioController{},
			),
		),
		beego.NSNamespace("/horario",
			beego.NSInclude(
				&controllers.HorarioController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
