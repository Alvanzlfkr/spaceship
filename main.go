package main

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	width  = 480
	height = 640
)

const (
	PLAY = iota
	GAMEOVER
)

type player struct {
	img  *ebiten.Image
	x, y float64
}

func (p *player) update(shoot *bool) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.x += 5
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.x -= 5
	}

	//control nembak
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		*shoot = true
	}
}

type bullet struct {
	img  *ebiten.Image
	x, y float64
}

type meteor struct {
	img  *ebiten.Image
	x, y float64
}

// Game implements ebiten.Game interface.
type Game struct {
	score int
	player
	bullet
	meteor
	shoot         bool
	meteors       []meteor
	current_scene int
	f             font.Face
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	g.player.update(&g.shoot)

	if g.shoot {
		g.bullet.y -= 15
	}

	if g.bullet.y < 0 {
		g.shoot = false
	}

	if !g.shoot {
		g.bullet.x, g.bullet.y = g.player.x+43, g.player.y
	}

	// push meteor
	for len(g.meteors) < 3 {
		rand.Seed(time.Now().UnixNano())
		g.meteors = append(g.meteors, meteor{
			g.meteor.img,
			float64(rand.Intn(width-28) + 28),
			-50,
		})
	}

	for i := 0; i < len(g.meteors); i++ {
		g.meteors[i].y += 3
		if g.meteors[i].y > height {
			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
		}
	}

	// check collision antara peluru player dengan meteor
	bw, bh := g.bullet.img.Size()
	for i := 0; i < len(g.meteors); i++ {
		mw, mh := g.meteors[i].img.Size()
		if g.bullet.x+float64(bw) >= g.meteors[i].x &&
			g.bullet.x <= g.meteors[i].x+float64(mw) &&
			g.bullet.y+float64(bh) >= g.meteors[i].y &&
			g.bullet.y <= g.meteors[i].y+float64(mh) && g.shoot {
			g.shoot = false
			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
			g.score += 1
		}
	}

	// check collision antara player dengan meteor
	pw, ph := g.player.img.Size()
	for i := 0; i < len(g.meteors); i++ {
		mw, mh := g.meteors[i].img.Size()
		if g.player.x+float64(pw) >= g.meteors[i].x &&
			g.player.x <= g.meteors[i].x+float64(mw) &&
			g.player.y+float64(ph) >= g.meteors[i].y &&
			g.player.y <= g.meteors[i].y+float64(mh) {
			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
			g.score = 0
			g.current_scene = GAMEOVER
		}
	}

	// jika game over tekan "r" untuk restart
	if ebiten.IsKeyPressed(ebiten.KeyR) && g.current_scene == GAMEOVER {
		g.current_scene = PLAY
		g.player.x, g.player.y = width/2, height-100
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.current_scene {
	case PLAY:
		//draw bullet
		bp := &ebiten.DrawImageOptions{}
		bp.GeoM.Translate(g.bullet.x, g.bullet.y)
		screen.DrawImage(g.bullet.img, bp)

		// draw meteor
		for i := 0; i < len(g.meteors); i++ {
			mp := &ebiten.DrawImageOptions{}
			mp.GeoM.Translate(g.meteors[i].x, g.meteors[i].y)
			screen.DrawImage(g.meteors[i].img, mp)
		}

		// draw player
		pp := &ebiten.DrawImageOptions{}
		pp.GeoM.Translate(g.player.x, g.player.y)
		screen.DrawImage(g.player.img, pp)

		//draw score
		text.Draw(screen, "Score: "+strconv.Itoa(g.score), g.f, 10, 20, color.White)
	case GAMEOVER:
		text.Draw(screen, "Game Over", g.f, width/2-77, height/4, color.White)
		text.Draw(screen, "Tekan \"R\" untuk restart", g.f, 0, height/2, color.White)
	}

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	// load font
	load_font, _ := opentype.Parse(fonts.PressStart2P_ttf)

	// load assets
	player_img, _, _ := ebitenutil.NewImageFromFile("playerShip3_red.png")
	bullet_img, _, _ := ebitenutil.NewImageFromFile("laserGreen09.png")
	meteor_img, _, _ := ebitenutil.NewImageFromFile("meteorBrown_small1.png")

	game := &Game{}
	//player
	game.player = player{
		img: player_img,
		x:   width / 2,
		y:   height - 100,
	}

	//bullet
	game.bullet.img = bullet_img
	game.shoot = false

	//meteor
	game.meteor.img = meteor_img

	//font
	game.f, _ = opentype.NewFace(load_font, &opentype.FaceOptions{
		Size:    20,
		DPI:     75,
		Hinting: font.HintingFull,
	})

	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Spaceship")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
