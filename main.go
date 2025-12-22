package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

const gridSize = 30
const ringBellsCount = 50

type GridSquare struct {
	Row             uint
	Col             uint
	Fleas           []*Flea
	AdjacentSquares []*GridSquare
}

func NewGridSquare(row, col uint) GridSquare {
	return GridSquare{
		Row: row,
		Col: col,
	}
}

func (s *GridSquare) ToString(withAdjacentSquares bool) string {
	if !withAdjacentSquares {
		return fmt.Sprintf("[%d, %d]: %d", s.Row, s.Col, len(s.Fleas))
	}

	var adjacentSquares []string

	for _, as := range s.AdjacentSquares {
		adjacentSquares = append(adjacentSquares, fmt.Sprintf("[%d, %d]", as.Row, as.Col))
	}

	return fmt.Sprintf("[%d, %d]: %d -> %s", s.Row, s.Col, len(s.Fleas), strings.Join(adjacentSquares, ", "))
}

type Grid struct {
	Squares [gridSize][gridSize]*GridSquare
	Fleas   []*Flea
}

type Flea struct {
	Square *GridSquare
}

func NewFlea(square *GridSquare) Flea {
	return Flea{
		Square: square,
	}
}

func NewGrid() Grid {
	var g Grid
	g.Fleas = []*Flea{}

	for row := uint(0); row < gridSize; row++ {
		for col := uint(0); col < gridSize; col++ {
			var square = NewGridSquare(row, col)
			var flea = NewFlea(&square)
			square.Fleas = []*Flea{&flea}
			g.Squares[row][col] = &square
			g.Fleas = append(g.Fleas, &flea)
		}
	}

	for row := uint(0); row < gridSize; row++ {
		for col := uint(0); col < gridSize; col++ {
			g.Squares[row][col].AdjacentSquares = getAdjacentSquares(&g, g.Squares[row][col])
		}
	}

	return g
}

func (g *Grid) Print(withAdjacentSquares bool) {
	for row := uint(0); row < gridSize; row++ {
		for col := uint(0); col < gridSize; col++ {
			fmt.Printf("%v\n", g.Squares[row][col].ToString(withAdjacentSquares))
		}
	}

	fmt.Println()
	fmt.Printf("Unoccupied fields: %d\n", countUnoccupiedFields(g))
}

type Move struct {
	X int
	Y int
}

func newMove(x, y int) Move {
	return Move{
		X: x,
		Y: y,
	}
}

var allMoves = []Move{
	newMove(-1, 0),
	newMove(0, -1),
	newMove(0, 1),
	newMove(1, 0),
}

func getAdjacentSquares(grid *Grid, square *GridSquare) []*GridSquare {
	var adjacentSquares []*GridSquare

	for _, m := range allMoves {
		var newRow = int(square.Row) + m.Y
		var newCol = int(square.Col) + m.X

		if newRow >= 0 && newRow < gridSize && newCol >= 0 && newCol < gridSize {
			adjacentSquares = append(adjacentSquares, grid.Squares[uint(newRow)][uint(newCol)])
		}
	}

	return adjacentSquares
}

func countUnoccupiedFields(g *Grid) uint {
	count := uint(0)

	for row := uint(0); row < gridSize; row++ {
		for col := uint(0); col < gridSize; col++ {
			if len(g.Squares[row][col].Fleas) == 0 {
				count++
			}
		}
	}

	return count
}

func removeFleaFromSquare(s *GridSquare, f *Flea) {
	var fleaIndex = -1

	for i, sf := range s.Fleas {
		if f == sf {
			fleaIndex = i
			break
		}
	}

	if fleaIndex >= 0 {
		s.Fleas = append(s.Fleas[:fleaIndex], s.Fleas[fleaIndex+1:]...)
	}
}

func ringBell(g *Grid) {
	for _, f := range g.Fleas {
		// remove flea from current square
		removeFleaFromSquare(f.Square, f)

		// move flea to a random adjacent square
		nextSquare := f.Square.AdjacentSquares[rand.IntN(len(f.Square.AdjacentSquares))]
		f.Square = nextSquare
		nextSquare.Fleas = append(nextSquare.Fleas, f)
	}
}

func checkFleasCount(g *Grid) {
	count := 0

	for row := uint(0); row < gridSize; row++ {
		for col := uint(0); col < gridSize; col++ {
			count += len(g.Squares[row][col].Fleas)
		}
	}

	const expectedCount = gridSize * gridSize
	if count != expectedCount {
		var msg = fmt.Sprintf("Expected %d fleas, but %d are present", expectedCount, count)
		panic(msg)
	}
}

func runSimulation() uint {
	var g = NewGrid()
	//g.Print(true)

	for rb := 1; rb <= ringBellsCount; rb++ {
		//checkFleasCount(&g)
		ringBell(&g)
		//g.Print(false)
	}

	return countUnoccupiedFields(&g)
}

func main() {
	const simulationsCount = 10000
	var sum uint

	var start = time.Now()

	for i := 1; i <= simulationsCount; i++ {
		sum += runSimulation()
	}

	var elapsedTime = time.Now().Sub(start)
	fmt.Printf("Average count of unoccupied fields: %f\n", float32(sum)/float32(simulationsCount))
	fmt.Printf("Elapsed time: %s", elapsedTime)
}
