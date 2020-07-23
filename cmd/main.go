package cmd

import (
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
}
