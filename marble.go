package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Marble represents a marble with physics properties
type Marble struct {
	X, Y     float64 // Position
	VX, VY   float64 // Velocity
	Radius   float64 // Radius of the marble
	Color    color.Color
	Friction float64 // Friction coefficient (0-1, where 1 = no friction)
}

// NewMarble creates a new marble at the specified position
func NewMarble(x, y, radius float64, c color.Color) *Marble {
	return &Marble{
		X:        x,
		Y:        y,
		VX:       0,
		VY:       0,
		Radius:   radius,
		Color:    c,
		Friction: 0.98, // Default friction
	}
}

// Update calculates the marble's new position based on its velocity and returns it
// The caller is responsible for checking collisions and applying the new position
func (m *Marble) Update() (newX, newY float64) {
	// Calculate new position based on velocity
	newX = m.X + m.VX
	newY = m.Y + m.VY

	// Apply friction to gradually slow down the marble
	m.VX *= m.Friction
	m.VY *= m.Friction

	// Stop very small movements to prevent jitter
	if math.Abs(m.VX) < 0.01 {
		m.VX = 0
	}
	if math.Abs(m.VY) < 0.01 {
		m.VY = 0
	}

	return newX, newY
}

// AddForce adds a force to the marble (for tilt mechanics)
func (m *Marble) AddForce(fx, fy float64) {
	m.VX += fx
	m.VY += fy
}

// SetPosition sets the marble's position
func (m *Marble) SetPosition(x, y float64) {
	m.X = x
	m.Y = y
}

// GetPosition returns the marble's current position
func (m *Marble) GetPosition() (float64, float64) {
	return m.X, m.Y
}

// SetVelocity sets the marble's velocity
func (m *Marble) SetVelocity(vx, vy float64) {
	m.VX = vx
	m.VY = vy
}

// GetVelocity returns the marble's current velocity
func (m *Marble) GetVelocity() (float64, float64) {
	return m.VX, m.VY
}

// Draw renders the marble to the screen
func (m *Marble) Draw(screen *ebiten.Image) {
	// Draw the marble as a filled circle
	vector.DrawFilledCircle(screen, float32(m.X), float32(m.Y), float32(m.Radius), m.Color, true)

	// Draw a subtle highlight to make it look more 3D
	highlightColor := color.RGBA{255, 255, 255, 100}
	highlightX := float32(m.X - m.Radius*0.3)
	highlightY := float32(m.Y - m.Radius*0.3)
	highlightRadius := float32(m.Radius * 0.3)
	vector.DrawFilledCircle(screen, highlightX, highlightY, highlightRadius, highlightColor, true)
}
