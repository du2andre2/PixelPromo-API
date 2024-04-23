package route

import "github.com/gin-gonic/gin"

type Server interface {
	Run()
}

type server struct {
	route Route
}

func NewServer(
	route Route,
) Server {
	return &server{
		route: route,
	}
}

func (s *server) Run() {
	router := gin.Default()
	s.route.Setup(router)
	router.Run("localhost:8080")
}
