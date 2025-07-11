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
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetEnabledProducts_DefaultValues(t *testing.T) {
	t.Setenv(ActivateSnykOssKey, "set it to anything to make sure it is reset")
	t.Setenv(ActivateSnykCodeKey, "set it to anything to make sure it is reset")
	t.Setenv(ActivateSnykIacKey, "set it to anything to make sure it is reset")
	t.Setenv(ActivateSnykAdvisorKey, "set it to anything to make sure it is reset")
	_ = os.Unsetenv(ActivateSnykOssKey)
	_ = os.Unsetenv(ActivateSnykCodeKey)
	_ = os.Unsetenv(ActivateSnykIacKey)
	_ = os.Unsetenv(ActivateSnykAdvisorKey)
	SetCurrentConfig(New())

	currentConfig.clientSettingsFromEnv()

	assert.Equal(t, true, CurrentConfig().IsSnykOssEnabled())
	assert.Equal(t, false, CurrentConfig().IsSnykCodeEnabled())
	assert.Equal(t, true, CurrentConfig().IsSnykIacEnabled())
	assert.Equal(t, false, CurrentConfig().IsSnykAdvisorEnabled())
}

func TestConfig_IsErrorReportingEnabledFromEnv_DefaultValues(t *testing.T) {
	t.Setenv(SendErrorReportsKey, "set it to anything to make sure it is reset")
	_ = os.Unsetenv(SendErrorReportsKey)
	SetCurrentConfig(New())
	CurrentConfig().clientSettingsFromEnv()

	assert.Equal(t, true, CurrentConfig().IsErrorReportingEnabled())
}
func TestConfig_IsErrorReportingEnabledFromEnv(t *testing.T) {
	t.Setenv(SendErrorReportsKey, "true")
	SetCurrentConfig(New())
	CurrentConfig().clientSettingsFromEnv()

	assert.Equal(t, true, CurrentConfig().IsErrorReportingEnabled())
}

func TestConfig_IsErrorReportingEnabledFromEnv_Error(t *testing.T) {
	t.Setenv(SendErrorReportsKey, "hurz")
	SetCurrentConfig(New())
	CurrentConfig().clientSettingsFromEnv()

	assert.Equal(t, true, CurrentConfig().IsErrorReportingEnabled())
}

func TestConfig_OrganizationFromEnv(t *testing.T) {
	orgUuid, _ := uuid.NewRandom()
	expectedOrgId := orgUuid.String()
	t.Setenv(Organization, expectedOrgId)
	SetCurrentConfig(New())
	CurrentConfig().clientSettingsFromEnv()

	assert.Equal(t, expectedOrgId, CurrentConfig().Organization())
}

func TestInitializeDefaultProductEnablement(t *testing.T) {
	t.Setenv(ActivateSnykOssKey, "false")
	t.Setenv(ActivateSnykCodeKey, "true")
	t.Setenv(ActivateSnykIacKey, "false")
	t.Setenv(ActivateSnykAdvisorKey, "true")

	SetCurrentConfig(New())

	assert.Equal(t, false, CurrentConfig().IsSnykOssEnabled())
	assert.Equal(t, true, CurrentConfig().IsSnykCodeEnabled())
	assert.Equal(t, false, CurrentConfig().IsSnykIacEnabled())
	assert.Equal(t, true, CurrentConfig().IsSnykAdvisorEnabled())
}

func TestGetEnabledProducts_Oss(t *testing.T) {
	t.Setenv(ActivateSnykOssKey, "false")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, false, CurrentConfig().isSnykOssEnabled)

	t.Setenv(ActivateSnykOssKey, "true")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, true, CurrentConfig().isSnykOssEnabled)
}

func TestGetEnabledProducts_Code(t *testing.T) {
	t.Setenv(ActivateSnykCodeKey, "false")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, false, CurrentConfig().IsSnykCodeEnabled())

	t.Setenv(ActivateSnykCodeKey, "true")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, true, CurrentConfig().IsSnykCodeEnabled())
}

func TestGetEnabledProducts_Iac(t *testing.T) {
	t.Setenv(ActivateSnykIacKey, "false")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, false, CurrentConfig().IsSnykIacEnabled())

	t.Setenv(ActivateSnykIacKey, "true")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, true, CurrentConfig().IsSnykIacEnabled())
}

func TestGetEnabledProducts_Advisor(t *testing.T) {
	t.Setenv(ActivateSnykAdvisorKey, "false")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, false, CurrentConfig().IsSnykAdvisorEnabled())

	t.Setenv(ActivateSnykAdvisorKey, "true")
	CurrentConfig().clientSettingsFromEnv()
	assert.Equal(t, true, CurrentConfig().IsSnykAdvisorEnabled())
}
