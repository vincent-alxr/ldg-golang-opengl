// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var (
	texture       uint32
	textureUpDown uint32
)

// Structure pour un sommet avec position et coordonnées de texture.
type Vertex struct {
	Position  [3]float32
	TexCoords [2]float32
}

// Vertices et indices pour le rendu.
var vertices []Vertex
var indices []uint32

var faceVertices = map[int][]float32{
	0: {-1, 1, 1, -1, -1, 1, 1, -1, 1, 1, 1, 1},     // Top
	1: {-1, -1, -1, 1, -1, -1, 1, -1, 1, -1, -1, 1}, // Bottom
	2: {-1, 1, -1, -1, -1, -1, -1, -1, 1, -1, 1, 1}, // Left
	3: {1, 1, 1, 1, -1, 1, 1, -1, -1, 1, 1, -1},     // Right
	4: {1, 1, -1, -1, 1, -1, -1, -1, -1, 1, -1, -1}, // Front
	5: {-1, 1, 1, 1, 1, 1, 1, -1, 1, -1, -1, 1},     // Back
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(VertexShader, FragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	const cameraSpeed = 0.05

	var cameraPos = mgl32.Vec3{0, 0, 3}    // Camera position
	var cameraFront = mgl32.Vec3{0, 0, -2} // Direction camera is looking at
	var cameraUp = mgl32.Vec3{0, 1, 0}

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
	viewUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	texture, err = NewTexture("textures/block/grass_block_side.png")
	if err != nil {
		log.Fatalln(err)
	}

	textureUpDown, err = NewTexture("textures/block/grass_block_top.png")
	if err != nil {
		log.Fatalln(err)
	}

	var chunk Chunk = Chunk{
		Size: mgl32.Vec3{16, 64, 16},
	}

	chunk.Initialize()

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(chunk.Vertices)*4, gl.Ptr(chunk.Vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// angle := 0.0
	// previousTime := glfw.GetTime()

	var lastX, lastY float64
	var firstMouse = true

	// Update camera direction here based on offsets
	var yaw float64 = -90
	var pitch float64

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		if window.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled {

			if firstMouse {
				lastX, lastY = xpos, ypos
				firstMouse = false
			}

			xoffset := xpos - lastX
			yoffset := lastY - ypos // Reversed since y-coordinates go from bottom to top
			lastX, lastY = xpos, ypos

			const sensitivity = 0.05
			xoffset *= sensitivity
			yoffset *= sensitivity

			yaw += xoffset // Assuming the initial yaw is -90.0 degrees
			pitch += yoffset

			// Constrain the pitch so the screen doesn't flip
			if pitch > 89.0 {
				pitch = 89.0
			}
			if pitch < -89.0 {
				pitch = -89.0
			}

			front := mgl32.Vec3{
				float32(math.Cos(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))),
				float32(math.Sin(float64(mgl32.DegToRad(float32(pitch))))),
				float32(math.Sin(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))),
			}
			cameraFront = front.Normalize()
		}
	})

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape {
			window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}

	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButton1 && window.GetInputMode(glfw.CursorMode) == glfw.CursorNormal {
			window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		}
	})

	lastTime := glfw.GetTime()
	nbFrames := 0

	for !window.ShouldClose() {
		currentTime := glfw.GetTime()
		nbFrames++
		if currentTime-lastTime >= 1.0 {
			log.Printf("%d ms/frame\n", 1000/nbFrames)
			nbFrames = 0
			lastTime += 1.0
		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if window.GetKey(glfw.KeyW) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
		}

		view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

		// Update
		// time := glfw.GetTime()
		// elapsed := time - previousTime
		// previousTime = time

		// angle += elapsed
		// model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(chunk.Vertices)))
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
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
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 0.0,

	// Back face
	-1.0, -1.0, -1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,

	// Left face
	-1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, -1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 0.0,

	// Right face
	1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 1.0,
	1.0, 1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
}
