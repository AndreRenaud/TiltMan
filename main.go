package main

import (
	"embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	tileSize = 32 // Size of each tile in pixels
)

//go:embed assets/*
var assetsFS embed.FS

// orientationChannel is a buffered channel for orientation events - from events_wasm.go
var orientationChannel = make(chan OrientationEvent, 10)

// OrientationEvent represents device orientation data
type OrientationEvent struct {
	Alpha float64 // Z-axis rotation (compass heading)
	Beta  float64 // X-axis rotation (front-to-back tilt)
	Gamma float64 // Y-axis rotation (left-to-right tilt)
}

// Game represents the main game state
type Game struct {
	marble                    *Marble
	gameMap                   *GameMap
	grassSpriteSheet          *SpriteSheet
	stoneSpriteSheet          *SpriteSheet
	screenWidth, screenHeight int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Handle device orientation events (for mobile/web)
	select {
	case event := <-orientationChannel:
		// Convert gamma (left-right tilt) to horizontal force
		// Gamma ranges from -90 to 90 degrees
		gammaForce := event.Gamma / 90.0 * 0.5 // Scale to reasonable force

		// Convert beta (front-back tilt) to vertical force
		// Beta ranges from -180 to 180 degrees, but we'll use -90 to 90
		betaForce := event.Beta / 90.0 * 0.5 // Scale to reasonable force

		// Apply orientation forces
		g.marble.AddForce(gammaForce, betaForce)
	default:
		// No orientation event, use keyboard input for tilt mechanics
	}

	// Handle keyboard input for tilt mechanics
	tiltForce := 0.2

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.marble.AddForce(-tiltForce, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.marble.AddForce(tiltForce, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.marble.AddForce(0, -tiltForce)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.marble.AddForce(0, tiltForce)
	}

	// Reset marble position if R is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.marble.SetPosition(640, 360) // Center of screen
		g.marble.SetVelocity(0, 0)
	}

	// Generate new maze if M is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.generateNewMaze()
	}

	// Update marble physics and get proposed new position
	proposedX, proposedY := g.marble.Update()

	// Apply map collision detection
	finalX, finalY := g.gameMap.CheckCollision(g.marble, proposedX, proposedY)
	g.marble.SetPosition(finalX, finalY)

	// Apply tile effects (speed changes)
	g.gameMap.ApplyTileEffects(g.marble)

	return nil
}

// createTileImage creates a 32x32 colored image for a tile
func createTileImage(tileColor color.Color) *ebiten.Image {
	img := ebiten.NewImage(tileSize, tileSize)
	img.Fill(tileColor)

	vector.DrawFilledRect(img, 0, 0, tileSize, tileSize, tileColor, false)

	return img
}

// getTileImageCallback returns the appropriate tile image for the given coordinates
func (g *Game) getTileImageCallback(m *GameMap, x, y int) *ebiten.Image {
	switch m.GetType(x, y) {
	case TileWall:
		// Use stone spritesheet with smart wall selection
		// Check if there's a non-wall to the right
		solidAround := []bool{
			m.IsSolid(x-1, y-1),
			m.IsSolid(x, y-1),
			m.IsSolid(x+1, y-1),
			m.IsSolid(x-1, y),
			m.IsSolid(x, y),
			m.IsSolid(x+1, y),
			m.IsSolid(x-1, y+1),
			m.IsSolid(x, y+1),
			m.IsSolid(x+1, y+1),
		}
		if !solidAround[1] && !solidAround[3] && !solidAround[5] && !solidAround[7] {
			return g.stoneSpriteSheet.GetTileImageByCoord(3, 3)
		} else if !solidAround[1] && !solidAround[3] && !solidAround[5] {
			return g.stoneSpriteSheet.GetTileImageByCoord(0, 3)
		} else if !solidAround[1] && !solidAround[5] && !solidAround[7] {
			return g.stoneSpriteSheet.GetTileImageByCoord(3, 2)
		} else if !solidAround[1] && !solidAround[3] && !solidAround[7] {
			return g.stoneSpriteSheet.GetTileImageByCoord(3, 0)
		} else if !solidAround[1] && !solidAround[7] {
			return g.stoneSpriteSheet.GetTileImageByCoord(3, 1)
		} else if !solidAround[3] && !solidAround[5] {
			return g.stoneSpriteSheet.GetTileImageByCoord(1, 3)
		} else if !solidAround[7] {
			return g.stoneSpriteSheet.GetTileImageByCoord(2, 1) // Wall with grass below
		} else if !solidAround[5] {
			return g.stoneSpriteSheet.GetTileImageByCoord(1, 2)
		} else if !solidAround[3] {
			return g.stoneSpriteSheet.GetTileImageByCoord(1, 0)
		} else if !solidAround[1] {
			return g.stoneSpriteSheet.GetTileImageByCoord(0, 1)
		}
		return g.stoneSpriteSheet.GetTileImageByCoord(1, 1) // Regular wall
	case TileSlow:
		return createTileImage(color.RGBA{100, 50, 50, 255}) // Red slow tile
	case TileFast:
		return createTileImage(color.RGBA{50, 100, 50, 255}) // Green fast tile
	case TileSlowMild:
		return createTileImage(color.RGBA{80, 50, 60, 255}) // Light red mild slow tile
	case TileFastMild:
		return createTileImage(color.RGBA{50, 80, 60, 255}) // Light green mild fast tile
	case TileFloor:
		fallthrough
	default:
		// Default to floor (grass)
		// These are all grass tiles, so grab one at random, but make sure it's the same for this x/y coordinate
		grassTiles := []struct {
			row, col int
		}{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 0}, {1, 1}, {2, 0}, {2, 1}, {3, 0}, {3, 1}, {3, 2}, {3, 3}}
		tileIndex := ((x + y*m.Width) * 289) % len(grassTiles)
		return g.grassSpriteSheet.GetTileImageByCoord(grassTiles[tileIndex].row, grassTiles[tileIndex].col)
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw the map
	g.gameMap.Draw(screen, g.getTileImageCallback)

	// Draw the marble
	g.marble.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

// generateNewMaze creates a new procedural maze and updates the game map
func (g *Game) generateNewMaze() {
	// Calculate maze dimensions based on screen size and tile size
	// Assume 32x32 tiles to match the existing setup
	mazeWidth := (g.screenWidth / tileSize) - 2 // Leave some border space
	mazeHeight := (g.screenHeight / tileSize) - 2

	// Ensure odd dimensions for proper maze structure
	if mazeWidth%2 == 0 {
		mazeWidth--
	}
	if mazeHeight%2 == 0 {
		mazeHeight--
	}

	mazeLines := CreateMazeWithSpecialTiles(mazeWidth, mazeHeight, 0.15)

	// Convert slice of strings to single string
	mazeStr := ""
	for _, line := range mazeLines {
		mazeStr += line + "\n"
	}

	// Update the game map with the new maze
	g.gameMap = NewGameMap(mazeStr, tileSize, g.screenWidth, g.screenHeight)

	// Reset marble to a safe starting position (top-left open area)
	startX := float64(g.gameMap.OffsetX + 2*tileSize)
	startY := float64(g.gameMap.OffsetY + 2*tileSize)
	g.marble.SetPosition(startX, startY)
	g.marble.SetVelocity(0, 0)
}

func main() {
	// Create a game instance with a marble and map
	game := &Game{
		screenWidth:  1280,
		screenHeight: 720,
	}

	// Print controls information
	log.Println("TiltMan Controls:")
	log.Println("- Arrow keys or WASD: Tilt the board")
	log.Println("- R: Reset marble position")
	log.Println("- M: Generate new random maze")
	log.Println("- On mobile: Tilt your device to control the marble!")

	// Load sprite sheets from embedded filesystem (assuming 32x32 tiles)
	game.grassSpriteSheet = NewSpriteSheetFromFS(assetsFS, "assets/grass.png", tileSize, tileSize)
	game.stoneSpriteSheet = NewSpriteSheetFromFS(assetsFS, "assets/stone.png", tileSize, tileSize)

	if game.grassSpriteSheet == nil {
		log.Fatalf("Warning: Failed to load grass sprite sheet")
	}
	if game.stoneSpriteSheet == nil {
		log.Fatalf("Warning: Failed to load stone sprite sheet")
	}

	// Create marble at starting position (adjust to be within the map)
	startX := float64(2 * tileSize)
	startY := float64(2 * tileSize)
	game.marble = NewMarble(startX, startY, 15, color.RGBA{255, 100, 100, 255})

	game.generateNewMaze()

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("TiltMan")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
