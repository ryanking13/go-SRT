package srt

func reverseMap(m map[string]string) map[string]string {
	_m := map[string]string{}
	for k, v := range m {
		_m[v] = k
	}

	return _m
}
