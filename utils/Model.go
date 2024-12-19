package utils

import (
	"fmt"
	"path"

	"github.com/bloeys/assimp-go/asig"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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
	material := scene.Materials[mesh.MaterialIndex]
	diffuseMaps := m.loadMaterialTextures(material, asig.TextureTypeDiffuse, "texture_diffuse")
	textures = append(textures, diffuseMaps...)

	specularMaps := m.loadMaterialTextures(material, asig.TextureTypeSpecular, "texture_specular")
	textures = append(textures, specularMaps...)

	normalMaps := m.loadMaterialTextures(material, asig.TextureTypeNormal, "texture_normal")
	textures = append(textures, normalMaps...)

	heightMaps := m.loadMaterialTextures(material, asig.TextureTypeHeight, "texture_height")
	textures = append(textures, heightMaps...)

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
	return New2DTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, filename)
}
