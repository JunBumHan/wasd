package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

const (
	gameWidth  = 20
	gameHeight = 10
)

type Player struct {
	x, y int
}

type Game struct {
	player Player
	width  int
	height int
}

func NewGame() *Game {
	return &Game{
		player: Player{x: gameWidth / 2, y: gameHeight / 2},
		width:  gameWidth,
		height: gameHeight,
	}
}

func (g *Game) clearScreen() {
	fmt.Print("\033[2J\033[H") // ANSI escape codes to clear screen and move cursor to top-left
}

func (g *Game) render() {
	g.clearScreen()
	
	// Draw the game area
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if x == g.player.x && y == g.player.y {
				fmt.Print("ᦂ")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	
	fmt.Println("\nUse W/A/S/D to move, Q to quit")
}

func (g *Game) movePlayer(dx, dy int) {
	newX := g.player.x + dx
	newY := g.player.y + dy
	
	// Boundary checking
	if newX >= 0 && newX < g.width && newY >= 0 && newY < g.height {
		g.player.x = newX
		g.player.y = newY
	}
}

func enableRawMode() {
	exec.Command("stty", "-echo", "cbreak").Run()
}

func disableRawMode() {
	exec.Command("stty", "echo", "-cbreak").Run()
}

func readChar() (byte, error) {
	reader := bufio.NewReader(os.Stdin)
	char, err := reader.ReadByte()
	return char, err
}

func main() {
	game := NewGame()
	
	// Set up signal handling to restore terminal on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		disableRawMode()
		fmt.Println("\nGoodbye!")
		os.Exit(0)
	}()
	
	fmt.Println("Terminal WASD Game")
	fmt.Println("Use W/A/S/D to move the player (ᦂ)")
	fmt.Println("Press Q to quit")
	fmt.Println("Press any key to start...")
	
	// Wait for any key to start
	readChar()
	
	// Enable raw mode for better input
	enableRawMode()
	defer disableRawMode()
	
	// Game loop
	game.render()
	
	for {
		char, err := readChar()
		if err != nil {
			break
		}
		
		switch char {
		case 'w', 'W':
			game.movePlayer(0, -1) // Move up
		case 's', 'S':
			game.movePlayer(0, 1)  // Move down
		case 'a', 'A':
			game.movePlayer(-1, 0) // Move left
		case 'd', 'D':
			game.movePlayer(1, 0)  // Move right
		case 'q', 'Q':
			disableRawMode()
			fmt.Println("Goodbye!")
			return
		}
		
		game.render()
	}
}