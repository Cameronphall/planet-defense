package main

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Asteroid struct {
	Pos      rl.Vector2
	SrcRec   rl.Rectangle
	DestRec  rl.Rectangle
	Radius   float32
	Texture  rl.Texture2D
	Alive    bool
	IsCargo  bool
	Velocity rl.Vector2
	Scale    float32
	Color    rl.Color
}

func NewAsteroid(pos rl.Vector2, src rl.Rectangle, dest rl.Rectangle, radius float32, text rl.Texture2D, velocity rl.Vector2, scale float32, color rl.Color) Asteroid {

	return Asteroid{
		Pos:      pos,
		SrcRec:   src,
		DestRec:  dest,
		Radius:   radius,
		Texture:  text,
		Alive:    true,
		IsCargo:  false,
		Velocity: velocity,
		Scale:    scale,
		Color:    color,
	}
}

func SpawnBigAsteroid(planet rl.Vector2, spawnRadius float32, speed float32, src rl.Rectangle, dest rl.Rectangle, texture rl.Texture2D, scale float32) Asteroid {
	angle := rand.Float64() * 2 * math.Pi

	spawnX := planet.X + float32(math.Cos(angle))*spawnRadius
	spawnY := planet.Y + float32(math.Sin(angle))*spawnRadius
	spawnPosition := rl.NewVector2(spawnX, spawnY)

	direction := rl.NewVector2(planet.X-spawnPosition.X, planet.Y-spawnPosition.Y)
	direction = rl.Vector2Normalize(direction)

	velocity := rl.NewVector2(direction.X*speed, direction.Y*speed)

	radius := (dest.Width / 2) * scale
	colorR := 100 + rand.IntN(156)
	colorG := 100 + rand.IntN(156)
	colorB := 100 + rand.IntN(156)
	color := rl.NewColor(uint8(colorR), uint8(colorG), uint8(colorB), 255)
	asteroid := NewAsteroid(spawnPosition, src, dest, radius, texture, velocity, scale, color)
	return asteroid

}

func (a Asteroid) DrawAsteroid(color rl.Color) {
	origin := rl.NewVector2(
		a.DestRec.Width*a.Scale/2,
		a.DestRec.Height*a.Scale/2,
	)
	scaledDest := rl.NewRectangle(
		a.Pos.X,
		a.Pos.Y,
		a.DestRec.Width*a.Scale,
		a.DestRec.Height*a.Scale,
	)
	if a.Alive {
		rl.DrawTexturePro(a.Texture, a.SrcRec, scaledDest, origin, 0, color)
	}
}

func (a *Asteroid) Update(dt float32, planet *Planet, sound rl.Sound) {
	if !a.Alive {
		return
	}
	a.Pos.X += a.Velocity.X * dt
	a.Pos.Y += a.Velocity.Y * dt

	a.DestRec.X = a.Pos.X - (a.DestRec.Width*a.Scale)/2
	a.DestRec.Y = a.Pos.Y - (a.DestRec.Height*a.Scale)/2

	distanceFromPlanet := rl.Vector2Distance(a.Pos, planet.Pos)
	if distanceFromPlanet <= planet.Radius+a.Radius {
		planet.CurrentHealth -= 10
		rl.PlaySound(sound)
		a.Alive = false
	}
}

func RenderAsteroids(asteroids []Asteroid) {
	for _, a := range asteroids {
		a.DrawAsteroid(a.Color)
	}
}

func UpdateAsteroids(asteroids []Asteroid, dt float32, planet *Planet, sound rl.Sound) []Asteroid {
	var alive []Asteroid
	for i := range asteroids {
		if asteroids[i].Alive {
			asteroids[i].Update(dt, planet, sound)
			if asteroids[i].Alive {
				alive = append(alive, asteroids[i])
			}
		}
	}
	return alive
}

func (a Asteroid) Split() []Asteroid {
	var children []Asteroid

	nextScale := float32(0)
	switch a.Scale {
	case 5:
		nextScale = 3

	case 3:
		nextScale = 1

	case 1:
		nextScale = 0.5
	default:
		return children
	}

	for i := 0; i < 2; i++ {
		angle := rand.Float64() * 2 * math.Pi
		direction := rl.NewVector2(float32(math.Cos(angle)), float32(math.Sin(angle)))
		speed := rl.Vector2Length(a.Velocity)
		velocity := rl.NewVector2(direction.X*speed, direction.Y*speed)

		radius := a.Radius * (nextScale / a.Scale)

		colorR := 100 + rand.IntN(156)
		colorG := 100 + rand.IntN(156)
		colorB := 100 + rand.IntN(156)
		color := rl.NewColor(uint8(colorR), uint8(colorG), uint8(colorB), 255)

		child := NewAsteroid(
			a.Pos,
			a.SrcRec,
			a.DestRec,
			radius,
			a.Texture,
			velocity,
			nextScale,
			color,
		)
		if nextScale == 0.5 {
			child.IsCargo = true
			child.Color = rl.NewColor(5, 235, 66, 255)
		}
		children = append(children, child)

	}
	return children
}
