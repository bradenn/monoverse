package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"io/ioutil"
	"strings"
)

func newProgramShader(vertexShaderPath, fragmentShaderPath string) (uint32, error) {
	vertexShader, err := loadShader(vertexShaderPath, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fragmentShader, err := loadShader(vertexShaderPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logSize int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logSize)
		log := strings.Repeat("\x00", int(logSize+1))
		gl.GetProgramInfoLog(program, logSize, nil, gl.Str(log))
		return 0, nil
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil

}

func loadShader(path string, shaderVariant uint32) (uint32, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	legacyCString := gl.Str(string(bytes) + "\x00")

	shader := gl.CreateShader(shaderVariant)
	gl.ShaderSource(shader, 1, &legacyCString, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logSize int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logSize)
		log := strings.Repeat("\x00", int(logSize+1))
		gl.GetShaderInfoLog(shader, logSize, nil, gl.Str(log))
		return 0, nil
	}
	return shader, nil
}
