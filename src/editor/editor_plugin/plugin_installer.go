/******************************************************************************/
/* plugin_installer.go                                                        */
/******************************************************************************/
/*                            This file is part of                            */
/*                                KAIJU ENGINE                                */
/*                          https://kaijuengine.com/                          */
/******************************************************************************/
/* MIT License                                                                */
/*                                                                            */
/* Copyright (c) 2023-present Kaiju Engine authors (AUTHORS.md).              */
/* Copyright (c) 2015-present Brent Farris.                                   */
/*                                                                            */
/* May all those that this source may reach be blessed by the LORD and find   */
/* peace and joy in life.                                                     */
/* Everyone who drinks of this water will be thirsty again; but whoever       */
/* drinks of the water that I will give him shall never thirst; John 4:13-14  */
/*                                                                            */
/* Permission is hereby granted, free of charge, to any person obtaining a    */
/* copy of this software and associated documentation files (the "Software"), */
/* to deal in the Software without restriction, including without limitation  */
/* the rights to use, copy, modify, merge, publish, distribute, sublicense,   */
/* and/or sell copies of the Software, and to permit persons to whom the      */
/* Software is furnished to do so, subject to the following conditions:       */
/*                                                                            */
/* The above copyright notice and this permission notice shall be included in */
/* all copies or substantial portions of the Software.                        */
/*                                                                            */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS    */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                 */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.     */
/* IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY       */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT  */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE      */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                              */
/******************************************************************************/

package editor_plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kaijuengine.com/platform/filesystem"
)

func AddGitPluginToStorage(modulePath string) error {
	plugFolder, err := PluginsFolder()
	if err != nil {
		return err
	}

	exists, err := gitPluginAlreadyStored(modulePath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	module := extractModule(modulePath)
	author, packageName := extractAuthorAndPackage(module)

	folderPath := filepath.Join(plugFolder, buildFolderName(modulePath))

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}

	cfg := buildPluginConfig(modulePath, module, author, packageName)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return filesystem.WriteFile(filepath.Join(folderPath, pluginConfigFile), data)
}

func RemoveGitPluginFromStorage(modulePath string) error {
	plugFolder, err := PluginsFolder()
	if err != nil {
		return err
	}

	dirs, err := os.ReadDir(plugFolder)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		cfg, err := loadPluginConfig(plugFolder, dir.Name())
		if err != nil {
			continue
		}

		if cfg.GitModule == modulePath {
			return os.RemoveAll(filepath.Join(plugFolder, dir.Name()))
		}
	}

	return nil
}

func GetStoredGitPlugins() ([]string, error) {
	plugFolder, err := PluginsFolder()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(plugFolder)
	if err != nil {
		return nil, err
	}

	var plugins []string

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		cfg, err := loadPluginConfig(plugFolder, dir.Name())
		if err != nil {
			continue
		}

		if cfg.GitModule != "" {
			plugins = append(plugins, cfg.GitModule)
		}
	}

	return plugins, nil
}

func loadPluginConfig(basePath, dirName string) (PluginConfig, error) {
	cfgPath := filepath.Join(basePath, dirName, pluginConfigFile)

	info, err := os.Stat(cfgPath)
	if err != nil || info.IsDir() {
		return PluginConfig{}, fmt.Errorf("invalid config")
	}

	data, err := filesystem.ReadFile(cfgPath)
	if err != nil {
		return PluginConfig{}, err
	}

	var cfg PluginConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return PluginConfig{}, err
	}

	return cfg, nil
}

func gitPluginAlreadyStored(modulePath string) (bool, error) {
	plugins, err := GetStoredGitPlugins()
	if err != nil {
		return false, err
	}

	for _, plugin := range plugins {
		if plugin == modulePath {
			return true, nil
		}
	}

	return false, nil
}

func extractModule(modulePath string) string {
	return strings.Split(modulePath, "@")[0]
}

// For a module path like "github.com/author/package@ref", this returns "author" and "package"
func extractAuthorAndPackage(module string) (author, packageName string) {
	parts := strings.Split(module, "/")
	n := len(parts)
	return parts[n-2], parts[n-1]
}

func buildFolderName(modulePath string) string {
	replacer := strings.NewReplacer("/", "_", "@", "_", ":", "_")
	return "git_" + replacer.Replace(modulePath)
}

func buildPluginConfig(modulePath, module, author, packageName string) PluginConfig {
	return PluginConfig{
		Name:        packageName,
		PackageName: packageName,
		Description: fmt.Sprintf("Git plugin from %s", module),
		Version:     0.1,
		Author:      author,
		Website:     "https://" + module,
		Enabled:     true,
		GitModule:   modulePath,
	}
}

func parseGitURL(gitURL string) (modulePath, ref string) {
	clean := strings.TrimSpace(gitURL)

	if idx := strings.IndexAny(clean, "?#"); idx != -1 {
		clean = clean[:idx]
	}

	if strings.HasPrefix(clean, "git@") {
		clean = strings.TrimPrefix(clean, "git@")
		clean = strings.Replace(clean, ":", "/", 1)
	}

	clean = strings.TrimPrefix(clean, "https://")
	clean = strings.TrimPrefix(clean, "http://")
	clean = strings.TrimPrefix(clean, "git://")

	clean = strings.TrimSuffix(clean, ".git")
	clean = strings.TrimSuffix(clean, "/")

	ref = "latest"

	if idx := strings.LastIndex(clean, "@"); idx != -1 {
		candidate := clean[idx+1:]
		clean = clean[:idx]

		if candidate != "" {
			ref = candidate
		}
	}

	modulePath = clean
	return
}

func AddPluginFromGit(gitURL string) (string, error) {
	modulePath, ref := parseGitURL(gitURL)

	if strings.Contains(modulePath, "github.com/KaijuEngine/kaiju") {
		modulePath = "kaijuengine.com"
		ref = ""
	}

	fullModuleRef := modulePath
	if ref != "" {
		fullModuleRef = fmt.Sprintf("%s@%s", modulePath, ref)
	}

	if err := AddGitPluginToStorage(fullModuleRef); err != nil {
		return "", fmt.Errorf("failed to save Git plugin to storage: %w", err)
	}

	return fullModuleRef, nil
}

func AddPluginFromGitHub(githubURL string) (string, error) {
	return AddPluginFromGit(githubURL)
}
