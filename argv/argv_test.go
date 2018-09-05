package argv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gokit/cmdkit/argv"
)

func TestParseArgsWithNoCommand(t *testing.T) {
	arg, err := argv.Parse("example --rack=20 --dirs=[drum flag kick] push git@ghu.com/fla.git")
	assert.NoError(t, err)
	assert.NotNil(t, arg.Sub)
	assert.NotEmpty(t, arg.Pairs)
	assert.Contains(t, arg.Pairs, "rack")
	assert.Contains(t, arg.Pairs, "dirs")
	assert.Equal(t, "push", arg.Sub.Name)
	assert.Equal(t, "git@ghu.com/fla.git", arg.Sub.Text)
}

func TestParseArgs(t *testing.T) {
	arg, err := argv.Parse("rocket  --name=wallet -rack=ball -h")
	assert.NoError(t, err)
	assert.Nil(t, arg.Sub)
	assert.NotEmpty(t, arg.Pairs)
	assert.Equal(t, "rocket", arg.Name)
	assert.Contains(t, arg.Pairs, "h")
	assert.Contains(t, arg.Pairs, "name")
	assert.Contains(t, arg.Pairs, "rack")
	assert.Contains(t, arg.Pairs["name"], "wallet")
	assert.Contains(t, arg.Pairs["rack"], "ball")
}

func TestParseArgsWithList(t *testing.T) {
	arg, err := argv.Parse("runket -w=323 -j danger ricker --name=[ bog willow crack ] -rack=ball -h renditions recka")
	assert.NoError(t, err)
	assert.NotNil(t, arg.Sub)
	assert.NotEmpty(t, arg.Pairs)
	assert.Equal(t, "runket", arg.Name)
	assert.Contains(t, arg.Pairs, "w")
	assert.Contains(t, arg.Pairs, "j")

	assert.Empty(t, arg.Sub.Pairs)
	assert.Equal(t, "danger", arg.Sub.Name)

	assert.NotNil(t, arg.Sub.Sub.Sub)
	assert.NotEmpty(t, arg.Sub.Sub.Pairs)
	assert.Contains(t, arg.Sub.Sub.Pairs, "h")
	assert.Contains(t, arg.Sub.Sub.Pairs, "name")
	assert.Equal(t, "ricker", arg.Sub.Sub.Name)
	assert.Contains(t, arg.Sub.Sub.Pairs["name"], "bog")
	assert.Contains(t, arg.Sub.Sub.Pairs["name"], "willow")
	assert.Contains(t, arg.Sub.Sub.Pairs["name"], "crack")

	assert.Equal(t, "renditions", arg.Sub.Sub.Sub.Name)
	assert.Equal(t, "recka", arg.Sub.Sub.Sub.Text)
}
