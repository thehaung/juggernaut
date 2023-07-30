package repository

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/thehaung/juggernaut/domain"
	"github.com/thehaung/juggernaut/internal/logger"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	_ipifyUrl = "https://api64.ipify.org?format=json"
	_seeIPUrl = "https://api.seeip.org/jsonip?"
)

const (
	_ipifyPartnerId = iota + 1
	_seeIpPartnerId
)

var (
	ipCache = cache.New(5*time.Minute, 10*time.Minute)
)

type juggernautRepository struct {
	logger logger.Interface
}

func NewJuggernautRepository() domain.JuggernautRepository {
	return &juggernautRepository{logger: logger.GetLogger()}
}

func (j *juggernautRepository) GetServerIPViaIpify(ctx context.Context, resultChan chan<- string) error {
	ip, err := j.getIPViaPartner(ctx, _ipifyPartnerId)
	if err != nil {
		return err
	}

	resultChan <- ip
	return nil
}

func (j *juggernautRepository) GetServerIPViaSeeIP(ctx context.Context, resultChan chan<- string) error {
	ip, err := j.getIPViaPartner(ctx, _seeIpPartnerId)
	if err != nil {
		return err
	}

	resultChan <- ip
	return nil
}

func (j *juggernautRepository) getIPViaPartner(ctx context.Context, partnerId int) (string, error) {
	resp, err := http.Get(j.getPartnerUrl(partnerId))
	if err != nil {
		j.logger.Errorf("Exec http.Get failed. Error: %s", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		j.logger.Errorf("Exec io.ReadAll failed. Error: %s", err.Error())
		return "", err
	}

	respString := string(respBytes)
	j.logger.Infof("PartnerId: %d Response: %s", partnerId, respString)
	ip := gjson.Get(respString, "ip").String()

	if strings.TrimSpace(ip) == "" {
		err = errors.New("exec gjson.GetBytes failed. A IP property is empty")
		j.logger.Error(err.Error())

		return "", err
	}

	if err = ipCache.Add("IP", ip, time.Minute); err != nil {
		j.logger.Errorf("Exec ipCache.Add failed. Error: %s", err.Error())
	}

	return ip, nil
}

func (j *juggernautRepository) getPartnerUrl(partnerId int) string {
	if partnerId == _ipifyPartnerId {
		return _ipifyUrl
	}

	return _seeIPUrl
}

func (j *juggernautRepository) GetCurrentIP(ctx context.Context) (string, error) {
	ip, found := ipCache.Get("IP")
	if !found {
		return "", errors.New("error not found value of key IP in memcache")
	}

	return ip.(string), nil
}
