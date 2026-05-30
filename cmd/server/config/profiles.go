package config

import (
	"speaking_hearts/models"
)

// GetDefaultRoutingRules returns the initial translation routing rules.
// This separates the configuration from the core processing loop, paving the way
// for dynamic, admin-controlled setups in the future.
func GetDefaultRoutingRules() []models.RoutingRule {
	return []models.RoutingRule{
		{SourceLang: "es", TargetLang: "en", Active: true},
		{SourceLang: "es", TargetLang: "ru", Active: true},
		{SourceLang: "es", TargetLang: "de", Active: true},
	}
}
