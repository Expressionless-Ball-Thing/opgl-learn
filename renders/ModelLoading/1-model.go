package ModelLoading

import (
	"fmt"
	"opgl-learn/utils"
	"path"
	"strconv"
	"unsafe"

	"github.com/bloeys/assimp-go/asig"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position, Normal, Tangent, Bitangent mgl32.Vec3
	TexCoords                            mgl32.Vec2
	m_BoneIDs                            [4]int
	m_Weights                            [4]float32
}

type Texture struct {
	id   uint32
	Type string
	Path string
}

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
	Textures []Texture

	vao, vbo, ebo uint32
}

func NewMesh(vertices []Vertex, indices []uint32, textures []Texture) Mesh {

	mesh := Mesh{
		Vertices: vertices,
		Indices:  indices,
		Textures: textures,
	}

	mesh.setupMesh()
	return mesh
}

func (m *Mesh) Draw(shader uint32) {
	var diffuseNr uint64 = 1
	var specularNr uint64 = 1
	var normalNr uint64 = 1
	var heightNr uint64 = 1
	for i := 1; i < len(m.Textures); i++ {
		number := ""
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		switch m.Textures[i].Type {
		case "texture_diffuse":
			number += strconv.FormatUint(diffuseNr, 10)
			diffuseNr++
		case "texture_specular":
			number += strconv.FormatUint(specularNr, 10)
			specularNr++
		case "texture_normal":
			number += strconv.FormatUint(normalNr, 10)
			normalNr++
		case "texture_height":
			number += strconv.FormatUint(heightNr, 10)
			heightNr++
		}
		utils.SetInt(shader, "material."+m.Textures[i].Type+number, int32(i))
		gl.BindTexture(gl.TEXTURE_2D, m.Textures[i].id)
	}

	gl.ActiveTexture(gl.TEXTURE0)

	//Draw Mesh
	gl.BindVertexArray(m.vao)
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, 0)
	gl.BindVertexArray(0)
}

func (m *Mesh) setupMesh() {
	// size of the Vertex struct
	var dummy Vertex
	structSize := int(unsafe.Sizeof(dummy))
	structSize32 := int32(structSize)

	// Create buffers/arrays
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	gl.BindVertexArray(m.vao)
	// Load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*structSize, gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*4, gl.Ptr(m.Indices), gl.STATIC_DRAW)

	// Vertex Positions
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, structSize32, (unsafe.Offsetof(dummy.Position)))

	// Vertex Normals
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, structSize32, (unsafe.Offsetof(dummy.Normal)))

	// Vertex Texture Coords
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, structSize32, (unsafe.Offsetof(dummy.TexCoords)))

	// Vertext Tangent
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(3, 3, gl.FLOAT, false, structSize32, (unsafe.Offsetof(dummy.Tangent)))

	// Vertex Bitangent
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointerWithOffset(4, 3, gl.FLOAT, false, structSize32, (unsafe.Offsetof(dummy.Bitangent)))

	// ids
	gl.EnableVertexAttribArray(5)
	gl.VertexAttribPointerWithOffset(5, 4, gl.INT, false, structSize32, (unsafe.Offsetof(dummy.m_BoneIDs)))

	// weights
	gl.EnableVertexAttribArray(6)
	gl.VertexAttribPointerWithOffset(6, 4, gl.FLOAT, false, structSize32, (unsafe.Offsetof(dummy.m_Weights)))

	gl.BindVertexArray(0)
}

type Model struct {
	Meshes         []Mesh
	Directory      string
	LoadedTextures map[string]Texture
}

func NewModel(path string) Model {
	model := Model{}
	model.LoadedTextures = make(map[string]Texture)
	model.loadModel(path)
	return model
}

func (m *Model) Draw(shader uint32) {
	for i := 0; i < len(m.Meshes); i++ {
		m.Meshes[i].Draw(shader)
	}
}

func (m *Model) loadModel(filepath string) {
	scene, _, err := asig.ImportFile(filepath, asig.PostProcessTriangulate|asig.PostProcessFlipUVs)

	if err != nil || scene.Flags&asig.SceneFlagIncomplete != 0 {
		fmt.Printf("ERROR::ASSIMP:: %d\n", scene.Flags)
		return
	}

	dir, _ := path.Split(filepath)
	m.Directory = dir
	m.processNode(scene.RootNode, scene)
}

func (m *Model) processNode(node *asig.Node, scene *asig.Scene) {

	for i := 0; i < len(node.MeshIndicies); i++ {
		mesh := scene.Meshes[node.MeshIndicies[i]]
		m.Meshes = append(m.Meshes, m.processMesh(mesh, scene))
	}

	for i := 0; i < len(node.Children); i++ {
		m.processNode(node.Children[i], scene)
	}
}

func (m *Model) processMesh(mesh *asig.Mesh, scene *asig.Scene) Mesh {

	vertices := []Vertex{}
	indices := []uint32{}
	textures := []Texture{}

	for i := 0; i < len(mesh.Vertices); i++ {

		vertex := Vertex{}
		// Process vertex Positions, normals and texture coords
		vertex.Position = mgl32.Vec3{mesh.Vertices[i].X(), mesh.Vertices[i].Y(), mesh.Vertices[i].Z()}
		vertex.Normal = mgl32.Vec3{mesh.Normals[i].X(), mesh.Normals[i].Y(), mesh.Normals[i].Z()}
		if len(mesh.TexCoords) > 0 {
			vertex.TexCoords = mgl32.Vec2{mesh.TexCoords[0][i].X(), mesh.TexCoords[0][i].Y()}
		} else {
			vertex.TexCoords = mgl32.Vec2{0, 0}
		}
		if len(mesh.Tangents) > 0 {
			vertex.Tangent = mgl32.Vec3{mesh.Tangents[i].X(), mesh.Tangents[i].Y(), mesh.Tangents[i].Z()}
		} else {
			vertex.Tangent = mgl32.Vec3{0, 0, 0}
		}
		if len(mesh.BitTangents) > 0 {
			vertex.Bitangent = mgl32.Vec3{mesh.BitTangents[i].X(), mesh.BitTangents[i].Y(), mesh.BitTangents[i].Z()}
		} else {
			vertex.Bitangent = mgl32.Vec3{0, 0, 0}
		}
		vertices = append(vertices, vertex)
	}

	// process indices
	for i := 0; i < len(mesh.Faces); i++ {
		face := mesh.Faces[i]
		for j := 0; j < len(face.Indices); j++ {
			indices = append(indices, uint32(face.Indices[j]))
		}
	}

	// process material
	if mesh.MaterialIndex >= 1 {
		material := scene.Materials[mesh.MaterialIndex]
		diffuseMaps := m.loadMaterialTextures(material, asig.TextureTypeDiffuse, "texture_diffuse")
		textures = append(textures, diffuseMaps...)

		specularMaps := m.loadMaterialTextures(material, asig.TextureTypeSpecular, "texture_specular")
		textures = append(textures, specularMaps...)

		normalMaps := m.loadMaterialTextures(material, asig.TextureTypeNormal, "texture_normal")
		textures = append(textures, normalMaps...)

		heightMaps := m.loadMaterialTextures(material, asig.TextureTypeHeight, "texture_height")
		textures = append(textures, heightMaps...)
	}

	return NewMesh(vertices, indices, textures)
}

func (m *Model) loadMaterialTextures(mat *asig.Material, matType asig.TextureType, typeName string) []Texture {
	textures := []Texture{}

	for i := 0; i < asig.GetMaterialTextureCount(mat, matType); i++ {

		info, _ := asig.GetMaterialTexture(mat, matType, uint(i))
		if val, ok := m.LoadedTextures[info.Path]; ok {
			textures = append(textures, val)
		} else {
			_, filename := path.Split(info.Path)
			texture := Texture{
				id:   TextureFromFile(filename, m.Directory, true),
				Type: typeName,
				Path: info.Path,
			}
			textures = append(textures, texture)
			m.LoadedTextures[info.Path] = texture
		}
	}

	return textures
}

func TextureFromFile(path string, directory string, gamma bool) uint32 {
	filename := directory + path

	var textureID uint32
	gl.GenTextures(1, &textureID)

	return utils.New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, filename)
}

type ModelLoad struct {
	ShaderProgram uint32
	camera        utils.Camera
	model         Model
}

var lastxPos float64 = 1920 / 2.0
var lastyPos float64 = 1080 / 2.0
var firstMouse bool = true
var lastFrame float64 = 0.0

func (ct *ModelLoad) InitGLPipeLine() {

	ct.camera = utils.NewCamera(mgl32.Vec3{0.0, 0.0, 3.0}, mgl32.Vec3{0, 1, 0}, utils.YAW, utils.PITCH)

	ct.ShaderProgram = utils.NewShader("./shaders/ModelLoading/1-ModelVert.glsl", "./shaders/ModelLoading/1-ModelFrag.glsl")
	ct.model = NewModel("./backpack/backpack.obj")

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
