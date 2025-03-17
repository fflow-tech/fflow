package compiler

/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * modified based on https://github.com/grafana/k6/blob/master/js/compiler/compiler.go
 */

import (
	_ "embed" // we need this for embedding Babel
	"sync"
	"time"

	"github.com/fflow-tech/fflow-sdk-go/faas"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

//go:embed babel.min.js
var babelSrc string //nolint:gochecknoglobals

var (
	DefaultOpts = map[string]interface{}{
		"plugins": []interface{}{
			[]interface{}{"transform-es2015-classes", map[string]interface{}{"loose": false}},
			"transform-es2015-object-super",
			[]interface{}{"transform-es2015-modules-commonjs", map[string]interface{}{"loose": false}},
			"transform-exponentiation-operator",
			"transform-regenerator",
			"transform-async-to-generator",
		},
		"ast":           false,
		"sourceMaps":    false,
		"babelrc":       false,
		"compact":       false,
		"retainLines":   true,
		"highlightCode": false,
	}

	onceBabelCode      sync.Once     // nolint:gochecknoglobals
	globalBabelCode    *goja.Program // nolint:gochecknoglobals
	globalBabelCodeErr error         // nolint:gochecknoglobals
	onceBabel          sync.Once     // nolint:gochecknoglobals
	globalBabel        *babel        // nolint:gochecknoglobals
)

// A Compiler compiles JavaScript source code (ES5.1 or ES6) into a goja.Program
type Compiler struct {
	ctx   faas.Context
	babel *babel
}

// New returns a new Compiler
func New(ctx faas.Context) *Compiler {
	return &Compiler{ctx: ctx}
}

// Transform the given code into ES5
func (c *Compiler) Transform(src, filename string) (code string, srcmap []byte, err error) {
	if c.babel == nil {
		onceBabel.Do(func() {
			globalBabel, err = newBabel()
		})
		c.babel = globalBabel
	}
	if err != nil {
		return
	}

	code, srcmap, err = c.babel.Transform(c.ctx, src, filename)
	return
}

// Compile the program in the given CompatibilityMode, wrapping it between pre and post code
func (c *Compiler) Compile(src, filename, pre, post string,
	strict bool) (*goja.Program, string, error) {
	code := pre + src + post
	ast, err := parser.ParseFile(nil, filename, code, 0, parser.WithDisableSourceMaps)
	if err != nil {
		// ES5 的代码 parseFile 失败则尝试 transform 一次
		code, _, err = c.Transform(code, filename)
		if err != nil {
			return nil, code, err
		}
		ast, err = parser.ParseFile(nil, filename, code, 0, parser.WithDisableSourceMaps)
		if err != nil {
			return nil, code, err
		}
	}
	pgm, err := goja.CompileAST(ast, strict)
	if err != nil {
		return nil, code, err
	}
	return pgm, code, err
}

type babel struct {
	vm        *goja.Runtime
	this      goja.Value
	transform goja.Callable
	m         sync.Mutex
}

func newBabel() (*babel, error) {
	onceBabelCode.Do(func() {
		globalBabelCode, globalBabelCodeErr = goja.Compile("babel.min.js", babelSrc, false)
	})
	if globalBabelCodeErr != nil {
		return nil, globalBabelCodeErr
	}
	vm := goja.New()
	_, err := vm.RunProgram(globalBabelCode)
	if err != nil {
		return nil, err
	}

	this := vm.Get("Babel")
	bObj := this.ToObject(vm)
	result := &babel{vm: vm, this: this}
	if err = vm.ExportTo(bObj.Get("transform"), &result.transform); err != nil {
		return nil, err
	}

	return result, err
}

// Transform the given code into ES5, while synchronizing to ensure only a single
// bundle instance / Goja VM is in use at a time.
func (b *babel) Transform(ctx faas.Context, src, filename string) (string, []byte, error) {
	b.m.Lock()
	defer b.m.Unlock()
	opts := make(map[string]interface{})
	for k, v := range DefaultOpts {
		opts[k] = v
	}
	opts["filename"] = filename

	startTime := time.Now()
	v, err := b.transform(b.this, b.vm.ToValue(src), b.vm.ToValue(opts))
	if err != nil {
		return "", nil, err
	}
	ctx.Logger().Infof("Babel: Transformed time %s", time.Since(startTime).String())

	vO := v.ToObject(b.vm)
	var code string
	if err = b.vm.ExportTo(vO.Get("code"), &code); err != nil {
		return code, nil, err
	}
	return code, nil, err
}
