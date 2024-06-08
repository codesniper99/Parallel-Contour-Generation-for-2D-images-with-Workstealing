package messagepackage

type Config struct {
	Mode        string // Represents which scheduler scheme to use
	ThreadCount int    // Runs parallel version with the specified number of threads
	Threshold   int    // Threshold for making contours
}
type Message struct {
	InPath  string
	OutPath string
	Effects []string
}
