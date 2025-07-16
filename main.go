package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game represents the main game state
type Game struct {
	marble                    *Marble
	gameMap                   *GameMap
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

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw the map
	g.gameMap.Draw(screen)

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
##################
#........#.......#
#...<........>...#
#........#.......#
####.#########.###
#........#.......#
#...>........<...#
#........#.......#
##################`

	// Create a game instance with a marble and map
	game := &Game{
		screenWidth:  1280,
		screenHeight: 720,
	}

	// Create the map with 40x40 pixel tiles
	game.gameMap = NewGameMap(sampleMap, 40, game.screenWidth, game.screenHeight)

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
