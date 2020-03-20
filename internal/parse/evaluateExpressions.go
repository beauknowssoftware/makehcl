package parse

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func evaluateStringMap(expr hcl.Expression, ctx *hcl.EvalContext) (map[string]string, error) {
	val, diag := expr.Value(ctx)
	if diag.HasErrors() {
		return nil, diag
	}

	t := val.Type()

	if !t.IsObjectType() {
		if et := t.MapElementType(); et == nil {
			return nil, fmt.Errorf("expected map of string, got %v", t.FriendlyName())
		} else if *et != cty.String {
			return nil, fmt.Errorf("expected list of string, but got map of %v", et.FriendlyName())
		}
	}

	m := val.AsValueMap()
	result := make(map[string]string, len(m))

	for k, v := range m {
		result[k] = v.AsString()
	}

	return result, nil
}

func evaluateStringArray(expr hcl.Expression, ctx *hcl.EvalContext) ([]string, error) {
	val, diag := expr.Value(ctx)
	if diag.HasErrors() {
		return nil, diag
	}

	t := val.Type()
	if t == cty.String {
		s, err := evaluateString(expr, ctx)
		if err != nil {
			return nil, err
		}

		return []string{s}, nil
	}

	if t.IsTupleType() {
		ets := t.TupleElementTypes()
		if ets == nil {
			return nil, fmt.Errorf("expected tuple of string, got %v", t.FriendlyName())
		}

		for i, et := range ets {
			if et != cty.String {
				return nil, fmt.Errorf("expected list of string, but element %v was %v", i, et.FriendlyName())
			}
		}

		sl := val.AsValueSlice()
		result := make([]string, 0, len(sl))

		for _, e := range sl {
			result = append(result, e.AsString())
		}

		return result, nil
	}

	if et := t.ListElementType(); et == nil {
		return nil, fmt.Errorf("expected list of string, got %v", t.FriendlyName())
	} else if *et != cty.String {
		return nil, fmt.Errorf("expected list of string, got list of %v", et.FriendlyName())
	}

	sl := val.AsValueSlice()
	result := make([]string, 0, len(sl))

	for _, e := range sl {
		result = append(result, e.AsString())
	}

	return result, nil
}

func evaluateIterable(expr hcl.Expression, ctx *hcl.EvalContext) ([]cty.Value, error) {
	val, diag := expr.Value(ctx)
	if diag.HasErrors() {
		return nil, diag
	}

	if !val.CanIterateElements() {
		return nil, fmt.Errorf("expected an iterable type, got %v", val.Type().FriendlyName())
	}

	return val.AsValueSlice(), nil
}

func evaluateValueArray(expr hcl.Expression, ctx *hcl.EvalContext) ([]cty.Value, error) {
	val, diag := expr.Value(ctx)
	if diag.HasErrors() {
		return nil, diag
	}

	t := val.Type()
	if t.IsTupleType() {
		return val.AsValueSlice(), nil
	}

	if !t.IsListType() {
		return nil, fmt.Errorf("expected list, got %v", t.FriendlyName())
	}

	return val.AsValueSlice(), nil
}

func evaluateString(expr hcl.Expression, ctx *hcl.EvalContext) (string, error) {
	val, diag := expr.Value(ctx)
	if diag.HasErrors() {
		return "", diag
	}

	t := val.Type()
	if t != cty.String {
		return "", fmt.Errorf("expected a string, got a %v", t.FriendlyName())
	}

	return val.AsString(), nil
}

func evaluateBool(expr hcl.Expression, ctx *hcl.EvalContext) (bool, error) {
	val, diag := expr.Value(ctx)
	if diag.HasErrors() {
		return false, diag
	}

	t := val.Type()
	if t != cty.Bool {
		return false, fmt.Errorf("expected a bool, got a %v", t.FriendlyName())
	}

	return val.True(), nil
}
