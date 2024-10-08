package config

import (
	"strconv"
	"time"
)

const (
	ldFlagDefaultReleaseTag  = "v0.0.2-4c3452b-100500"
	ldFlagDefaultCommit      = "0000000000000000000000000000000000000000"
	ldFlagDefaultShortCommit = "00000000"
	ldFlagDefaultBuildNumber = 100500
)

var _ ldFlagManagerService = (*ldFlagManager)(nil)

type ldFlagManager struct {
	buildDateAt   time.Time
	releaseTag    string
	commitID      string
	shortCommitID string
	buildNumber   uint64
	buildDateTS   uint64
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
	return m.buildDateAt.Unix()
}

func (m *ldFlagManager) GetBuildDate() time.Time {
	return m.buildDateAt
}

func newDefaultLdFlagManager() *ldFlagManager {
	currTime := time.Now()

	return &ldFlagManager{
		releaseTag:    ldFlagDefaultReleaseTag,
		commitID:      ldFlagDefaultCommit,
		shortCommitID: ldFlagDefaultShortCommit,
		buildNumber:   ldFlagDefaultBuildNumber,
		buildDateTS:   uint64(currTime.Unix()),
		buildDateAt:   time.Now(),
	}
}

//nolint:gochecknoglobals
var ldFlagsSrv *ldFlagManager

func NewLdFlagsManager(
	errFmtSvc errorFormatterService,
	releaseTag,
	commitID,
	shortCommitID,
	buildNumber,
	buildDateTS string,
) (*ldFlagManager, error) {
	if ldFlagsSrv != nil {
		return ldFlagsSrv, nil
	}

	buildDateTSRaw, err := strconv.ParseUint(buildDateTS, 10, 0)
	if err != nil {
		return nil, errFmtSvc.ErrorOnly(err)
	}

	buildNumberRaw, err := strconv.ParseUint(buildNumber, 10, 0)
	if err != nil {
		return nil, errFmtSvc.ErrorOnly(err)
	}

	ldFlagsSrv = &ldFlagManager{
		releaseTag:    releaseTag,
		commitID:      commitID,
		shortCommitID: shortCommitID,
		buildNumber:   buildNumberRaw,
		buildDateTS:   buildDateTSRaw,
		buildDateAt:   time.Unix(int64(buildDateTSRaw), 0),
	}

	return ldFlagsSrv, nil
}
