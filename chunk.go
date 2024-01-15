package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Chunk struct {
	Size mgl32.Vec3

	Blocks [][][]*Block

	Vertices []float32
}

const TOP = 0
const BOTTOM = 1
const LEFT = 2
const RIGHT = 3
const FRONT = 4
const BACK = 5

func (chunk *Chunk) Initialize() {
	// Initialise les blocks et leur faces visible
	chunk.Blocks = make([][][]*Block, int(chunk.Size[0]))
	chunk.Vertices = []float32{}

	addVertices := func(face int, x, y, z float32) []float32 {
		switch face {
		case 0:
			{
				return []float32{
					-1.0 + x, 1.0 + y, -1.0 + z, 0.0, 0.0,
					-1.0 + x, 1.0 + y, 1.0 + z, 0.0, 1.0,
					1.0 + x, 1.0 + y, -1.0 + z, 1.0, 0.0,
					1.0 + x, 1.0 + y, -1.0 + z, 1.0, 0.0,
					-1.0 + x, 1.0 + y, 1.0 + z, 0.0, 1.0,
					1.0 + x, 1.0 + y, 1.0 + z, 1.0, 1.0,
				}
			}
		case 1:
			{
				return []float32{
					-1.0 + x, -1.0 + y, -1.0 + z, 0.0, 0.0,
					1.0 + x, -1.0 + y, -1.0 + z, 1.0, 0.0,
					-1.0 + x, -1.0 + y, 1.0 + z, 0.0, 1.0,
					1.0 + x, -1.0 + y, -1.0 + z, 1.0, 0.0,
					1.0 + x, -1.0 + y, 1.0 + z, 1.0, 1.0,
					-1.0 + x, -1.0 + y, 1.0 + z, 0.0, 1.0,
				}
			}
		case 2:
			{
				return []float32{
					-1.0 + x, -1.0 + y, 1.0 + z, 1.0, 1.0,
					-1.0 + x, -1.0 + y, -1.0 + z, 0.0, 1.0,
					-1.0 + x, 1.0 + y, 1.0 + z, 1.0, 0.0,
					-1.0 + x, -1.0 + y, -1.0 + z, 0.0, 1.0,
					-1.0 + x, 1.0 + y, -1.0 + z, 0.0, 0.0,
					-1.0 + x, 1.0 + y, 1.0 + z, 1.0, 0.0,
				}
			}
		case 3:
			{
				return []float32{
					1.0 + x, -1.0 + y, 1.0 + z, 0.0, 1.0,
					1.0 + x, 1.0 + y, 1.0 + z, 0.0, 0.0,
					1.0 + x, -1.0 + y, -1.0 + z, 1.0, 1.0,
					1.0 + x, -1.0 + y, -1.0 + z, 1.0, 1.0,
					1.0 + x, 1.0 + y, 1.0 + z, 0.0, 0.0,
					1.0 + x, 1.0 + y, -1.0 + z, 1.0, 0.0,
				}
			}
		case 4:
			{
				return []float32{
					-1.0 + x, -1.0 + y, 1.0 + z, 0.0, 1.0,
					1.0 + x, -1.0 + y, 1.0 + z, 1.0, 1.0,
					-1.0 + x, 1.0 + y, 1.0 + z, 0.0, 0.0,
					1.0 + x, -1.0 + y, 1.0 + z, 1.0, 1.0,
					1.0 + x, 1.0 + y, 1.0 + z, 1.0, 0.0,
					-1.0 + x, 1.0 + y, 1.0 + z, 0.0, 0.0,
				}
			}
		case 5:
			{
				return []float32{
					-1.0 + x, -1.0 + y, -1.0 + z, 1.0, 1.0,
					-1.0 + x, 1.0 + y, -1.0 + z, 1.0, 0.0,
					1.0 + x, -1.0 + y, -1.0 + z, 0.0, 1.0,
					1.0 + x, -1.0 + y, -1.0 + z, 0.0, 1.0,
					-1.0 + x, 1.0 + y, -1.0 + z, 1.0, 0.0,
					1.0 + x, 1.0 + y, -1.0 + z, 0.0, 0.0,
				}
			}
		default:
			{
				return []float32{}
			}
		}

	}

	// Initialisez chaque slice de la deuxième dimension.
	for x := range chunk.Blocks {
		chunk.Blocks[x] = make([][]*Block, int(chunk.Size[1]))

		// Initialisez chaque slice de la troisième dimension.
		for y := range chunk.Blocks[x] {
			chunk.Blocks[x][y] = make([]*Block, int(chunk.Size[2]))

			for z := range chunk.Blocks[x][y] {

				if x > 10 && y > 20 && y < 40 {
					chunk.Blocks[x][y][z] = &Block{
						Position:     mgl32.Vec3{float32(x), float32(y), float32(z)},
						Solid:        false,
						FacesVisible: [6]bool{false, false, false, false, false, false},
					}
				} else {
					chunk.Blocks[x][y][z] = &Block{
						Position:     mgl32.Vec3{float32(x), float32(y), float32(z)},
						Solid:        true,
						FacesVisible: [6]bool{false, false, false, false, false, false},
					}
				}
			}
		}
	}

	for x := 0; x < int(chunk.Size[0]); x++ {

		for y := 0; y < int(chunk.Size[1]); y++ {

			for z := 0; z < int(chunk.Size[2]); z++ {
				// neighbors := chunk.GetNeighbors(x, y, z)
				if chunk.Blocks[x][y][z].Solid {
					chunk.Blocks[x][y][z].FacesVisible[0] = !chunk.HasSolidNeighbor(TOP, x, y, z)       // Haut
					chunk.Blocks[x][y][z].FacesVisible[1] = !chunk.HasSolidNeighbor(BOTTOM, x, y, z)    // Bas
					chunk.Blocks[x][y][z].FacesVisible[LEFT] = !chunk.HasSolidNeighbor(LEFT, x, y, z)   // Gauche
					chunk.Blocks[x][y][z].FacesVisible[RIGHT] = !chunk.HasSolidNeighbor(RIGHT, x, y, z) // Droite
					chunk.Blocks[x][y][z].FacesVisible[4] = !chunk.HasSolidNeighbor(FRONT, x, y, z)     // Avant
					chunk.Blocks[x][y][z].FacesVisible[5] = !chunk.HasSolidNeighbor(BACK, x, y, z)      // Arrière

					block := chunk.Blocks[x][y][z]
					for face := 0; face < 6; face++ {
						if block.FacesVisible[face] {
							chunk.Vertices = append(chunk.Vertices, addVertices(face, block.Position.X()*2, block.Position.Y()*2, block.Position.Z()*2)...)
						}
					}
				}
			}
		}
	}

}

func (chunk *Chunk) Render() {

}

func (chunk *Chunk) GetVertices() {

}

func (chunk *Chunk) HasNeighbor(dir, x, y, z int) bool {
	isValidCoordinate := func(coord, maxCoord int) bool {
		return coord >= 0 && coord < maxCoord
	}

	switch dir {
	case LEFT:
		return isValidCoordinate(x-1, int(chunk.Size.X())) && chunk.Blocks[x-1][y][z] != nil
	case RIGHT:
		return isValidCoordinate(x+1, int(chunk.Size.X())) && chunk.Blocks[x+1][y][z] != nil
	case BOTTOM:
		return isValidCoordinate(y-1, int(chunk.Size.Y())) && chunk.Blocks[x][y-1][z] != nil
	case TOP:
		return isValidCoordinate(y+1, int(chunk.Size.Y())) && chunk.Blocks[x][y+1][z] != nil
	case BACK:
		return isValidCoordinate(z-1, int(chunk.Size.Z())) && chunk.Blocks[x][y][z-1] != nil
	case FRONT:
		return isValidCoordinate(z+1, int(chunk.Size.Z())) && chunk.Blocks[x][y][z+1] != nil
	}

	return false
}

func (chunk *Chunk) HasSolidNeighbor(dir, x, y, z int) bool {
	if chunk.HasNeighbor(dir, x, y, z) {
		switch dir {
		case LEFT:
			return chunk.Blocks[x-1][y][z].Solid
		case RIGHT:
			return chunk.Blocks[x+1][y][z].Solid
		case BOTTOM:
			return chunk.Blocks[x][y-1][z].Solid
		case TOP:
			return chunk.Blocks[x][y+1][z].Solid
		case BACK:
			return chunk.Blocks[x][y][z-1].Solid
		case FRONT:
			return chunk.Blocks[x][y][z+1].Solid
		}
	}
	return false
}
