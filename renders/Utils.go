package renders

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Given the filepath to the vertex and fragment shader file, make the shader program
func NewShader(vertexFile, fragmentFile string) (program uint32) {

	vertSrc, err := os.ReadFile(vertexFile)
	if err != nil {
		fmt.Println("Error in reading vertex shader file")
		fmt.Println(err)
		os.Exit(1)
	}

	fragSrc, err2 := os.ReadFile(fragmentFile)
	if err != nil {
		fmt.Println("Error in reading fragment shader file")
		fmt.Println(err2)
		os.Exit(1)
	}

	return CreateShaderProgram(string(vertSrc)+"\x00", string(fragSrc)+"\x00")

}

// build and compile our shader program
func CreateShaderProgram(vertexShaderSource, fragmentShaderSource string) (program uint32) {
	// build and compile our shader program
	// ------------------------------------
	// vertex shader
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	// fragment shader
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	// link shaders
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	// check for linking errors
	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var infoLen int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &infoLen)

		log := strings.Repeat("\x00", int(infoLen+1))
		gl.GetProgramInfoLog(shaderProgram, infoLen, nil, gl.Str(log))

		fmt.Printf("failed to link shader program %v: %v\n", shaderProgram, log)
		os.Exit(1)
	}

	// Deleting shaders because we loaded them into the program.
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram
}

func compileShader(ShaderSource string, ShaderType uint32) uint32 {
	shader := gl.CreateShader(ShaderType)
	csources, free := gl.Strs(ShaderSource)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	// check for shader compile errors
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var infoLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &infoLen)

		log := strings.Repeat("\x00", int(infoLen+1))
		gl.GetShaderInfoLog(shader, infoLen, nil, gl.Str(log))

		fmt.Printf("failed to compile shader %v: %v\n", shader, log)
		os.Exit(1)
	}

	return shader
}

// load and create a texture
func New2DTexture(wrap_s, wrap_t, min_filter, max_filter int32, texturePath string) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture) // all upcoming GL_TEXTURE_2D operations now have effect on this texture object
	// set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, wrap_s) // set texture wrapping to GL_REPEAT (default wrapping method)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, wrap_t)
	// set texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, min_filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, max_filter)

	// load image, create texture and generate mipmaps
	file, err := os.Open(texturePath)
	if err != nil {
		fmt.Println("Cannot open texture file")
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	img, err2 := png.Decode(file)
	if err2 != nil {
		fmt.Println("error decoding the image")
		fmt.Println(err2)
		os.Exit(1)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		fmt.Println("unsupported stride")
		os.Exit(1)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Bounds().Max.X),
		int32(img.Bounds().Max.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}
