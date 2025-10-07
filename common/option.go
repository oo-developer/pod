package common

func HasOption(opt string, args ...string) bool {
	for _, o := range args {
		if o == opt {
			return true
		}
	}
	return false
}
