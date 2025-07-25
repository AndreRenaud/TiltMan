package main

import (
	"image"
	"io/fs"
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

// NewSpriteSheetFromFS creates a new sprite sheet from an embedded filesystem
func NewSpriteSheetFromFS(filesystem fs.FS, imagePath string, tileWidth, tileHeight int) *SpriteSheet {
	file, err := filesystem.Open(imagePath)
	if err != nil {
		log.Printf("Failed to open sprite sheet from embedded FS %s: %v", imagePath, err)
		return nil
	}
	defer file.Close()

	img, _, err := ebitenutil.NewImageFromReader(file)
	if err != nil {
		log.Printf("Failed to create image from embedded FS %s: %v", imagePath, err)
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

// GetTileImageByCoord returns a sub-image of the sprite at the specified row and column
func (s *SpriteSheet) GetTileImageByCoord(row, col int) *ebiten.Image {
	if s == nil || s.image == nil {
		return nil
	}

	// Check bounds
	if row >= s.tilesPerCol || col >= s.tilesPerRow || row < 0 || col < 0 {
		log.Printf("Sprite coordinates (%d, %d) out of bounds (max: %d, %d)", row, col, s.tilesPerCol-1, s.tilesPerRow-1)
		return nil
	}

	// Calculate source rectangle
	srcX := col * s.tileWidth
	srcY := row * s.tileHeight

	// Create sub-image
	return s.image.SubImage(image.Rect(srcX, srcY, srcX+s.tileWidth, srcY+s.tileHeight)).(*ebiten.Image)
}
