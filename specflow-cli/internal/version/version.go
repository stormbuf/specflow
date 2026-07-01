package version

// 这些变量通过 GoReleaser ldflags 注入：
// -X github.com/stormbuf/specflow/internal/version.Version={{.Version}}
// -X github.com/stormbuf/specflow/internal/version.Commit={{.ShortCommit}}
// -X github.com/stormbuf/specflow/internal/version.Date={{.Date}}

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
