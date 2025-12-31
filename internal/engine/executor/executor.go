package executor

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/issueye/macro_operation/pkg/bindings"
)

// Executor 脚本执行器
type Executor struct {
	vm        *goja.Runtime
	isRunning bool
	bindings  *bindings.APIBindings
}

// NewExecutor 创建执行器
func NewExecutor() *Executor {
	vm := goja.New()

	executor := &Executor{
		vm:       vm,
		bindings: bindings.NewAPIBindings(),
	}

	// 注册API
	executor.registerAPIs()

	return executor
}

// registerAPIs 注册API函数
func (e *Executor) registerAPIs() {
	// 鼠标操作API
	e.vm.Set("mouseMove", func(call goja.FunctionCall) goja.Value {
		x := call.Argument(0).ToInteger()
		y := call.Argument(1).ToInteger()
		err := e.bindings.MouseMove(int(x), int(y))
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("mouseClick", func(call goja.FunctionCall) goja.Value {
		button := call.Argument(0).String()
		err := e.bindings.MouseClick(button)
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("mouseDrag", func(call goja.FunctionCall) goja.Value {
		x := call.Argument(0).ToInteger()
		y := call.Argument(1).ToInteger()
		err := e.bindings.MouseDrag(int(x), int(y))
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("mouseScroll", func(call goja.FunctionCall) goja.Value {
		delta := call.Argument(0).ToInteger()
		err := e.bindings.MouseScroll(int(delta))
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	// 键盘操作API
	e.vm.Set("keyDown", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		err := e.bindings.KeyDown(key)
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("keyUp", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		err := e.bindings.KeyUp(key)
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("keyType", func(call goja.FunctionCall) goja.Value {
		text := call.Argument(0).String()
		err := e.bindings.KeyType(text)
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("keyTap", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		err := e.bindings.KeyTap(key)
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	// 系统操作API
	e.vm.Set("sleep", func(call goja.FunctionCall) goja.Value {
		ms := call.Argument(0).ToInteger()
		err := e.bindings.Sleep(int(ms))
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("screenshot", func(call goja.FunctionCall) goja.Value {
		path := call.Argument(0).String()
		err := e.bindings.Screenshot(path)
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	e.vm.Set("log", func(call goja.FunctionCall) goja.Value {
		msg := call.Argument(0).String()
		fmt.Printf("[LOG] %s\n", msg)
		return goja.Undefined()
	})
}

// Execute 执行脚本
func (e *Executor) Execute(script string) error {
	if e.isRunning {
		return fmt.Errorf("executor already running")
	}

	e.isRunning = true
	defer func() {
		e.isRunning = false
	}()

	// 编译脚本
	program, err := goja.Compile("", script, true)
	if err != nil {
		return fmt.Errorf("compile error: %w", err)
	}

	// 执行脚本
	_, err = e.vm.RunProgram(program)
	if err != nil {
		return fmt.Errorf("runtime error: %w", err)
	}

	return nil
}

// IsRunning 检查是否正在运行
func (e *Executor) IsRunning() bool {
	return e.isRunning
}
