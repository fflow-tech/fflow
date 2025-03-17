package sql

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	insertAppSQL = "INSERT INTO `app` *"
	deletAppSQL  = "DELETE FROM `app` (.*)"
	selectAppSQL = "SELECT \\* FROM `app` (.*)"
	countAppSQL  = "SELECT count\\(\\*\\) FROM `app` (.*)"
	updateAppSQL = "UPDATE `app` (.*)"
)

func Test_AppDAO_Create(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.CreateAppDTO
		wantErr bool
	}{
		{
			name:    "invalid param",
			req:     &dto.CreateAppDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(insertAppSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(errors.New("invalid app"))
				mmock.ExpectRollback()
			},
			req: &dto.CreateAppDTO{
				Name:    "test",
				Creator: "weixxxu",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(insertAppSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.CreateAppDTO{
				Name:    "test",
				Creator: "weixxxu",
			},
		},
	}

	mockDAO := NewAppDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.Create(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("AppDAO Create() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_AppDefDAO_Delete(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.DeleteAppDTO
		wantErr bool
	}{
		{
			name:    "invalid param",
			req:     &dto.DeleteAppDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(deletAppSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(errors.New("invalid app"))
				mmock.ExpectRollback()
			},
			req: &dto.DeleteAppDTO{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(deletAppSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.DeleteAppDTO{
				Name: "test",
			},
		},
	}

	mockDAO := NewAppDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if err = mockDAO.Delete(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("AppDAO Delete() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_AppDAO_PageQuery(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.PageQueryAppDTO
		wantErr bool
	}{
		{
			name: "fail",
			mock: func() {
				mmock.ExpectQuery(selectAppSQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(1, "test", "weixxxu"),
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(2, "test2", "weixxxu"),
				).WillReturnError(errors.New("error"))
			},
			req: &dto.PageQueryAppDTO{
				Name:      "test",
				Creator:   "weixxxu",
				PageQuery: &constants.PageQuery{},
				Order:     &constants.Order{},
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectQuery(selectAppSQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(1, "test", "weixxxu"),
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(2, "test2", "weixxxu"),
				).WillReturnError(nil)
			},
			req: &dto.PageQueryAppDTO{
				Name:      "test",
				Creator:   "weixxxu",
				PageQuery: &constants.PageQuery{},
				Order:     &constants.Order{},
			},
		},
	}

	mockDAO := NewAppDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.PageQuery(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("AppDAO PageQueryTimeList() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_AppDAO_Get(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.GetAppDTO
		wantErr bool
	}{
		{
			name:    "invalid param",
			req:     &dto.GetAppDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectQuery(selectAppSQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(1, "test", "weixxxu"),
				).WillReturnError(errors.New("error"))
			},
			req: &dto.GetAppDTO{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectQuery(selectAppSQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnCreator}).
						AddRow(1, "test", "weixxxu"),
				).WillReturnError(nil)
			},
			req: &dto.GetAppDTO{
				Name: "test",
			},
		},
	}

	mockDAO := NewAppDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.Get(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("AppDAO Get() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_AppDefDAO_Count(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	mmock.ExpectQuery(countAppSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = NewAppDAO(mysql.NewClient(mockdb)).Count(&dto.CountAppDTO{
		Creator: "weixxxu",
		Name:    "test",
	})
	assert.Nil(t, err)
}
