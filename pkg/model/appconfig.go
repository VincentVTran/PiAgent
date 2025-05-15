package config

var ApplicationConfig struct {
	RabbitMQ struct {
		Local      string `json:"local"`
		Prod       string `json:"prod"`
		Exchange   string `json:"exchange"`
		RoutingKey string `json:"routingKey"`
	} `json:"rabbitmq"`
}
