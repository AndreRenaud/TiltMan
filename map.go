package main

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// TileType represents different types of tiles in the map
type TileType int

const (
	TileWall TileType = iota
	TileFloor
	TileSlow
	TileFast
	TileSlowMild
	TileFastMild
)

// Tile represents a single tile in the map
type Tile struct {
	Type   TileType
	X, Y   int // Grid coordinates
	Solid  bool
	Effect float64 // Speed multiplier for special tiles
}

// GameMap represents the game map
type GameMap struct {
	Tiles    [][]Tile
	Width    int // Number of tiles horizontally
	Height   int // Number of tiles vertically
	TileSize int // Size of each tile in pixels
	OffsetX  int // X offset for centering the map
	OffsetY  int // Y offset for centering the map
}

// NewGameMap creates a new game map from an ASCII string
func NewGameMap(asciiMap string, tileSize int, screenWidth, screenHeight int) *GameMap {
	lines := strings.Split(strings.TrimSpace(asciiMap), "\n")
	height := len(lines)
	width := 0

	// Find the maximum width
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	// Calculate offsets to center the map
	offsetX := (screenWidth - width*tileSize) / 2
	offsetY := (screenHeight - height*tileSize) / 2

	gameMap := &GameMap{
		Tiles:    make([][]Tile, height),
		Width:    width,
		Height:   height,
		TileSize: tileSize,
		OffsetX:  offsetX,
		OffsetY:  offsetY,
	}

	// Parse the ASCII map
	for y, line := range lines {
		gameMap.Tiles[y] = make([]Tile, width)
		for x, char := range line {
			tile := Tile{X: x, Y: y}

			switch char {
			case '#':
				tile.Type = TileWall
				tile.Solid = true
				tile.Effect = 1.0
			case '.':
				tile.Type = TileFloor
				tile.Solid = false
				tile.Effect = 1.0
			case '<':
				tile.Type = TileSlow
				tile.Solid = false
				tile.Effect = 0.5 // Slow down marble
			case '>':
				tile.Type = TileFast
				tile.Solid = false
				tile.Effect = 1.5 // Speed up marble
			case '(':
				tile.Type = TileSlowMild
				tile.Solid = false
				tile.Effect = 0.75 // Mildly slow down marble
			case ')':
				tile.Type = TileFastMild
				tile.Solid = false
				tile.Effect = 1.25 // Mildly speed up marble
			default:
				// Default to floor for unknown characters
				tile.Type = TileFloor
				tile.Solid = false
				tile.Effect = 1.0
			}

			gameMap.Tiles[y][x] = tile
		}
	}

	return gameMap
}

// GetTileAt returns the tile at the given pixel coordinates
func (m *GameMap) GetTileAt(pixelX, pixelY float64) *Tile {
	// Convert pixel coordinates to grid coordinates
	gridX := int((pixelX - float64(m.OffsetX)) / float64(m.TileSize))
	gridY := int((pixelY - float64(m.OffsetY)) / float64(m.TileSize))

	// Check bounds
	if gridX < 0 || gridX >= m.Width || gridY < 0 || gridY >= m.Height {
		return nil
	}

	return &m.Tiles[gridY][gridX]
}

// IsWallAt checks if there's a wall at the given pixel coordinates
func (m *GameMap) IsWallAt(pixelX, pixelY float64) bool {
	tile := m.GetTileAt(pixelX, pixelY)
	return tile != nil && tile.Solid
}

// GetEffectAt returns the speed effect at the given pixel coordinates
func (m *GameMap) GetEffectAt(pixelX, pixelY float64) float64 {
	tile := m.GetTileAt(pixelX, pixelY)
	if tile != nil {
		return tile.Effect
	}
	return 1.0 // Default effect
}

// CheckCollision checks for collision with walls and returns corrected position
func (m *GameMap) CheckCollision(marble *Marble, newX, newY float64) (float64, float64) {
	radius := marble.Radius

	// Check collision in multiple points around the marble's circumference
	points := []struct{ dx, dy float64 }{
		{-radius, 0},                   // Left
		{radius, 0},                    // Right
		{0, -radius},                   // Top
		{0, radius},                    // Bottom
		{-radius * 0.7, -radius * 0.7}, // Top-left
		{radius * 0.7, -radius * 0.7},  // Top-right
		{-radius * 0.7, radius * 0.7},  // Bottom-left
		{radius * 0.7, radius * 0.7},   // Bottom-right
	}

	for _, point := range points {
		checkX := marble.X + point.dx
		checkY := marble.Y + point.dy

		if m.IsWallAt(checkX, checkY) {
			// Simple collision response - stop movement in the direction of collision
			if point.dx < 0 && marble.VX < 0 { // Left collision
				newX = marble.X
				marble.VX = -marble.VX * 0.3 // Bounce with damping
			} else if point.dx > 0 && marble.VX > 0 { // Right collision
				newX = marble.X
				marble.VX = -marble.VX * 0.3
			}

			if point.dy < 0 && marble.VY < 0 { // Top collision
				newY = marble.Y
				marble.VY = -marble.VY * 0.3
			} else if point.dy > 0 && marble.VY > 0 { // Bottom collision
				newY = marble.Y
				marble.VY = -marble.VY * 0.3
			}
		}
	}

	return newX, newY
}

// ApplyTileEffects applies the effects of the tile the marble is on
func (m *GameMap) ApplyTileEffects(marble *Marble) {
	effect := m.GetEffectAt(marble.X, marble.Y)

	// Apply speed effect
	if effect != 1.0 {
		marble.VX *= effect
		marble.VY *= effect
	}
}

// TileImageCallback is a function type that returns an image for a given tile coordinate
type TileImageCallback func(tiles [][]Tile, x, y int) *ebiten.Image

// Draw renders the map to the screen using a callback to get tile images
func (m *GameMap) Draw(screen *ebiten.Image, getTileImage TileImageCallback) {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			// Calculate pixel position
			pixelX := float64(m.OffsetX + x*m.TileSize)
			pixelY := float64(m.OffsetY + y*m.TileSize)

			// Get the tile image from the callback
			tileImage := getTileImage(m.Tiles, x, y)

			if tileImage != nil {
				// Draw the tile image, scaling it to fit the tile size
				options := &ebiten.DrawImageOptions{}

				options.GeoM.Translate(pixelX, pixelY)

				screen.DrawImage(tileImage, options)
			}
		}
	}
}
