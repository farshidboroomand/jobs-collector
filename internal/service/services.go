package service

// Services is a container for all service instances, providing a centralized access point for business logic across the application.
type Services struct {
	Bot       *BotService
	Continent *ContinentService
}
