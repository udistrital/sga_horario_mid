package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"],
        beego.ControllerComments{
            Method: "GetColocacionesDeGrupoEstudio",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"],
        beego.ControllerComments{
            Method: "DeleteColocacionEspacioAcademico",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"],
        beego.ControllerComments{
            Method: "GetSobreposicionEnGrupoEstudio",
            Router: "/grupo-estudio/sobreposicion",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"],
        beego.ControllerComments{
            Method: "GetColocacionInfoAdicional",
            Router: "/info-adicional/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"],
        beego.ControllerComments{
            Method: "GetColocacionesGrupoSinDetalles",
            Router: "/sin-detalles",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:ColocacionEspacioAcademicoController"],
        beego.ControllerComments{
            Method: "GetSobreposicionColocacion",
            Router: "/sobreposicion",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:DocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:DocenteController"],
        beego.ControllerComments{
            Method: "GetPreasignacionesSegunDocenteYPeriodo",
            Router: "/pre-asignacion",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:DocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:DocenteController"],
        beego.ControllerComments{
            Method: "GetDocenteYVincuaciones",
            Router: "/vinculaciones",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:EspacioFisicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:EspacioFisicoController"],
        beego.ControllerComments{
            Method: "GetEspaciosOCupadoSegunPeriodo",
            Router: "/ocupados",
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

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"],
        beego.ControllerComments{
            Method: "PostGrupoEstudio",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"],
        beego.ControllerComments{
            Method: "DeleteGrupoEstudio",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:GrupoEstudioController"],
        beego.ControllerComments{
            Method: "PutGrupoEstudio",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:HorarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_horario_mid/controllers:HorarioController"],
        beego.ControllerComments{
            Method: "GetActividadesParaHorarioYPlanDocente",
            Router: "/calendario",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
