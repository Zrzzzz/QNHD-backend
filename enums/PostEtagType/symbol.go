package PostEtagType

var msgSymbol = map[Enum]string{
	NONE:      "",
	RECOMMEND: "recommend",
	THEME:     "theme",
	TOP:       "top",
}

func (code Enum) GetSymbol() string {
	return msgSymbol[code]
}

func Contains(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range msgSymbol {
		if v == s {
			return true
		}
	}
	return false
}
