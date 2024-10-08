/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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
