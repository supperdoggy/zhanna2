package structs

type Card struct {
	Suit  string
	Value string
}

const (
	// Suit
	SuitDiamonds string = "diamonds"
	SuitHearts   string = "hearts"
	SuitSpades   string = "spades"
	SuitClubs    string = "clubs"

	// Card Value
	ValueSix   string = "6"
	ValueSeven string = "7"
	ValueEight string = "8"
	ValueNine  string = "9"
	ValueTen   string = "10"
	ValueJack  string = "j"
	ValueQueen string = "q"
	ValueKing  string = "k"
	ValueAce   string = "a"
)

var (
	// DO NOT CHANGE COZ I CANT MAKE IT CONSTANT
	Cards = [...]Card{
		{Suit: SuitClubs, Value: ValueSix},
		{Suit: SuitClubs, Value: ValueSeven},
		{Suit: SuitClubs, Value: ValueEight},
		{Suit: SuitClubs, Value: ValueNine},
		{Suit: SuitClubs, Value: ValueTen},
		{Suit: SuitClubs, Value: ValueJack},
		{Suit: SuitClubs, Value: ValueQueen},
		{Suit: SuitClubs, Value: ValueKing},
		{Suit: SuitClubs, Value: ValueAce},

		{Suit: SuitDiamonds, Value: ValueSix},
		{Suit: SuitDiamonds, Value: ValueSeven},
		{Suit: SuitDiamonds, Value: ValueEight},
		{Suit: SuitDiamonds, Value: ValueNine},
		{Suit: SuitDiamonds, Value: ValueTen},
		{Suit: SuitDiamonds, Value: ValueJack},
		{Suit: SuitDiamonds, Value: ValueQueen},
		{Suit: SuitDiamonds, Value: ValueKing},
		{Suit: SuitDiamonds, Value: ValueAce},

		{Suit: SuitSpades, Value: ValueSix},
		{Suit: SuitSpades, Value: ValueSeven},
		{Suit: SuitSpades, Value: ValueEight},
		{Suit: SuitSpades, Value: ValueNine},
		{Suit: SuitSpades, Value: ValueTen},
		{Suit: SuitSpades, Value: ValueJack},
		{Suit: SuitSpades, Value: ValueQueen},
		{Suit: SuitSpades, Value: ValueKing},
		{Suit: SuitSpades, Value: ValueAce},

		{Suit: SuitHearts, Value: ValueSix},
		{Suit: SuitHearts, Value: ValueSeven},
		{Suit: SuitHearts, Value: ValueEight},
		{Suit: SuitHearts, Value: ValueNine},
		{Suit: SuitHearts, Value: ValueTen},
		{Suit: SuitHearts, Value: ValueJack},
		{Suit: SuitHearts, Value: ValueQueen},
		{Suit: SuitHearts, Value: ValueKing},
		{Suit: SuitHearts, Value: ValueAce},
	}
)
