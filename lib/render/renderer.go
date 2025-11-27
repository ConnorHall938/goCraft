package render

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	Window  *glfw.Window
	Program uint32
	Atlas   uint32
	VAO     uint32
	VBO     uint32
	EBO     uint32
}

func NewRenderer(window *glfw.Window, program, atlas uint32) *Renderer {
	r := &Renderer{
		Window:  window,
		Program: program,
		Atlas:   atlas,
	}

	// depth test should be enabled exactly once.
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// Optional but recommended
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	return r
}

// Called once after mesh generation
func (r *Renderer) UploadMesh(vertices []float32, indices []uint32) {
	gl.GenVertexArrays(1, &r.VAO)
	gl.GenBuffers(1, &r.VBO)
	gl.GenBuffers(1, &r.EBO)

	gl.BindVertexArray(r.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, r.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Vertex layout: pos(3), uv(2), tint(3)
	stride := int32(8 * 4)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(5*4))

	gl.BindVertexArray(0)
}

func (r *Renderer) Render(indicesCount int, mvp mgl32.Mat4) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(r.Program)

	loc := gl.GetUniformLocation(r.Program, gl.Str("uMVP\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &mvp[0])

	atlasLoc := gl.GetUniformLocation(r.Program, gl.Str("uAtlas\x00"))
	gl.Uniform1i(atlasLoc, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.Atlas)

	gl.BindVertexArray(r.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(indicesCount), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}
