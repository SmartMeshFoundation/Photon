package v1

import (
	"strconv"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Abort(code int) {
	this.Controller.Abort(strconv.Itoa(code))
}
