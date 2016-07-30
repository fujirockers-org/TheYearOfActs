package main

import (
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/fujirockers-org/TheYearOfActs/model"
	"github.com/fujirockers-org/TheYearOfActs/util"
)

const YEAR = "2016"

func main() {
	f, err := os.Open("resource/2016/artist.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}

	list := doc.Find(".artistlist td a").Map(func(i int, s *goquery.Selection) string { return s.Text() })
	list = util.Filter(list, func(s string) bool { return len(s) != 0 })

	s := toStruct(list)
	util.WriteJson("resource/2016/artists.json", s)
}

func toStruct(names []string) *model.FujiRockers {
	rockers := new(model.FujiRockers)
	rockers.Year = YEAR
	rockers.Artists = []model.Artist{}

	for _, v := range names {
		a := new(model.Artist)
		a.Name = v
		a.Reference = append(a.Reference, v)

		rockers.Artists = append(rockers.Artists, *a)
	}
	return rockers
}
