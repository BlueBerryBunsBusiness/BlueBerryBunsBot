package cmd

import (
	"log"

	"github.com/Necroforger/dgrouter/exrouter"
)

// CommandController is used to batch manage routed commands
type CommandController struct {
	Router *exrouter.Route
}

// New is a CommandController constuctor. Use this instead of creating an instance directly.
func New(r *exrouter.Route) *CommandController {
	c := &CommandController{r}
	c.Router.Default = c.Router.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range c.Router.Routes {
			text += v.Name + " : \t" + v.Description + "\n"
		}
		ctx.Reply("```" + text + "```")
	}).Desc("prints this help menu")
	return c
}

// Add adds set and get for testing management of context variables
func Add(c *CommandController) {
	c.Router.On("set", func(ctx *exrouter.Context) {
		ctx.Set(ctx.Args.Get(1), ctx.Args.Get(2))
		reply := "set called"
		reply += "\n Variable is: " + ctx.Args.Get(1)
		reply += "\n Value is: " + ctx.Args.Get(2)
		ctx.Reply(reply)
	}).Desc("Sets a variable to a value")

	c.Router.On("get", func(ctx *exrouter.Context) {
		val := ctx.Get(ctx.Args.Get(1))
		log.Printf(val.(string))
		reply := "get called"
		reply += "\nVarible is: " + ctx.Args.Get(1)
		ctx.Reply(reply)
	}).Desc("Gets the value of a variable")
}
