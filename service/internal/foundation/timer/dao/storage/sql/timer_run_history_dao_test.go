package sql

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	insertRunHistorySQL       = "INSERT INTO `run_history` *"
	deletRunHistorySQL        = "DELETE FROM `run_history` (.*)"
	selectRunHistorySQL       = "SELECT \\* FROM `run_history` (.*)"
	selectIDFromRunHistorySQL = "SELECT `id` FROM `run_history` (.*)"
	countRunHistorySQL        = "SELECT count\\(\\*\\) FROM `run_history` (.*)"
	updateRunHistorySQL       = "UPDATE `run_history` (.*)"
)

func Test_RunHistoryDAO_Create(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.CreateRunHistoryDTO
		wantErr bool
	}{
		{
			name:    "invalid param",
			req:     &dto.CreateRunHistoryDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(insertRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(fmt.Errorf("invalid history"))
				mmock.ExpectRollback()
			},
			req: &dto.CreateRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(insertRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.CreateRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
		},
	}

	mockDAO := NewRunHistoryDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.Create(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("RunHistoryDAO Create() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_RunHistoryDAO_Get(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.GetRunHistoryDTO
		wantErr bool
	}{
		{
			name:    "invalid param",
			req:     &dto.GetRunHistoryDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectQuery(selectRunHistorySQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName}).
						AddRow(1, "test"),
				).WillReturnError(fmt.Errorf("error"))
			},
			req: &dto.GetRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectQuery(selectRunHistorySQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName}).
						AddRow(1, "test"),
				).WillReturnError(nil)
			},
			req: &dto.GetRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
		},
	}

	mockDAO := NewRunHistoryDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.Get(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("RunHistoryDAO Get() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_RunHistoryDAO_Delete(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.DeleteRunHistoryDTO
		wantErr bool
	}{
		{
			name:    "invalid defID",
			req:     &dto.DeleteRunHistoryDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(deletRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(fmt.Errorf("invalid def"))
				mmock.ExpectRollback()
			},
			req: &dto.DeleteRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(deletRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.DeleteRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
		},
	}

	mockDAO := NewRunHistoryDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if err = mockDAO.Delete(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("RunHistoryDAO Delete() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_RunHistoryDAO_Update(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.UpdateRunHistoryDTO
		wantErr bool
	}{
		{
			name:    "invalid defID",
			req:     &dto.UpdateRunHistoryDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(updateRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(fmt.Errorf("error"))
				mmock.ExpectRollback()
			},
			req: &dto.UpdateRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(updateRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.UpdateRunHistoryDTO{
				DefID:    "111",
				RunTimer: "2022-06-01",
			},
		},
	}

	mockDAO := NewRunHistoryDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if err = mockDAO.Update(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("RunHistoryApp Update() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_RunHistoryDAO_PageQuery(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.PageQueryRunHistoryDTO
		wantErr bool
	}{
		{
			name: "fail",
			mock: func() {
				mmock.ExpectQuery(selectRunHistorySQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(1, "test", "weixxxu"),
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(2, "test2", "weixxxu"),
				).WillReturnError(fmt.Errorf("error"))
			},
			req: &dto.PageQueryRunHistoryDTO{
				Name:      "test",
				PageQuery: &constants.PageQuery{},
				Order:     &constants.Order{},
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectQuery(selectRunHistorySQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(1, "test", "weixxxu"),
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(2, "test2", "weixxxu"),
				).WillReturnError(nil)
			},
			req: &dto.PageQueryRunHistoryDTO{
				Name:      "test",
				PageQuery: &constants.PageQuery{},
				Order:     &constants.Order{},
			},
		},
	}

	mockDAO := NewRunHistoryDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.PageQuery(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("RunHistoryDAO PageQuery() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_RunHistoryDAO_Count(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	mmock.ExpectQuery(countRunHistorySQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = NewRunHistoryDAO(mysql.NewClient(mockdb)).Count(&dto.PageQueryRunHistoryDTO{
		Name: "test",
	})
	assert.Nil(t, err)
}

func Test_RunHistoryDAO_DeleteByRunTime(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		time    time.Time
		wantErr bool
	}{
		{
			name: "select fail",
			mock: func() {
				mmock.ExpectQuery(selectIDFromRunHistorySQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID}).
						AddRow(1),
					sqlmock.NewRows([]string{columnID}).
						AddRow(2),
				).WillReturnError(fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "delete fail",
			mock: func() {
				mmock.ExpectQuery(selectIDFromRunHistorySQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID}).
						AddRow(1),
					sqlmock.NewRows([]string{columnID}).
						AddRow(2),
				).WillReturnError(nil)

				mmock.ExpectBegin()
				mmock.ExpectExec(deletRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(fmt.Errorf("invalid def"))
				mmock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "delete success",
			mock: func() {
				mmock.ExpectQuery(selectIDFromRunHistorySQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID}).
						AddRow(1),
					sqlmock.NewRows([]string{columnID}).
						AddRow(2),
				).WillReturnError(nil)

				mmock.ExpectBegin()
				mmock.ExpectExec(deletRunHistorySQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
		},
	}

	mockDAO := NewRunHistoryDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if err = mockDAO.DeleteByRunTime(tt.time); (err != nil) != tt.wantErr {
			t.Errorf("RunHistoryDAO DeleteByRunTime() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}
