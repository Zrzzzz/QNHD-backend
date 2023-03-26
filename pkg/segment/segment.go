package segment

import (
	"qnhd/pkg/setting"
	"regexp"
	"strings"

	"github.com/go-ego/gse"
	"mvdan.cc/xurls/v2"
)

var (
	seg       gse.Segmenter
	urlFilter *regexp.Regexp
)

func Setup() {
	load()
}

func load() {
	seg.LoadDict()
	if setting.EnvironmentSetting.RELEASE == "1" {
		seg.LoadDict("dict/zh/s_1.txt, dict/zh/t_1.txt, dict/jp/dict.txt")
		seg.LoadStop("dict/zh/stop_word.txt, dict/zh/stop_tokens.txt")
	} else {
		seg.LoadDict("zh", "jp", "en")
		seg.LoadStop("zh", "jp")
	}

	urlFilter = xurls.Relaxed()
}

func Cut(text string, sep string) string {
	// 过滤网址
	urls := urlFilter.FindAllString(text, -1)
	for _, x := range urls {
		text = strings.ReplaceAll(text, x, "")
	}
	return strings.Join(seg.CutSearch(text), sep)
}
