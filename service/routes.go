package service

func (s *Service) registerRoutes() {
	//s.r.GET("/", s.Homepage)

	s.r.POST("/users", s.CreateUser)

	authorized := s.r.Group("/")
	//authorized.Use(s.ValidateToken)
	{
		authorized.GET("/users/:id", s.GetUser)
		authorized.POST("/users/:id", s.UpdateUser)
	}
}
