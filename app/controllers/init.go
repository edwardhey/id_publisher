package controllers

import (

	//"encoding/json"
	"fmt"

	"github.com/revel/revel"
)

type ControllerError struct {
	Msg  string
	Code int
}

func (e ControllerError) Error() string {
	return fmt.Sprintf("%v: %v", e.Msg, e.Code)
}

var (
//Mc = app.Mc
)

const (
	ERROR_PARAMS         = iota + 1
	ERROR_DATA_NOTEXISTS = 404
	ERROR_SYSTEM         = 500
)

type BaseController struct {
	*revel.Controller
	Data  map[string]interface{}
	Code  int
	ERROR string
	//RD            Redispool
}

func (self *BaseController) beforeAction() revel.Result {
	self.Data = make(map[string]interface{})
	return nil
}

func (c *BaseController) afterAction() revel.Result {
	output := make(map[string]interface{})
	output["code"] = c.Code
	output["data"] = c.Data
	if c.Code > 0 {
		output["error"] = c.ERROR
	}
	return c.RenderJson(output)

}

func init() {
	revel.OnAppStart(func() {
		//Mq = app.Mq
		//Mc = app.Mc
		//TTT = app.TTT
		//fmt.Println("controllers", app.Mq)
	})
	//Mq.Publish("/test", 1, false, "aaa")

	revel.InterceptMethod((*BaseController).beforeAction, revel.BEFORE)
	revel.InterceptMethod((*BaseController).afterAction, revel.AFTER)
}
