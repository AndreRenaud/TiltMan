package main

import (
	"math/rand"
	"time"
)

// MazeGenerator creates ASCII mazes of arbitrary size
type MazeGenerator struct {
	width  int
	height int
	maze   [][]rune
	rng    *rand.Rand
}

// Direction represents movement directions for maze generation
type Direction struct {
	dx, dy int
}

var directions = []Direction{
	{0, -2}, // Up
	{2, 0},  // Right
	{0, 2},  // Down
	{-2, 0}, // Left
}

// NewMazeGenerator creates a new maze generator with the specified dimensions
// Width and height should be odd numbers for proper maze structure
func NewMazeGenerator(width, height int) *MazeGenerator {
	// Ensure odd dimensions for proper maze structure
	if width%2 == 0 {
		width++
	}
	if height%2 == 0 {
		height++
	}

	mg := &MazeGenerator{
		width:  width,
		height: height,
		maze:   make([][]rune, height),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Initialize maze grid
	for y := 0; y < height; y++ {
		mg.maze[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			mg.maze[y][x] = '#' // Start with all walls
		}
	}

	return mg
}

// GenerateMaze creates a maze using recursive backtracking algorithm
func (mg *MazeGenerator) GenerateMaze() []string {
	// Start from position (1, 1) - first open cell
	mg.maze[1][1] = '.'
	mg.carvePassages(1, 1)

	// Ensure border is all walls
	mg.ensureBorder()

	// Convert maze to string array
	result := make([]string, mg.height)
	for y := 0; y < mg.height; y++ {
		result[y] = string(mg.maze[y])
	}

	return result
}

// carvePassages uses recursive backtracking to carve maze passages
func (mg *MazeGenerator) carvePassages(x, y int) {
	// Randomize direction order
	dirs := make([]Direction, len(directions))
	copy(dirs, directions)
	mg.shuffle(dirs)

	for _, dir := range dirs {
		nx := x + dir.dx
		ny := y + dir.dy

		// Check if the new position is valid and unvisited
		if mg.isValidCell(nx, ny) && mg.maze[ny][nx] == '#' {
			// Carve the wall between current and new cell
			mg.maze[y+dir.dy/2][x+dir.dx/2] = '.'
			// Mark new cell as passage
			mg.maze[ny][nx] = '.'
			// Recursively carve from new cell
			mg.carvePassages(nx, ny)
		}
	}
}

// isValidCell checks if coordinates are within maze bounds and on valid grid positions
func (mg *MazeGenerator) isValidCell(x, y int) bool {
	return x > 0 && x < mg.width-1 && y > 0 && y < mg.height-1 && x%2 == 1 && y%2 == 1
}

// ensureBorder makes sure all border cells are walls
func (mg *MazeGenerator) ensureBorder() {
	for x := 0; x < mg.width; x++ {
		mg.maze[0][x] = '#'           // Top border
		mg.maze[mg.height-1][x] = '#' // Bottom border
	}
	for y := 0; y < mg.height; y++ {
		mg.maze[y][0] = '#'          // Left border
		mg.maze[y][mg.width-1] = '#' // Right border
	}
}

// shuffle randomizes the order of directions for random maze generation
func (mg *MazeGenerator) shuffle(dirs []Direction) {
	for i := len(dirs) - 1; i > 0; i-- {
		j := mg.rng.Intn(i + 1)
		dirs[i], dirs[j] = dirs[j], dirs[i]
	}
}

// AddSpecialTiles adds special tiles to the maze (speed tiles, etc.)
func (mg *MazeGenerator) AddSpecialTiles(maze []string, density float64) []string {
	if density <= 0 || density > 1 {
		return maze
	}

	result := make([]string, len(maze))
	specialTiles := []rune{'<', '>', '(', ')'}

	for y, row := range maze {
		runes := []rune(row)
		for x, cell := range runes {
			if cell == '.' && mg.rng.Float64() < density {
				// Replace with random special tile
				runes[x] = specialTiles[mg.rng.Intn(len(specialTiles))]
			}
		}
		result[y] = string(runes)
	}

	return result
}

// CreateSimpleMaze creates a basic maze without complex algorithms (for smaller mazes)
func CreateSimpleMaze(width, height int) []string {
	mg := NewMazeGenerator(width, height)
	return mg.GenerateMaze()
}

// CreateMazeWithSpecialTiles creates a maze and adds special speed tiles
func CreateMazeWithSpecialTiles(width, height int, specialTileDensity float64) []string {
	mg := NewMazeGenerator(width, height)
	maze := mg.GenerateMaze()
	return mg.AddSpecialTiles(maze, specialTileDensity)
}
