package polling

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWorkerPool struct{}

func (m *mockWorkerPool) Submit(f func()) error {
	f()
	return fmt.Errorf("submit fail")
}

type mockPollingProxy struct{}

func (m *mockPollingProxy) GetPollingTaskWorkLock() (string, error) {
	return "", nil
}

func (m *mockPollingProxy) SendPollingTaskWork(timeSlice string) error {
	return nil
}

func Test_Polling_Type(t *testing.T) {
	service := &Polling{}
	assert.Equal(t, "Polling", service.Type())
}

func Test_Polling_Restart(t *testing.T) {
	mockService := &Polling{
		command: &mockPollingProxy{},
		pool:    &mockWorkerPool{},
		workChan: []chan bool{
			make(chan bool, 1),
		},
	}
	assert.Nil(t, mockService.Restart())
}
