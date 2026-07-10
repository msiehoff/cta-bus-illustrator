package fake

import "context"

// PipelineRouteProvider returns a small fixed route set for local pipeline testing.
type PipelineRouteProvider struct{}

func NewPipelineRouteProvider() *PipelineRouteProvider {
	return &PipelineRouteProvider{}
}

func (p *PipelineRouteProvider) GetRouteIDs(_ context.Context) ([]string, error) {
	return []string{"8", "66"}, nil
}
