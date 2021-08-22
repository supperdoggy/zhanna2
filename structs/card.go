package structs

type Card struct {
	Suit  Suit
	Value Value
}

type Suit string
type Value string

const (
	// Suit
	SuitDiamonds Suit = "diamonds"
	SuitHearts   Suit = "hearts"
	SuitSpades   Suit = "spades"
	SuitClubs    Suit = "clubs"

	// Card Value
	ValueSix   Value = "6"
	ValueSeven Value = "7"
	ValueEight Value = "8"
	ValueNine  Value = "9"
	ValueTen   Value = "10"
	ValueJack  Value = "j"
	ValueQueen Value = "q"
	ValueKing  Value = "k"
	ValueAce   Value = "a"
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
