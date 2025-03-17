package sql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync"
	"testing"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	drivermysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once = sync.Once{}
	gdb  *gorm.DB
	mock sqlmock.Sqlmock
)

const (
	insertTimerDefSQL = "INSERT INTO `timer_def` *"
	deletTimerDefSQL  = "DELETE FROM `timer_def` (.*)"
	selectTimerDefSQL = "SELECT \\* FROM `timer_def` (.*)"
	countTimerDefSQL  = "SELECT count\\(\\*\\) FROM `timer_def` (.*)"
	updateTimerDefSQL = "UPDATE `timer_def` (.*)"

	columnID         = "id"
	columnDefID      = "def_id"
	columnName       = "name"
	columnApp        = "app"
	columnCreator    = "creator"
	columnStatus     = "status"
	columnCron       = "cron"
	columnNotifyType = "notify_type"
	columnTimerType  = "timer_type"
)

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	var (
		db  *sql.DB
		err error
	)

	once.Do(func() {
		db, mock, err = sqlmock.New()
		if err != nil {
			return
		}
		mock.ExpectQuery("SELECT(.*)").WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7.18-txsql-log"))

		gdb, err = gorm.Open(drivermysql.New(drivermysql.Config{Conn: db}), &gorm.Config{})
	})

	if err != nil {
		return nil, nil, err
	}

	return gdb, mock, nil
}

func Test_TimerDefDAO_Create(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.CreateTimerDefDTO
		wantErr bool
	}{
		{
			name: "fail",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(insertTimerDefSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(errors.New("invalid def"))
				mmock.ExpectRollback()
			},
			req:     &dto.CreateTimerDefDTO{},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(insertTimerDefSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.CreateTimerDefDTO{},
		},
	}

	mockDAO := NewTimerDefDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.Create(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("TimerDefDAO Create() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_TimerDefDAO_Delete(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.DeleteTimerDefDTO
		wantErr bool
	}{
		{
			name:    "invalid defID",
			req:     &dto.DeleteTimerDefDTO{},
			wantErr: true,
		},
		{
			name: "fail",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(deletTimerDefSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(errors.New("invalid def"))
				mmock.ExpectRollback()
			},
			req: &dto.DeleteTimerDefDTO{
				DefID: "111",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(deletTimerDefSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.DeleteTimerDefDTO{
				DefID: "111",
			},
		},
	}

	mockDAO := NewTimerDefDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if err = mockDAO.Delete(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("TimerDefDAO Delete() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_TimerDefDAO_PageQueryTimeList(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.PageQueryTimeDefDTO
		wantErr bool
	}{
		{
			name: "fail",
			mock: func() {
				mmock.ExpectQuery(selectTimerDefSQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnDefID, columnCreator, columnApp, columnStatus, columnCron, columnNotifyType, columnTimerType}).
						AddRow(1, "test", "111", "weixxxu", "test", 1, "0 */2 * * * ? *", 3, 2),
					sqlmock.NewRows([]string{columnID, columnName, columnDefID, columnCreator, columnApp, columnStatus, columnCron, columnNotifyType, columnTimerType}).
						AddRow(2, "test", "222", "weixxxu", "test", 1, "0 */5 * * * ? *", 3, 2),
				).WillReturnError(errors.New("error"))
			},
			req: &dto.PageQueryTimeDefDTO{
				Name:    "test",
				Creator: "weixxxu",
			},
			wantErr: true,
		},
		{
			name: "success",
			mock: func() {
				mmock.ExpectQuery(selectTimerDefSQL).WillReturnRows(
					sqlmock.NewRows([]string{columnID, columnName, columnDefID, columnCreator, columnApp, columnStatus, columnCron, columnNotifyType, columnTimerType}).
						AddRow(1, "test", "111", "weixxxu", "test", 1, "0 */2 * * * ? *", 3, 2),
					sqlmock.NewRows([]string{columnID, columnName, columnDefID, columnCreator, columnApp, columnStatus, columnCron, columnNotifyType, columnTimerType}).
						AddRow(2, "test", "222", "weixxxu", "test", 1, "0 */5 * * * ? *", 3, 2),
				).WillReturnError(nil)
			},
			req: &dto.PageQueryTimeDefDTO{
				Name:    "test",
				Creator: "weixxxu",
			},
		},
	}

	mockDAO := NewTimerDefDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if _, err = mockDAO.PageQueryTimeList(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("TimerDefDAO PageQueryTimeList() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}

func Test_TimerDefDAO_Count(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	mmock.ExpectQuery(countTimerDefSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = NewTimerDefDAO(mysql.NewClient(mockdb)).Count(&dto.CountTimerDefDTO{
		Creator: "weixxxu",
		Name:    "test",
	})
	assert.Nil(t, err)
}

func Test_TimerDefDAO_CountByStatus(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	mmock.ExpectQuery(countTimerDefSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = NewTimerDefDAO(mysql.NewClient(mockdb)).CountByStatus(1)
	assert.Nil(t, err)
}

func Test_TimerDefDAO_UpdateStatus(t *testing.T) {
	mockdb, mmock, err := getDBMock()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		mock    func()
		req     *dto.UpdateTimerDefDTO
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mmock.ExpectBegin()
				mmock.ExpectExec(updateTimerDefSQL).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)
				mmock.ExpectCommit()
			},
			req: &dto.UpdateTimerDefDTO{},
		},
	}

	mockDAO := NewTimerDefDAO(mysql.NewClient(mockdb))
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		if err = mockDAO.UpdateStatus(tt.req); (err != nil) != tt.wantErr {
			t.Errorf("UpdateStatus() err got = %v, expect = %t", err, tt.wantErr)
		}
	}
}
