package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"strings"

	"goCraft/lib/atlas"
	"goCraft/lib/cube"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
)

func init() {
	// OpenGL requires the main thread
	runtime.LockOSThread()
}

func main() {
	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// Configure OpenGL version 4.6 Core Profile
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Textured Cube", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("OpenGL version:", gl.GoStr(gl.GetString(gl.VERSION)))

	// --- Compile shaders ---
	program := initShaderProgram("assets/shaders/cube.vert", "assets/shaders/cube.frag")

	// --- Load texture ---
	texture := loadAllTextures([]string{
		"assets/textures/grass_block_top.png",
		"assets/textures/grass_block_side.png",
		"assets/textures/dirt.png",
		"assets/textures/stone.png",
		"assets/textures/oak_log.png",
		"assets/textures/oak_log_top.png",
	})

	// --- Cube vertices ---
	verts, inds := cube.MakeCube(
		atlas.TextureIndex["wood_end"],
		atlas.TextureIndex["wood_end"],
		atlas.TextureIndex["wood_side"],
	)

	// --- Setup OpenGL buffers ---
	var VAO, VBO, EBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	// Upload vertex data
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	// Upload indices
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(inds)*4, gl.Ptr(inds), gl.STATIC_DRAW)
	// Vertex attribute 0 = position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	// Vertex attribute 1 = texture coordinates
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.BindVertexArray(0)
	gl.Enable(gl.DEPTH_TEST)

	// --- Main loop ---
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		time := float32(glfw.GetTime())

		// Camera
		eye := mgl32.Vec3{2, 2, 4}
		center := mgl32.Vec3{0, 0, 0}
		up := mgl32.Vec3{0, 1, 0}

		view := mgl32.LookAtV(eye, center, up)

		// Projection
		proj := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)

		// Model rotation (rotate on two axes to expose all faces)
		model := mgl32.HomogRotate3DY(time * 0.7).Mul4(
			mgl32.HomogRotate3DX(time * 0.4),
		)

		mvp := proj.Mul4(view).Mul4(model)

		// Upload to shader
		loc := gl.GetUniformLocation(program, gl.Str("uMVP\x00"))
		gl.UniformMatrix4fv(loc, 1, false, &mvp[0])

		// Bind texture
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		// Draw cube
		gl.BindVertexArray(VAO)
		gl.DrawElements(gl.TRIANGLES, int32(len(inds)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.BindVertexArray(0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// --- Shader loading/compiling ---
func initShaderProgram(vertexPath, fragmentPath string) uint32 {
	vertexShader := compileShader(readFile(vertexPath), gl.VERTEX_SHADER)
	fragmentShader := compileShader(readFile(fragmentPath), gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		logInfo := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logInfo))
		panic("failed to link program: " + logInfo)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return program
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logInfo := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logInfo))
		panic("failed to compile shader: " + logInfo)
	}

	return shader
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	return string(data)
}

// Map texture names to atlas index
var TextureIndex = map[string]int{
	"grass_top":  0,
	"grass_side": 1,
	"dirt":       2,
	"stone":      3,
	"wood_side":  4,
	"wood_end":   5,
}

const (
	AtlasCols = 4 // 4 textures per row
	AtlasRows = 2 // 2 textures per column
)

func loadAllTextures(paths []string) uint32 {
	rgba := image.NewRGBA(image.Rect(0, 0, 1024*4, 1024*2))
	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		x := (i % 4) * 1024
		y := (i / 4) * 1024
		draw.Draw(rgba, image.Rect(x, y, x+1024, y+1024), img, image.Point{}, draw.Src)
		file.Close()
	}
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8,
		int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	// Wrap and filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	return tex
}
