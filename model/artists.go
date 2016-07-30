package model

type FujiRockers struct {
	Year    string
	Artists []Artist
}

type Artist struct {
	Name      string
	Reference []string
}

func NewAggregate() *Aggregate {
	a := new(Aggregate)
	a.BestAct = make(map[string]int)
	a.GoodAct = make(map[string]int)
	a.WorstAct = make(map[string]int)

	return a
}

type Aggregate struct {
	BestAct  map[string]int
	GoodAct  map[string]int
	WorstAct map[string]int
}
