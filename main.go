package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1200, 800, "Asteroid Defense")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	gameOver := false

	rl.SetTargetFPS(60)
	planetText := rl.LoadTexture("Assets/AssetJam/Earth.png")
	asteroidText := rl.LoadTexture("Assets/Asteroids/Asteroid.png")
	shipText := rl.LoadTexture("Assets/Ships/New Piskel.png")
	pew := rl.LoadSound("Assets/Sounds/347171__davidsraba__shoot-sound.wav")
	projectileCollusion := rl.LoadSound("Assets/Sounds/538647__speakwithanimals__vox_glitch_1.wav")
	projectileCargoCollision := rl.LoadSound("Assets/Sounds/327736__distillerystudio__error_03.wav")
	pickupCargoSound := rl.LoadSound("Assets/Sounds/703541__yoshicakes77__coin.ogg")
	dropoffSound := rl.LoadSound("Assets/Sounds/346116__lulyc__retro-game-heal-sound.wav")
	planetAsteroidCollision := rl.LoadSound("Assets/Sounds/421877__sventhors__ouch_1.wav")
	rl.SetSoundVolume(pew, 0.3)
	backgroundMusic := rl.LoadMusicStream("Assets/Music/Space Music Pack/slow-travel.wav")
	rl.PlayMusicStream(backgroundMusic)
	backgroundMusic.Looping = true
	defer rl.UnloadMusicStream(backgroundMusic)
	scale := float32(3)
	var projectiles []Projectile
	var asteroids []Asteroid

	shipSrc := rl.NewRectangle(0, 0, float32(shipText.Width), float32(shipText.Height))
	planetSrc := rl.NewRectangle(0, 0, float32(planetText.Width), float32(planetText.Height))
	asteroidSrc := rl.NewRectangle(0, 0, float32(asteroidText.Width), float32(asteroidText.Height))
	shipW := float32(shipText.Width) * scale
	shipH := float32(shipText.Height) * scale

	shipOrigin := rl.NewVector2(shipW/2, shipH/2)

	asteroidDest := rl.NewRectangle(100, 200, float32(asteroidText.Width*int32(scale)), float32(asteroidText.Height*int32(scale)))
	planetDest := rl.NewRectangle(650, 500, float32(planetText.Width*int32(scale)), float32(planetText.Height*int32(scale)))

	planetCenter := rl.NewVector2(
		planetDest.X+planetDest.Width/2,
		planetDest.Y+planetDest.Height/2,
	)
	orbitDistance := float32(200)
	shipCenter := rl.NewVector2(planetCenter.X+orbitDistance, planetCenter.Y)
	shipDest := rl.NewRectangle(shipCenter.X, shipCenter.Y, shipW, shipH)
	ship := NewShip(shipCenter, 300, shipText, shipSrc, shipDest, shipOrigin, 0, rl.White)
	planet := NewPlanet(planetCenter, planetSrc, planetDest, planetText, planetDest.Width/2*0.55, 100, 100)
	camera := rl.NewCamera2D(
		rl.NewVector2(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2),
		ship.Pos,
		0.0,
		1.0,
	)

	spawnTimer := float32(0)
	spawnDelay := float32(4)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		if !gameOver {
			spawnTimer += dt
			if spawnTimer >= spawnDelay {
				spawnTimer = 0
				asteroids = append(asteroids, SpawnBigAsteroid(planet.Pos, 800, 150, asteroidSrc, asteroidDest, asteroidText, 3))
			}

			UpdateAsteroids(asteroids, dt, &planet, planetAsteroidCollision)
			UpdateProjectiles(projectiles, dt)
			CheckCargoPickup(&ship, &asteroids, pickupCargoSound)
			CheckCargoDropoff(&ship, &planet, dropoffSound)
			DetectAsteroidCollision(projectiles, &asteroids, projectileCargoCollision, projectileCollusion)
			if planet.CurrentHealth <= 0 {
				gameOver = true
			}
			ship.Update(dt, camera)
			camera.Target = ship.Pos
		} else {
			if rl.IsKeyPressed(rl.KeyR) {
				gameOver = false
				planet.CurrentHealth = planet.MaxHealth
				ship.Cargo = 0
				ship.Pos = rl.NewVector2(planet.Pos.X+orbitDistance, planet.Pos.Y)
				ship.DestRec.X = ship.Pos.X
				ship.DestRec.Y = ship.Pos.Y
				ship.Rotation = 0
				asteroids = nil
				projectiles = nil
				spawnTimer = 0

			}
		}

		rl.UpdateMusicStream(backgroundMusic)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(camera)

		RenderProjectiles(projectiles)
		RenderAsteroids(asteroids)

		planet.DrawPlanet()
		ship.DrawShip()
		if rl.IsKeyPressed(rl.KeySpace) {
			projectile := ship.Shoot(camera)
			projectiles = append(projectiles, projectile)
			rl.PlaySound(pew)
		}

		rl.EndMode2D()
		DrawHealthBar(50, 30, 300, 25, planet.CurrentHealth, planet.MaxHealth)
		rl.DrawText(fmt.Sprintf("Planet Health: %.0f / %.0f", planet.CurrentHealth, planet.MaxHealth), 50, 10, 10, rl.White)
		rl.DrawText(fmt.Sprintf("Cargo: %.0f", ship.Cargo), 200, 10, 10, rl.White)
		if gameOver {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Black)
			rl.DrawText("Game Over!", 450, 350, 50, rl.White)
			rl.DrawText("Press R to Restart", 350, 420, 50, rl.White)
		}
		rl.EndDrawing()
	}
}
