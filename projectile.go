package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct{
	Pos rl.Vector2
	StartPos rl.Vector2
	Vel rl.Vector2
	Alive bool
	MaxDist float32
	Radius float32
}

func NewProjectile(pos rl.Vector2, dir rl.Vector2) Projectile{
	return Projectile{
		Pos: rl.Vector2{X: pos.X, Y: pos.Y},
		StartPos: rl.Vector2{X: pos.X, Y: pos.Y},
		Vel: rl.Vector2Scale(dir, 500),
		MaxDist: 800,
		Alive: true,
		Radius: 1.5,
	}
}

func (p *Projectile) Update(dt float32){
	if !p.Alive{
		return
	}

	p.Pos.X += p.Vel.X * dt
	p.Pos.Y += p.Vel.Y * dt

	dx := p.Pos.X - p.StartPos.X
	dy := p.Pos.Y - p.StartPos.Y

	dist := float32(math.Hypot(float64(dx), float64(dy)))

	if dist >= p.MaxDist{
		p.Alive = false
	}

}

func (p Projectile) DrawProjectile(){
	if p.Alive {
		rl.DrawCircleV(p.Pos, 3, rl.Yellow)
	}
}

func RenderProjectiles(projectiles []Projectile){
	for _, p := range projectiles {
		p.DrawProjectile()
	}
}

func UpdateProjectiles(projectiles []Projectile, dt float32)[]Projectile{
	var alive []Projectile
	for i := range projectiles {
		if projectiles[i].Alive {
			projectiles[i].Update(dt)
			if projectiles[i].Alive {
				alive = append(alive, projectiles[i])
			}
		}
	}
	return alive
}

func DetectAsteroidCollision(projectiles []Projectile, asteroids *[]Asteroid, cargoSound rl.Sound, asteroidSound rl.Sound){
	for i := range *asteroids {
		a := &(*asteroids)[i]
		if !a.Alive{
			continue
		}
		for j := range projectiles {
			p := &projectiles[j]
			if !p.Alive {
				continue
			}

			distance := rl.Vector2Distance(p.Pos, a.Pos)
			if distance <= a.Radius + p.Radius {
				a.Alive = false
				p.Alive = false
				if a.IsCargo{
					rl.PlaySound(cargoSound)
				}else{
					rl.PlaySound(asteroidSound)	
				}
				children := a.Split()
				*asteroids = append(*asteroids, children...)
				break
			}
		}
	}
}

