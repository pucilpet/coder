package codersdk_test

import (
	"testing"

	"github.com/coder/coder/codersdk"
)

type exclusion struct {
	flag bool
	env  bool
	yaml bool
}

func TestDeploymentValues_HighlyConfigurable(t *testing.T) {
	t.Parallel()

	// This test ensures that every deployment option has
	// a corresponding Flag, Env, and YAML name, unless explicitly excluded.

	excludes := map[string]exclusion{
		// These are used to configure YAML support itself, so
		// they make no sense within the YAML file.
		"Config Path": {
			yaml: true,
		},
		"Write Config": {
			yaml: true,
		},
		// Dangerous values? Not sure we should help users
		// persistent their configuration.
		"DANGEROUS: Allow Path App Sharing": {
			yaml: true,
		},
		"DANGEROUS: Allow Site Owners to Access Path Apps": {
			yaml: true,
		},
		// Secrets
		"Trace Honeycomb API Key": {
			yaml: true,
		},
		"OAuth2 GitHub Client Secret": {
			yaml: true,
		},
		"OIDC Client Secret": {
			yaml: true,
		},
		"Postgres Connection URL": {
			yaml: true,
		},
		"SCIM API Key": {
			yaml: true,
		},
		// These complex objects should be configured through YAML.
		"Support Links": {
			flag: true,
			env:  true,
		},
		"Git Auth Providers": {
			// Technically Git Auth Providers can be provided through the env,
			// but bypassing clibase. See cli.ReadGitAuthProvidersFromEnv.
			flag: true,
			env:  true,
		},
	}

	set := (&codersdk.DeploymentValues{}).Options()
	for _, opt := range set {
		// These are generally for development, so their configurability is
		// not relevant.
		if opt.Hidden {
			delete(excludes, opt.Name)
			continue
		}

		if codersdk.IsSecretDeploymentOption(opt) && opt.YAML != "" {
			// Secrets should not be written to YAML and instead should continue
			// to be read from the environment.
			//
			// Unfortunately, secrets are still accepted through flags for
			// legacy purposes. Eventually, we should prevent that.
			t.Errorf("Option %q is a secret but has a YAML name", opt.Name)
		}

		excluded := excludes[opt.Name]
		switch {
		case opt.YAML == "" && !excluded.yaml:
			t.Errorf("Option %q should have a YAML name", opt.Name)
		case opt.YAML != "" && excluded.yaml:
			t.Errorf("Option %q is excluded but has a YAML name", opt.Name)
		case opt.Flag == "" && !excluded.flag:
			t.Errorf("Option %q should have a flag name", opt.Name)
		case opt.Flag != "" && excluded.flag:
			t.Errorf("Option %q is excluded but has a flag name", opt.Name)
		case opt.Env == "" && !excluded.env:
			t.Errorf("Option %q should have an env name", opt.Name)
		case opt.Env != "" && excluded.env:
			t.Errorf("Option %q is excluded but has an env name", opt.Name)
		}

		delete(excludes, opt.Name)
	}

	for opt := range excludes {
		t.Errorf("Excluded option %q is not in the deployment config. Remove it?", opt)
	}
}
