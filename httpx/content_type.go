package httpx

import (
	"net/http"
	"strings"
)

func ContentType(req *http.Request) string {
	ct := req.Header.Get("Content-Type")
	if ct == "" {
		ct = req.Header.Get("Accept")
	}

	ct = strings.TrimSpace(ct)
	var cts []string
	if strings.Contains(ct, ",") {
		cts = strings.Split(ct, ",")
	} else {
		cts = strings.Split(ct, ";")
	}
	for _, c := range cts {
		c = strings.TrimSpace(c)
		if strings.HasPrefix(c, "*/*") {
			continue
		}
		return strings.ToLower(c)
	}
	if ct == "*/*" {
		return ""
	}
	return ct
}
