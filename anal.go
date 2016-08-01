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

const ATTENDANT = "【参加日】"

var counter util.Counter

func main() {
	counter.Initialize()

	files, err := filepath.Glob("resource/2016/thread/*.html")
	if err != nil {
		panic(err)
	}

	votes := []Vote{}
	for _, v := range files {
		votes = append(votes, parseThread(v)...)
	}
	util.WriteJson("resource/2016/award/vote.json", votes)

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
	return strings.Contains(t, ATTENDANT)
}

func parseVote(t string) Vote {

	BEST_ACT := []string{"【ベストアクト】"}
	GOOD_ACT := []string{"【良かったアクト】", "【目当てじゃなかったが良かったアクト】"}
	WORST_ACT := []string{"【ワーストアクト】"}
	BEST_FOOD := []string{"【ベスト飯】"}
	WORST_FOOD := []string{"【ワースト飯】"}
	ETC_COMMENT := []string{"【その他一言】", "【一言】"}

	v := new(Vote)
	v.Id = counter.Get()

	if strings.Contains(t, ATTENDANT) {
		t = strings.Split(t, ATTENDANT)[1]
	}
	v.Attend, t = parse(t, BEST_ACT)
	v.BestActs, t = parse(t, GOOD_ACT)
	v.GoodActs, t = parse(t, WORST_ACT)
	v.WorstActs, t = parse(t, BEST_FOOD)
	v.BestFoods, t = parse(t, WORST_FOOD)
	v.WorstFoods, t = parse(t, ETC_COMMENT)
	v.Comment = strings.Trim(t, " 　\r\n")

	return *v
}

func parse(t string, h []string) (string, string) {
	for _, v := range h {
		s := strings.Split(t, v)
		if len(s) > 1 {
			return strings.Trim(s[0], " 　\r\n"), strings.Join(s[1:], "")
		}
	}
	return t, ""
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
		agg.BestAct = count(agg.BestAct, v.Id, v.BestActs, artists)
		agg.GoodAct = count(agg.GoodAct, v.Id, v.GoodActs, artists)
		agg.WorstAct = count(agg.WorstAct, v.Id, v.WorstActs, artists)
	}
	return *agg
}

func count(m map[string]model.Result, id int, t string, artists []model.Artist) map[string]model.Result {
	list := []string{}
	for _, a := range artists {
		if util.Find(a.Reference, t) {
			list = append(list, a.Name)
		}
	}
	return collect(m, list, id)
}

func collect(m map[string]model.Result, l []string, id int) map[string]model.Result {
	for _, v := range l {
		if _, ok := m[v]; ok {
			tmp := m[v]
			tmp.Count++
			tmp.Ids = append(m[v].Ids, id)
			m[v] = tmp
		} else {
			tmp := new(model.Result)
			tmp.Count = 1
			tmp.Ids = []int{id}
			m[v] = *tmp
		}
	}
	return m
}

type Vote struct {
	Id         int
	Attend     string `json:"attend"`
	BestActs   string `json:"best_acts"`
	GoodActs   string `json:"good_acts"`
	WorstActs  string `json:"worst_acts"`
	BestFoods  string `json:"best_food"`
	WorstFoods string `json:"worst_food"`
	Comment    string `json:"comment"`
}
