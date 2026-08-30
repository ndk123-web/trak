package config

import "fmt"

type Config struct {
	RawBaseUrl       string
	RegistryName     string
	RepositoryBranch string
	GithubUsername   string
}

var TrakConfig = &Config{
	RawBaseUrl:       "",
	RegistryName:     "trak-registry",
	RepositoryBranch: "main",
	GithubUsername:   "ndk123-web",
}

func UpdateBaseUrl() {
	var BaseUrl string = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/refs/heads/%s/", TrakConfig.GithubUsername, TrakConfig.RegistryName, TrakConfig.RepositoryBranch)
	TrakConfig.RawBaseUrl = BaseUrl
}
