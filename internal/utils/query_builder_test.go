package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// テスト用のモデル構造体
type TestModel struct {
	ID         uint
	CreatedAt  string
	UpdatedAt  string
	gorm.Model `gorm:"-"`
}

// テーブル名を指定するメソッドを追加
func (TestModel) TableName() string {
	return "test"
}

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string]string
		wantQuery   string
		wantVars    []interface{}
	}{
		{
			name: "created_at__gteとcreated_at__lteが指定された場合",
			queryParams: map[string]string{
				"created_at__gte": "2023-01-01",
				"created_at__lte": "2023-01-03",
			},
			wantQuery: "SELECT * FROM `test` WHERE created_at >= ? AND created_at <= ?",
			wantVars:  []interface{}{"2023-01-01", "2023-01-03"},
		},
		{
			name: "created_at__gtとcreated_at__ltが指定された場合",
			queryParams: map[string]string{
				"created_at__gt": "2023-01-01",
				"created_at__lt": "2023-01-03",
			},
			wantQuery: "SELECT * FROM `test` WHERE created_at > ? AND created_at < ?",
			wantVars:  []interface{}{"2023-01-01", "2023-01-03"},
		},
		{
			name: "updated_at__gteとupdated_at__lteが指定された場合",
			queryParams: map[string]string{
				"updated_at__gte": "2023-01-01",
				"updated_at__lte": "2023-01-03",
			},
			wantQuery: "SELECT * FROM `test` WHERE updated_at >= ? AND updated_at <= ?",
			wantVars:  []interface{}{"2023-01-01", "2023-01-03"},
		},
		{
			name: "updated_at__gtとupdated_at__ltが指定された場合",
			queryParams: map[string]string{
				"updated_at__gt": "2023-01-01",
				"updated_at__lt": "2023-01-03",
			},
			wantQuery: "SELECT * FROM `test` WHERE updated_at > ? AND updated_at < ?",
			wantVars:  []interface{}{"2023-01-01", "2023-01-03"},
		},
		{
			name: "ページネーションが指定された場合",
			queryParams: map[string]string{
				"page":      "2",
				"page_size": "10",
			},
			wantQuery: "SELECT * FROM `test` LIMIT 10 OFFSET 10",
			wantVars:  nil,
		},
		{
			name: "無効な日付が指定された場合",
			queryParams: map[string]string{
				"created_at__gte": "invalid_date",
				"created_at__lte": "2023-01-01",
			},
			wantQuery: "SELECT * FROM `test`",
			wantVars:  []interface{}{"invalid_date", "2023-01-01"},
		},
		{
			name: "無効なページネーションが指定された場合",
			queryParams: map[string]string{
				"page":      "invalid_page",
				"page_size": "invalid_page_size",
			},
			wantQuery: "SELECT * FROM `test` LIMIT 20",
			wantVars:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			// SQLiteのインメモリデータベースを設定
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			assert.NoError(t, err)

			// テスト用のテーブルを作成
			err = db.AutoMigrate(&TestModel{})
			assert.NoError(t, err)

			// テストデータを作成
			testData := []TestModel{
				{
					ID:        1,
					CreatedAt: "2023-01-01",
					UpdatedAt: "2023-01-01",
				},
				{
					ID:        2,
					CreatedAt: "2023-01-02",
					UpdatedAt: "2023-01-02",
				},
				{
					ID:        3,
					CreatedAt: "2023-01-03",
					UpdatedAt: "2023-01-03",
				},
			}

			// テストデータを挿入
			for _, data := range testData {
				result := db.Create(&data)
				assert.NoError(t, result.Error)
			}

			// クエリを構築
			query := BuildQuery(c, db)

			// 実際のSQLクエリを取得するために DryRun を実行
			stmt := query.Session(&gorm.Session{DryRun: true}).Find(&TestModel{}).Statement
			assert.NotNil(t, stmt)

			// 生成されたSQLクエリを取得
			sql := stmt.SQL.String()
			assert.Contains(t, sql, tt.wantQuery)
			assert.Equal(t, tt.wantVars, stmt.Vars)
		})
	}
}
