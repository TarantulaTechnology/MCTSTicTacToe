package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Node struct {
	state       GameState
	parent      *Node
	children    []*Node
	visits      int
	totalReward float64
}

type GameState interface {
	GetPossibleActions() []Action
	PerformAction(Action) GameState
	IsTerminal() bool
	GetReward() float64
}

type Action interface{}

func NewNode(state GameState, parent *Node) *Node {
	return &Node{
		state:    state,
		parent:   parent,
		children: []*Node{},
	}
}

func (n *Node) UCTValue() float64 {
	if n.visits == 0 {
		return math.Inf(1)
	}
	return n.totalReward/float64(n.visits) + 1.41*math.Sqrt(math.Log(float64(n.parent.visits))/float64(n.visits))
}

func (n *Node) Select() *Node {
	if n.IsLeaf() {
		return n
	}
	bestValue := math.Inf(-1)
	var bestChild *Node
	for _, child := range n.children {
		childValue := child.UCTValue()
		if childValue > bestValue {
			bestValue = childValue
			bestChild = child
		}
	}
	return bestChild.Select()
}

func (n *Node) Expand() {
	actions := n.state.GetPossibleActions()
	for _, action := range actions {
		newState := n.state.PerformAction(action)
		newNode := NewNode(newState, n)
		n.children = append(n.children, newNode)
	}
}

func (n *Node) Backpropagate(reward float64) {
	n.visits++
	n.totalReward += reward
	if n.parent != nil {
		n.parent.Backpropagate(reward)
	}
}

func (n *Node) IsLeaf() bool {
	return len(n.children) == 0
}

func MCTS(rootState GameState, iterations int) *Node {
	root := NewNode(rootState, nil)
	for i := 0; i < iterations; i++ {
		node := root.Select()
		if !node.state.IsTerminal() {
			node.Expand()
		}
		child := node.Select()
		reward := child.state.GetReward()
		child.Backpropagate(reward)
	}
	return root
}

///////////////////////////////////////

const (
	Empty = iota
	X
	O
)

type Board [3][3]int

type TicTacToeState struct {
	board  Board
	player int
}

func (s *TicTacToeState) GetPossibleActions() []Action {
	var actions []Action
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.board[i][j] == Empty {
				actions = append(actions, [2]int{i, j})
			}
		}
	}
	return actions
}

func (s *TicTacToeState) PerformAction(action Action) GameState {
	move := action.([2]int)
	newBoard := s.board
	newBoard[move[0]][move[1]] = s.player
	return &TicTacToeState{board: newBoard, player: 3 - s.player}
}

func (s *TicTacToeState) IsTerminal() bool {
	return s.checkWin(X) || s.checkWin(O) || s.checkDraw()
}

func (s *TicTacToeState) GetReward() float64 {
	if s.checkWin(X) {
		return 1
	} else if s.checkWin(O) {
		return -1
	}
	return 0
}

func (s *TicTacToeState) checkWin(player int) bool {
	for i := 0; i < 3; i++ {
		if (s.board[i][0] == player && s.board[i][1] == player && s.board[i][2] == player) || (s.board[0][i] == player && s.board[1][i] == player && s.board[2][i] == player) {
			return true
		}
	}
	return (s.board[0][0] == player && s.board[1][1] == player && s.board[2][2] == player) || (s.board[0][2] == player && s.board[1][1] == player && s.board[2][0] == player)
}

func (s *TicTacToeState) checkDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.board[i][j] == Empty {
				return false
			}
		}
	}
	return true
}

func (s *TicTacToeState) PrintBoard() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch s.board[i][j] {
			case X:
				fmt.Print("X ")
			case O:
				fmt.Print("O ")
			default:
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Main
func main() {
	rand.Seed(time.Now().UnixNano())
	state := &TicTacToeState{player: X}

	for !state.IsTerminal() {
		state.PrintBoard()
		var row, col int
		if state.player == O {
			fmt.Println("Enter your move (row col):")
			_, err := fmt.Scanf("%d %d", &row, &col)
			if err != nil || row < 0 || row > 2 || col < 0 || col > 2 || state.board[row][col] != Empty {
				fmt.Println("Invalid input. Please try again.")
				continue
			}
			state.board[row][col] = O
			state.player = 3 - state.player
		} else {
			node := MCTS(state, 1000)
			var bestChild *Node
			bestVisits := -1
			for _, child := range node.children {
				if child.visits > bestVisits {
					bestVisits = child.visits
					bestChild = child
				}
			}
			state = bestChild.state.(*TicTacToeState)
		}
	}

	state.PrintBoard()
	result := state.GetReward()
	if result == 1 {
		fmt.Println("Player X wins!")
	} else if result == -1 {
		fmt.Println("Player O wins!")
	} else {
		fmt.Println("It's a draw!")
	}
}
