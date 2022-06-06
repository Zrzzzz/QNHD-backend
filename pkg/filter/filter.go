package filter

import (
	"encoding/json"
	"qnhd/pkg/logging"
	"qnhd/pkg/setting"
	"strings"

	"github.com/Baidu-AIP/golang-sdk/aip/censor"
	"github.com/importcjj/sensitive"
)

var (
	CommonFilter   WordFilter
	NicknameFilter WordFilter
)

func Setup() {
	censor.NewClient(setting.AppSetting.BaiduAIAppKey, setting.AppSetting.BaiduAISecretKey)
	//使用百度云API认证机制
	aiClient := censor.NewCloudClient(setting.AppSetting.BaiduAIAppKey, setting.AppSetting.BaiduAISecretKey)
	CommonFilter = WordFilter{dictPlace: "conf/sensitive.txt", aiFilter: aiClient}
	NicknameFilter = WordFilter{dictPlace: "conf/nickname-sensitive.txt", aiFilter: aiClient}
	CommonFilter.Reload()
	NicknameFilter.Reload()
}

func (c *WordFilter) Reload() error {
	c.filter = sensitive.New()
	return c.filter.LoadWordDict(c.dictPlace)
}

func (c *WordFilter) Filter(s string) string {
	var ret = s
	aiRes := c.aiFilter.TextCensor(ret)
	var r AIFilterResult
	if err := json.Unmarshal([]byte(aiRes), &r); err != nil {
		logging.Error(err.Error())
	}
	if r.ConclusionType == 2 || r.ConclusionType == 3 {
		for _, data := range r.Data {
			for _, hit := range data.Hits {
				for _, word := range hit.Words {
					ret = strings.Replace(ret, word, strings.Repeat("*", len([]rune(word))), -1)
				}
			}
		}
	}

	return c.filter.Replace(s, '*')
}

func (c *WordFilter) Validate(s string) (bool, string) {
	aiRes := c.aiFilter.TextCensor(s)
	var r AIFilterResult
	if err := json.Unmarshal([]byte(aiRes), &r); err != nil {
		logging.Error(err.Error())
	}
	if r.ConclusionType == 2 || r.ConclusionType == 3 {
		if len(r.Data[0].Hits[0].Words) == 0 {
			return false, r.Data[0].Msg
		}
		return false, strings.Join(r.Data[0].Hits[0].Words, "|")
	}
	return c.filter.Validate(s)
}
