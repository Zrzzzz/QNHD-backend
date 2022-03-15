package filter

import "github.com/importcjj/sensitive"

var filter *sensitive.Filter

func Setup() {
	Reload()
}

func Reload() error {
	filter = sensitive.New()
	return filter.LoadWordDict("conf/sensitive.txt")
}

func Filter(s string) string {
	return filter.Replace(s, '*')
}

func Validate(s string) (bool, string) {
	return filter.Validate(s)
}
