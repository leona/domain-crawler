package crawler

import (
	"fmt"
	_ "errors"
)

func isTopLevel(component string) bool {
    switch component {
    case
		"com",
		"io",
		"net",
		"org",
        "cn":
        return true
    }
    return false
}

func isInvalidExtension(extension string) bool {
    switch extension {
    case
        ".jpg",
        ".jpeg",
		".gif",
		".mp4",
		".png",
		".mp3",
		".pdf",
		".css",
		".js",
		".webp",
        ".svg":
        return true
    }
    return false
}

func Info(level int, args ...interface{}) {
	if *InputOptions.Verbose >= level {
		fmt.Println(args)
	}
}

func reverse(list []string) []string {
	reversed := []string{}
	
	for i := len(list) - 1;i >= 0;i-- {
		reversed = append(reversed, list[i])
	}

	return reversed
}