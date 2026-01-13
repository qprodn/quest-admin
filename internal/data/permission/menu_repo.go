package permission

import (
	"context"
	"database/sql"
	"time"

	biz "quest-admin/internal/biz/permission"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/idgen"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type Menu struct {
	bun.BaseModel `bun:"table:qa_menu,alias:m"`

	ID            string     `bun:"id,pk"`
	Name          string     `bun:"name,notnull"`
	Permission    string     `bun:"permission,default:''"`
	Type          int32      `bun:"type,notnull"`
	Sort          int32      `bun:"sort,default:0"`
	ParentID      string     `bun:"parent_id,default:'0'"`
	Path          string     `bun:"path,default:''"`
	Icon          string     `bun:"icon,default:'#'"`
	Component     string     `bun:"component"`
	ComponentName string     `bun:"component_name"`
	Status        int32      `bun:"status,default:0"`
	Visible       bool       `bun:"visible,default:true"`
	KeepAlive     bool       `bun:"keep_alive,default:true"`
	AlwaysShow    bool       `bun:"always_show,default:true"`
	CreateBy      string     `bun:"create_by"`
	CreateAt      time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy      string     `bun:"update_by"`
	UpdateAt      time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	DeleteAt      *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type menuRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewMenuRepo(data *data.Data, logger log.Logger) biz.MenuRepo {
	return &menuRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *menuRepo) Create(ctx context.Context, menu *biz.Menu) error {
	now := time.Now()
	dbMenu := &Menu{
		ID:            idgen.GenerateID(),
		Name:          menu.Name,
		Permission:    menu.Permission,
		Type:          menu.Type,
		Sort:          menu.Sort,
		ParentID:      menu.ParentID,
		Path:          menu.Path,
		Icon:          menu.Icon,
		Component:     menu.Component,
		ComponentName: menu.ComponentName,
		Status:        menu.Status,
		Visible:       menu.Visible,
		KeepAlive:     menu.KeepAlive,
		AlwaysShow:    menu.AlwaysShow,
		CreateBy:      "todo",
		CreateAt:      now,
		UpdateBy:      "todo",
		UpdateAt:      now,
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbMenu).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *menuRepo) FindByID(ctx context.Context, id string) (*biz.Menu, error) {
	dbMenu := &Menu{ID: id}
	err := r.data.DB(ctx).NewSelect().Model(dbMenu).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizMenu(dbMenu), nil
}

func (r *menuRepo) FindByName(ctx context.Context, name string) (*biz.Menu, error) {
	dbMenu := &Menu{}
	err := r.data.DB(ctx).NewSelect().Model(dbMenu).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizMenu(dbMenu), nil
}

func (r *menuRepo) List(ctx context.Context) ([]*biz.Menu, error) {
	var dbMenus []*Menu
	err := r.data.DB(ctx).NewSelect().
		Model(&dbMenus).
		Where("status = ?", 1).
		Order("sort ASC, create_at DESC").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*biz.Menu{}, nil
		}
		return nil, err
	}

	menus := make([]*biz.Menu, 0, len(dbMenus))
	for _, dbMenu := range dbMenus {
		menus = append(menus, r.toBizMenu(dbMenu))
	}
	return menus, nil
}

func (r *menuRepo) FindByParentID(ctx context.Context, parentID string) ([]*biz.Menu, error) {
	var dbMenus []*Menu
	err := r.data.DB(ctx).NewSelect().
		Model(&dbMenus).
		Where("parent_id = ?", parentID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*biz.Menu{}, nil
		}
		return nil, err
	}

	menus := make([]*biz.Menu, 0, len(dbMenus))
	for _, dbMenu := range dbMenus {
		menus = append(menus, r.toBizMenu(dbMenu))
	}
	return menus, nil
}

func (r *menuRepo) Update(ctx context.Context, menu *biz.Menu) error {
	dbMenu := &Menu{
		ID:            menu.ID,
		Name:          menu.Name,
		Permission:    menu.Permission,
		Type:          menu.Type,
		Sort:          menu.Sort,
		ParentID:      menu.ParentID,
		Path:          menu.Path,
		Icon:          menu.Icon,
		Component:     menu.Component,
		ComponentName: menu.ComponentName,
		Status:        menu.Status,
		Visible:       menu.Visible,
		KeepAlive:     menu.KeepAlive,
		AlwaysShow:    menu.AlwaysShow,
		UpdateBy:      "todo",
		UpdateAt:      time.Now(),
	}

	_, err := r.data.DB(ctx).NewUpdate().Model(dbMenu).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *menuRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*Menu)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *menuRepo) toBizMenu(dbMenu *Menu) *biz.Menu {
	return &biz.Menu{
		ID:            dbMenu.ID,
		Name:          dbMenu.Name,
		Permission:    dbMenu.Permission,
		Type:          dbMenu.Type,
		Sort:          dbMenu.Sort,
		ParentID:      dbMenu.ParentID,
		Path:          dbMenu.Path,
		Icon:          dbMenu.Icon,
		Component:     dbMenu.Component,
		ComponentName: dbMenu.ComponentName,
		Status:        dbMenu.Status,
		Visible:       dbMenu.Visible,
		KeepAlive:     dbMenu.KeepAlive,
		AlwaysShow:    dbMenu.AlwaysShow,
		CreateBy:      dbMenu.CreateBy,
		CreateAt:      dbMenu.CreateAt,
		UpdateBy:      dbMenu.UpdateBy,
		UpdateAt:      dbMenu.UpdateAt,
		Children:      []*biz.Menu{},
	}
}
