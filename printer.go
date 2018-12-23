package main

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
)

func CPrint(tag string, tagColor aurora.Color, text string) {
	coloredTag := aurora.Bold(aurora.Colorize("["+tag+"]:", tagColor))
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:03"), coloredTag, text)
}
