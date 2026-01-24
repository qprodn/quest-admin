package common

import (
	"context"
	"time"

	bizOrganization "quest-admin/internal/biz/organization"
	bizPermission "quest-admin/internal/biz/permission"
	bizTenant "quest-admin/internal/biz/tenant"
	bizUser "quest-admin/internal/biz/user"

	"github.com/go-kratos/kratos/v2/log"
)

func NewTestContext() context.Context {
	return context.Background()
}

func NewTestLogger() log.Logger {
	return log.DefaultLogger
}

func MockNowTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func CreateTestBizUser() *bizUser.User {
	now := time.Now()
	return &bizUser.User{
		ID:       "test_user_001",
		Username: "testuser",
		Password: "test_password_hashed",
		Nickname: "Test User",
		Email:    "test@example.com",
		Mobile:   "13800138000",
		Sex:      1,
		Avatar:   "",
		Status:   1,
		Remark:   "Test user",
		CreateBy: "system",
		CreateAt: now,
		UpdateBy: "system",
		UpdateAt: now,
		TenantID: "test_tenant",
	}
}

func CreateTestBizRole() *bizPermission.Role {
	now := time.Now()
	return &bizPermission.Role{
		ID:               "test_role_001",
		Name:             "Test Role",
		Code:             "TEST_ROLE",
		Sort:             1,
		DataScope:        1,
		DataScopeDeptIDs: "",
		Status:           1,
		Type:             1,
		Remark:           "Test role",
		CreateBy:         "system",
		CreateAt:         now,
		UpdateBy:         "system",
		UpdateAt:         now,
		TenantID:         "test_tenant",
	}
}

func CreateTestBizDepartment() *bizOrganization.Department {
	now := time.Now()
	return &bizOrganization.Department{
		ID:       "test_dept_001",
		Name:     "Test Department",
		ParentID: "",
		Sort:     1,
		Status:   1,
		CreateBy: "system",
		CreateAt: now,
		UpdateBy: "system",
		UpdateAt: now,
		TenantID: "test_tenant",
	}
}

func CreateTestBizTenant() *bizTenant.Tenant {
	now := time.Now()
	return &bizTenant.Tenant{
		ID:            "test_tenant_001",
		Name:          "Test Company",
		ContactUserID: "",
		ContactName:   "Test Contact",
		ContactMobile: "13800138000",
		Status:        1,
		Website:       "https://test.example.com",
		PackageID:     "",
		ExpireTime:    now.AddDate(1, 0, 0),
		AccountCount:  10,
		CreateBy:      "system",
		CreateAt:      now,
		UpdateBy:      "system",
		UpdateAt:      now,
	}
}

func CreateTestBizMenu() *bizPermission.Menu {
	now := time.Now()
	return &bizPermission.Menu{
		ID:            "test_menu_001",
		Name:          "Test Menu",
		Permission:    "test:permission",
		Type:          1,
		Sort:          1,
		ParentID:      "",
		Path:          "/test",
		Icon:          "test-icon",
		Component:     "TestComponent",
		ComponentName: "test-component",
		Status:        1,
		Visible:       true,
		KeepAlive:     true,
		AlwaysShow:    true,
		CreateBy:      "system",
		CreateAt:      now,
		UpdateBy:      "system",
		UpdateAt:      now,
	}
}
