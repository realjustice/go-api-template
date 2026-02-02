package database

import (
	"context"

	"go-api-template/pkg/errors"

	"gorm.io/gorm"
)

// BaseRepository 基础 Repository，提供通用的 CRUD 操作
// 其他 Repository 可以嵌入此结构体，复用基础方法
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础 Repository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// DB 获取数据库连接（用于复杂查询）
func (r *BaseRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// ========== 查询操作 ==========

// FindByID 根据 ID 查询单条记录
func (r *BaseRepository) FindByID(ctx context.Context, id interface{}, dest interface{}) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).First(dest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, "query by id failed")
	}
	return nil
}

// FindOne 根据条件查询单条记录
func (r *BaseRepository) FindOne(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error {
	err := r.db.WithContext(ctx).Where(query, args...).First(dest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, "query one failed")
	}
	return nil
}

// FindAll 查询所有记录
func (r *BaseRepository) FindAll(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error {
	err := r.db.WithContext(ctx).Where(query, args...).Find(dest).Error
	if err != nil {
		return errors.Wrap(err, "query all failed")
	}
	return nil
}

// FindPage 分页查询
func (r *BaseRepository) FindPage(ctx context.Context, dest interface{}, page, pageSize int, query interface{}, args ...interface{}) (int64, error) {
	var total int64

	db := r.db.WithContext(ctx).Model(dest)
	if query != nil {
		db = db.Where(query, args...)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return 0, errors.Wrap(err, "count failed")
	}

	// 查询分页数据
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(dest).Error
	if err != nil {
		return 0, errors.Wrap(err, "query page failed")
	}

	return total, nil
}

// Count 统计数量
func (r *BaseRepository) Count(ctx context.Context, model interface{}, query interface{}, args ...interface{}) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(model)
	if query != nil {
		db = db.Where(query, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, "count failed")
	}
	return count, nil
}

// Exists 判断记录是否存在
func (r *BaseRepository) Exists(ctx context.Context, model interface{}, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(model).Where(query, args...).Limit(1).Count(&count).Error
	if err != nil {
		return false, errors.Wrap(err, "check exists failed")
	}
	return count > 0, nil
}

// ========== 创建操作 ==========

// Create 创建记录
func (r *BaseRepository) Create(ctx context.Context, value interface{}) error {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return errors.Wrap(err, "create failed")
	}
	return nil
}

// CreateInBatches 批量创建
func (r *BaseRepository) CreateInBatches(ctx context.Context, value interface{}, batchSize int) error {
	err := r.db.WithContext(ctx).CreateInBatches(value, batchSize).Error
	if err != nil {
		return errors.Wrap(err, "create in batches failed")
	}
	return nil
}

// ========== 更新操作 ==========

// Update 更新记录（全部字段）
func (r *BaseRepository) Update(ctx context.Context, value interface{}) error {
	err := r.db.WithContext(ctx).Save(value).Error
	if err != nil {
		return errors.Wrap(err, "update failed")
	}
	return nil
}

// UpdateFields 更新指定字段
func (r *BaseRepository) UpdateFields(ctx context.Context, model interface{}, query interface{}, updates map[string]interface{}, args ...interface{}) error {
	err := r.db.WithContext(ctx).Model(model).Where(query, args...).Updates(updates).Error
	if err != nil {
		return errors.Wrap(err, "update fields failed")
	}
	return nil
}

// UpdateColumn 更新单个字段（不触发钩子）
func (r *BaseRepository) UpdateColumn(ctx context.Context, model interface{}, query interface{}, column string, value interface{}, args ...interface{}) error {
	err := r.db.WithContext(ctx).Model(model).Where(query, args...).Update(column, value).Error
	if err != nil {
		return errors.Wrap(err, "update column failed")
	}
	return nil
}

// ========== 删除操作 ==========

// Delete 删除记录
func (r *BaseRepository) Delete(ctx context.Context, model interface{}, id interface{}) error {
	err := r.db.WithContext(ctx).Delete(model, id).Error
	if err != nil {
		return errors.Wrap(err, "delete failed")
	}
	return nil
}

// DeleteWhere 根据条件删除
func (r *BaseRepository) DeleteWhere(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error {
	err := r.db.WithContext(ctx).Where(query, args...).Delete(model).Error
	if err != nil {
		return errors.Wrap(err, "delete where failed")
	}
	return nil
}

// ========== 事务操作 ==========

// Transaction 执行事务
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// ========== 原生 SQL ==========

// Exec 执行原生 SQL
func (r *BaseRepository) Exec(ctx context.Context, sql string, values ...interface{}) error {
	err := r.db.WithContext(ctx).Exec(sql, values...).Error
	if err != nil {
		return errors.Wrap(err, "exec sql failed")
	}
	return nil
}

// Raw 执行原生查询
func (r *BaseRepository) Raw(ctx context.Context, dest interface{}, sql string, values ...interface{}) error {
	err := r.db.WithContext(ctx).Raw(sql, values...).Scan(dest).Error
	if err != nil {
		return errors.Wrap(err, "raw query failed")
	}
	return nil
}
