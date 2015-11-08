package handler

type Context struct {
	Destination string `json:"destination"`
	Parameters  `json:"parameters"`
	Error       error `json:"error"`
	HandlerFunc
}

func NewContext() *Context {
	return &Context{Parameters: Parameters(make(map[string]interface{}))}
}

func (ctx *Context) Handle(data interface{}) {
	ctx.Error = ctx.HandlerFunc(data, ctx.Destination, ctx.Parameters)
}
