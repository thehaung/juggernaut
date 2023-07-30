package usecase

import (
	"context"
	"github.com/thehaung/juggernaut/domain"
	"github.com/thehaung/juggernaut/internal/logger"
	"time"
)

type juggernautUseCase struct {
	juggernautRepository domain.JuggernautRepository
	logger               logger.Interface
}

func NewJuggernautUseCase(repository domain.JuggernautRepository) domain.JuggernautUseCase {
	return &juggernautUseCase{juggernautRepository: repository, logger: logger.GetLogger()}
}

func (j *juggernautUseCase) GetServerPublicIP(ctx context.Context) (string, error) {
	resultChan := make(chan string, 1)
	ctx, cancel := context.WithTimeout(ctx, time.Minute)

	defer cancel()
	go func() {
		err := j.juggernautRepository.GetServerIPViaIpify(ctx, resultChan)
		if err != nil {
			j.logger.Errorf("Exec GetServerIPViaIpify. Error: %s", err.Error())
			return
		}
	}()

	go func() {
		err := j.juggernautRepository.GetServerIPViaSeeIP(ctx, resultChan)
		if err != nil {
			j.logger.Errorf("Exec GetServerIPViaSeeIP. Error: %s", err.Error())
			return
		}
	}()

	res := <-resultChan
	j.logger.Infof("Exec GetServerPublicIP. Success IP: %s", res)

	return res, nil
}

func (j *juggernautUseCase) GetCurrentIP(ctx context.Context) (string, error) {
	return j.juggernautRepository.GetCurrentIP(ctx)
}
