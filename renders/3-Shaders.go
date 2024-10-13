package renders

import (
	"math"
	"opgl-learn/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type ChangingTriangle struct {
	ShaderProgram uint32
	VAO, VBO      uint32
}

func (ct *ChangingTriangle) InitGLPipeLine() {

	var vertexShaderSource = `
	#version 330 core
	layout (location = 0) in vec3 aPos;
	
	out vec4 vertexColor;
	
	void main()
	{
		gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
		vertexColor = vec4(0.5, 0.0, 0.0, 1.0);
	}
	` + "\x00"

	var fragmentShaderSource = `
	#version 330 core
	out vec4 FragColor;
	uniform vec4 ourColor;
	void main()
	{
		FragColor = ourColor;
	}
	` + "\x00"

	ct.ShaderProgram = utils.CreateShaderProgram(vertexShaderSource, fragmentShaderSource)

	var vertices = []float32{
		-0.5, -0.5, 0.0, // left
		0.5, -0.5, 0.0, // right
		0.0, 0.5, 0.0, // top
	}

	gl.GenVertexArrays(1, &ct.VAO)
	gl.GenBuffers(1, &ct.VBO)

	gl.BindVertexArray(ct.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

}

func (ct *ChangingTriangle) Draw() {

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// activate shader
	gl.UseProgram(ct.ShaderProgram)

	// update the uniform color
	timeValue := glfw.GetTime()
	greenValue := (math.Sin(timeValue) / (2.0)) + 0.5
	vertexColorLocation := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("ourColor"+"\x00"))
	gl.Uniform4f(vertexColorLocation, 0.0, float32(greenValue), 0.0, 1.0)

	// Render the triangle
	gl.BindVertexArray(ct.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

type RainbowTriangle struct {
	ShaderProgram uint32
	VAO, VBO      uint32
}

func (rt *RainbowTriangle) InitGLPipeLine() {

	rt.ShaderProgram = utils.NewShader("./shaders/3-rainbowTriangleVert.glsl", "./shaders/3-rainbowTriangleFrag.glsl")

	var vertices = []float32{
		// positions    // colors
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom left
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // top
	}

	gl.GenVertexArrays(1, &rt.VAO)
	gl.GenBuffers(1, &rt.VBO)

	gl.BindVertexArray(rt.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, rt.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position Attribute
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
	gl.EnableVertexAttribArray(0)

	// Color Attribute
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

}

func (rt *RainbowTriangle) Draw() {

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// activate shader
	gl.UseProgram(rt.ShaderProgram)

	// Render the triangle
	gl.BindVertexArray(rt.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
