package main

import (
	"fmt"
	"opgl-learn/renders"
	lighting "opgl-learn/renders/Lighting"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var WIDTH, HEIGHT = 800, 600

func main() {

	if err := glfw.Init(); err != nil {
		fmt.Println("Something went wrong with the GLFW Init process")
		os.Exit(1)
	}
	defer glfw.Terminate()
	window := setup()

	// configure global opengl state
	// -----------------------------
	gl.Enable(gl.DEPTH_TEST) // Tell opengl to enable depth testing

	// Load Render here
	var render renders.Render = &lighting.Materials{}
	// set up vertex data (and buffer(s)) and configure vertex attributes
	render.InitGLPipeLine()
	window.SetCursorPosCallback(render.MouseCallback)
	window.SetScrollCallback(render.ScrollCallback)
	// render loop
	// -----------
	for !window.ShouldClose() {

		// input
		render.KeyboardCallback(window)

		render.Draw()

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers() // Swap the front and back buffers
		glfw.PollEvents()
	}

	glfw.Terminate()
}

func framebuffer_size_callback(window *glfw.Window, width, height int) {
	// set the viewport
	gl.Viewport(0, 0, int32(width), int32(height))
}

// Setup GLFW and OpenGL function loaders
func setup() *glfw.Window {
	// glfw: initialize and configure
	// ------------------------------
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	// glfw window creation
	// --------------------
	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "OpenGL", nil, nil)
	if err != nil {
		fmt.Println("Failed to create GLFW window")
		os.Exit(1)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebuffer_size_callback)

	// load all OpenGL function pointers
	// --------------------
	// OpenGL function loader must be initalised after the context (GLFW in this case).
	if err := gl.Init(); err != nil {
		fmt.Println("Something went wrong with initalising function loaders:", err)
		os.Exit(1)
	}

	return window
}
