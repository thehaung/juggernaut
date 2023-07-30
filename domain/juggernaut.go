package domain

import "context"

type Juggernaut struct {
	Command int
}

type JuggernautUseCase interface {
	GetServerPublicIP(ctx context.Context) (string, error)
	GetCurrentIP(ctx context.Context) (string, error)
}

type JuggernautRepository interface {
	GetServerIPViaIpify(ctx context.Context, resultChan chan<- string) error
	GetServerIPViaSeeIP(ctx context.Context, resultChan chan<- string) error
	GetCurrentIP(ctx context.Context) (string, error)
}
