package cmdkit_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/gokit/cmdkit"
)

func TestFlagParsing(t *testing.T) {
	var suite = []struct {
		MustFail bool
		Expected interface{}
		Type     cmdkit.FlagType
		Value    []string
	}{
		{
			Type:     cmdkit.String,
			Expected: "wallet",
			Value:    []string{"wallet"},
		},
		{
			Type:     cmdkit.StringList,
			Value:    []string{"wallet", "river"},
			Expected: []string{"wallet", "river"},
		},
		{
			Type:     cmdkit.BoolList,
			Value:    []string{"false", "true"},
			Expected: []bool{false, true},
		},
		{
			Type:     cmdkit.Bool,
			Value:    []string{"true"},
			Expected: true,
		},
		{
			Type:     cmdkit.Bool,
			Value:    []string{"false"},
			Expected: false,
		},
		{
			Type:     cmdkit.UInt64List,
			Value:    []string{"1", "2"},
			Expected: []uint64{1, 2},
		},
		{
			Type:     cmdkit.IntList,
			Value:    []string{"1", "2"},
			Expected: []int{1, 2},
		},
		{
			Type:     cmdkit.Int64List,
			Value:    []string{"1", "2"},
			Expected: []int64{1, 2},
		},
		{
			Type:     cmdkit.IntList,
			Value:    []string{"1", "2"},
			Expected: []int{1, 2},
		},
		{
			Type:     cmdkit.Int8,
			Value:    []string{"1"},
			Expected: int8(1),
		},
		{
			Type:     cmdkit.Int16,
			Value:    []string{"1"},
			Expected: int16(1),
		},
		{
			Type:     cmdkit.Int32,
			Value:    []string{"1"},
			Expected: int32(1),
		},
		{
			Type:     cmdkit.Int64,
			Value:    []string{"1"},
			Expected: int64(1),
		},
		{
			Type:     cmdkit.Float32,
			Value:    []string{"1.32"},
			Expected: float32(1.32),
		},
		{
			Type:     cmdkit.Float64,
			Value:    []string{"1.32"},
			Expected: float64(1.32),
		},
		{
			Type:     cmdkit.Duration,
			Value:    []string{"2s"},
			Expected: time.Second * 2,
		},
		{
			Type:     cmdkit.DurationList,
			Value:    []string{"2s", "3m", "30h"},
			Expected: []time.Duration{time.Second * 2, time.Minute * 3, time.Hour * 30},
		},
	}

	for _, tcase := range suite {
		switch tcase.Type {
		case cmdkit.UInt:
			flag := cmdkit.UIntFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Int:
			flag := cmdkit.IntFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Int8:
			flag := cmdkit.Int8Flag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Int32:
			flag := cmdkit.Int32Flag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.UInt64:
			flag := cmdkit.UInt64Flag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Int16:
			flag := cmdkit.Int16Flag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Int64:
			flag := cmdkit.Int64Flag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Bool:
			flag := cmdkit.BoolFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.TBool:
			flag := cmdkit.TBoolFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.String:
			flag := cmdkit.StringFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Float32:
			flag := cmdkit.Float32Flag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Float64:
			flag := cmdkit.Float64Flag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Duration:
			flag := cmdkit.DurationFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.UIntList:
			flag := cmdkit.UIntListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.IntList:
			flag := cmdkit.IntListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.UInt64List:
			flag := cmdkit.UInt64ListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Int64List:
			flag := cmdkit.Int64ListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.BoolList:
			flag := cmdkit.BoolListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.StringList:
			flag := cmdkit.StringListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.Float64List:
			flag := cmdkit.Float64ListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		case cmdkit.DurationList:
			flag := cmdkit.DurationListFlag()
			received, err := flag.Parse(tcase.Value[0], tcase.Value[1:]...)
			if tcase.MustFail && err != nil {
				continue
			}
			if !tcase.MustFail && err != nil {
				t.Fatalf("Should not have failed: %+q\n", err)
			}
			if tcase.Expected != nil {
				if !reflect.DeepEqual(tcase.Expected, received) {
					t.Logf("Recieved: %#v\n", received)
					t.Logf("Expected: %#v\n", tcase.Expected)
					t.Fatal("Should match expected")
				}
			}
		}
	}
}
