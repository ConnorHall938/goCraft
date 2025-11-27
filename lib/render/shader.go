package render

import (
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func LoadShader(path string, shaderType uint32) uint32 {
	srcBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read shader: %v", err)
	}

	src := string(srcBytes) + "\x00"
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, csources, nil)
	free()

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)
		logInfo := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(shader, logLen, nil, gl.Str(logInfo))
		log.Fatalf("shader compile error: %s", logInfo)
	}
	return shader
}

func LinkProgram(vert, frag uint32) uint32 {
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vert)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLen)
		logInfo := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(prog, logLen, nil, gl.Str(logInfo))
		log.Fatalf("failed to link program: %s", logInfo)
	}

	gl.DeleteShader(vert)
	gl.DeleteShader(frag)
	return prog
}
