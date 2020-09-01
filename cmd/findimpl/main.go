package main

import (
	"github.com/nu50218/findimpl"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(findimpl.Analyzer) }
