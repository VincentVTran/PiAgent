package config

var ApplicationConfig struct {
	Local struct {
		RabbitMQLink string `json:"rabbitMQLink"`
		Exchange     string `json:"exchange"`
		RoutingKey   string `json:"routingKey"`
	} `json:"local"`
	Prod struct {
		RabbitMQLink string `json:"rabbitMQLink"`
		Exchange     string `json:"exchange"`
		RoutingKey   string `json:"routingKey"`
	} `json:"prod"`
}
