package expr

import (
	"fmt"
	"math"
	"reflect"
)

type Env map[Var]interface{}

type runtimePanic string

func SafetyEvalBool(expr Expr, env Env) (value bool, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// no panic
		case runtimePanic:
			value = false
			err = fmt.Errorf("%s", x)
		default:
			// unexpected panic: resume state of panic.
			panic(x)
		}
	}()
	if expr == nil {
		return false, nil
	}
	value = ConvertToBool(expr.Eval(env))
	return
}

func (v Var) Eval(env Env) reflect.Value {
	switch v {
	case "true":
		return reflect.ValueOf(true)
	case "false":
		return reflect.ValueOf(false)
	default:
		if i, ok := env[v]; ok {
			return reflect.ValueOf(i)
		}
		panic(runtimePanic(fmt.Sprintf("undefined variable: %s", v)))
	}
}

func (l literal) Eval(_ Env) reflect.Value {
	return reflect.ValueOf(l.value)
}

func (u unary) Eval(env Env) reflect.Value {
	switch u.op {
	case "+":
		return unaryPlus(u.x.Eval(env))
	case "-":
		return unaryMinus(u.x.Eval(env))
	case "!":
		return logicalNegation(u.x.Eval(env))
	case "~":
		return bitwiseComplement(u.x.Eval(env))
	}
	panic(runtimePanic(fmt.Sprintf("unsupported unary operator: %q", u.op)))
}

func (b binary) Eval(env Env) reflect.Value {
	switch b.op {
	case "+":
		return addition(b.x.Eval(env), b.y.Eval(env))
	case "-":
		return subtraction(b.x.Eval(env), b.y.Eval(env))
	case "*":
		return multiplication(b.x.Eval(env), b.y.Eval(env))
	case "/":
		return division(b.x.Eval(env), b.y.Eval(env))
	case "%":
		return modulus(b.x.Eval(env), b.y.Eval(env))
	case "&":
		return bitwiseAnd(b.x.Eval(env), b.y.Eval(env))
	case "&&":
		return logicalAnd(b.x.Eval(env), b.y.Eval(env))
	case "|":
		return bitwiseOr(b.x.Eval(env), b.y.Eval(env))
	case "||":
		return logicalOr(b.x.Eval(env), b.y.Eval(env))
	case "=", "==":
		return comparisonEqual(b.x.Eval(env), b.y.Eval(env))
	case ">":
		return comparisonGreater(b.x.Eval(env), b.y.Eval(env))
	case ">=":
		return comparisonGreaterOrEqual(b.x.Eval(env), b.y.Eval(env))
	case "<":
		return comparisonLess(b.x.Eval(env), b.y.Eval(env))
	case "<=":
		return comparisonLessOrEqual(b.x.Eval(env), b.y.Eval(env))
	case "!=":
		return comparisonNotEqual(b.x.Eval(env), b.y.Eval(env))
	}
	panic(runtimePanic(fmt.Sprintf("unsupported binary operator: %q", b.op)))
}

func (c call) Eval(env Env) reflect.Value {
	switch c.fn {
	case "pow":
		return reflect.ValueOf(math.Pow(ConvertToFloat(c.args[0].Eval(env)), ConvertToFloat(c.args[1].Eval(env))))
	case "sin":
		return reflect.ValueOf(math.Sin(ConvertToFloat(c.args[0].Eval(env))))
	case "sqrt":
		return reflect.ValueOf(math.Sqrt(ConvertToFloat(c.args[0].Eval(env))))
	}
	panic(runtimePanic(fmt.Sprintf("unsupported function call: %s", c.fn)))
}

func ConvertToBool(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	default:
		panic(runtimePanic(fmt.Sprintf("cannot convert data type: %s to bool", v.Kind().String())))
	}
}

func ConvertToInt(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(v.Float())
	default:
		panic(runtimePanic(fmt.Sprintf("cannot convert data type: %s to int", v.Kind().String())))
	}
}

func ConvertToUint(v reflect.Value) uint64 {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return uint64(v.Float())
	default:
		panic(runtimePanic(fmt.Sprintf("cannot convert data type: %s to uint", v.Kind().String())))
	}
}

func ConvertToFloat(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	default:
		panic(runtimePanic(fmt.Sprintf("cannot convert data type: %s to float", v.Kind().String())))
	}
}

func unaryPlus(v reflect.Value) reflect.Value {
	return v
}

func unaryMinus(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Bool:
		return v
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(-v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(-v.Uint())
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(-v.Float())
	default:
		panic(runtimePanic(fmt.Sprintf("unary minus not support type: %s", v.Kind().String())))
	}
}

func logicalNegation(v reflect.Value) reflect.Value {
	return reflect.ValueOf(!ConvertToBool(v))
}

func bitwiseComplement(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(!v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(^v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(^v.Uint())
	case reflect.Float32, reflect.Float64:
		panic(runtimePanic("cannot eval ~ for float"))
	default:
		panic(runtimePanic(fmt.Sprintf("bitwise complement not support type: %s", v.Kind().String())))
	}
}

func typeLevel(k reflect.Kind) int {
	switch k {
	case reflect.Float32, reflect.Float64:
		return 4
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return 3
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 2
	case reflect.Bool:
		return 1
	default:
		return 0
	}
}

func typeAscend(a reflect.Kind, b reflect.Kind) reflect.Kind {
	if typeLevel(a) >= typeLevel(b) {
		return a
	} else {
		return b
	}
}

func addition(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) + ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) + ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) + ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) + ConvertToInt(right)
		return reflect.ValueOf(r != 0)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support addition", left.Kind().String(), right.Kind().String())))
	}
}

func subtraction(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) - ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) - ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) - ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) - ConvertToInt(right)
		return reflect.ValueOf(r != 0)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support subtraction", left.Kind().String(), right.Kind().String())))
	}
}

func multiplication(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) * ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) * ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) * ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) * ConvertToInt(right)
		return reflect.ValueOf(r != 0)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support multiplication", left.Kind().String(), right.Kind().String())))
	}
}

func division(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) / ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) / ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) / ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) / ConvertToInt(right)
		return reflect.ValueOf(r != 0)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support division", left.Kind().String(), right.Kind().String())))
	}
}

func modulus(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) % ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) % ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) % ConvertToInt(right)
		return reflect.ValueOf(r != 0)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support division", left.Kind().String(), right.Kind().String())))
	}
}

func bitwiseAnd(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) & ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) & ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToBool(left) && ConvertToBool(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support bitwise and", left.Kind().String(), right.Kind().String())))
	}
}

func bitwiseOr(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) | ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) | ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToBool(left) || ConvertToBool(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support bitwise or", left.Kind().String(), right.Kind().String())))
	}
}

func logicalAnd(left reflect.Value, right reflect.Value) reflect.Value {
	r := ConvertToBool(left) && ConvertToBool(right)
	return reflect.ValueOf(r)
}

func logicalOr(left reflect.Value, right reflect.Value) reflect.Value {
	r := ConvertToBool(left) || ConvertToBool(right)
	return reflect.ValueOf(r)
}

func comparisonEqual(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) == ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) == ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) == ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) == ConvertToInt(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support comparison equal", left.Kind().String(), right.Kind().String())))
	}
}

func comparisonNotEqual(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) != ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) != ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) != ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) != ConvertToInt(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support comparison not equal", left.Kind().String(), right.Kind().String())))
	}
}

func comparisonGreater(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) > ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) > ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) > ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) > ConvertToInt(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support comparison greater", left.Kind().String(), right.Kind().String())))
	}
}

func comparisonGreaterOrEqual(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) >= ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) >= ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) >= ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) >= ConvertToInt(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support comparison greater or equal", left.Kind().String(), right.Kind().String())))
	}
}

func comparisonLess(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) < ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) < ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) < ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) < ConvertToInt(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support comparison less", left.Kind().String(), right.Kind().String())))
	}
}

func comparisonLessOrEqual(left reflect.Value, right reflect.Value) reflect.Value {
	k := typeAscend(left.Kind(), right.Kind())
	switch k {
	case reflect.Float32, reflect.Float64:
		r := ConvertToFloat(left) <= ConvertToFloat(right)
		return reflect.ValueOf(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := ConvertToUint(left) <= ConvertToUint(right)
		return reflect.ValueOf(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := ConvertToInt(left) <= ConvertToInt(right)
		return reflect.ValueOf(r)
	case reflect.Bool:
		r := ConvertToInt(left) <= ConvertToInt(right)
		return reflect.ValueOf(r)
	default:
		panic(runtimePanic(fmt.Sprintf("type %s and %s not support comparison less or equal", left.Kind().String(), right.Kind().String())))
	}
}
