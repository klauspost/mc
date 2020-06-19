/*
 * MinIO Client (C) 2015 MinIO, Inc.
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

package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/minio/mc/pkg/probe"

	"github.com/mitchellh/go-homedir"
)

// mcCustomConfigDir contains the whole path to config dir. Only access via get/set functions.
var mcCustomConfigDir string

// setMcConfigDir - set a custom MinIO Client config folder.
func setMcConfigDir(configDir string) {
	mcCustomConfigDir = configDir
}

// getMcConfigDir - construct MinIO Client config folder.
func getMcConfigDir() (string, *probe.Error) {
	if mcCustomConfigDir != "" {
		return mcCustomConfigDir, nil
	}
	homeDir, e := homedir.Dir()
	if e != nil {
		return "", probe.NewError(e)
	}
	configDir := filepath.Join(homeDir, defaultMCConfigDir())
	return configDir, nil
}

// Return default default mc config directory.
// Generally you want to use getMcConfigDir which returns custom overrides.
func defaultMCConfigDir() string {
	if runtime.GOOS == "windows" {
		// For windows the path is slightly different
		cmd := filepath.Base(os.Args[0])
		if strings.HasSuffix(strings.ToLower(cmd), ".exe") {
			cmd = cmd[:strings.LastIndex(cmd, ".")]
		}
		return fmt.Sprintf("%s\\", cmd)
	}
	return fmt.Sprintf(".%s/", filepath.Base(os.Args[0]))
}

// mustGetMcConfigDir - construct MinIO Client config folder or fail
func mustGetMcConfigDir() (configDir string) {
	configDir, err := getMcConfigDir()
	fatalIf(err.Trace(), "Unable to get mcConfigDir.")

	return configDir
}

// createMcConfigDir - create MinIO Client config folder
func createMcConfigDir() *probe.Error {
	p, err := getMcConfigDir()
	if err != nil {
		return err.Trace()
	}
	if e := os.MkdirAll(p, 0700); e != nil {
		return probe.NewError(e)
	}
	return nil
}

// getMcConfigPath - construct MinIO Client configuration path
func getMcConfigPath() (string, *probe.Error) {
	if mcCustomConfigDir != "" {
		return filepath.Join(mcCustomConfigDir, globalMCConfigFile), nil
	}
	dir, err := getMcConfigDir()
	if err != nil {
		return "", err.Trace()
	}
	return filepath.Join(dir, globalMCConfigFile), nil
}

// mustGetMcConfigPath - similar to getMcConfigPath, ignores errors
func mustGetMcConfigPath() string {
	path, err := getMcConfigPath()
	fatalIf(err.Trace(), "Unable to get mcConfigPath.")

	return path
}

// newMcConfig - initializes a new version '9' config.
func newMcConfig() *configV9 {
	cfg := newConfigV9()
	cfg.loadDefaults()
	return cfg
}

// loadMcConfigCached - returns loadMcConfig with a closure for config cache.
func loadMcConfigFactory() func() (*configV9, *probe.Error) {
	// Load once and cache in a closure.
	cfgCache, err := loadConfigV9()

	// loadMcConfig - reads configuration file and returns config.
	return func() (*configV9, *probe.Error) {
		return cfgCache, err
	}
}

// loadMcConfig - returns configuration, initialized later.
var loadMcConfig func() (*configV9, *probe.Error)

// saveMcConfig - saves configuration file and returns error if any.
func saveMcConfig(config *configV9) *probe.Error {
	if config == nil {
		return errInvalidArgument().Trace()
	}

	err := createMcConfigDir()
	if err != nil {
		return err.Trace(mustGetMcConfigDir())
	}

	// Save the config.
	if err := saveConfigV9(config); err != nil {
		return err.Trace(mustGetMcConfigPath())
	}

	// Refresh the config cache.
	loadMcConfig = loadMcConfigFactory()
	return nil
}

// isMcConfigExists returns err if config doesn't exist.
func isMcConfigExists() bool {
	configFile, err := getMcConfigPath()
	if err != nil {
		return false
	}
	if _, e := os.Stat(configFile); e != nil {
		return false
	}
	return true
}

// cleanAlias removes any forbidden trailing slashes or backslashes
// before any validation to avoid annoying mc complaints.
func cleanAlias(s string) string {
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSuffix(s, "\\")
	return s
}

// isValidAlias - Check if alias valid.
func isValidAlias(alias string) bool {
	return regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-_]+$").MatchString(alias)
}

// getHostConfig retrieves host specific configuration such as access keys, signature type.
func getHostConfig(alias string) (*hostConfigV9, *probe.Error) {
	mcCfg, err := loadMcConfig()
	if err != nil {
		return nil, err.Trace(alias)
	}

	// if host is exact return quickly.
	if _, ok := mcCfg.Hosts[alias]; ok {
		hostCfg := mcCfg.Hosts[alias]
		return &hostCfg, nil
	}

	// return error if cannot be matched.
	return nil, errNoMatchingHost(alias).Trace(alias)
}

// mustGetHostConfig retrieves host specific configuration such as access keys, signature type.
func mustGetHostConfig(alias string) *hostConfigV9 {
	hostCfg, _ := getHostConfig(alias)
	// If alias is not found,
	// look for it in the environment variable.
	if hostCfg == nil {
		if envConfig, ok := os.LookupEnv(mcEnvHostPrefix + alias); ok {
			hostCfg, _ = expandAliasFromEnv(envConfig)
		}
	}
	if hostCfg == nil {
		if envConfig, ok := os.LookupEnv(mcEnvHostsDeprecatedPrefix + alias); ok {
			errorIf(errInvalidArgument().Trace(mcEnvHostsDeprecatedPrefix+alias), "`MC_HOSTS_<alias>` environment variable is deprecated. Please use `MC_HOST_<alias>` instead for the same functionality.")
			hostCfg, _ = expandAliasFromEnv(envConfig)
		}
	}
	return hostCfg
}

var (
	hostKeys      = regexp.MustCompile("^(https?://)(.*?):(.*?)@(.*?)$")
	hostKeyTokens = regexp.MustCompile("^(https?://)(.*?):(.*?):(.*?)@(.*?)$")
)

// parse url usually obtained from env.
func parseEnvURLStr(envURL string) (*url.URL, string, string, string, *probe.Error) {
	var accessKey, secretKey, sessionToken string
	var parsedURL string
	if hostKeyTokens.MatchString(envURL) {
		parts := hostKeyTokens.FindStringSubmatch(envURL)
		if len(parts) != 6 {
			return nil, "", "", "", errInvalidArgument().Trace(envURL)
		}
		accessKey = parts[2]
		secretKey = parts[3]
		sessionToken = parts[4]
		parsedURL = fmt.Sprintf("%s%s", parts[1], parts[5])
	} else if hostKeys.MatchString(envURL) {
		parts := hostKeys.FindStringSubmatch(envURL)
		if len(parts) != 5 {
			return nil, "", "", "", errInvalidArgument().Trace(envURL)
		}
		accessKey = parts[2]
		secretKey = parts[3]
		parsedURL = fmt.Sprintf("%s%s", parts[1], parts[4])
	}
	var u *url.URL
	var e error
	if parsedURL != "" {
		u, e = url.Parse(parsedURL)
	} else {
		u, e = url.Parse(envURL)
	}
	if e != nil {
		return nil, "", "", "", probe.NewError(e)
	}
	// Look for if URL has invalid values and return error.
	if !((u.Scheme == "http" || u.Scheme == "https") &&
		(u.Path == "/" || u.Path == "") && u.Opaque == "" &&
		!u.ForceQuery && u.RawQuery == "" && u.Fragment == "") {
		return nil, "", "", "", errInvalidArgument().Trace(u.String())
	}
	if accessKey == "" && secretKey == "" {
		if u.User != nil {
			accessKey = u.User.Username()
			secretKey, _ = u.User.Password()
		}
	}
	return u, accessKey, secretKey, sessionToken, nil
}

const (
	mcEnvHostPrefix            = "MC_HOST_"
	mcEnvHostsDeprecatedPrefix = "MC_HOSTS_"
)

func expandAliasFromEnv(envURL string) (*hostConfigV9, *probe.Error) {
	u, accessKey, secretKey, sessionToken, err := parseEnvURLStr(envURL)
	if err != nil {
		return nil, err.Trace(envURL)
	}

	return &hostConfigV9{
		URL:          u.String(),
		API:          "S3v4",
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
	}, nil
}

// expandAlias expands aliased URL if any match is found, returns as is otherwise.
func expandAlias(aliasedURL string) (alias string, urlStr string, hostCfg *hostConfigV9, err *probe.Error) {
	// Extract alias from the URL.
	alias, path := url2Alias(aliasedURL)

	var envConfig string
	var ok bool

	if envConfig, ok = os.LookupEnv(mcEnvHostPrefix + alias); !ok {
		envConfig, ok = os.LookupEnv(mcEnvHostsDeprecatedPrefix + alias)
		if ok {
			errorIf(errInvalidArgument().Trace(mcEnvHostsDeprecatedPrefix+alias), "`MC_HOSTS_<alias>` environment variable is deprecated. Please use `MC_HOST_<alias>` instead for the same functionality.")
		}
	}

	if ok {
		hostCfg, err = expandAliasFromEnv(envConfig)
		if err != nil {
			return "", "", nil, err.Trace(aliasedURL)
		}
		return alias, urlJoinPath(hostCfg.URL, path), hostCfg, nil
	}

	// Find the matching alias entry and expand the URL.
	if hostCfg = mustGetHostConfig(alias); hostCfg != nil {
		return alias, urlJoinPath(hostCfg.URL, path), hostCfg, nil
	}

	snapshots, err := listSnapshots()
	if err == nil {
		for _, s := range snapshots {
			if s.Name() == alias {
				return alias, aliasedURL, nil, nil
			}
		}
	}

	return "", aliasedURL, nil, nil // No matching entry found. Return original URL as is.
}

// mustExpandAlias expands aliased URL if any match is found, returns as is otherwise.
func mustExpandAlias(aliasedURL string) (alias string, urlStr string, hostCfg *hostConfigV9) {
	alias, urlStr, hostCfg, _ = expandAlias(aliasedURL)
	return alias, urlStr, hostCfg
}
