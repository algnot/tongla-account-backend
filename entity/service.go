package entity

import (
	"gorm.io/gorm"
	"time"
	"tongla-account/di/config"
)

type Service struct {
	Name         string         `json:"name" gorm:"type:varchar(255)"`
	ClientId     string         `json:"client_id" gorm:"type:varchar(255);primarykey"`
	ClientSecret string         `json:"client_secret" gorm:"type:varchar(255)"`
	RedirectUri  string         `json:"redirect_uri" gorm:"type:varchar(255)"`
	Issuer       string         `json:"issuer" gorm:"type:varchar(255)"`
	Owner        string         `json:"owner" gorm:"type:varchar(255);index"`
	Scope        string         `json:"scope" gorm:"type:varchar(255)"`
	GrantType    string         `json:"grant_type" gorm:"type:varchar(255)"`
	ResponseType string         `json:"response_type" gorm:"type:varchar(255)"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ServicesResponse struct {
	Name                   string `json:"name"`
	ClientId               string `json:"client_id"`
	ClientSecret           string `json:"client_secret"`
	RedirectUri            string `json:"redirect_uri"`
	Issuer                 string `json:"issuer"`
	Scopes                 string `json:"scopes"`
	GrantTypes             string `json:"grant_types"`
	ResponseType           string `json:"response_type"`
	OpenidConfigurationUri string `json:"openid_configuration_uri"`
}

func (service Service) ToResponse() ServicesResponse {
	appConfig := config.GetConfig()
	return ServicesResponse{
		Name:                   service.Name,
		ClientId:               service.ClientId,
		ClientSecret:           service.ClientSecret,
		RedirectUri:            service.RedirectUri,
		Issuer:                 service.Issuer,
		Scopes:                 service.Scope,
		GrantTypes:             service.GrantType,
		ResponseType:           service.ResponseType,
		OpenidConfigurationUri: appConfig.ServerConfig.BackendPath + "/openid/.well-known/configuration",
	}
}
