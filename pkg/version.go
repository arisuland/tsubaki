package pkg

var (
	Version    string
	CommitHash string
	BuildDate  string
)

func SetVersion(version string, commitHash string, buildDate string) {
	Version = version
	CommitHash = commitHash
	BuildDate = buildDate
}
