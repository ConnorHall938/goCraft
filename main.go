package main

import (
	"fmt"
	"goCraft/lib/camera"
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
	fmt.Printf("Chunk mesh: %d vertices, %d indices\n", len(mesh.Vertices)/8, len(mesh.Indices))

	// Build renderer
	r := render.NewRenderer(window, prog, atlas)
	r.UploadMesh(mesh.Vertices, mesh.Indices)

	camera := camera.NewCamera(mgl32.Vec3{0, 130, 0})
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		camera.ProcessMouse(xpos, ypos)
	})

	lastTime := glfw.GetTime()

	fmt.Println("Ready!")

	// Render loop
	for !window.ShouldClose() {

		now := glfw.GetTime()
		dt := float32(now - lastTime)
		lastTime = now

		// --- keyboard movement ---
		if window.GetKey(glfw.KeyW) == glfw.Press {
			camera.MoveForward(dt)
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			camera.MoveBackward(dt)
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			camera.MoveLeft(dt)
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			camera.MoveRight(dt)
		}
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		w, h := window.GetFramebufferSize()

		view := camera.ViewMatrix()
		proj := mgl32.Perspective(mgl32.DegToRad(60), float32(w)/float32(h), 0.1, 2000.0)

		mvp := proj.Mul4(view)

		r.Render(len(mesh.Indices), mvp)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
