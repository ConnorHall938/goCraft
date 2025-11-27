package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position mgl32.Vec3
	Yaw      float32 // left–right
	Pitch    float32 // up–down

	Speed      float32
	MouseSens  float32
	FirstMouse bool
	LastX      float64
	LastY      float64
}

func NewCamera(pos mgl32.Vec3) *Camera {
	return &Camera{
		Position:   pos,
		Yaw:        -90.0, // looking forward on Z–
		Pitch:      0.0,
		Speed:      10.0,
		MouseSens:  0.1,
		FirstMouse: true,
	}
}

func (c *Camera) ViewMatrix() mgl32.Mat4 {
	dir := c.Direction()
	return mgl32.LookAtV(c.Position, c.Position.Add(dir), mgl32.Vec3{0, 1, 0})
}

func (c *Camera) Direction() mgl32.Vec3 {
	yawRad := mgl32.DegToRad(c.Yaw)
	pitchRad := mgl32.DegToRad(c.Pitch)

	x := float32(math.Cos(float64(yawRad)) * math.Cos(float64(pitchRad)))
	y := float32(math.Sin(float64(pitchRad)))
	z := float32(math.Sin(float64(yawRad)) * math.Cos(float64(pitchRad)))

	return mgl32.Vec3{x, y, z}.Normalize()
}

func (c *Camera) MoveForward(dt float32) {
	c.Position = c.Position.Add(c.Direction().Mul(dt * c.Speed))
}
func (c *Camera) MoveBackward(dt float32) {
	c.Position = c.Position.Sub(c.Direction().Mul(dt * c.Speed))
}

func (c *Camera) MoveRight(dt float32) {
	right := c.Direction().Cross(mgl32.Vec3{0, 1, 0}).Normalize()
	c.Position = c.Position.Add(right.Mul(dt * c.Speed))
}

func (c *Camera) MoveLeft(dt float32) {
	right := c.Direction().Cross(mgl32.Vec3{0, 1, 0}).Normalize()
	c.Position = c.Position.Sub(right.Mul(dt * c.Speed))
}

func (c *Camera) ProcessMouse(x, y float64) {
	if c.FirstMouse {
		c.LastX = x
		c.LastY = y
		c.FirstMouse = false
		return
	}

	dx := (x - c.LastX) * float64(c.MouseSens)
	dy := (c.LastY - y) * float64(c.MouseSens)

	c.LastX = x
	c.LastY = y

	c.Yaw += float32(dx)
	c.Pitch += float32(dy)

	if c.Pitch > 89.0 {
		c.Pitch = 89.0
	}
	if c.Pitch < -89.0 {
		c.Pitch = -89.0
	}
}
