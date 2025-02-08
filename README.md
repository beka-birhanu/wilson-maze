# Wilson Maze  

A maze generation package using Wilson's algorithm . Provides functionality to generate mazes with probabilistic rewards on paths. 

## Usage  

### Creating a Maze  
```go
New(width, height int) (game_i.Maze, error)
```
Creates a new maze with the given dimensions. 

### Populating Rewards  
```go
func (m *WillsonMaze) PopulateReward(r struct {
	RewardOne      int32
	RewardTwo      int32
	RewardTypeProb float32
}) error 
```
Assigns rewards to the maze cells based on probability. 

 Note:
- The function takes two reward values and a probability for selecting the second one. 
- It assumes the second reward is larger and reduces its probability toward the edges of the maze using Manhattan distance from center. 

### Maze Interaction  
- Check if a position is in bounds:
 ```go
 InBound(row, col int) bool
 ```
- Validate a move:
 ```go
 IsValidMove(move game_i.Move) bool
 ```
- Move within the maze:
 ```go
 Move(move game_i.Move) (int32, error)
 ```
- Get the total remaining reward:
 ```go
 GetTotalReward() int32
 ```
- Generate a new valid move:
 ```go
 NewValidMove(curPos game_i.CellPosition, dir string) (game_i.Move, error)
 ```
