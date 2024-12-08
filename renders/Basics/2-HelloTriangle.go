package basics

import (
	"opgl-learn/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
)

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
in vec4 vertexColor;
void main()
{
	FragColor = vertexColor;
}
` + "\x00"

type Triangle struct {
	ShaderProgram uint32
	VAO, VBO      uint32
}

func (t *Triangle) InitGLPipeLine() {
	var vertices = []float32{
		-0.5, -0.5, 0.0, // left
		0.5, -0.5, 0.0, // right
		0.0, 0.5, 0.0, // top
	}

	t.ShaderProgram = utils.CreateShaderProgram(vertexShaderSource, fragmentShaderSource)

	// Generate Buffers
	gl.GenVertexArrays(1, &t.VAO)
	gl.GenBuffers(1, &t.VBO)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(t.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, t.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// set the vertex attributes pointers
	// Note: we could've also specified the stride as 0 to let OpenGL determine the stride (this only works when values are tightly packed).
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	// note that this is allowed, the call to glVertexAttribPointer registered VBO
	// as the vertex attribute's bound vertex buffer object so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// You can unbind the VAO afterwards so other VAO calls won't accidentally modify this VAO, but this rarely happens. Modifying other
	// VAOs requires a call to glBindVertexArray anyways so we generally don't unbind VAOs (nor VBOs) when it's not directly necessary.
	gl.BindVertexArray(0)

}

func (t *Triangle) Draw() {

	gl.ClearColor(0.5, 0.5, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// draw Triangle
	gl.UseProgram(t.ShaderProgram)
	gl.BindVertexArray(t.VAO) // seeing as we only have a single VAO there's no need to bind it every time, but we'll do so to keep things a bit more organized
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

type TwoTriangles struct {
	ShaderProgram uint32
	VAO, VBO, EBO uint32
}

func (tt *TwoTriangles) InitGLPipeLine() {

	var TwoTrianglesvertices = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		-0.5, 0.5, 0.0, // top left
	}

	var TwoTrianglesindices = []uint32{
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	tt.ShaderProgram = utils.CreateShaderProgram(vertexShaderSource, fragmentShaderSource)

	// Generate Buffers
	gl.GenVertexArrays(1, &tt.VAO)
	gl.GenBuffers(1, &tt.VBO)
	gl.GenBuffers(1, &tt.EBO)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(tt.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, tt.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(TwoTrianglesvertices), gl.Ptr(TwoTrianglesvertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, tt.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(TwoTrianglesindices), gl.Ptr(TwoTrianglesindices), gl.STATIC_DRAW)

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
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

}

func (tt *TwoTriangles) Draw() {

	gl.ClearColor(0.5, 0.5, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// draw Triangle
	gl.UseProgram(tt.ShaderProgram)
	gl.BindVertexArray(tt.VAO)                                     // seeing as we only have a single VAO there's no need to bind it every time, but we'll do so to keep things a bit more organized
	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0) // Draw a triangle base on the indicies specified in the EBO.
}
