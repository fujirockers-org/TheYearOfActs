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
	a.BestAct = make(map[string]Result)
	a.GoodAct = make(map[string]Result)
	a.WorstAct = make(map[string]Result)

	return a
}

type Aggregate struct {
	BestAct  map[string]Result
	GoodAct  map[string]Result
	WorstAct map[string]Result
}

type Result struct {
	Count int
	Ids   []int
}
