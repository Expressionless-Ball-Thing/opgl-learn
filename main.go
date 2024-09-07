package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var WIDTH, HEIGHT = 800, 600

var vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;
void main()
{
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"

var fragmentShaderSource = `
#version 330 core
out vec4 FragColor;
void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"

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

	// build and compile our shader program
	// ------------------------------------
	// vertex shader
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)

	// fragment shader
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	// link shaders
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	// check for linking errors
	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var infoLen int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &infoLen)

		log := strings.Repeat("\x00", int(infoLen+1))
		gl.GetProgramInfoLog(shaderProgram, infoLen, nil, gl.Str(log))

		fmt.Printf("failed to link shader program %v: %v\n", shaderProgram, log)
		os.Exit(1)
	}

	// Deleting shaders because we loaded them into the program.
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	var vertices = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		-0.5, 0.5, 0.0, // top left
	}

	var indices = []uint32{
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	var VBO, VAO, EBO uint32
	// Generate Buffers
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// set the vertex attributes pointers
	// Note: we could've also specified the stride as 0 to let OpenGL determine the stride (this only works when values are tightly packed).
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	// note that this is allowed, the call to glVertexAttribPointer registered VBO
	// as the vertex attribute's bound vertex buffer object so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// remember: do NOT unbind the EBO while a VAO is active as the bound element buffer object IS stored in the VAO; keep the EBO bound.
	//glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, 0);

	// You can unbind the VAO afterwards so other VAO calls won't accidentally modify this VAO, but this rarely happens. Modifying other
	// VAOs requires a call to glBindVertexArray anyways so we generally don't unbind VAOs (nor VBOs) when it's not directly necessary.
	gl.BindVertexArray(0)

	// wireframe polygons.
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// render loop
	// -----------
	for !window.ShouldClose() {

		// input
		processInput(window)

		// render
		gl.ClearColor(0.5, 0.5, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// draw Triangle
		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(VAO) // seeing as we only have a single VAO there's no need to bind it every time, but we'll do so to keep things a bit more organized
		// gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers() // Swap the front and back buffers
		glfw.PollEvents()
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
	gl.DeleteProgram(shaderProgram)
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

func compileShader(ShaderSource string, ShaderType uint32) uint32 {
	shader := gl.CreateShader(ShaderType)
	csources, free := gl.Strs(ShaderSource)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// check for shader compile errors
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var infoLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &infoLen)

		log := strings.Repeat("\x00", int(infoLen+1))
		gl.GetShaderInfoLog(shader, infoLen, nil, gl.Str(log))

		fmt.Printf("failed to compile shader %v: %v\n", shader, log)
		os.Exit(1)
	}

	return shader
}
