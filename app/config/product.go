package config

var productionMode bool = false

func IsProductionEnabled() bool {
	return productionMode
}
