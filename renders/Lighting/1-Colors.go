package lighting

import (
	"opgl-learn/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var cube_vertices = []float32{
	-0.5, -0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, 0.5, -0.5,
	0.5, 0.5, -0.5,
	-0.5, 0.5, -0.5,
	-0.5, -0.5, -0.5,

	-0.5, -0.5, 0.5,
	0.5, -0.5, 0.5,
	0.5, 0.5, 0.5,
	0.5, 0.5, 0.5,
	-0.5, 0.5, 0.5,
	-0.5, -0.5, 0.5,

	-0.5, 0.5, 0.5,
	-0.5, 0.5, -0.5,
	-0.5, -0.5, -0.5,
	-0.5, -0.5, -0.5,
	-0.5, -0.5, 0.5,
	-0.5, 0.5, 0.5,

	0.5, 0.5, 0.5,
	0.5, 0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, -0.5, 0.5,
	0.5, 0.5, 0.5,

	-0.5, -0.5, -0.5,
	0.5, -0.5, -0.5,
	0.5, -0.5, 0.5,
	0.5, -0.5, 0.5,
	-0.5, -0.5, 0.5,
	-0.5, -0.5, -0.5,

	-0.5, 0.5, -0.5,
	0.5, 0.5, -0.5,
	0.5, 0.5, 0.5,
	0.5, 0.5, 0.5,
	-0.5, 0.5, 0.5,
	-0.5, 0.5, -0.5,
}

type Lighting struct {
	ShaderProgram, LightCubeShader uint32
	VBO, cubeVAO                   uint32
	lightCubeVAO                   uint32
	camera                         utils.Camera
}

func (ct *Lighting) InitGLPipeLine() {

	ct.camera = utils.NewCamera(mgl32.Vec3{0.0, 0.0, 3.0}, mgl32.Vec3{0, 1, 0}, utils.YAW, utils.PITCH)

	ct.ShaderProgram = utils.NewShader("./shaders/Lighting/1-ColorsVert.glsl", "./shaders/Lighting/1-ColorsFrag.glsl")
	ct.LightCubeShader = utils.NewShader("./shaders/Lighting/1-LightVert.glsl", "./shaders/Lighting/1-LightFrag.glsl")

	vertices := cube_vertices

	gl.GenVertexArrays(1, &ct.cubeVAO)
	gl.GenBuffers(1, &ct.VBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(ct.cubeVAO)

	// Position attribs
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	// second, configure the light's VAO (VBO stays the same; the vertices are the same for the light object which is also a 3D cube)
	gl.GenVertexArrays(1, &ct.lightCubeVAO)
	gl.BindVertexArray(ct.lightCubeVAO)

	// we only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it; the VBO's data already contains all we need
	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	// tell opengl for each sampler to which texture unit it belongs to (only has to be done once)
	// -------------------------------------------------------------------------------------------
	gl.UseProgram(ct.ShaderProgram) // don't forget to activate/use the shader before setting uniforms!

}

func (ct *Lighting) Draw() {
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// activate shader
	gl.UseProgram(ct.ShaderProgram)
	utils.SetVec3(ct.ShaderProgram, "objectColor", &mgl32.Vec3{1.0, 0.5, 0.31})
	utils.SetVec3(ct.ShaderProgram, "lightColor", &mgl32.Vec3{1.0, 1.0, 1.0})

	// camera/view transformation
	projection := mgl32.Perspective(mgl32.DegToRad(float32(ct.camera.Zoom)), float32(800/600), 0.1, 100)
	view := ct.camera.GetViewMatrix()
	utils.SetMat4(ct.ShaderProgram, "view", &view)
	utils.SetMat4(ct.ShaderProgram, "projection", &projection)

	// world transforms
	model := mgl32.Ident4()
	utils.SetMat4(ct.ShaderProgram, "model", &model)

	// render the cube
	gl.BindVertexArray(ct.cubeVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	// Draw the Lamp
	gl.UseProgram(ct.LightCubeShader)
	utils.SetMat4(ct.LightCubeShader, "view", &view)
	utils.SetMat4(ct.LightCubeShader, "projection", &projection)
	model = mgl32.Ident4()
	model = model.Mul4(mgl32.Translate3D(1.2, 1.0, 2.0)).Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
	utils.SetMat4(ct.LightCubeShader, "model", &model)

	gl.BindVertexArray(ct.lightCubeVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
}

var lastFrame float64 = 0.0

func (ct *Lighting) KeyboardCallback(window *glfw.Window) {

	currentFrame := glfw.GetTime()
	deltaTime := currentFrame - lastFrame
	lastFrame = currentFrame

	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		ct.camera.ProcessKeyboard(utils.FORWARD, deltaTime)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		ct.camera.ProcessKeyboard(utils.BACKWARD, deltaTime)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		ct.camera.ProcessKeyboard(utils.LEFT, deltaTime)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		ct.camera.ProcessKeyboard(utils.RIGHT, deltaTime)
	}
}

var (
	lastxPos   float64 = 800.0 / 2.0
	lastyPos   float64 = 600.0 / 20.0
	firstMouse bool    = true
)

func (ct *Lighting) MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	if firstMouse {
		firstMouse = false
		lastxPos = xpos
		lastyPos = ypos
	}

	xoffset := xpos - lastxPos
	yoffset := lastyPos - ypos
	lastxPos = xpos
	lastyPos = ypos

	ct.camera.ProcessMouseMovement(xoffset, yoffset, true)
}

func (ct *Lighting) ScrollCallback(window *glfw.Window, xoff float64, yoff float64) {
	ct.camera.ProcessMouseScroll(yoff)
}
