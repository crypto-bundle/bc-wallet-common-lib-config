package envconfig

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
	buildDateTS   uint64
	builtDateAt   time.Time
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
	return m.builtDateAt.Unix()
}

func (m *ldFlagManager) GetBuildDate() time.Time {
	return m.builtDateAt
}

func newDefaultLdFlagManager() *ldFlagManager {
	return &ldFlagManager{
		builtDateAt:   time.Now(),
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
	buildTime := time.Now()
	return &ldFlagManager{
		builtDateAt:   buildTime,
		buildDateTS:   uint64(buildTime.Unix()),
		version:       appVersion,
		releaseTag:    releaseTag,
		commitID:      commitID,
		shortCommitID: shortCommitID,
		buildNumber:   buildNumber,
	}
}

var ldFlagsSrv *ldFlagManager

func NewLdFlagsManager(
	version,
	releaseTag,
	commitID,
	shortCommitID string,
	buildNumber,
	buildDateTS uint64,
) *ldFlagManager {
	if ldFlagsSrv != nil {
		return ldFlagsSrv
	}

	ldFlagsSrv = &ldFlagManager{
		version:       version,
		releaseTag:    releaseTag,
		commitID:      commitID,
		shortCommitID: shortCommitID,
		buildNumber:   buildNumber,
		buildDateTS:   buildDateTS,
		builtDateAt:   time.Unix(int64(buildDateTS), 0),
	}

	return ldFlagsSrv
}
