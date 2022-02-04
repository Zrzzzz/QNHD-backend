package segment

import (
	"strings"

	"github.com/go-ego/gse"
)

var (
	seg gse.Segmenter
)

func Setup() {
	load()
}

func load() {
	seg.LoadDict("zh", "jp", "en")
	seg.LoadStop("zh", "jp")
}

func Cut(text string, sep string) string {
	str := seg.CutAll(text)
	res := make([]string, 0, len(str))
	for _, t := range str {
		if t != " " {

			res = append(res, t)
		}
	}
	return strings.Join(res, sep)
}
