package basics

import "github.com/go-gl/gl/v3.3-core/gl"

type EmptyWindow struct {
}

func (ew *EmptyWindow) InitGLPipeLine() {

}

func (ew *EmptyWindow) Draw() {
	gl.ClearColor(0.5, 0.5, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
