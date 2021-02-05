package erst

import "fmt"

const (
	NoUriToDownload = "No URI Found to Download"
)


func GenMissingFlagErr(flag string, synopsis string) string {
	return fmt.Sprintf("Important Flag `%s` (%s) missing", flag, synopsis)
}