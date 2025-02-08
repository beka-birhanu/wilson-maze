package maze

import game_i "github.com/beka-birhanu/vinom-interfaces/game"

// Cell represents a single cell in a maze grid.
// It includes properties for walls on each side and an associated reward value.
type Cell struct {
	northWall bool
	southWall bool
	eastWall  bool
	westWall  bool
	reward    int32
}

// HasNorthWall returns true if there is a wall on the north side of the cell.
func (c *Cell) HasNorthWall() bool {
	return c.northWall
}

// HasSouthWall returns true if there is a wall on the south side of the cell.
func (c *Cell) HasSouthWall() bool {
	return c.southWall
}

// HasEastWall returns true if there is a wall on the east side of the cell.
func (c *Cell) HasEastWall() bool {
	return c.eastWall
}

// HasWestWall returns true if there is a wall on the west side of the cell.
func (c *Cell) HasWestWall() bool {
	return c.westWall
}

// GetReward returns the reward value assigned to the cell.
func (c *Cell) GetReward() int32 {
	return c.reward
}

// SetNorthWall updates the presence of a wall on the north side of the cell.
func (c *Cell) SetNorthWall(hasWall bool) {
	c.northWall = hasWall
}

// SetSouthWall updates the presence of a wall on the south side of the cell.
func (c *Cell) SetSouthWall(hasWall bool) {
	c.southWall = hasWall
}

// SetEastWall updates the presence of a wall on the east side of the cell.
func (c *Cell) SetEastWall(hasWall bool) {
	c.eastWall = hasWall
}

// SetWestWall updates the presence of a wall on the west side of the cell.
func (c *Cell) SetWestWall(hasWall bool) {
	c.westWall = hasWall
}

// SetReward assigns a new reward value to the cell.
func (c *Cell) SetReward(reward int32) {
	c.reward = reward
}

// CellPosition represents the position of a cell in the maze grid using row and column indices.
type CellPosition struct {
	row int32 // The row index of the cell.
	col int32 // The column index of the cell.
}

// GetRow returns the row index of the cell.
func (cp *CellPosition) GetRow() int32 {
	return cp.row
}

// GetCol returns the column index of the cell.
func (cp *CellPosition) GetCol() int32 {
	return cp.col
}

// SetRow updates the row index of the cell.
func (cp *CellPosition) SetRow(row int32) {
	cp.row = row
}

// SetCol updates the column index of the cell.
func (cp *CellPosition) SetCol(col int32) {
	cp.col = col
}

// Move represents a movement from one cell to another in a specific direction.
type Move struct {
	from game_i.CellPosition
	to   game_i.CellPosition
}

// From returns the starting position of the move.
func (m *Move) From() game_i.CellPosition {
	return m.from
}

// To returns the destination position of the move.
func (m *Move) To() game_i.CellPosition {
	return m.to
}

// SetFrom updates the starting position of the move.
func (m *Move) SetFrom(from game_i.CellPosition) {
	m.from = from
}

// SetTo updates the destination position of the move.
func (m *Move) SetTo(to game_i.CellPosition) {
	m.to = to
}
