package erst

import "fmt"

const (
	NoUriToDownload = "No URI Found to Download"
	NoGidFound = "No such GID was found"
)


func GenMissingFlagErr(flag string, synopsis string) string {
	return fmt.Sprintf("Important Flag `%s` (%s) missing", flag, synopsis)
}