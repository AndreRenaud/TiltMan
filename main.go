package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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
	img := ebiten.NewImage(32, 32)
	img.Fill(tileColor)

	vector.DrawFilledRect(img, 0, 0, 32, 32, tileColor, false)

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

		if !solidAround[7] { // Check if there's grass (floor) below
			return g.stoneSpriteSheet.GetTileImageByCoord(2, 1) // Wall with grass below
		} else if !solidAround[5] { // Check if there's a non-wall to the right
			return g.stoneSpriteSheet.GetTileImageByCoord(1, 2) // Wall with right edge
		}
		return g.stoneSpriteSheet.GetTileImageByCoord(1, 1) // Regular wall
	case TileFloor:
		// Use grass spritesheet at position (0,0)
		return g.grassSpriteSheet.GetTileImageByCoord(0, 0)
	case TileSlow:
		return createTileImage(color.RGBA{100, 50, 50, 255}) // Red slow tile
	case TileFast:
		return createTileImage(color.RGBA{50, 100, 50, 255}) // Green fast tile
	case TileSlowMild:
		return createTileImage(color.RGBA{80, 50, 60, 255}) // Light red mild slow tile
	case TileFastMild:
		return createTileImage(color.RGBA{50, 80, 60, 255}) // Light green mild fast tile
	default:
		// Default to floor (grass)
		return g.grassSpriteSheet.GetTileImageByCoord(0, 0)
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

func main() {
	// Define a sample map
	sampleMap := `
######################
#........#...........#
#...<....(...)..>....#
#........#...........#
####.#########.....###
#........#...........#
#...)....(...<..>....#
#........#...........#
######################`

	// Create a game instance with a marble and map
	game := &Game{
		screenWidth:  1280,
		screenHeight: 720,
	}

	// Load sprite sheets (assuming 32x32 tiles)
	game.grassSpriteSheet = NewSpriteSheet("assets/grass.png", 32, 32)
	game.stoneSpriteSheet = NewSpriteSheet("assets/stone.png", 32, 32)

	if game.grassSpriteSheet == nil {
		log.Fatalf("Warning: Failed to load grass sprite sheet")
	}
	if game.stoneSpriteSheet == nil {
		log.Fatalf("Warning: Failed to load stone sprite sheet")
	}

	// Create the map with 32x32 pixel tiles (matching sprite size)
	game.gameMap = NewGameMap(sampleMap, 32, game.screenWidth, game.screenHeight)

	// Create marble at starting position (adjust to be within the map)
	startX := float64(game.gameMap.OffsetX + 2*game.gameMap.TileSize) // Start in an open area
	startY := float64(game.gameMap.OffsetY + 2*game.gameMap.TileSize)
	game.marble = NewMarble(startX, startY, 15, color.RGBA{255, 100, 100, 255})

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("TiltMan")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
