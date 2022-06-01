package filter

import "github.com/importcjj/sensitive"

type WordFilter struct {
	dictPlace string
	filter    *sensitive.Filter
}

var (
	CommonFilter   WordFilter
	NicknameFilter WordFilter
)

func Setup() {
	CommonFilter = WordFilter{dictPlace: "conf/sensitive.txt"}
	NicknameFilter = WordFilter{dictPlace: "conf/nickname-sensitive.txt"}
	CommonFilter.Reload()
	NicknameFilter.Reload()
}

func (c *WordFilter) Reload() error {
	c.filter = sensitive.New()
	return c.filter.LoadWordDict(c.dictPlace)
}

func (c *WordFilter) Filter(s string) string {
	return c.filter.Replace(s, '*')
}

func (c *WordFilter) Validate(s string) (bool, string) {
	return c.filter.Validate(s)
}
