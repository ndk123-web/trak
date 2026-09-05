package config

import "fmt"

type Config struct {
	RawBaseUrl       string
	RegistryName     string
	RepositoryBranch string
	GithubUsername   string
	Version          string
}

var AppVersion = "v1.3.0"

var TrakConfig = &Config{
	RawBaseUrl:       "",
	RegistryName:     "trak-registry",
	RepositoryBranch: "main",
	GithubUsername:   "ndk123-web",
	Version:          "v1.3.0",
}

func UpdateBaseUrl() {
	if AppVersion != "" {
		TrakConfig.Version = AppVersion
	}
	var BaseUrl string = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/refs/heads/%s/", TrakConfig.GithubUsername, TrakConfig.RegistryName, TrakConfig.RepositoryBranch)
	TrakConfig.RawBaseUrl = BaseUrl
}
