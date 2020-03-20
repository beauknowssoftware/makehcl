package functions

import (
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var ShellSpec = function.Spec{
	Params: []function.Parameter{
		{Type: cty.String},
	},
	Type: function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		a := args[0].AsString()
		c := exec.Command("sh", "-c", a)

		outPipe, err := c.StdoutPipe()
		if err != nil {
			err = errors.Wrap(err, "failed to create stdout pipe")
			return cty.UnknownVal(cty.DynamicPseudoType), err
		}

		errPipe, err := c.StderrPipe()
		if err != nil {
			err = errors.Wrap(err, "failed to create stderr pipe")
			return cty.UnknownVal(cty.DynamicPseudoType), err
		}

		if err = c.Start(); err != nil {
			err = errors.Wrap(err, "failed to start")
			return cty.UnknownVal(cty.DynamicPseudoType), err
		}

		resBytes, err := ioutil.ReadAll(outPipe)
		if err != nil {
			err = errors.Wrap(err, "failed to read stdout")
			return cty.UnknownVal(cty.DynamicPseudoType), err
		}

		errBytes, err := ioutil.ReadAll(errPipe)
		if err != nil {
			err = errors.Wrap(err, "failed to read stderr")
			return cty.UnknownVal(cty.DynamicPseudoType), err
		}

		if err := c.Wait(); err != nil {
			err = errors.Wrap(err, string(errBytes))
			return cty.UnknownVal(cty.DynamicPseudoType), err
		}

		res := strings.TrimSpace(string(resBytes))

		return cty.StringVal(res), nil
	},
}
