package main

import (
	"fmt"
	"goCraft/lib/chunk"
	"goCraft/lib/render"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() { runtime.LockOSThread() }

func main() {
	glfw.Init()
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(1200, 800, "GoCraft", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	if err := gl.Init(); err != nil {
		panic(err)
	}

	vert := render.LoadShader("assets/shaders/cube.vert", gl.VERTEX_SHADER)
	frag := render.LoadShader("assets/shaders/cube.frag", gl.FRAGMENT_SHADER)
	prog := render.LinkProgram(vert, frag)

	atlas := render.LoadAtlas([]string{
		"assets/textures/grass_top.png",
		"assets/textures/grass_side.png",
		"assets/textures/dirt.png",
		"assets/textures/stone.png",
		"assets/textures/log_oak.png",
		"assets/textures/log_oak_top.png",
	})

	// Build a chunk mesh
	ch := chunk.NewChunk()
	ch.Fill(1) // fill with grass
	mesh := ch.BuildMesh()

	// Build renderer
	r := render.NewRenderer(window, prog, atlas)
	r.UploadMesh(mesh.Vertices, mesh.Indices)

	fmt.Println("Ready!")

	// Render loop
	for !window.ShouldClose() {
		w, h := window.GetFramebufferSize()

		eye := mgl32.Vec3{60, 80, 60}
		center := mgl32.Vec3{16, 64, 16}
		up := mgl32.Vec3{0, 1, 0}

		view := mgl32.LookAtV(eye, center, up)
		proj := mgl32.Perspective(mgl32.DegToRad(60), float32(w)/float32(h), 0.1, 2000.0)

		mvp := proj.Mul4(view)

		r.Render(len(mesh.Indices), mvp)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
