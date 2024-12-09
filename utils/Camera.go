package utils

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// Camera movement enums
const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
	UP
	DOWN
)

// Camera default values

const YAW float64 = -90.0
const PITCH float64 = 0.0
const SPEED float64 = 2.5
const SENSITIVITY float64 = 0.1
const ZOOM float64 = 45.0

type Camera struct {
	// Camera attributes
	Position, Front, Up, Right, WorldUp mgl32.Vec3

	// euler Angles
	Yaw, Pitch float64

	// camera options
	MovementSpeed, MouseSensitivity, Zoom float64
}

// constructor with vectors
func NewCamera(position, up mgl32.Vec3, yaw, pitch float64) Camera {
	c := Camera{
		Position:         position,
		WorldUp:          up,
		Yaw:              yaw,
		Pitch:            pitch,
		Front:            mgl32.Vec3{0.0, 0.0, -1.0},
		MovementSpeed:    SPEED,
		MouseSensitivity: SENSITIVITY,
		Zoom:             ZOOM,
	}
	c.updateCameraVectors()
	return c
}

// constructor with scalar values
func NewCameraWithScalars(posX, posY, posZ, upX, upY, upZ float32, yaw, pitch float64) Camera {
	c := Camera{
		Position:         mgl32.Vec3{posX, posY, posZ},
		WorldUp:          mgl32.Vec3{upX, upY, upZ},
		Yaw:              yaw,
		Pitch:            pitch,
		Front:            mgl32.Vec3{0.0, 0.0, -1.0},
		MovementSpeed:    SPEED,
		MouseSensitivity: SENSITIVITY,
		Zoom:             ZOOM,
	}
	c.updateCameraVectors()
	return c
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

// Processes input received from any keyboard-like input system. Accepts input parameter in the form of camera defined ENUM (to abstract it from windowing systems)
func (c *Camera) ProcessKeyboard(direction int, deltaTime float64) {
	velocity := float32(c.MovementSpeed * deltaTime)

	switch direction {
	case FORWARD:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case BACKWARD:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case LEFT:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case RIGHT:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	case UP:
		c.Position = c.Position.Add(c.Up.Mul(velocity))
	case DOWN:
		c.Position = c.Position.Sub(c.Up.Mul(velocity))
	}
}

// Processes input received from a mouse input system. Expects the offset value in both the x and y direction.
func (c *Camera) ProcessMouseMovement(xoffset, yoffset float64, constrainPitch bool) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	// Make sure that when pitch is out of bounds, screen doesn't get flipped
	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}
	// Update Front, Right and Up Vectors using the updated Eular angles
	c.updateCameraVectors()
}

// Processes input received from a mouse scroll-wheel event. Only requires input on the vertical wheel-axis
func (c *Camera) ProcessMouseScroll(yoffset float64) {
	if c.Zoom >= 1.0 && c.Zoom <= 45.0 {
		c.Zoom -= yoffset
	}
	if c.Zoom <= 1.0 {
		c.Zoom = 1.0
	}
	if c.Zoom >= 45.0 {
		c.Zoom = 45.0
	}
}

// calculates the front vector from the Camera's (updated) Euler Angles
func (c *Camera) updateCameraVectors() {
	x := float32(math.Cos(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch)))
	y := float32(math.Sin(mgl64.DegToRad(c.Pitch)))
	z := float32(math.Sin(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch)))
	front := mgl32.Vec3{x, y, z}
	c.Front = front.Normalize()
	// Also re-calculate the Right and Up vector
	// Normalize the vectors, because their length gets closer to 0 the more you look up or down which results in slower movement.
	c.Right = front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}
