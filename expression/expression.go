package expression

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/robertkrimen/otto"
	"liuyi-test/utils"
	"strconv"
	"strings"
)

func assert(exp string) (bool, error) {
	vm := otto.New()
	line := "result = " + exp + ";"
	vm.Run(line)
	r, e := vm.Get("result")
	utils.PanicError(e)
	return r.ToBoolean()
}

/**
@code = #code
*/
func Assign(ctx map[string]interface{}, right string, exp string) {
	if ctx == nil {
		utils.PanicError(utils.NewError1(ctxIsEmpty))
	}

	ss := strings.Split(exp, "=")
	l := ss[0]
	r := strings.TrimSpace(ss[1])

	if strings.HasPrefix(r, SHARP) {
		_, v, e := getVariableFromString(right, r)()
		utils.PanicError(e)

		if strings.HasPrefix(l, AT) {
			key := strings.TrimSpace(strings.Split(l, AT)[1])
			ctx[key] = v
		}
	}
}

/**
paramName = @code
 */
func Param(ctx map[string]interface{}, script string) (interface{}, error) {
	if strings.Contains(script,AT){
		origin := strings.Split(script,AT)
		_, v, e := getVariableFromMap(ctx, origin[1])()
		utils.PanicError(e)
		return origin[0]+fmt.Sprintf("%v",v),nil
	}else{
		return script,nil
	}
}

/**
#code == 200 && #data.firstPage == true && #data.list.0.content == 'sdsjdks'
*/
func Assert(ctx map[string]interface{}, right string, exp string) (bool, error) {

	s := strings.Split(exp, " ")

	for i, v := range s {
		if strings.HasPrefix(v, AT) {
			_, v, e := getVariableFromMap(ctx, v)()
			utils.PanicError(e)
			s[i] = afterGetVariable(fmt.Sprintf("%v", v))
		}

		if strings.HasPrefix(v, SHARP) {
			_, v, e := getVariableFromString(right, v)()
			utils.PanicError(e)
			s[i] = afterGetVariable(fmt.Sprintf("%v", v))
		}
	}

	_exp := strings.Join(s, " ")
	return assert(_exp)
}

func afterGetVariable(v string) string {
	_, e1 := cast2Int(v)
	if e1 != nil {
		_, e2 := cast2Bool(v)
		if e2 != nil {
			return "'" + v + "'"
		} else {
			return v
		}
	} else {
		return v
	}
}

func getVariableFromMap(ctx map[string]interface{}, script string) func() (string /*key*/, interface{} /*value*/, error) {

	if !strings.HasPrefix(script, AT) {
		return func() (s string, i interface{}, err error) {
			s = script
			i = ctx[script]
			err = nil
			return s, i, err
		}
	}

	if ctx == nil {
		return func() (s string, i interface{}, err error) {
			return "", nil, utils.NewError1(ctxIsEmpty)
		}
	}

	variables := strings.Split(script, AT)
	if !utils.SliceNotEmptyAndLength(variables, 2) {
		return func() (s string, i interface{}, err error) {
			return "", nil, utils.NewError2(invalidScript, script)
		}

	}

	if utils.StringIsBlank(variables[1]) {
		return func() (s string, i interface{}, err error) {
			return "", nil, utils.NewError2(invalidScript, script)
		}
	}

	key := variables[1]

	return func() (string, interface{}, error) {
		return key, ctx[key], nil
	}

}

func getVariableFromString(src string, script string) func() (string /*key*/, interface{} /*value*/, error) {

	if !strings.HasPrefix(script, SHARP) {
		return func() (string, interface{}, error) {
			return script, script, nil
		}
	}

	if utils.StringIsBlank(src) {
		return func() (s string, i interface{}, err error) {
			return "", nil, utils.NewError1(strIsEmpty)
		}
	}

	if utils.StringIsBlank(script) {
		return func() (s string, i interface{}, err error) {
			return "", nil, utils.NewError2(invalidScript, script)
		}
	}

	variables := strings.Split(script, SHARP)

	if !utils.SliceNotEmptyAndLength(variables, 2) {
		return func() (s string, i interface{}, err error) {
			return "", nil, utils.NewError2(invalidScript, script)
		}
	}

	if utils.StringIsBlank(variables[1]) {
		return func() (s string, i interface{}, err error) {
			return "", nil, utils.NewError2(invalidScript, script)
		}
	}

	dot := strings.Split(variables[1], ".")

	var t []interface{} = make([]interface{}, len(dot))
	for i, v := range dot {
		index, err := cast2Int(v)
		if err == nil {
			t[i] = index
		} else {
			t[i] = v
		}
	}

	return func() (string /*key*/, interface{} /*value*/, error) {
		v := jsoniter.Get([]byte(src), t...).ToString()
		return variables[1], v, nil
	}

}

func cast2Int(s string) (int, error) {
	return strconv.Atoi(s)
}

func cast2Bool(s string) (bool, error) {
	switch s {
	case "true", "TRUE", "True":
		return true, nil
	case "false", "FALSE", "False":
		return false, nil
	}
	return false, utils.NewError2(couldNotToBool,s)
}

const (
	noVarInCtx        		string = "No [%s] exist in context"
	noVarInString     	string = "No [%s] exist in string"
	ctxIsEmpty        	string = "context is empty or nil"
	strIsEmpty        		string = "string is empty or nil"
	invalidExpression 	string = "invalid expression [%s]"
	invalidScript     		string = "invalid script [%s]"
	couldNotToBool 		string = "could not cast to bool [%s]"

	AT    string = "@"
	SHARP string = "#"
)
