package util

func GetIP(addr string) string {
	i := 0
	end := -1
	for i < len(addr) {
		if addr[i] == ':' {
			end = i
		}
		i++
	}
	if end == -1 {
		return ""
	}
	return addr[:end]
}

func GetPort(addr string) string {
	i := 0
	end := -1
	for i < len(addr) {
		if addr[i] == ':' {
			end = i
		}
		i++
	}
	if end == -1 {
		return ""
	}
	return addr[end+1:]
}
