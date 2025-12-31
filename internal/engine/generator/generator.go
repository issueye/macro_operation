package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/issueye/macro_operation/internal/model"
)

// Generator 脚本生成器
type Generator struct {
	indentLevel  int
	addComments  bool
	optimizeCode bool
}

// NewGenerator 创建脚本生成器
func NewGenerator() *Generator {
	return &Generator{
		indentLevel:  1,
		addComments:  true,
		optimizeCode: true,
	}
}

// Generate 生成JavaScript脚本
func (g *Generator) Generate(name string, events []model.Event) (string, error) {
	if len(events) == 0 {
		return "", fmt.Errorf("no events to generate")
	}

	var sb strings.Builder

	// 生成操作代码
	for i, event := range events {
		code, err := g.generateEventCode(event)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)

		// 添加延迟(除了最后一个操作)
		if i < len(events)-1 {
			nextEvent := events[i+1]
			delay := nextEvent.GetTimestamp().Sub(event.GetTimestamp())
			delayMs := delay.Milliseconds()

			// 只添加有意义的延迟(>20ms)
			if delayMs > 20 {
				sb.WriteString(g.indent(fmt.Sprintf("sleep(%d);\n", delayMs)))
			}
		}
	}

	// 应用模板
	script := fmt.Sprintf(`// 宏名称: %s
// 录制时间: %s
// 操作数量: %d

function main() {
%s
}

// 执行宏
main();
`, name, time.Now().Format("2006-01-02 15:04:05"), len(events), sb.String())

	return script, nil
}

// generateEventCode 生成单个事件的代码
func (g *Generator) generateEventCode(event model.Event) (string, error) {
	switch e := event.(type) {
	case model.KeyboardEvent:
		return g.generateKeyboardCode(e)
	case model.MouseEvent:
		return g.generateMouseCode(e)
	default:
		return "", fmt.Errorf("unsupported event type: %v", event.GetType())
	}
}

// generateKeyboardCode 生成键盘事件代码
func (g *Generator) generateKeyboardCode(event model.KeyboardEvent) (string, error) {
	var code strings.Builder

	switch event.Type {
	case model.KeyDown:
		if g.addComments {
			code.WriteString(g.indent(fmt.Sprintf("// 按下键盘: %s\n", event.Keyname)))
		}
		code.WriteString(g.indent(fmt.Sprintf("keyDown('%s');\n", event.Keyname)))

	case model.KeyUp:
		if g.addComments {
			code.WriteString(g.indent(fmt.Sprintf("// 释放键盘: %s\n", event.Keyname)))
		}
		code.WriteString(g.indent(fmt.Sprintf("keyUp('%s');\n", event.Keyname)))

	case model.KeyPress:
		if g.addComments {
			code.WriteString(g.indent(fmt.Sprintf("// 输入文本: %s\n", event.Keyname)))
		}
		code.WriteString(g.indent(fmt.Sprintf("keyTap('%s');\n", event.Keyname)))

	default:
		return "", fmt.Errorf("unsupported keyboard event type: %v", event.Type)
	}

	return code.String(), nil
}

// generateMouseCode 生成鼠标事件代码
func (g *Generator) generateMouseCode(event model.MouseEvent) (string, error) {
	var code strings.Builder

	switch event.Type {
	case model.MouseMove:
		if g.addComments {
			code.WriteString(g.indent(fmt.Sprintf("// 鼠标移动到 (%d, %d)\n", event.X, event.Y)))
		}
		code.WriteString(g.indent(fmt.Sprintf("mouseMove(%d, %d);\n", event.X, event.Y)))

	case model.MouseClick:
		if g.addComments {
			code.WriteString(g.indent(fmt.Sprintf("// 鼠标点击 (%d, %d) 按钮: %s\n",
				event.X, event.Y, event.Button)))
		}
		code.WriteString(g.indent(fmt.Sprintf("mouseClick('%s');\n", event.Button)))

	case model.MouseDrag:
		if g.addComments {
			code.WriteString(g.indent(fmt.Sprintf("// 鼠标拖拽到 (%d, %d)\n", event.X, event.Y)))
		}
		code.WriteString(g.indent(fmt.Sprintf("mouseDrag(%d, %d);\n", event.X, event.Y)))

	case model.MouseScroll:
		if g.addComments {
			code.WriteString(g.indent(fmt.Sprintf("// 鼠标滚动: %d\n", event.Delta)))
		}
		code.WriteString(g.indent(fmt.Sprintf("mouseScroll(%d);\n", event.Delta)))

	default:
		return "", fmt.Errorf("unsupported mouse event type: %v", event.Type)
	}

	return code.String(), nil
}

// indent 生成缩进
func (g *Generator) indent(code string) string {
	indent := strings.Repeat("  ", g.indentLevel)
	return indent + code
}

// OptimizeEvents 优化事件序列
func (g *Generator) OptimizeEvents(events []model.Event) []model.Event {
	if !g.optimizeCode {
		return events
	}

	var optimized []model.Event
	var lastMouseX, lastMouseY int = -1, -1

	for _, event := range events {
		// 过滤微小鼠标移动
		if me, ok := event.(model.MouseEvent); ok && me.Type == model.MouseMove {
			// 计算移动距离
			if lastMouseX >= 0 && lastMouseY >= 0 {
				dist := (me.X-lastMouseX)*(me.X-lastMouseX) + (me.Y-lastMouseY)*(me.Y-lastMouseY)
				if dist < 25 { // 移动距离小于5像素
					continue
				}
			}
			lastMouseX, lastMouseY = me.X, me.Y
		} else {
			lastMouseX, lastMouseY = -1, -1
		}

		optimized = append(optimized, event)
	}

	return optimized
}
