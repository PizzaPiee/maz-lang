package environment

import (
	"maz-lang/object"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvironment(t *testing.T) {
	env := New()

	intObj := object.Integer{Value: 5}
	boolObj := object.Boolean{Value: false}

	env.Set("num", &intObj)
	env.Set("isBuzz", &boolObj)

	num := env.Get("num")
	isBuzz := env.Get("isBuzz")

	if !cmp.Equal(num, &intObj) {
		t.Errorf("expected 'num' to be %+v, instead got %+v\n", intObj, num)
	}

	if !cmp.Equal(isBuzz, &boolObj) {
		t.Errorf("expected 'num' to be %+v, instead got %+v\n", boolObj, isBuzz)
	}

	childEnv := New()
	childIntObj := object.Integer{Value: 100}
	childEnv.Set("child_num", &childIntObj)

	env.Extend(&childEnv)

	childNum := env.Get("child_num")

	if !cmp.Equal(childNum, &childIntObj) {
		t.Errorf("expected 'num' to be %+v, instead got %+v\n", childIntObj, childNum)
	}

	none := env.Get("doesnotexists")

	if !cmp.Equal(none, nil) {
		t.Errorf("expected 'num' to be %+v, instead got %+v\n", nil, none)
	}
}
