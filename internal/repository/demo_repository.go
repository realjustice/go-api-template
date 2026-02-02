package repository

import (
	"context"

	"go-api-template/internal/model"
	"go-api-template/pkg/database"
	"go-api-template/pkg/errors"

	"gorm.io/gorm"
)

// DemoRepository Demo 数据访问层
type DemoRepository struct {
	*database.BaseRepository // 嵌入 BaseRepository，复用基础方法
	db                       *gorm.DB
}

// NewDemoRepository 创建 Demo Repository
func NewDemoRepository(db *gorm.DB) *DemoRepository {
	return &DemoRepository{
		BaseRepository: database.NewBaseRepository(db),
		db:             db,
	}
}

// ========== 使用 BaseRepository 的通用方法 ==========

// FindByID 根据 ID 查询（使用基类方法）
func (r *DemoRepository) FindByID(ctx context.Context, id uint) (*model.Demo, error) {
	var demo model.Demo
	err := r.BaseRepository.FindByID(ctx, id, &demo)
	if err != nil {
		return nil, errors.Wrapf(err, "demo not found, id: %d", id)
	}
	return &demo, nil
}

// FindAll 查询所有（使用基类方法）
func (r *DemoRepository) FindAll(ctx context.Context) ([]*model.Demo, error) {
	var demos []*model.Demo
	err := r.BaseRepository.FindAll(ctx, &demos, "1 = 1") // 查询所有
	if err != nil {
		return nil, err
	}
	return demos, nil
}

// Create 创建（使用基类方法）
func (r *DemoRepository) Create(ctx context.Context, demo *model.Demo) error {
	return r.BaseRepository.Create(ctx, demo)
}

// Update 更新（使用基类方法）
func (r *DemoRepository) Update(ctx context.Context, demo *model.Demo) error {
	return r.BaseRepository.Update(ctx, demo)
}

// Delete 删除（使用基类方法）
func (r *DemoRepository) Delete(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, &model.Demo{}, id)
}

// ========== 业务特定的复杂查询（直接使用 GORM）==========

// FindByStatus 根据状态查询（复杂查询示例）
func (r *DemoRepository) FindByStatus(ctx context.Context, status int) ([]*model.Demo, error) {
	var demos []*model.Demo
	// 直接使用 GORM，保留灵活性
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&demos).Error
	if err != nil {
		return nil, errors.Wrap(err, "query by status failed")
	}
	return demos, nil
}

// FindPage 分页查询（使用基类方法）
func (r *DemoRepository) FindPage(ctx context.Context, page, pageSize int) ([]*model.Demo, int64, error) {
	var demos []*model.Demo
	total, err := r.BaseRepository.FindPage(ctx, &demos, page, pageSize, "1 = 1")
	if err != nil {
		return nil, 0, err
	}
	return demos, total, nil
}

// UpdateStatus 更新状态（使用基类方法）
func (r *DemoRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.BaseRepository.UpdateColumn(ctx, &model.Demo{}, "id = ?", "status", status, id)
}

// CountByStatus 统计指定状态的数量（使用基类方法）
func (r *DemoRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	return r.BaseRepository.Count(ctx, &model.Demo{}, "status = ?", status)
}

// ExistsByTitle 检查标题是否存在（使用基类方法）
func (r *DemoRepository) ExistsByTitle(ctx context.Context, title string) (bool, error) {
	return r.BaseRepository.Exists(ctx, &model.Demo{}, "title = ?", title)
}

// ========== 高级查询（直接使用 GORM，展示灵活性）==========

// Search 搜索（支持多条件）
func (r *DemoRepository) Search(ctx context.Context, keyword string, status *int, page, pageSize int) ([]*model.Demo, int64, error) {
	var demos []*model.Demo
	var total int64

	// 构建查询（直接使用 GORM 的链式调用）
	query := r.db.WithContext(ctx).Model(&model.Demo{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "count search results failed")
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&demos).Error
	if err != nil {
		return nil, 0, errors.Wrap(err, "search failed")
	}

	return demos, total, nil
}

// BatchUpdateStatus 批量更新状态（直接使用 GORM）
func (r *DemoRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	err := r.db.WithContext(ctx).
		Model(&model.Demo{}).
		Where("id IN ?", ids).
		Update("status", status).Error
	if err != nil {
		return errors.Wrap(err, "batch update status failed")
	}
	return nil
}

// ========== 事务支持 ==========

// CreateWithTx 在事务中创建（供 Service 层使用）
func (r *DemoRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, demo *model.Demo) error {
	err := tx.WithContext(ctx).Create(demo).Error
	if err != nil {
		return errors.Wrap(err, "create with tx failed")
	}
	return nil
}

// UpdateWithTx 在事务中更新（供 Service 层使用）
func (r *DemoRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, demo *model.Demo) error {
	err := tx.WithContext(ctx).Save(demo).Error
	if err != nil {
		return errors.Wrap(err, "update with tx failed")
	}
	return nil
}
