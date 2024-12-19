package utils

import (
	"strconv"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
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
	for i := 0; i < len(m.Textures); i++ {
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
		SetInt(shader, "material."+m.Textures[i].Type+number, int32(i))
		gl.BindTexture(gl.TEXTURE_2D, m.Textures[i].id)
	}

	//Draw Mesh
	gl.BindVertexArray(m.vao)
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, 0)
	gl.BindVertexArray(0)

	gl.ActiveTexture(gl.TEXTURE0)
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
	gl.VertexAttribIPointerWithOffset(5, 4, gl.INT, structSize32, (unsafe.Offsetof(dummy.m_BoneIDs)))

	// weights
	gl.EnableVertexAttribArray(6)
	gl.VertexAttribPointerWithOffset(6, 4, gl.FLOAT, false, structSize32, (unsafe.Offsetof(dummy.m_Weights)))

	gl.BindVertexArray(0)
}
