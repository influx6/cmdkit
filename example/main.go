package main

import "fmt"
import "github.com/gokit/cmdkit"

func main() {
	cmdkit.Run("example", cmdkit.Flags(
		cmdkit.IntFlag(cmdkit.FlagName("age")),
		cmdkit.StringFlag(cmdkit.FlagName("name")),
	), cmdkit.Commands(
		cmdkit.Cmd(
			"add",
			cmdkit.Desc("displays a add message"),
			cmdkit.WithAction(func(ctx cmdkit.Context) error {
				fmt.Printf("Welcome to add: %q -> %d \n", ctx.String("name"), ctx.Int("age"))
				return nil
			})),
	))
}
