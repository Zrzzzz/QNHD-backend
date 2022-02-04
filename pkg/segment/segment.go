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
	return strings.Join(seg.CutAll(text), sep)
}
