package engine

type Move int

type Board [][]Move

type Config struct {
	Board Board
    Player1 Move
    Player2 Move
}

const (
	EMPTY Move = iota
	X 
	O
)

func (m Move) String() string {
    switch m {
    case X:
        return "X"
    case O:
        return "O"
    default:
        return "-" 
    }
}