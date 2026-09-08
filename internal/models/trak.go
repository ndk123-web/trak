package models

import (
	"time"
)

type WorkspaceMetadata struct {
	Id              string          `json:"id"`
	Template        string          `json:"template"`
	TemplateVersion string          `json:"template_version"`
	CreatedAt       time.Time       `json:"created_at"`
	Author          string          `json:"author"`
	Source          string          `json:"source"`
	Version         string          `json:"version"`
	Repository      string          `json:"repository"`
	Name            string          `json:"name"`
	ModuleBreakdown map[string]bool `json:"module_breakdown"`
}

type UserConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SystemConfigWorkspaceModel struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	LastOpened string `json:"lastOpened"`
	TrackId    string `json:"trackId"`
}

type TrakUserConfig struct {
	UserName   string                       `json:"username"`
	Email      string                       `json:"email"`
	Password   string                       `json:"password"`
	Workspaces []SystemConfigWorkspaceModel `json:"workspaces"`
}

func (cf *TrakUserConfig) SetEmail(email string) *TrakUserConfig {
	if email != "" {
		cf.Email = email
	}
	return cf
}

func (cf *TrakUserConfig) SetPassword(password string) *TrakUserConfig {
	if password != "" {
		cf.Password = password
	}
	return cf
}

func (cf *TrakUserConfig) SetUsername(username string) *TrakUserConfig {
	if username != "" {
		cf.UserName = username
	}
	return cf
}
