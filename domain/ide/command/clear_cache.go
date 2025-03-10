/*
 * © 2024 Snyk Limited
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

package command

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/sourcegraph/go-lsp"

	"github.com/snyk/snyk-ls/application/config"
	"github.com/snyk/snyk-ls/internal/types"
)

type clearCache struct {
	command types.CommandData
	c       *config.Config
}

func (cmd *clearCache) Command() types.CommandData {
	return cmd.command
}

// Execute Deletes persisted, inMemory Cache or both.
// Parameters: folderUri either folder Uri or empty for all folders
// cacheType: either inMemory or persisted or empty for both.
func (cmd *clearCache) Execute(_ context.Context) (any, error) {
	logger := cmd.c.Logger().With().Str("method", "clearCache.Execute").Logger()
	args := cmd.command.Arguments
	var parsedFolderUri *lsp.DocumentURI
	folderURI, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid folder URI")
	}

	if folderURI != "" {
		decodedPath, err := url.PathUnescape(folderURI)
		if err != nil {
			logger.Error().Err(err).Msgf("could not decode folder Uri %s", folderURI)
			return nil, err
		}
		uri := lsp.DocumentURI(decodedPath)
		parsedFolderUri = &uri
	}

	cacheType, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid cache URI")
	}

	if cacheType == "" {
		cmd.purgeInMemoryCache(&logger, parsedFolderUri)
		cmd.purgePersistedCache(&logger, parsedFolderUri)
	} else if cacheType == "inMemory" {
		cmd.purgeInMemoryCache(&logger, parsedFolderUri)
	} else if cacheType == "persisted" {
		cmd.purgePersistedCache(&logger, parsedFolderUri)
	}

	return nil, nil
}

func (cmd *clearCache) purgeInMemoryCache(logger *zerolog.Logger, folderUri *lsp.DocumentURI) {
	ws := cmd.c.Workspace()
	trusted, _ := ws.GetFolderTrust()
	for _, folder := range trusted {
		if folderUri != nil && *folderUri != folder.Uri() {
			continue
		}
		logger.Info().Msgf("deleting in-memory cache for folder %s", folder.Path())
		folder.Clear()
		if config.CurrentConfig().IsAutoScanEnabled() {
			go folder.ScanFolder(context.Background())
		}
	}
}

func (cmd *clearCache) purgePersistedCache(logger *zerolog.Logger, folderUri *lsp.DocumentURI) {
	var folderList []types.FilePath
	ws := cmd.c.Workspace()
	clearerExister := ws.GetScanSnapshotClearerExister()
	if clearerExister == nil {
		logger.Error().Msgf("could not find scan persister")
		return
	}
	for _, folder := range ws.Folders() {
		if folderUri != nil && *folderUri != folder.Uri() {
			continue
		}
		folderList = append(folderList, folder.Path())
	}
	logger.Info().Msgf("deleting perrsisted cache for folders %v", folderList)
	clearerExister.Clear(folderList, false)
}
