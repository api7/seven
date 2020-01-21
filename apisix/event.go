package apisix

// define event for workflow

type Event struct {
	Method string               // ADD UPDATE DELETE
	Kind   string               // route service upstream
	Func   func(...interface{}) // callback
}
