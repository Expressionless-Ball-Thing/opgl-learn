package renders

import (
	"fmt"
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
