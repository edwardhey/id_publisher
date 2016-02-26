package controllers

import (
	"id_publisher/app"

	"github.com/revel/revel"
)

type App struct {
	BaseController
}

func (c *App) Index() revel.Result {

	var (
		t int
	)
	c.Params.Bind(&t, "t")
	var err error
	if t > 184 {
		c.ERROR = "t overflow!"
		c.Code = ERROR_PARAMS
		return nil
	} else if t == 0 {
		c.ERROR = "t must not be 0!"
		c.Code = ERROR_PARAMS
		return nil
	}

	orderIndex := app.GetByKey(t)
	/*
		orderIndex, err := getOrderIndex(key)
	*/

	if err != nil {
		c.ERROR = err.Error()
		c.Code = ERROR_SYSTEM
	} else {
		c.Data["oid"] = orderIndex
	}

	//app.Mc.Delete(key)
	return nil
}
