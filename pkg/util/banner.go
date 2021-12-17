package util

import (
	"html/template"
	"io/ioutil"
	"os"
)

type BannerTemplateData struct {
	Telemetry  bool
	CommitHash string
	Version    string
	BuildDate  string
}

func fallbackEnvString(envString string, fallback string) string {
	if envString == "" {
		return fallback
	} else {
		return envString
	}
}

func PrintBanner(version string, commitHash string, buildDate string) {
	contents, err := ioutil.ReadFile("./assets/banner.txt")
	if err != nil {
		panic(err)
	}

	value := string(contents)
	t := template.New("banner template")
	telemetry := fallbackEnvString(os.Getenv("TSUBAKI_TELEMETRY"), "no")

	var enabled bool
	if telemetry == "no" || telemetry == "false" {
		enabled = false
	} else {
		enabled = true
	}

	data := BannerTemplateData{
		Version:    version,
		CommitHash: commitHash,
		Telemetry:  enabled,
		BuildDate:  buildDate,
	}

	ta, err := t.Parse(value)
	if err != nil {
		panic(err)
	}

	if err := ta.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
