package lighting

import (
	"math"
	"opgl-learn/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type LightingMaps struct {
	ShaderProgram, LightCubeShader uint32
	VBO, cubeVAO                   uint32
	lightCubeVAO                   uint32
	camera                         utils.Camera
	lightPos                       mgl32.Vec3
	diffuseMap, specularMap        uint32
}

func (ct *LightingMaps) InitGLPipeLine() {

	ct.lightPos = mgl32.Vec3{1.2, 1.0, 2.0}

	ct.camera = utils.NewCamera(mgl32.Vec3{0.5, 1.0, 4.0}, mgl32.Vec3{0, 1, 0}, utils.YAW, utils.PITCH)

	ct.ShaderProgram = utils.NewShader("./shaders/Lighting/4-LightingMapsVert.glsl", "./shaders/Lighting/4-LightingMapsFrag.glsl")
	ct.LightCubeShader = utils.NewShader("./shaders/Lighting/1-LightVert.glsl", "./shaders/Lighting/1-LightFrag.glsl")

	vertices := []float32{
		// positions          // normals           // texture coords
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	}

	gl.GenVertexArrays(1, &ct.cubeVAO)
	gl.GenBuffers(1, &ct.VBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(ct.cubeVAO)

	// Position attribs
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	// Normal Attribs, it's the third float of each 6 float block.
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)

	// TexCoords
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	// second, configure the light's VAO (VBO stays the same; the vertices are the same for the light object which is also a 3D cube)
	gl.GenVertexArrays(1, &ct.lightCubeVAO)
	gl.BindVertexArray(ct.lightCubeVAO)

	// we only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it; the VBO's data already contains all we need
	gl.BindBuffer(gl.ARRAY_BUFFER, ct.VBO)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	// Texture Stuff
	ct.diffuseMap = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/container2.png")
	ct.specularMap = utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "./assets/container2_specular.png")
	gl.UseProgram(ct.ShaderProgram)
	utils.SetInt(ct.ShaderProgram, "material.diffuse", 0)
	utils.SetInt(ct.ShaderProgram, "material.specular", 1)

}

func (ct *LightingMaps) Draw() {

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	ct.lightPos[0] = float32(1.0 + math.Sin(float64(glfw.GetTime()))*2)
	ct.lightPos[1] = float32(math.Sin(float64(glfw.GetTime()) / 2.0))

	// activate shader
	gl.UseProgram(ct.ShaderProgram)
	utils.SetVec3(ct.ShaderProgram, "light.position", &ct.lightPos)
	utils.SetVec3(ct.ShaderProgram, "viewPos", &(ct.camera.Position))

	// material stuff
	utils.SetFloat(ct.ShaderProgram, "material.shininess", 64.0)

	// Time varing light color
	lightColor := mgl32.Vec3{1, 1, 1}

	diffuseColor := lightColor.Mul(0.5)
	ambientColor := lightColor.Mul(0.2)

	utils.SetVec3(ct.ShaderProgram, "light.ambient", &ambientColor)
	utils.SetVec3(ct.ShaderProgram, "light.diffuse", &diffuseColor)
	utils.SetVec3(ct.ShaderProgram, "light.specular", &mgl32.Vec3{1, 1, 1})

	// camera/view transformation
	projection := mgl32.Perspective(mgl32.DegToRad(float32(ct.camera.Zoom)), float32(800/600), 0.1, 100)
	view := ct.camera.GetViewMatrix()
	utils.SetMat4(ct.ShaderProgram, "view", &view)
	utils.SetMat4(ct.ShaderProgram, "projection", &projection)

	// world transforms
	model := mgl32.Ident4()
	// model = model.Mul4(mgl32.Scale3D(3, 3, 3))
	utils.SetMat4(ct.ShaderProgram, "model", &model)

	// texture stuff
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ct.diffuseMap)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ct.specularMap)

	// render the cube
	gl.BindVertexArray(ct.cubeVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	// Draw the Lamp
	gl.UseProgram(ct.LightCubeShader)
	utils.SetMat4(ct.LightCubeShader, "view", &view)
	utils.SetMat4(ct.LightCubeShader, "projection", &projection)
	model = mgl32.Ident4()
	model = model.Mul4(mgl32.Translate3D(ct.lightPos.X(), ct.lightPos.Y(), ct.lightPos.Z())).Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
	utils.SetMat4(ct.LightCubeShader, "model", &model)

	gl.BindVertexArray(ct.lightCubeVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
}

func (ct *LightingMaps) KeyboardCallback(window *glfw.Window) {

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

func (ct *LightingMaps) MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
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

func (ct *LightingMaps) ScrollCallback(window *glfw.Window, xoff float64, yoff float64) {
	ct.camera.ProcessMouseScroll(yoff)
}
