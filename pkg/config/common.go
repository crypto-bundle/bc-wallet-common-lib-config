package config

type ldFlagManagerService interface {
	GetVersion() string
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
}
