package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	_ "github.com/sijms/go-ora/v2"
)

// TestRequest 连通性测试请求
type TestRequest struct {
	DBType   string `json:"db_type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// TestResult 连通性测试结果
type TestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TestConnectivity 测试数据库连通性
func TestConnectivity(req TestRequest) TestResult {
	switch req.DBType {
	case "mysql":
		return testMySQL(req)
	case "postgresql":
		return testPostgreSQL(req)
	case "oracle":
		return testOracle(req)
	case "redis":
		return testRedis(req)
	case "sqlserver":
		return testSQLServer(req)
	// case "sqlite":
	// 	return testSQLite(req)
	default:
		return TestResult{Success: false, Message: "不支持的数据源类型"}
	}
}

// testMySQL 测试MySQL连通性
func testMySQL(req TestRequest) TestResult {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=5s",
		req.Username, req.Password, req.Host, req.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}

	return TestResult{Success: true, Message: "MySQL 连接成功!"}
}

// testPostgreSQL 测试PostgreSQL连通性
func testPostgreSQL(req TestRequest) TestResult {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable connect_timeout=5",
		req.Host, req.Port, req.Username, req.Password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}

	return TestResult{Success: true, Message: "PostgreSQL 连接成功!"}
}

// testOracle 测试Oracle连通性
func testOracle(req TestRequest) TestResult {
	connStr := fmt.Sprintf("oracle://%s:%s@%s:%d/?service_name=ORCL&ssl=false&sslmode=disable",
		req.Username, req.Password, req.Host, req.Port)

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}

	return TestResult{Success: true, Message: "Oracle 连接成功!"}
}

// testRedis 测试Redis连通性
func testRedis(req TestRequest) TestResult {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", req.Host, req.Port),
		Password: req.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}

	return TestResult{Success: true, Message: "Redis 连接成功!"}
}

// testSQLServer 测试SQL Server连通性
func testSQLServer(req TestRequest) TestResult {
	connStr := fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=master;connection timeout=5",
		req.Host, req.Port, req.Username, req.Password)

	db, err := sql.Open("mssql", connStr)
	if err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
	}

	return TestResult{Success: true, Message: "SQL Server 连接成功!"}
}

// testSQLite 测试SQLite连通性（已移除SQLite支持）
// func testSQLite(req TestRequest) TestResult {
// 	// SQLite是本地文件，host作为文件路径
// 	dbPath := req.Host
// 	if dbPath == "" {
// 		dbPath = "./test.db"
// 	}
// 
// 	db, err := sql.Open("sqlite", dbPath)
// 	if err != nil {
// 		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
// 	}
// 	defer db.Close()
// 
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 
// 	if err := db.PingContext(ctx); err != nil {
// 		return TestResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}
// 	}
// 
// 	return TestResult{Success: true, Message: "SQLite 连接成功!"}
// }
