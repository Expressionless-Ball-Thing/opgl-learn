package renders

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Transformation struct {
	ShaderProgram      uint32
	VAO, VBO, EBO      uint32
	texture1, texture2 uint32
}

func (ct *Transformation) InitGLPipeLine() {
	ct.ShaderProgram = NewShader("./shaders/5-TransformVert.glsl", "./shaders/5-TransformFrag.glsl")

	// Eight per vertex, 3 position, 3 color, 2 texture coords.
	var vertices = []float32{
		// positions    // colors     // texture coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
	}

	var indices = []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	gl.GenVertexArrays(1, &ct.VAO)
	gl.GenBuffers(1, &ct.VBO)
	gl.GenBuffers(1, &ct.EBO)

	gl.BindVertexArray(ct.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ct.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	// color
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)

	// texture
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	ct.texture1 = New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/container.png")
	ct.texture2 = New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/face.png")

	// tell opengl for each sampler to which texture unit it belongs to (only has to be done once)
	// -------------------------------------------------------------------------------------------
	gl.UseProgram(ct.ShaderProgram) // don't forget to activate/use the shader before setting uniforms!
	// set it manually like so:
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture1"+"\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture2"+"\x00")), 1)
}

func (ct *Transformation) Draw() {

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// bind Texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture2)

	// create transformations
	trans := mgl32.Ident4()
	trans = trans.Mul4(mgl32.Translate3D(0.5, -0.5, 0))
	trans = trans.Mul4(mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0.0, 0.0, 1.0}))

	transformLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("transform"+"\x00"))
	gl.UniformMatrix4fv(transformLoc, 1, false, &trans[0])

	// render container
	gl.BindVertexArray(ct.VAO)
	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

}
