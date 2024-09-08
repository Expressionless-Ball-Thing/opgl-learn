package renders

type Render interface {
	Draw()

	// set up vertex data (and buffer(s)) and configure vertex attributes
	InitGLPipeLine()
}
