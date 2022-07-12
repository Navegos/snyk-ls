package server

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	"github.com/rs/zerolog/log"

	"github.com/snyk/snyk-ls/application/config"
	"github.com/snyk/snyk-ls/application/di"
	"github.com/snyk/snyk-ls/presentation/lsp"
)

func WorkspaceDidChangeConfiguration() jrpc2.Handler {
	return handler.New(func(ctx context.Context, params lsp.DidChangeConfigurationParams) (interface{}, error) {
		log.Info().Str("method", "WorkspaceDidChangeConfiguration").Interface("params", params).Msg("RECEIVED")
		defer log.Info().Str("method", "WorkspaceDidChangeConfiguration").Interface("params", params).Msg("DONE")
		UpdateSettings(ctx, params.Settings)
		return nil, nil
	})
}

func UpdateSettings(ctx context.Context, settings lsp.Settings) {
	emptySettings := lsp.Settings{}
	if settings == emptySettings {
		return
	}
	updateProductEnablement(settings)
	updateCliConfig(ctx, settings)
	updateEnvironment(settings)
	updatePath(settings)
	updateTelemetry(settings)
	updateOrganization(settings)
	manageBinariesAutomatically(settings)
}

func updateOrganization(settings lsp.Settings) {
	org := strings.TrimSpace(settings.Organization)
	if org != "" {
		config.CurrentConfig().SetOrganization(org)
	}
}

func updateTelemetry(settings lsp.Settings) {
	parseBool, err := strconv.ParseBool(settings.SendErrorReports)
	if err != nil {
		log.Err(err).Msgf("couldn't read send error reports %s", settings.SendErrorReports)
	}
	config.CurrentConfig().SetErrorReportingEnabled(parseBool)

	parseBool, err = strconv.ParseBool(settings.EnableTelemetry)
	if err != nil {
		log.Err(err).Msgf("couldn't read send error reports %s", settings.SendErrorReports)
	}
	config.CurrentConfig().SetTelemetryEnabled(parseBool)
}

func manageBinariesAutomatically(settings lsp.Settings) {
	parseBool, err := strconv.ParseBool(settings.ManageBinariesAutomatically)
	if err != nil {
		log.Err(err).Msgf("couldn't read manage binaries automatically %s", settings.ManageBinariesAutomatically)
	} else {
		config.CurrentConfig().SetManageBinariesAutomatically(parseBool)
	}
}

// TODO store in config, move parsing to CLI
func updatePath(settings lsp.Settings) {
	err := os.Setenv("PATH", os.Getenv("PATH")+string(os.PathSeparator)+settings.Path)
	if err != nil {
		log.Err(err).Msgf("couldn't add path %s", settings.Path)
	}
}

// TODO store in config, move parsing to CLI
func updateEnvironment(settings lsp.Settings) {
	envVars := strings.Split(settings.AdditionalEnv, ";")
	for _, envVar := range envVars {
		v := strings.Split(envVar, "=")
		if len(v) != 2 {
			continue
		}
		err := os.Setenv(v[0], v[1])
		if err != nil {
			log.Err(err).Msgf("couldn't set env variable %s", envVar)
		}
	}
}

func updateCliConfig(ctx context.Context, settings lsp.Settings) {
	var err error
	cliSettings := config.CliSettings{}
	cliSettings.Insecure, err = strconv.ParseBool(settings.Insecure)
	if err != nil {
		log.Err(err).Msg("couldn't parse insecure setting")
	}
	cliSettings.Endpoint = strings.Trim(settings.Endpoint, " ")
	cliSettings.AdditionalParameters = strings.Split(settings.AdditionalParams, " ")

	if cliSettings.Endpoint != config.CurrentConfig().CliSettings().Endpoint {
		// Reset CLI token
		err = di.Authenticator().ClearAuthentication(ctx)
		if err != nil {
			log.Err(err).Msg("couldn't reset token")
		}
	}

	config.CurrentConfig().SetCliSettings(cliSettings)
}

func updateProductEnablement(settings lsp.Settings) {
	parseBool, err := strconv.ParseBool(settings.ActivateSnykCode)
	if err != nil {
		log.Err(err).Msg("couldn't parse code setting")
	} else {
		config.CurrentConfig().SetSnykCodeEnabled(parseBool)
	}
	parseBool, err = strconv.ParseBool(settings.ActivateSnykOpenSource)
	if err != nil {
		log.Err(err).Msg("couldn't parse open source setting")
	} else {
		config.CurrentConfig().SetSnykOssEnabled(parseBool)
	}
	parseBool, err = strconv.ParseBool(settings.ActivateSnykIac)
	if err != nil {
		log.Err(err).Msg("couldn't parse iac setting")
	} else {
		config.CurrentConfig().SetSnykIacEnabled(parseBool)
	}
}
