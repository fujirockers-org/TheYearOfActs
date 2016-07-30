package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fujirockers-org/TheYearOfActs/model"
	"github.com/fujirockers-org/TheYearOfActs/util"
)

const (
	ATTENDANT   = "【参加日】"
	BEST_ACT    = "【ベストアクト】"
	GOOD_ACT    = "【良かったアクト】"
	WORST_ACT   = "【ワーストアクト】"
	BEST_FOOD   = "【ベスト飯】"
	WORST_FOOD  = "【ワースト飯】"
	ETC_COMMENT = "【その他一言】"
)

func main() {
	files, err := filepath.Glob("resource/2016/thread/*.html")
	if err != nil {
		panic(err)
	}

	votes := []Vote{}
	for _, v := range files {
		votes = append(votes, parseThread(v)...)
	}

	artists := decodeArtists()
	agg := aggregate(votes, artists)

	util.WriteJson("resource/2016/award/acts.json", agg)
}

func parseThread(path string) (votes []Vote) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}

	doc.Find(".post").Each(func(i int, s *goquery.Selection) {
		t := s.Find(".message").Text()
		if isVotingPost(t) {
			votes = append(votes, parseVote(t))
		}
	})
	return votes
}

func isVotingPost(t string) bool {
	return strings.Contains(t, BEST_ACT)
}

func parseVote(t string) Vote {
	v := new(Vote)
	v.Attend, t = parse(t, BEST_ACT)
	v.BestActs, t = parse(t, GOOD_ACT)
	v.GoodActs, t = parse(t, WORST_ACT)
	v.WorstActs, t = parse(t, BEST_FOOD)
	v.BestFoods, t = parse(t, WORST_FOOD)
	v.WorstFoods, v.Comment = parse(t, ETC_COMMENT)

	return *v
}

func parse(t, h string) (string, string) {
	s := strings.Split(t, h)
	return s[0], strings.Join(s[1:], "")
}

func decodeArtists() []model.Artist {
	f, err := ioutil.ReadFile("resource/2016/artists.json")
	if err != nil {
		panic(err)
	}

	var org model.FujiRockers
	json.Unmarshal([]byte(f), &org)
	return org.Artists
}

func aggregate(votes []Vote, artists []model.Artist) model.Aggregate {
	agg := model.NewAggregate()
	for _, v := range votes {
		agg.BestAct = count(agg.BestAct, v.BestActs, artists)
		agg.GoodAct = count(agg.GoodAct, v.GoodActs, artists)
		agg.WorstAct = count(agg.WorstAct, v.WorstActs, artists)
	}
	return *agg
}

func count(m map[string]int, t string, artists []model.Artist) map[string]int {
	list := []string{}
	for _, a := range artists {
		if util.Find(a.Reference, t) {
			list = append(list, a.Name)
		}
	}
	return collect(m, list)
}

func collect(m map[string]int, l []string) map[string]int {
	for _, v := range l {
		if _, ok := m[v]; ok {
			m[v]++
		} else {
			m[v] = 1
		}
	}
	return m
}

type Vote struct {
	Attend     string
	BestActs   string
	GoodActs   string
	WorstActs  string
	BestFoods  string
	WorstFoods string
	Comment    string
}
