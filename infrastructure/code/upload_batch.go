/*
 * © 2022 Snyk Limited All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package code

import "github.com/snyk/snyk-ls/internal/types"

const (
	maxFileSize               = 1024 * 1024
	maxUploadBatchSize        = 1024*1024*4 - 1024 // subtract 1k for potential headers
	jsonOverheadRequest       = "{\"files\":{}}"
	jsonOverHeadRequestLength = len(jsonOverheadRequest)
	jsonUriOverhead           = "\"\":{}"
	jsonHashSizePerFile       = "\"hash\":\"0123456789012345678901234567890123456789012345678901234567890123\""
	jsonContentOverhead       = ",\"content\":\"\""
	jsonOverheadPerFile       = jsonUriOverhead + jsonContentOverhead
)

type UploadBatch struct {
	hash      string
	documents map[types.FilePath]BundleFile
	size      int
}

func NewUploadBatch() *UploadBatch {
	return &UploadBatch{
		documents: map[types.FilePath]BundleFile{},
	}
}

// todo simplify the size computation
// maybe consider an addFile / canFitFile interface with proper error handling
func (b *UploadBatch) canFitFile(uri types.FilePath, size int) bool {
	docPayloadSize := b.getTotalDocPayloadSize(uri, size)
	newSize := docPayloadSize + b.getSize()
	b.size += docPayloadSize
	return newSize < maxUploadBatchSize
}

func (b *UploadBatch) getTotalDocPayloadSize(filePath types.FilePath, size int) int {
	return len(jsonHashSizePerFile) + len(jsonOverheadPerFile) + len([]byte(filePath)) + size
}

func (b *UploadBatch) getSize() int {
	if len(b.documents) == 0 {
		return 0
	}
	jsonCommasForFiles := len(b.documents) - 1
	var size = jsonOverHeadRequestLength + jsonCommasForFiles // if more than one file, they are separated by commas in the req
	return size + b.size
}

func (b *UploadBatch) hasContent() bool {
	return len(b.documents) > 0
}

type BundleFile struct {
	Hash    string `json:"hash"`
	Content string `json:"content"`
	Size    int    `json:"-"`
}
