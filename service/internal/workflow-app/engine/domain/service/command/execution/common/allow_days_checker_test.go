package common

import (
	"reflect"
	"testing"
	"time"

	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/stretchr/testify/suite"
)

// TestAllowDaysSuite 测试
func TestAllowDaysSuite(t *testing.T) {
	suite.Run(t, new(allowDaysSuite))
}

type allowDaysSuite struct {
	suite.Suite

	allowDaysChecker AllowDaysChecker
}

// SetupTest 执行用例执行前准备工作
func (s *allowDaysSuite) SetupTest() {
	s.allowDaysChecker, _ = NewDefaultAllowDaysChecker()
}

// TestCheckAllowDays CheckAllowDays测试用例
func (s *allowDaysSuite) TestCheckAllowDays() {
	type args struct {
		allowDaysPolicy entity.AllowDaysPolicy
		triggerName     string
		startDate       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    bool
	}{
		{
			name: "triggerName any success",
			args: args{
				allowDaysPolicy: entity.Any,
				triggerName:     "any",
				startDate:       "2021-08-11",
			},
			wantErr: nil,
			want:    true,
		},
		{
			name: "triggerName week success",
			args: args{
				allowDaysPolicy: entity.Week,
				triggerName:     "week",
				startDate:       "2021-08-11",
			},
			wantErr: nil,
			want:    true,
		},
		{
			name: "triggerName weekend success",
			args: args{
				allowDaysPolicy: entity.Weekend,
				triggerName:     "weekend",
				startDate:       "2021-08-14",
			},
			wantErr: nil,
			want:    true,
		},
		{
			name: "triggerName week fail",
			args: args{
				allowDaysPolicy: entity.Week,
				triggerName:     "week",
				startDate:       "2021-10-01",
			},
			wantErr: nil,
			want:    false,
		},
		{
			name: "triggerName weekend fail",
			args: args{
				allowDaysPolicy: entity.Weekend,
				triggerName:     "weekend",
				startDate:       "2021-08-11",
			},
			wantErr: nil,
			want:    false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			dateTime, _ := time.Parse("2006-01-02", tt.args.startDate)
			got, err := s.allowDaysChecker.Check(dateTime, tt.args.allowDaysPolicy)
			s.Equal(tt.wantErr, err, "CheckAllowDays() got = %v, want %v", err, tt.wantErr)
			if !reflect.DeepEqual(got, tt.want) {
				s.T().Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
