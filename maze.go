/*
Package maze provides tools for creating and managing rectangular mazes.

It defines the `Maze` structure, composed of `Cell` objects that include wall configurations
and optional rewards.

The package includes functionality for random maze generation using Wilson's algorithm, wall manipulation,
and reward assignment. Rewards can be dynamically distributed based on proximity to the maze center.

Utility functions enable neighbor detection, move validation, and ASCII visualization of the maze.
*/
package main

import (
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"strings"

	game_i "github.com/beka-birhanu/vinom-interfaces/game"
)

const (
	// maxMazeDimenssion defines the maximum allowable maze size (width or height).
	maxMazeDimenssion = 20
)

var (
	// Directions maps movement directions (North, South, East, West) to row and column deltas.
	Directions = map[string]CellPosition{
		"North": {row: -1, col: 0},
		"South": {row: 1, col: 0},
		"East":  {row: 0, col: 1},
		"West":  {row: 0, col: -1},
	}

	// ErrInvalidMove is returned when a move request is not valid in the current maze state.
	ErrInvalidMove = errors.New("invalid move request")
)

// WillsonMaze represents a rectangular maze consisting of cells with walls and optional rewards.
type WillsonMaze struct {
	width      int             // The number of columns in the maze.
	height     int             // The number of rows in the maze.
	grid       [][]game_i.Cell // The 2D grid of cells that form the maze.
	totalRward int32           // total reward in the maze.
}

// New initializes a new maze with the given dimensions and generates its layout.
func New(width, height int) (game_i.Maze, error) {
	if min(width, height) <= 0 || max(width, height) > maxMazeDimenssion {
		return nil, fmt.Errorf("invalid maze dimensions")
	}

	// Create a grid of cells with all walls initially intact.
	grid := make([][]game_i.Cell, height)
	for k := range grid {
		grid[k] = make([]game_i.Cell, width)
		for j := range grid[k] {
			grid[k][j] = &Cell{
				northWall: true,
				southWall: true,
				eastWall:  true,
				westWall:  true,
				reward:    0,
			}
		}
	}

	maze := &WillsonMaze{
		width:  width,
		height: height,
		grid:   grid,
	}
	maze.generateMaze()
	return maze, nil
}

// randomCellPosition generates a random position within the maze bounds.
func (m *WillsonMaze) randomCellPosition() game_i.CellPosition {
	return &CellPosition{row: int32(rand.Intn(m.height)), col: int32(rand.Intn(m.width))}
}

// randomUnvisitedCellPosition selects a random cell position that has not been visited.
func (m *WillsonMaze) randomUnvisitedCellPosition(visited map[string]struct{}) game_i.CellPosition {
	for {
		pos := m.randomCellPosition()
		key := fmt.Sprintf("%d,%d", pos.GetRow(), pos.GetCol())
		if _, included := visited[key]; !included {
			return pos
		}
	}
}

// inBoundMoves returns all valid neighboring positions and the moves required to reach them.
func (m *WillsonMaze) inBoundMoves(pos game_i.CellPosition) []game_i.Move {
	var result []game_i.Move
	for _, delta := range Directions {
		neighbor := &CellPosition{row: pos.GetRow() + delta.row, col: pos.GetCol() + delta.col}
		if m.InBound(int(neighbor.row), int(neighbor.col)) {
			result = append(result, &Move{from: pos, to: neighbor})
		}
	}
	return result
}

// getDirection returns direction string to move from pos1 to pos2.
func (m *WillsonMaze) getDirection(pos1, pos2 game_i.CellPosition) string {
	horizontalDir := map[int32]string{1: "West", -1: "East"}
	verticalDir := map[int32]string{-1: "South", 1: "North"}

	if pos1.GetCol() != pos2.GetCol() && pos1.GetRow() != pos2.GetRow() {
		return "Direction Not Allowed"
	}

	if pos1.GetCol() == pos2.GetCol() { // vertical move
		delta := pos1.GetRow() - pos2.GetRow()
		return verticalDir[delta]
	} else {
		delta := pos1.GetCol() - pos2.GetCol()
		return horizontalDir[delta]
	}
}

// openWall removes the wall between two adjacent cells in the specified direction.
func (m *WillsonMaze) openWall(move game_i.Move) error {
	dir := m.getDirection(move.From(), move.To())
	switch dir {
	case "North":
		m.grid[move.From().GetRow()][move.From().GetCol()].SetNorthWall(false)
		m.grid[move.To().GetRow()][move.To().GetCol()].SetSouthWall(false)
	case "South":
		m.grid[move.From().GetRow()][move.From().GetCol()].SetSouthWall(false)
		m.grid[move.To().GetRow()][move.To().GetCol()].SetNorthWall(false)
	case "East":
		m.grid[move.From().GetRow()][move.From().GetCol()].SetEastWall(false)
		m.grid[move.To().GetRow()][move.To().GetCol()].SetWestWall(false)
	case "West":
		m.grid[move.From().GetRow()][move.From().GetCol()].SetWestWall(false)
		m.grid[move.To().GetRow()][move.To().GetCol()].SetEastWall(false)
	}
	return nil
}

// randomWalk performs a random walk starting from an unvisited cell, recording moves taken.
func (m *WillsonMaze) randomWalk(visited map[string]struct{}) []game_i.Move {
	start := m.randomUnvisitedCellPosition(visited)
	visits := make(map[string]game_i.Move)
	cell := start

	for {
		moves := m.inBoundMoves(cell)
		randomNeighbor := moves[rand.Intn(len(moves))].To()
		cellKey := fmt.Sprintf("%d,%d", cell.GetRow(), cell.GetCol())
		visits[cellKey] = &Move{from: cell, to: randomNeighbor}
		neighborKey := fmt.Sprintf("%d,%d", randomNeighbor.GetRow(), randomNeighbor.GetCol())
		if _, included := visited[neighborKey]; included {
			break
		}
		cell = randomNeighbor
	}

	return slices.Collect(maps.Values(visits))
}

// generateMaze creates a maze using Wilson's algorithm.
func (m *WillsonMaze) generateMaze() {
	visited := make(map[string]struct{})
	start := m.randomCellPosition()
	visited[fmt.Sprintf("%d,%d", start.GetRow(), start.GetCol())] = struct{}{}

	for len(visited) < m.width*m.height {
		for _, move := range m.randomWalk(visited) {
			from := move.From()
			to := move.To()
			_ = m.openWall(&Move{from: from, to: to})
			visited[fmt.Sprintf("%d,%d", from.GetRow(), from.GetCol())] = struct{}{}
		}
	}
}

// InBound checks if a position is within the maze bounds.
func (m *WillsonMaze) InBound(row, col int) bool {
	return row >= 0 && row < m.height && col >= 0 && col < m.width
}

// IsValidMove determines if a move is valid, ensuring walls are open and positions are in bounds.
func (m *WillsonMaze) IsValidMove(move game_i.Move) bool {
	if !m.InBound(int(move.From().GetRow()), int(move.From().GetCol())) || !m.InBound(int(move.To().GetRow()), int(move.To().GetCol())) {
		return false
	}

	dir := m.getDirection(move.From(), move.To())
	switch dir {
	case "North":
		return !m.grid[move.From().GetRow()][move.From().GetCol()].HasNorthWall() && !m.grid[move.To().GetRow()][move.To().GetCol()].HasSouthWall()
	case "South":
		return !m.grid[move.From().GetRow()][move.From().GetCol()].HasSouthWall() && !m.grid[move.To().GetRow()][move.To().GetCol()].HasNorthWall()
	case "East":
		return !m.grid[move.From().GetRow()][move.From().GetCol()].HasEastWall() && !m.grid[move.To().GetRow()][move.To().GetCol()].HasWestWall()
	case "West":
		return !m.grid[move.From().GetRow()][move.From().GetCol()].HasWestWall() && !m.grid[move.To().GetRow()][move.To().GetCol()].HasEastWall()
	default:
		return false
	}
}

// Move executes a move in the maze and returns the reward of the destination cell.
func (m *WillsonMaze) Move(move game_i.Move) (int32, error) {
	if !m.IsValidMove(move) {
		return 0, ErrInvalidMove
	}

	reward := m.grid[move.To().GetRow()][move.To().GetCol()].GetReward()
	m.grid[move.To().GetRow()][move.To().GetCol()].SetReward(0)
	return reward, nil
}

// String provides an ASCII visualization of the maze.
func (m *WillsonMaze) String() string {
	var output strings.Builder

	// Top boundary
	output.WriteString("+" + strings.Repeat("---+", m.width) + "\n")

	for row := 0; row < m.height; row++ {
		cellRow := "|"
		for col := 0; col < m.width; col++ {
			cell := m.grid[row][col]

			if cell.GetReward() != 0 {
				cellRow += fmt.Sprintf(" %d ", cell.GetReward())
			} else {
				cellRow += "   "
			}

			if cell.HasEastWall() {
				cellRow += "|"
			} else {
				cellRow += " "
			}
		}
		output.WriteString(cellRow + "\n")

		// Bottom walls
		wallRow := "+"
		for col := 0; col < m.width; col++ {
			cell := m.grid[row][col]
			if cell.HasSouthWall() {
				wallRow += "---+"
			} else {
				wallRow += "   +"
			}
		}
		output.WriteString(wallRow + "\n")
	}

	return output.String()
}

// Height returns the number of rows in the maze.
func (m *WillsonMaze) Height() int {
	return m.height
}

// Width returns the number of columns in the maze.
func (m *WillsonMaze) Width() int {
	return m.width
}

// RemoveReward clears the reward from a specified cell position.
func (m *WillsonMaze) RemoveReward(pos game_i.CellPosition) error {
	if pos.GetRow() < 0 || pos.GetRow() >= int32(m.height) || pos.GetCol() < 0 || pos.GetCol() >= int32(m.width) {
		return fmt.Errorf("position out of bounds")
	}
	m.grid[pos.GetRow()][pos.GetCol()].SetReward(0)
	return nil
}

// GetTotalReward implements game.Maze.
func (m *WillsonMaze) GetTotalReward() int32 {
	return m.totalRward
}

// NewValidMove implements game.Maze.
func (m *WillsonMaze) NewValidMove(curPos game_i.CellPosition, dir string) (game_i.Move, error) {
	delta, ok := Directions[dir]
	if !ok {
		return nil, errors.New("invalid direction")
	}

	nextPos := &CellPosition{
		row: curPos.GetRow() + delta.row,
		col: curPos.GetCol() + delta.col,
	}

	move := &Move{
		from: curPos,
		to:   nextPos,
	}

	if !m.IsValidMove(move) {
		return nil, errors.New("invalid move")
	}

	return move, nil
}

// RetriveGrid returns the 2D grid of cells that make up the maze (implements game.Maze).
func (m *WillsonMaze) RetriveGrid() [][]game_i.Cell {
	return m.grid
}

// SetGrid sets the 2D grid of cells for the maze (implements game.Maze).
func (m *WillsonMaze) SetGrid(g [][]game_i.Cell) {
	m.grid = g
}
