package requests

import (
	"atlas-cos/rest/response"
)

const (
	BaseRequest string = "http://atlas-nginx:80"
)

type configuration struct {
	retries int
	mappers []response.ConditionalMapperProvider
}

type Configurator func(c *configuration)

//goland:noinspection GoUnusedExportedFunction
func SetRetries(amount int) Configurator {
	return func(c *configuration) {
		c.retries = amount
	}
}

//goland:noinspection GoUnusedExportedFunction
func AddMappers(mappers []response.ConditionalMapperProvider) Configurator {
	return func(c *configuration) {
		c.mappers = append(c.mappers, mappers...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func AddMapper(mapper response.ConditionalMapperProvider) Configurator {
	return func(c *configuration) {
		c.mappers = append(c.mappers, mapper)
	}
}
