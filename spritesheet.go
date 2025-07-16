package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// SpriteSheet represents a collection of sprites in a single image
type SpriteSheet struct {
	image       *ebiten.Image
	tileWidth   int
	tileHeight  int
	tilesPerRow int
	tilesPerCol int
}

// NewSpriteSheet creates a new sprite sheet from an image file
func NewSpriteSheet(imagePath string, tileWidth, tileHeight int) *SpriteSheet {
	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Printf("Failed to load sprite sheet from %s: %v", imagePath, err)
		return nil
	}
	return NewSpriteSheetFromImage(img, tileWidth, tileHeight)
}

// NewSpriteSheetFromImage creates a new sprite sheet from an existing ebiten.Image
func NewSpriteSheetFromImage(img *ebiten.Image, tileWidth, tileHeight int) *SpriteSheet {
	bounds := img.Bounds()
	tilesPerRow := bounds.Dx() / tileWidth
	tilesPerCol := bounds.Dy() / tileHeight

	return &SpriteSheet{
		image:       img,
		tileWidth:   tileWidth,
		tileHeight:  tileHeight,
		tilesPerRow: tilesPerRow,
		tilesPerCol: tilesPerCol,
	}
}

// GetTileImage returns a sub-image of the sprite at the specified index
// Index is calculated as: row * tilesPerRow + col
func (s *SpriteSheet) GetTileImage(index int) *ebiten.Image {
	if s == nil || s.image == nil {
		return nil
	}

	// Calculate row and column from index
	row := index / s.tilesPerRow
	col := index % s.tilesPerRow

	// Check bounds
	if row >= s.tilesPerCol || col >= s.tilesPerRow || row < 0 || col < 0 {
		log.Printf("Sprite index %d out of bounds (max: %d)", index, s.tilesPerRow*s.tilesPerCol-1)
		return nil
	}

	// Calculate source rectangle
	srcX := col * s.tileWidth
	srcY := row * s.tileHeight

	// Create sub-image
	return s.image.SubImage(image.Rect(srcX, srcY, srcX+s.tileWidth, srcY+s.tileHeight)).(*ebiten.Image)
}

// DrawTile draws a specific tile from the sprite sheet at the given screen position
func (s *SpriteSheet) DrawTile(screen *ebiten.Image, tileIndex int, x, y float64, options *ebiten.DrawImageOptions) {
	tileImg := s.GetTileImage(tileIndex)
	if tileImg == nil {
		return
	}

	if options == nil {
		options = &ebiten.DrawImageOptions{}
	}

	// Apply translation
	options.GeoM.Translate(x, y)

	screen.DrawImage(tileImg, options)
}
