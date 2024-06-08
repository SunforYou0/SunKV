package utils

// ToCmdLine conver str -> [][]byte
func ToCmdLine(cmd ...string) [][]byte {
	args := make([][]byte, len(cmd))
	for i, s := range cmd {
		args[i] = []byte(s)
	}
	return args
}
func ToCmdLine2(cmdName string, args ...[]byte) [][]byte {
	res := make([][]byte, len(args)+1)
	res[0] = []byte(cmdName)
	for i, s := range args {
		res[i+1] = s
	}
	return res
}
