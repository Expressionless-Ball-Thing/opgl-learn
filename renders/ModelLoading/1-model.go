package ModelLoading

import (
	"opgl-learn/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type ModelLoad struct {
	ShaderProgram uint32
	camera        utils.Camera
	model         utils.Model
}

var lastxPos float64 = 1920 / 2.0
var lastyPos float64 = 1080 / 2.0
var firstMouse bool = true
var lastFrame float64 = 0.0

func (ct *ModelLoad) InitGLPipeLine() {

	ct.camera = utils.NewCamera(mgl32.Vec3{0.0, 0.0, 3.0}, mgl32.Vec3{0, 1, 0}, utils.YAW, utils.PITCH)

	ct.ShaderProgram = utils.NewShader("./shaders/ModelLoading/1-ModelVert.glsl", "./shaders/ModelLoading/1-ModelFrag.glsl")
	ct.model = utils.NewModel("./backpack/backpack.obj")

	gl.UseProgram(ct.ShaderProgram)
}

func (ct *ModelLoad) Draw() {

	gl.ClearColor(0.05, 0.05, 0.05, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(ct.ShaderProgram)

	// view/projection transformations
	projection := mgl32.Perspective(mgl32.DegToRad(float32(ct.camera.Zoom)), float32(800/600), 0.1, 100)
	view := ct.camera.GetViewMatrix()
	utils.SetMat4(ct.ShaderProgram, "view", &view)
	utils.SetMat4(ct.ShaderProgram, "projection", &projection)

	// render the loaded model
	model := mgl32.Ident4()
	model = mgl32.Translate3D(0, 0, 0).Mul4(model)
	model = mgl32.Scale3D(1, 1, 1).Mul4(model)
	utils.SetMat4(ct.ShaderProgram, "model", &model)
	ct.model.Draw(ct.ShaderProgram)
}

func (ct *ModelLoad) KeyboardCallback(window *glfw.Window) {

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
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		ct.camera.ProcessKeyboard(utils.UP, deltaTime)
	}
	if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		ct.camera.ProcessKeyboard(utils.DOWN, deltaTime)
	}
}

func (ct *ModelLoad) MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
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

func (ct *ModelLoad) ScrollCallback(window *glfw.Window, xoff float64, yoff float64) {
	ct.camera.ProcessMouseScroll(yoff)
}
