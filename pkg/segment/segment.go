package segment

import (
	"qnhd/pkg/setting"
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
	if setting.EnvironmentSetting.RELEASE == "1" {
		seg.LoadDict("dict/zh/s_1.txt, dict/zh/t_1.txt, dict/jp/dict.txt")
		seg.LoadStop("dict/zh/stop_word.txt, dict/zh/stop_tokens.txt")
	} else {
		seg.LoadDict("zh", "jp", "en")
		seg.LoadStop("zh", "jp")
	}
}

func Cut(text string, sep string) string {
	return strings.Join(seg.CutAll(text), sep)
}
