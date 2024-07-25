package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"],
        beego.ControllerComments{
            Method: "GetColocacionesSegunGrupoEstudio",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:DocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:DocenteController"],
        beego.ControllerComments{
            Method: "GetDocente",
            Router: "/vinculaciones",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"],
        beego.ControllerComments{
            Method: "GetGruposEstudioSegunHorarioYSemestre",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
