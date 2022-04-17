package template

import (
	"fmt"
	"net/url"
	"strings"
)

func GeneArgs(keys, values []string) string {
	var ret []string
	if len(keys) != len(values) {
		return ""
	}
	for i := range keys {
		ret = append(ret, fmt.Sprintf("%s=%s", url.QueryEscape(keys[i]), url.QueryEscape(values[i])))
	}
	return strings.Join(ret, "&")
}

func GeneTemplateString(temp string, args string) (string, error) {
	argsMap, err := splitArgs(args)
	if err != nil {
		return temp, err
	}
	return fillTemplate(temp, argsMap), nil
}

func splitArgs(args string) (map[string]string, error) {
	if args == "" {
		return nil, nil
	}
	maps := make(map[string]string)
	var err error
	for _, a := range strings.Split(args, "&") {
		splits := strings.Split(a, "=")
		k, err := url.QueryUnescape(splits[0])
		if err != nil {
			return nil, err
		}
		v, err := url.QueryUnescape(splits[1])
		if err != nil {
			return nil, err
		}
		maps[k] = v
	}
	return maps, err
}

func fillTemplate(temp string, args map[string]string) string {
	for k, v := range args {
		temp = strings.ReplaceAll(temp, fmt.Sprintf("<%s>", k), v)
	}
	temp = strings.ReplaceAll(temp, "\\<", "<")
	temp = strings.ReplaceAll(temp, "\\>", ">")
	return temp
}
