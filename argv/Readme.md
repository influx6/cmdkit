Argv
----------
Argv provides a commandline argument parser using a subcommand approach where users can layer subcommands
as single line calls. Be aware that the position of the flags (i.e those with `--` or `-` prefix) matters.

## Install

```
go get -u github.com/gokit/cmdkit
```

## Examples

```
$ ./example --rack=20 --dirs=[drum flag kick] push git@ghu.com/fla.git
```

Produces a `argv.Argv`, if marshalled in json:

```json
{
	"Name": "example",
	"Sub": {
		"Name": "push",
		"Sub": null,
		"Text": "git@ghu.com/fla.git",
		"Flags": null,
		"Pairs": {}
	},
	"Text": "push git@ghu.com/fla.git",
	"Flags": null,
	"Pairs": {
		"dirs": [
			"drum",
			"flag",
			"kick"
		],
		"rack": [
			"20"
		]
	}
}
```


