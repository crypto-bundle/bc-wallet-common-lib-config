package config

import "time"

const (
	ldFlagDefaultVersion     = "v0.0.0"
	ldFlagDefaultReleaseTag  = "v0.0.0"
	ldFlagDefaultCommit      = "0000000000000000000000000000000000000000"
	ldFlagDefaultShortCommit = "00000000"
	ldFlagDefaultBuildNumber = 0
)

type ldFlagManager struct {
	version       string
	releaseTag    string
	commitID      string
	shortCommitID string
	buildNumber   uint64
	builtAt       time.Time
}

func (m *ldFlagManager) GetVersion() string {
	return m.version
}

func (m *ldFlagManager) GetReleaseTag() string {
	return m.releaseTag
}

func (m *ldFlagManager) GetCommitID() string {
	return m.commitID
}

func (m *ldFlagManager) GetShortCommitID() string {
	return m.shortCommitID
}

func (m *ldFlagManager) GetBuildNumber() uint64 {
	return m.buildNumber
}

func (m *ldFlagManager) GetBuildDateTS() int64 {
	return m.builtAt.Unix()
}

func (m *ldFlagManager) GetBuildDate() time.Time {
	return m.builtAt
}

func newDefaultLdFlagManager() *ldFlagManager {
	return &ldFlagManager{
		builtAt:       time.Now(),
		version:       ldFlagDefaultVersion,
		releaseTag:    ldFlagDefaultReleaseTag,
		commitID:      ldFlagDefaultCommit,
		shortCommitID: ldFlagDefaultShortCommit,
		buildNumber:   ldFlagDefaultBuildNumber,
	}
}

func newMockLdFlagManager(appVersion string,
	releaseTag string,
	commitID string,
	shortCommitID string,
	buildNumber uint64,
) *ldFlagManager {
	return &ldFlagManager{
		builtAt:       time.Now(),
		version:       appVersion,
		releaseTag:    releaseTag,
		commitID:      commitID,
		shortCommitID: shortCommitID,
		buildNumber:   buildNumber,
	}
}
