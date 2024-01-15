package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Block struct {
	Position     mgl32.Vec3
	Solid        bool
	FacesVisible [6]bool
}

func (block *Block) Render(modelUniform int32) {

	model := mgl32.Translate3D(block.Position.X(), block.Position.Y(), block.Position.Z())
	// Add any other transformations here (rotation, scaling, etc.)

	// Send the model matrix to the shader
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// Bind texture for the cube
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.DrawArrays(gl.TRIANGLES, 1*2*6, 2*2*6)

	gl.BindTexture(gl.TEXTURE_2D, textureUpDown)
	gl.DrawArrays(gl.TRIANGLES, 0, 1*2*6)
}
