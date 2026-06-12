package application

func (a *App) registerRoutes() {
	v1 := a.Router.Group("v1")

	collector := v1.Group("/jobs")

	collector.POST("/linkedin", CollectByLinkedin())
}
