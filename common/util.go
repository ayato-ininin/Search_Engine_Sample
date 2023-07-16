package common

func Contains(arr []string, target string) bool {
	for _, item := range arr {
			if item == target {
					return true
			}
	}
	return false
}
