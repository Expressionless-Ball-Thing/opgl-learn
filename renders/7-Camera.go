package renders

import (
	"math"
	"opgl-learn/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	ShaderProgram                    uint32
	VAO, VBO, EBO                    uint32
	texture1, texture2               uint32
	cube_positions                   []mgl32.Vec3
	cameraPos, cameraUp, cameraFront mgl32.Vec3
}

var (
	firstMouse bool    = true
	lastxPos   float64 = 800.0 / 2.0
	lastyPos   float64 = 600.0 / 20.0
	pitch      float64 = 0.0
	yaw        float64 = -90.0
	fov        float64 = 45
)

func (ct *Camera) InitGLPipeLine() {
	ct.cameraPos = mgl32.Vec3{0, 0, 3}
	ct.cameraFront = mgl32.Vec3{0, 0, -1}
	ct.cameraUp = mgl32.Vec3{0, 1, 0}

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
	utils.SetInt(ct.ShaderProgram, "texture1", 0)
	utils.SetInt(ct.ShaderProgram, "texture2", 1)
}

func (ct *Camera) Draw() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// bind Texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, ct.texture2)

	// activate shader
	gl.UseProgram(ct.ShaderProgram)

	// camera/view transformation
	projection := mgl32.Ident4()
	// Look AT consists of a matrix with the right, up and direction vector, multiplied by a matrix consists of the camera's position.
	view := mgl32.LookAtV(ct.cameraPos, ct.cameraPos.Add(ct.cameraFront), ct.cameraUp)
	utils.SetMat4(ct.ShaderProgram, "view", &view)

	projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(16.0/9.0), 0.1, 100)
	utils.SetMat4(ct.ShaderProgram, "projection", &projection)

	// render boxes
	gl.BindVertexArray(ct.VAO)
	for i := 0; i < len(ct.cube_positions); i++ {
		// calculate the model matrix for each object and pass it to shader before drawing
		model := mgl32.Translate3D(ct.cube_positions[i][0], ct.cube_positions[i][1], ct.cube_positions[i][2])
		angle := i * 20.0
		model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(float32(angle)), mgl32.Vec3{1, 0.3, 0.5}))
		utils.SetMat4(ct.ShaderProgram, "model", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
}

var lastFrame float64 = 0.0

func (ct *Camera) KeyboardCallback(window *glfw.Window) {

	currentFrame := glfw.GetTime()
	deltaTime := currentFrame - lastFrame
	lastFrame = currentFrame

	cameraSpeed := float32(2.5 * deltaTime)
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		ct.cameraPos = ct.cameraPos.Add(ct.cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		ct.cameraPos = ct.cameraPos.Sub(ct.cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		ct.cameraPos = ct.cameraPos.Sub(ct.cameraFront.Cross(ct.cameraUp).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		ct.cameraPos = ct.cameraPos.Add(ct.cameraFront.Cross(ct.cameraUp).Normalize().Mul(cameraSpeed))
	}
}

func (ct *Camera) MouseCallback(window *glfw.Window, xpos float64, ypos float64) {
	if firstMouse {
		firstMouse = false
		lastxPos = xpos
		lastyPos = ypos
	}

	xoffset := xpos - lastxPos
	yoffset := lastyPos - ypos
	lastxPos = xpos
	lastyPos = ypos

	sensitivity := 0.5
	xoffset *= sensitivity
	yoffset *= sensitivity

	yaw += xoffset
	pitch += yoffset

	// make sure that when pitch is out of bounds, screen doesn't get flipped
	if pitch > 89.0 {
		pitch = 89.0
	}
	if pitch < -89.0 {
		pitch = -89.0
	}

	ct.cameraFront = mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))),
		float32(math.Sin(float64(mgl32.DegToRad(float32(pitch))))),
		float32(math.Sin(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))),
	}
}

func (ct *Camera) ScrollCallback(window *glfw.Window, xoff float64, yoff float64) {
	fov -= yoff
	if fov < 1 {
		fov = 1
	}
	if fov > 45 {
		fov = 45
	}
}
