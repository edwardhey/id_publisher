package app

import (
	"fmt"
	"sync"
	"time"

	"strconv"

	"github.com/revel/revel"
)

var mutex sync.Mutex
var m map[int]map[string]int = make(map[int]map[string]int, 3)
var current int
var timeout time.Duration = 3

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	current = int(time.Now().Unix())
	m[current] = make(map[string]int, 184)
	// register startup functions with OnAppStart
	// ( order dependent )
	//revel.OnAppStart(InitMc)

	time.AfterFunc(timeout*time.Second, FlushAll)
	time.AfterFunc(1*time.Second, ChangeCurrentTime)
	//revel.InterceptMethod((*BaseController).beforeAction, revel.BEFORE)
	//revel.InterceptMethod((*BaseController).afterAction, revel.AFTER)
	// revel.OnAppStart(FillCache)
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func GetByKey(t int) string {
	mutex.Lock()
	defer mutex.Unlock()

	key := fmt.Sprintf("%d_%d", t, current)
	n, ok := m[current][key]
	if !ok {
		if _, ok := m[current]; !ok {
			m[current] = make(map[string]int, 184)
		}
		m[current][key] = 0
		n = 0
	}
	m[current][key] += 1
	tm := time.Unix(int64(current), 0)
	date, _ := strconv.Atoi(tm.Format("060102"))
	return fmt.Sprintf("%d%d%06d%05d%05d", t, 0, date, current%86400, n)
}

func FlushAll() {
	mutex.Lock()
	defer mutex.Unlock()
	//fmt.Println("flushall!")
	for t, _ := range m {
		if t != current {
			delete(m, t)
			//fmt.Printf("Delete %d:%d\r\n", current, t)
		}

	}
	time.AfterFunc(timeout*time.Second, FlushAll)
}

func ChangeCurrentTime() {
	//.Println("chanrge current time!")
	current = int(time.Now().Unix())
	time.AfterFunc(1*time.Second, ChangeCurrentTime)
}
