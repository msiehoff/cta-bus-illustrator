package app

import "github.com/msiehoff/cta-bus-illustrator/backend/business"

type RouteRepository interface {
	GetRoutes() ([]business.Route, error)
}

type RouteService struct {
	repo RouteRepository
}

func NewRouteService(repo RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) GetRoutes() ([]business.Route, error) {
	return s.repo.GetRoutes()
}
