package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ship struct {
	Pos      rl.Vector2
	Speed    float32
	Texture  rl.Texture2D
	SrcRec   rl.Rectangle
	DestRec  rl.Rectangle
	Origin   rl.Vector2
	Rotation float32
	Color    rl.Color
	Cargo    float32
}

func (s *Ship) Update(dt float32, camera rl.Camera2D) {

	mousePosition := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
	dx := mousePosition.X - s.Pos.X
	dy := mousePosition.Y - s.Pos.Y
	s.Rotation = float32(math.Atan2(float64(dy), float64(dx)))*rl.Rad2deg + 90

	toMouse := rl.Vector2{
		X: mousePosition.X - s.Pos.X,
		Y: mousePosition.Y - s.Pos.Y,
	}

	len := float32(math.Hypot(float64(toMouse.X), float64(toMouse.Y)))
	if len < 5 {
		toMouse.X = 0
		toMouse.Y = 0
	}
	if len != 0 {
		toMouse.X /= len
		toMouse.Y /= len
	}

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		s.Pos.X += toMouse.X * s.Speed * dt
		s.Pos.Y += toMouse.Y * s.Speed * dt
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		s.Pos.X -= toMouse.X * s.Speed * dt
		s.Pos.Y -= toMouse.Y * s.Speed * dt
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		s.Pos.X += toMouse.Y * s.Speed * dt
		s.Pos.Y += -toMouse.X * s.Speed * dt
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		s.Pos.X += -toMouse.Y * s.Speed * dt
		s.Pos.Y += toMouse.X * s.Speed * dt
	}
	s.DestRec.X = s.Pos.X
	s.DestRec.Y = s.Pos.Y

}

func NewShip(pos rl.Vector2, speed float32, texture rl.Texture2D, src rl.Rectangle, dest rl.Rectangle, origin rl.Vector2, rotation float32, color rl.Color) Ship {
	dest.X = pos.X
	dest.Y = pos.Y
	return Ship{
		Pos:      pos,
		Speed:    speed,
		Texture:  texture,
		SrcRec:   src,
		DestRec:  dest,
		Origin:   origin,
		Rotation: rotation,
		Color:    color,
		Cargo:    0,
	}
}

func (s Ship) DrawShip() {
	rl.DrawTexturePro(s.Texture, s.SrcRec, s.DestRec, s.Origin, s.Rotation, s.Color)
}

func (s *Ship) Shoot(camera rl.Camera2D) Projectile {
	mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)

	dir := rl.Vector2{
		X: mousePos.X - s.Pos.X,
		Y: mousePos.Y - s.Pos.Y,
	}
	if l := float32(rl.Vector2Length(dir)); l != 0 {
		dir.X /= l
		dir.Y /= l
	}

	nose := rl.Vector2Add(s.Pos, rl.Vector2Scale(dir, s.DestRec.Height*0.5))

	return NewProjectile(nose, dir)
}

func CheckCargoPickup(ship *Ship, asteroids *[]Asteroid, pickupSound rl.Sound) {
	for i := range *asteroids {
		a := &(*asteroids)[i]
		if !a.Alive || !a.IsCargo {
			continue
		}

		distance := rl.Vector2Distance(ship.Pos, a.Pos)
		pickupRange := (a.Radius * a.Scale) + 20

		if distance <= pickupRange {
			ship.Cargo += 1
			a.Alive = false
			rl.PlaySound(pickupSound)
		}
	}
}

func CheckCargoDropoff(ship *Ship, planet *Planet, dropoffSound rl.Sound) {
	if planet.CurrentHealth == planet.MaxHealth || ship.Cargo == 0{
		return
	}

	distance := rl.Vector2Distance(ship.Pos, planet.Pos)
	dropoffRange := (planet.Radius) + 20

	if distance <=  dropoffRange{
		amountToHeal := planet.MaxHealth - planet.CurrentHealth
		if amountToHeal >= ship.Cargo{
			planet.CurrentHealth += ship.Cargo
			ship.Cargo = 0
		}else{
			planet.CurrentHealth += amountToHeal
			ship.Cargo -= amountToHeal
		}
		
		rl.PlaySound(dropoffSound)
	}

}
