//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	err := healthCheck(defaultAttempts)
	if err != nil {
		panic(fmt.Errorf("integration tests: host %s is not available: %w", basePath, err))
	}

	log.Infof("integration tests: host %s is available", basePath)
	os.Exit(m.Run())
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Infof("integration tests: host %s is not available, attempts left: %d", basePath, attempts)
		time.Sleep(time.Second)
		attempts--
	}

	return err
}
