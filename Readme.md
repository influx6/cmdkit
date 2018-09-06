Cmdkit
--------------
Cmdkit provides a simple commandline library with zero outside dependency, using 
internal libraries as much as possible. It's geared towards a simple API.

If you need broader functionality then I would encourage you to use more robust libraries.

## Install

```
go get -u github.com/gokit/cmdkit
```


## Usage

```go
import "fmt"
import "github.com/gokit/cmdkit"

func main() {
	cmdkit.Run(
		"example", 
		cmdkit.Flags(
			cmdkit.IntFlag(cmdkit.FlagName("age")),
			cmdkit.StringFlag(cmdkit.FlagName("name")),
		), 
		cmdkit.Commands(
			cmdkit.Cmd(
				"add",
				cmdkit.Desc("displays a add message"),
				cmdkit.WithAction(func(ctx cmdkit.Context) error {
					fmt.Printf("Welcome to add: %q -> %d \n", ctx.String("name"), ctx.Int("age"))
					return nil
				}),
				cmdkit.SubCommands(
					cmdkit.Cmd(
						"broc",
						cmdkit.Desc("displays a broc message"),
						cmdkit.WithAction(func(ctx cmdkit.Context) error {
							fmt.Printf("Welcome to bro adder: %q -> %d \n", ctx.String("name"), ctx.Int("age"))
							return nil
						}),
					),
				),
			),
		))
}
```

The code above produces the ff:

````bash
Usage: example [flags] [command] 

⡿ COMMANDS:

	⠙ add        displays a add message

⡿ HELP:

	Run [command] --help to print this message
	Run example --flags to print all flags of all commands.

⡿ Flags:
	
	⠙ age  (int)          
	
	⠙ name  (string)          
	
	⠙ help  (bool)          
	
	⠙ flags  (bool)          
	
	⠙ timeout  (time.Duration)          
	
````