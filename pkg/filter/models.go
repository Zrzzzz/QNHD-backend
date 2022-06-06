package filter

import (
	"github.com/Baidu-AIP/golang-sdk/aip/censor"
	"github.com/importcjj/sensitive"
)

type WordFilter struct {
	dictPlace string
	filter    *sensitive.Filter
	aiFilter  *censor.ContentCensorClient
}

type AIFilterResult struct {
	LogID          int64          `json:"log_id"`
	Conclusion     string         `json:"conclusion"`
	ConclusionType int64          `json:"conclusionType"`
	Data           []AIFilterData `json:"data"`
}

type AIFilterData struct {
	Type           int64         `json:"type"`
	SubType        int64         `json:"subType"`
	Conclusion     string        `json:"conclusion"`
	ConclusionType int64         `json:"conclusionType"`
	Msg            string        `json:"msg"`
	Hits           []AIFilterHit `json:"hits"`
}

type AIFilterHit struct {
	DatasetName string   `json:"datasetName"`
	Words       []string `json:"words"`
	Probability *float64 `json:"probability,omitempty"`
}
