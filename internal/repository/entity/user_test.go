package entity

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestUserEntity_Create(t *testing.T) {
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB

		// 输入
		ctx  context.Context
		user User

		/*结果error*/
		wantError error
	}{
		{
			name: "创建用户成功",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.Equal(t, err, nil)
				// 对于创建、删除、改，主要是rowsAffteced和lastinsertid
				mockRes := sqlmock.NewResult(int64(1), 1)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(mockRes)
				return db
			},

			ctx:       context.Background(),
			user:      User{},
			wantError: nil,
		},
		{
			name: "创建用户重复",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.Equal(t, err, nil)
				// 唯一索引冲突的error
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(&mysql2.MySQLError{Number: 1062})
				return db
			},

			ctx:       context.Background(),
			user:      User{},
			wantError: ErrUserDuplciate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
				// 不调用show version
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// 禁止ping数据库
				DisableAutomaticPing: true,
				// 不开启默认的事务
				SkipDefaultTransaction: true,
			})
			// 初始化数据库必须没有Error
			assert.Equal(t, err, nil)
			entity := NewUserEntity(db)
			err = entity.Create(tc.ctx, tc.user)
			assert.Equal(t, err, tc.wantError)
		})
	}
}
