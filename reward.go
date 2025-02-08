package maze

import (
	"fmt"
	"math"
	"math/rand"

	game_i "github.com/beka-birhanu/vinom-interfaces/game"
)

// PopulateReward assigns rewards to maze cells based on the RewardModel.
// The probability of assigning RewardTwo decreases as cells are closer
// to the center of the maze.
func (m *WillsonMaze) PopulateReward(r struct {
	RewardOne      int32
	RewardTwo      int32
	RewardTypeProb float32
}) error {
	if r.RewardTypeProb > 1 || r.RewardTypeProb < 0 || min(r.RewardOne, r.RewardTwo) < 0 {
		return fmt.Errorf("invalid RewardModel")
	}

	visited := map[string]struct{}{}
	stack := []game_i.CellPosition{&CellPosition{row: 0, col: 0}}
	startPosKey := "0,0"
	visited[startPosKey] = struct{}{}
	var totalRewardPlaced int32 = 0

	for len(stack) > 0 {
		cell := pop(&stack)

		// Assign RewardOne as a base reward
		reward := r.RewardOne
		// Adjust probability dynamically and potentially assign RewardTwo
		if rand.Float32() > calcProb(r.RewardTypeProb, cell, m.width, m.height) {
			reward = r.RewardTwo
		}
		m.grid[cell.GetRow()][cell.GetCol()].SetReward(reward)
		totalRewardPlaced += reward

		for _, nbr := range m.inBoundMoves(cell) {
			key := fmt.Sprintf("%d,%d", nbr.To().GetRow(), nbr.To().GetCol())
			if _, seen := visited[key]; !seen {
				visited[key] = struct{}{}
				stack = append(stack, nbr.To())
			}
		}
	}

	m.totalRward = totalRewardPlaced
	return nil
}

// pop removes and returns the last element of a stack of CellPositions.
// The stack is represented as a slice of CellPosition.
func pop(s *[]game_i.CellPosition) game_i.CellPosition {
	lastIndex := len(*s) - 1
	popped := (*s)[lastIndex]
	*s = (*s)[:lastIndex] // Remove the last element
	return popped
}

// calcProb calculates the adjusted probability of assigning RewardTwo
// based on the cell's distance from the center of the maze.
// As the distance to the center decreases, the probability of RewardTwo increases.
func calcProb(baseProb float32, cell game_i.CellPosition, mazeWidth, mazeHeight int) float32 {
	midRow, midCol := int32(mazeHeight/2), int32(mazeWidth/2)

	// Calculate the Manhattan distance to the maze center
	distToMid := math.Abs(float64(cell.GetRow()-midRow)) + math.Abs(float64(cell.GetCol()-midCol))
	maxDist := float64(midRow + midCol)

	// Normalize the distance and invert it for proximity scoring
	normalizedDist := 1.0 - distToMid/maxDist

	// Scale the probability using the base value
	return baseProb + (1-baseProb)*float32(normalizedDist)/10
}
