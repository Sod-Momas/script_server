package models

import (
	"fmt"
	"sync"
	"time"
)

// 内存存储的全局变量
var (
	datasourceMap = make(map[int64]*Datasource)
	currentID     int64 = 0
	mutex         sync.RWMutex
)

// Datasource 数据源模型
type Datasource struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	DBType    string    `json:"db_type"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// Save 保存数据源到内存
func (d *Datasource) Save() error {
	mutex.Lock()
	defer mutex.Unlock()

	// 生成新ID
	currentID++
	d.ID = currentID
	d.CreatedAt = time.Now()

	// 存储到内存map
	datasourceMap[d.ID] = d

	return nil
}

// GetAll 获取所有数据源
func GetAll() ([]Datasource, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	var datasources []Datasource
	for _, ds := range datasourceMap {
		datasources = append(datasources, *ds)
	}

	return datasources, nil
}

// GetByID 根据ID获取单个数据源
func GetByID(id int64) (*Datasource, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	ds, exists := datasourceMap[id]
	if !exists {
		return nil, fmt.Errorf("数据源不存在")
	}

	return ds, nil
}
