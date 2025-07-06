package main

import (
	"os"

	"codeberg.org/anaseto/goal"
	"codeberg.org/anaseto/goal/cmd"
	"codeberg.org/anaseto/goal/help"
	gos "codeberg.org/anaseto/goal/os"
)

func main() {
	ctx := goal.NewContext() // new evaluation context for Goal code
	ctx.Log = os.Stderr      // configure logging with \x to stderr
	gos.Import(ctx, "")      // register all IO/OS primitives with prefix ""
	ctx.AssignGlobal("http.get", ctx.RegisterMonad("http.get", vfHttpGet))
	ctx.AssignGlobal("http.serve", ctx.RegisterMonad("http.serve", vfHttpServe))
	ctx.AssignGlobal("http.register", ctx.RegisterMonad("http.register", vfHttpRegister))
	ctx.AssignGlobal("db.save", ctx.RegisterDyad("db.save", vfDbSave))
	ctx.AssignGlobal("db.get", ctx.RegisterDyad("db.get", vfDbGet))
	ctx.AssignGlobal("date.fs", ctx.RegisterDyad("date.fs", vfTimeStrDate))
	ctx.AssignGlobal("date.sf", ctx.RegisterDyad("date.sf", vfTimeDateStr))
	cmd.Exit(cmd.Run(ctx, cmd.Config{Help: help.HelpFunc(), Man: "goal"}))
}
