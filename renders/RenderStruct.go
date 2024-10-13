package renders

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Render interface {
	Draw()

	// set up vertex data (and buffer(s)) and configure vertex attributes
	InitGLPipeLine()

	// callback function for handling keyboard inputs
	KeyboardCallback(window *glfw.Window)

	// callback function for handling Mouse cursor movement
	MouseCallback(window *glfw.Window, xpos float64, ypos float64)

	// callback function for handling Mouse Scroll movement
	ScrollCallback(window *glfw.Window, xoff float64, yoff float64)
}

type BaseRender struct {
	Render
}
