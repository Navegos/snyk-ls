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

package config

import (
	"os"
	"strconv"
)

const (
	ActivateSnykOssKey     = "ACTIVATE_SNYK_OPEN_SOURCE"
	ActivateSnykCodeKey    = "ACTIVATE_SNYK_CODE"
	ActivateSnykIacKey     = "ACTIVATE_SNYK_IAC"
	ActivateSnykAdvisorKey = "ACTIVATE_SNYK_ADVISOR"
	SendErrorReportsKey    = "SEND_ERROR_REPORTS"
	Organization           = "SNYK_CFG_ORG"
)

func (c *Config) clientSettingsFromEnv() {
	c.productEnablementFromEnv()
	c.errorReportsEnablementFromEnv()
	c.orgFromEnv()
	c.path = os.Getenv("PATH")
}

func (c *Config) orgFromEnv() {
	org := os.Getenv(Organization)
	if org != "" {
		c.SetOrganization(org)
	}
}

func (c *Config) errorReportsEnablementFromEnv() {
	errorReports := os.Getenv(SendErrorReportsKey)
	if errorReports == "false" {
		c.SetErrorReportingEnabled(false)
	} else {
		c.SetErrorReportingEnabled(true)
	}
}

func (c *Config) productEnablementFromEnv() {
	oss := os.Getenv(ActivateSnykOssKey)
	code := os.Getenv(ActivateSnykCodeKey)
	iac := os.Getenv(ActivateSnykIacKey)
	advisor := os.Getenv(ActivateSnykAdvisorKey)

	if oss != "" {
		parseBool, err := strconv.ParseBool(oss)
		if err != nil {
			c.Logger().Debug().Err(err).Str("method", "clientSettingsFromEnv").Msgf("couldn't parse oss config %s", oss)
		}
		c.SetSnykOssEnabled(parseBool)
	}

	if code != "" {
		parseBool, err := strconv.ParseBool(code)
		if err != nil {
			c.Logger().Debug().Err(err).Str("method", "clientSettingsFromEnv").Msgf("couldn't parse code config %s", code)
		}
		c.SetSnykCodeEnabled(parseBool)
	}

	if iac != "" {
		parseBool, err := strconv.ParseBool(iac)
		if err != nil {
			c.Logger().Debug().Err(err).Str("method", "clientSettingsFromEnv").Msgf("couldn't parse iac config %s", iac)
		}
		c.SetSnykIacEnabled(parseBool)
	}

	if advisor != "" {
		parseBool, err := strconv.ParseBool(advisor)
		if err != nil {
			c.Logger().Debug().Err(err).Str("method", "clientSettingsFromEnv").Msgf("couldn't parse advisor config %s", advisor)
		}
		c.SetSnykAdvisorEnabled(parseBool)
	}
}
