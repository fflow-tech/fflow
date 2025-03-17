package repo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPollingTaskDAO struct {
}

func (m *mockPollingTaskDAO) GetBucketNum() int {
	return 0
}
func (m *mockPollingTaskDAO) GetTaskBucketID(hashID string) (string, error) {
	if hashID == "error" {
		return "", errors.New("invalid hashID")
	}
	return "abcd", nil
}
func (m *mockPollingTaskDAO) SetBucketNum(num int) error {
	if num == -1 {
		return errors.New("invalid num ")
	}
	return nil
}
func (m *mockPollingTaskDAO) SetTimeSlice(timeDuration string) error {
	if timeDuration == "error" {
		return errors.New("invalid timeSlice")
	}
	return nil
}
func (m *mockPollingTaskDAO) GetTimeSlice(timeDuration string) error {
	if timeDuration == "error" {
		return errors.New("invalid timeSlice")
	}
	return nil
}
func (m *mockPollingTaskDAO) SuccessTimeSlice(timeDuration string) error {
	if timeDuration == "error" {
		return errors.New("invalid timeSlice")
	}
	return nil
}

func Test_PollingTaskRepo_GetTaskBucketID(t *testing.T) {
	tests := []struct {
		name    string
		hashID  string
		wantErr bool
	}{
		{
			name:    "error",
			hashID:  "error",
			wantErr: true,
		},
		{
			name:   "success",
			hashID: "success",
		},
	}

	mockRepo := &PollingTaskRepo{&mockPollingTaskDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockRepo.GetTaskBucketID(tt.hashID); (err != nil) != tt.wantErr {
				t.Errorf("GetTaskBucketID() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_PollingTaskRepo_SetBucketNum(t *testing.T) {
	tests := []struct {
		name    string
		num     int
		wantErr bool
	}{
		{
			name:    "error",
			num:     -1,
			wantErr: true,
		},
		{
			name: "success",
			num:  2,
		},
	}

	mockRepo := &PollingTaskRepo{&mockPollingTaskDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.SetBucketNum(tt.num); (err != nil) != tt.wantErr {
				t.Errorf("SetBucketNum() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_PollingTaskRepo_SetTimeSlice(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		wantErr  bool
	}{
		{
			name:     "error",
			duration: "error",
			wantErr:  true,
		},
		{
			name:     "success",
			duration: "30 seconds",
		},
	}

	mockRepo := &PollingTaskRepo{&mockPollingTaskDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.SetTimeSlice(tt.duration); (err != nil) != tt.wantErr {
				t.Errorf("SetTimeSlice() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_PollingTaskRepo_GetTimeSlice(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		wantErr  bool
	}{
		{
			name:     "error",
			duration: "error",
			wantErr:  true,
		},
		{
			name:     "success",
			duration: "30 seconds",
		},
	}

	mockRepo := &PollingTaskRepo{&mockPollingTaskDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.GetTimeSlice(tt.duration); (err != nil) != tt.wantErr {
				t.Errorf("GetTimeSlice() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_PollingTaskRepo_SuccessTimeSlice(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		wantErr  bool
	}{
		{
			name:     "error",
			duration: "error",
			wantErr:  true,
		},
		{
			name:     "success",
			duration: "30 seconds",
		},
	}

	mockRepo := &PollingTaskRepo{&mockPollingTaskDAO{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockRepo.SuccessTimeSlice(tt.duration); (err != nil) != tt.wantErr {
				t.Errorf("SuccessTimeSlice() got err: %v, expect err: %t", err, tt.wantErr)
			}
		})
	}
}

func Test_PollingTaskRepo_GetBucketNum(t *testing.T) {
	mockRepo := &PollingTaskRepo{&mockPollingTaskDAO{}}
	assert.Equal(t, 0, mockRepo.GetBucketNum())
}
