package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"],
        beego.ControllerComments{
            Method: "GetGruposEstudio",
            Router: "/getAll",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
