package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var WIDTH, HEIGHT = 800, 600

func main() {

	// glfw: initialize and configure
	// ------------------------------
	if err := glfw.Init(); err != nil {
		fmt.Println("Something went wrong with the GLFW Init process")
		os.Exit(1)
	}
	defer glfw.Terminate()
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

	// render loop
	// -----------
	for !window.ShouldClose() {

		glfw.PollEvents()

		// input
		processInput(window)

		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers() // Swap the front and back buffers
	}
}

func framebuffer_size_callback(window *glfw.Window, width, height int) {
	// set the viewport
	gl.Viewport(0, 0, int32(width), int32(height))
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	window.SwapBuffers()
}
