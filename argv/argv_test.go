package argv_test

import (
	"reflect"
	"testing"

	"github.com/gokit/cmdkit/argv"
)

func TestParseArgsWithNoCommand(t *testing.T) {
	arg, err := argv.Parse("example --rack=20 --dirs=[drum flag kick] push git@ghu.com/fla.git")
	noError(t, err)
	notNil(t, arg.Sub)
	notEmpty(t, arg.Pairs)
	contains(t, arg.Pairs, "rack")
	contains(t, arg.Pairs, "dirs")
	equal(t, "push", arg.Sub.Name)
	equal(t, "git@ghu.com/fla.git", arg.Sub.Text)
}

func TestParseArgs(t *testing.T) {
	arg, err := argv.Parse("rocket  --name=wallet -rack=ball -h")
	noError(t, err)
	isNil(t, arg.Sub)
	notEmpty(t, arg.Pairs)
	equal(t, "rocket", arg.Name)
	contains(t, arg.Pairs, "h")
	contains(t, arg.Pairs, "name")
	contains(t, arg.Pairs, "rack")
	contains(t, arg.Pairs["name"], "wallet")
	contains(t, arg.Pairs["rack"], "ball")
}

func TestParseArgsWithList(t *testing.T) {
	arg, err := argv.Parse("runket -w=323 -j danger ricker --name=[ bog willow crack ] -rack=ball -h renditions recka")
	noError(t, err)
	notNil(t, arg.Sub)
	notEmpty(t, arg.Pairs)
	equal(t, "runket", arg.Name)
	contains(t, arg.Pairs, "w")
	contains(t, arg.Pairs, "j")

	isEmpty(t, arg.Sub.Pairs)
	equal(t, "danger", arg.Sub.Name)

	notNil(t, arg.Sub.Sub.Sub)
	notEmpty(t, arg.Sub.Sub.Pairs)
	contains(t, arg.Sub.Sub.Pairs, "h")
	contains(t, arg.Sub.Sub.Pairs, "name")
	equal(t, "ricker", arg.Sub.Sub.Name)
	contains(t, arg.Sub.Sub.Pairs["name"], "bog")
	contains(t, arg.Sub.Sub.Pairs["name"], "willow")
	contains(t, arg.Sub.Sub.Pairs["name"], "crack")

	equal(t, "renditions", arg.Sub.Sub.Sub.Name)
	equal(t, "recka", arg.Sub.Sub.Sub.Text)
}

func noError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error occured: %#v\n", err)
	}
}

func equal(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Logf("Expected: %#v\n", expected)
		t.Logf("Actual: %#v\n", actual)
		t.Fatalf("Actual is not equal to expected\n")
	}
}

func contains(t *testing.T, item interface{}, k string) {
	switch tm := item.(type) {
	case []string:
		for _, elem := range tm {
			if elem == k {
				return
			}
		}
		t.Fatalf("Unable to find key %q in %#v\n", k, item)
	case map[string][]string:
		if _, ok := tm[k]; !ok {
			t.Fatalf("Unable to find key %q in %#v\n", k, item)
		}
	}
}

func isNil(t *testing.T, item interface{}) {
	tm := reflect.ValueOf(item)
	if !tm.IsNil() {
		t.Fatalf("Expected value is not nil: %#v\n", item)
	}
}

func notNil(t *testing.T, item interface{}) {
	tm := reflect.ValueOf(item)
	if tm.IsNil() {
		t.Fatal("Expected nil value")
	}
}

func isEmpty(t *testing.T, item interface{}) {
	ft := reflect.ValueOf(item)
	switch ft.Kind() {
	case reflect.Map:
		if ft.Len() != 0 {
			t.Fatal("Value is  empty")
		}
	case reflect.Chan:
		if ft.Len() != 0 {
			t.Fatal("Value is  empty")
		}
	case reflect.Array:
		if ft.Len() != 0 {
			t.Fatal("Value is  empty")
		}
	case reflect.Slice:
		if ft.Len() != 0 {
			t.Fatal("Value is  empty")
		}
	case reflect.String:
		if ft.Len() != 0 {
			t.Fatal("Value is  empty")
		}
	default:
		t.Fatal("Value is not a map, slice, array, string or channel")
	}
}

func notEmpty(t *testing.T, item interface{}) {
	ft := reflect.ValueOf(item)
	switch ft.Kind() {
	case reflect.Map:
		if ft.Len() == 0 {
			t.Fatal("Value is not empty")
		}
	case reflect.Chan:
		if ft.Len() == 0 {
			t.Fatal("Value is not empty")
		}
	case reflect.Array:
		if ft.Len() == 0 {
			t.Fatal("Value is not empty")
		}
	case reflect.Slice:
		if ft.Len() == 0 {
			t.Fatal("Value is not empty")
		}
	case reflect.String:
		if ft.Len() == 0 {
			t.Fatal("Value is not empty")
		}
	default:
		t.Fatal("Value is not a map, slice, array, string or channel")
	}
}
