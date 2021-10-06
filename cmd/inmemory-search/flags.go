package main

import (
	"flag"
	"os"
)

type flags struct {
	http    string
	logFile string
}

func getFlags() *flags {
	flg := &flags{
		http: ":8080",
	}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&flg.http, "http", flg.http, "HTTP ")
	fs.StringVar(&flg.logFile, "log-file", flg.logFile, "log file")
	_ = fs.Parse(os.Args[1:]) // Ignore error, because it exits on error
	return flg
}
