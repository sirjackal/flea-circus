package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

const gridSize = 30
const ringBellsCount = 50
const simulationsCount = 10000
const workersCount = 128

type GridSquare struct {
	Row             int
	Col             int
	Fleas           []*Flea
	AdjacentSquares []*GridSquare
}

func NewGridSquare(row, col int) GridSquare {
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

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			var square = NewGridSquare(row, col)
			var flea = NewFlea(&square)
			square.Fleas = []*Flea{&flea}
			g.Squares[row][col] = &square
			g.Fleas = append(g.Fleas, &flea)
		}
	}

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			g.Squares[row][col].AdjacentSquares = getAdjacentSquares(&g, g.Squares[row][col])
		}
	}

	return g
}

func (g *Grid) Print(withAdjacentSquares bool) {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
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
			adjacentSquares = append(adjacentSquares, grid.Squares[newRow][newCol])
		}
	}

	return adjacentSquares
}

func countUnoccupiedFields(g *Grid) int {
	count := 0

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
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

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			count += len(g.Squares[row][col].Fleas)
		}
	}

	const expectedCount = gridSize * gridSize
	if count != expectedCount {
		var msg = fmt.Sprintf("Expected %d fleas, but %d are present", expectedCount, count)
		panic(msg)
	}
}

func runSingleSimulation() int {
	var g = NewGrid()
	//g.Print(true)

	for rb := 1; rb <= ringBellsCount; rb++ {
		//checkFleasCount(&g)
		ringBell(&g)
		//g.Print(false)
	}

	return countUnoccupiedFields(&g)
}

func worker(jobs <-chan int, results chan<- int) {
	for range jobs {
		results <- runSingleSimulation()
	}
}

func runParallelSimulation(simulationsCount int) {
	var start = time.Now()

	jobs := make(chan int, simulationsCount)
	results := make(chan int, simulationsCount)

	for w := 1; w <= workersCount; w++ {
		go worker(jobs, results)
	}

	for j := 1; j <= simulationsCount; j++ {
		jobs <- j
	}
	close(jobs)

	var sum int

	for a := 1; a <= simulationsCount; a++ {
		sum += <-results
	}
	close(results)

	var elapsedTime = time.Since(start)
	fmt.Printf("Average count of unoccupied fields: %.6f\n", float32(sum)/float32(simulationsCount))
	fmt.Printf("Elapsed time: %s", elapsedTime)
}

func runSequentialSimulation(simulationsCount int) {
	var start = time.Now()
	var sum int

	for i := 1; i <= simulationsCount; i++ {
		sum += runSingleSimulation()
	}

	var elapsedTime = time.Since(start)
	fmt.Printf("Average count of unoccupied fields: %.6f\n", float32(sum)/float32(simulationsCount))
	fmt.Printf("Elapsed time: %s", elapsedTime)
}

func main() {
	//runSequentialSimulation(simulationsCount)
	runParallelSimulation(simulationsCount)
}
