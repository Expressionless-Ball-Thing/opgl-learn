package renders

import (
	"opgl-learn/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Coordinates struct {
	ShaderProgram      uint32
	VAO, VBO, EBO      uint32
	texture1, texture2 uint32
}

func (ct *Coordinates) InitGLPipeLine() {
	ct.ShaderProgram = utils.NewShader("./shaders/6-CoordinatesVert.glsl", "./shaders/6-CoordinatesFrag.glsl")

	// Eight per vertex, 3 position, 2 texture coords.
	var vertices = []float32{
		// positions    // texture Coords.
		0.5, 0.5, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 0.0, 1.0, // top left
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
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)

	// texture
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	ct.texture1 = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/container.png")
	ct.texture2 = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/face.png")

	// tell opengl for each sampler to which texture unit it belongs to (only has to be done once)
	// -------------------------------------------------------------------------------------------
	gl.UseProgram(ct.ShaderProgram) // don't forget to activate/use the shader before setting uniforms!
	// set it manually like so:
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture1"+"\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture2"+"\x00")), 1)
}

func (ct *Coordinates) Draw() {

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// bind Texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture2)

	// create transformations
	model, view, projection := mgl32.Ident4(), mgl32.Ident4(), mgl32.Ident4()
	model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(-55), mgl32.Vec3{1, 0, 0}))
	view = view.Mul4(mgl32.Translate3D(0, 0, -3))
	projection = mgl32.Perspective(mgl32.DegToRad(45), 16.0/9.0, 0.1, 100)

	// retrieve the matrix uniform locations
	gl.UseProgram(ct.ShaderProgram)
	modelLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("model"+"\x00"))
	viewLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("view"+"\x00"))
	projectionLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("projection"+"\x00"))

	// pass them to the shaders (3 different ways)
	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

	// render container
	gl.BindVertexArray(ct.VAO)
	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

}

var cube_vertices = []float32{
	-0.5, -0.5, -0.5, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,

	-0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,

	-0.5, 0.5, -0.5, 0.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
}

type Cube struct {
	ShaderProgram      uint32
	VAO, VBO, EBO      uint32
	texture1, texture2 uint32
}

func (ct *Cube) InitGLPipeLine() {

	ct.ShaderProgram = utils.NewShader("./shaders/6-CoordinatesVert.glsl", "./shaders/6-CoordinatesFrag.glsl")

	var vertices = cube_vertices

	gl.GenVertexArrays(1, &ct.VAO)
	gl.GenBuffers(1, &ct.VBO)

	gl.BindVertexArray(ct.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	ct.texture1 = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/container.png")
	ct.texture2 = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/face.png")

	// tell opengl for each sampler to which texture unit it belongs to (only has to be done once)
	// -------------------------------------------------------------------------------------------
	gl.UseProgram(ct.ShaderProgram) // don't forget to activate/use the shader before setting uniforms!
	// set it manually like so:
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture1"+"\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture2"+"\x00")), 1)
}

func (ct *Cube) Draw() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// bind Texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture2)

	// create transformations
	model, view, projection := mgl32.Ident4(), mgl32.Ident4(), mgl32.Ident4()
	model = model.Mul4(mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0.5, 1.0, 0.0}))
	view = view.Mul4(mgl32.Translate3D(0, 0, -3))
	projection = mgl32.Perspective(mgl32.DegToRad(45), 16.0/9.0, 0.1, 100)

	// retrieve the matrix uniform locations
	gl.UseProgram(ct.ShaderProgram)
	modelLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("model"+"\x00"))
	viewLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("view"+"\x00"))
	projectionLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("projection"+"\x00"))

	// pass them to the shaders (3 different ways)
	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

	// render container
	gl.BindVertexArray(ct.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

}

type MoreCubes struct {
	ShaderProgram      uint32
	VAO, VBO, EBO      uint32
	texture1, texture2 uint32
	cube_positions     []mgl32.Vec3
}

func (ct *MoreCubes) InitGLPipeLine() {

	ct.ShaderProgram = utils.NewShader("./shaders/6-CoordinatesVert.glsl", "./shaders/6-CoordinatesFrag.glsl")

	var vertices = cube_vertices

	ct.cube_positions = []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}

	gl.GenVertexArrays(1, &ct.VAO)
	gl.GenBuffers(1, &ct.VBO)

	gl.BindVertexArray(ct.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	ct.texture1 = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/container.png")
	ct.texture2 = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/face.png")

	// tell opengl for each sampler to which texture unit it belongs to (only has to be done once)
	// -------------------------------------------------------------------------------------------
	gl.UseProgram(ct.ShaderProgram) // don't forget to activate/use the shader before setting uniforms!
	// set it manually like so:
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture1"+"\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(ct.ShaderProgram, gl.Str("texture2"+"\x00")), 1)
}

func (ct *MoreCubes) Draw() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// bind Texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture2)

	// create transformations
	view, projection := mgl32.Ident4(), mgl32.Ident4()
	view = view.Mul4(mgl32.Translate3D(0, 0, -3))
	projection = mgl32.Perspective(mgl32.DegToRad(45), float32(800/600), 0.1, 100)

	// retrieve the matrix uniform locations
	modelLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("model"+"\x00"))
	viewLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("view"+"\x00"))
	projectionLoc := gl.GetUniformLocation(ct.ShaderProgram, gl.Str("projection"+"\x00"))

	// pass them to the shaders (3 different ways)
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

	// render container
	gl.BindVertexArray(ct.VAO)
	for i := 0; i < len(ct.cube_positions); i++ {
		model := mgl32.Ident4()
		model = model.Mul4(mgl32.Translate3D(ct.cube_positions[i][0], ct.cube_positions[i][1], ct.cube_positions[i][2]))
		angle := i * 20
		model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(float32(angle)), mgl32.Vec3{1, 0.3, 0.5}))
		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
}
