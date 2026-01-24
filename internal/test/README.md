# 测试目录结构说明

本目录包含项目所有业务代码的测试文件，完全按照 internal 目录下的业务代码层级进行组织。

## 目录结构

```
internal/test/
├── common/                        # 公共测试工具
│   ├── context.go                 # 测试上下文管理
│   ├── mock.go                    # Mock 对象定义
│   └── fixtures/                  # 测试数据
│       ├── test_data.sql            # SQL 测试数据
│       └── test_data.json          # JSON 测试数据
│
├── data/                          # Data 层测试
│   ├── idgen/
│   │   └── sonyflake_test.go
│   ├── user/
│   │   ├── user_repo_test.go
│   │   ├── user_role_repo_test.go
│   │   ├── user_dept_repo_test.go
│   │   └── user_post_repo_test.go
│   ├── permission/
│   │   ├── role_repo_test.go
│   │   ├── menu_repo_test.go
│   │   └── role_menu_repo_test.go
│   ├── organization/
│   │   ├── department_repo_test.go
│   │   └── post_repo_test.go
│   └── tenant/
│       ├── tenant_repo_test.go
│       └── package_repo_test.go
│
├── biz/                           # Biz 层测试
│   ├── user/
│   │   ├── user_biz_test.go
│   │   ├── user_role_biziz_test.go
│   │   ├── user_dept_biz_test.go
│   │   └── user_post_biz_test.go
│   ├── permission/
│   │   ├── role_biz_test.go
│   │   └── menu_biz_test.go
│   ├── organization/
│   │   ├── department_biz_test.go
│   │   └── post_biz_test.go
│   ├── tenant/
│   │   ├── tenant_biz_test.go
│   │   └── package_biz_test.go
│   └── auth/
│       └── auth_biz_test.go.go
│
└── service/                       # Service 层测试
    ├── user/
    │   ├── user_service_test.go
    │   ├── user_role_service_test.go
    │   ├── user_dept_service_test.go
    │   └── user_post_service_test.go
    ├── permission/
    │   ├── role_service_test.go
    │   └── menu_service_test.go
    ├── organization/
    │   ├── department_service_test.go
    │   └── post_service_test.go
    ├── tenant/
    │   ├── tenant_service_test.go
    │   └── package_service_test.go
    └── auth/
        └── auth_service_test.go
```

## 运行测试

### 运行所有测试
```bash
make test
```

### 按层级运行测试
```bash
# 只测试 data 层
make test-data

# 只测试 biz 层
make test-biz

# 只测试 service 层
make test-service
```

### 生成测试覆盖率报告
```bash
make test-coverage       # 生成文本覆盖率报告
make test-coverage-html  # 生成 HTML 覆盖率报告
```

### 清理测试产物
```bash
make test-clean
```

## 测试工具说明

### TestContext
提供测试上下文管理，包含：
- Ctx: context.Context - 测试上下文
- Logger: 日志记录器
- cleanupFuncs: 清理函数列表

### Mock 对象
提供各种业务接口的 Mock 实现：
- MockUserRepo
- MockRoleRepo
- MockDepartmentRepo
- MockTenantRepo
- MockIDGenerator
- MockTransactionManager

### 测试数据
- test_data.sql: PostgreSQL 初始化数据
- test_data.json: JSON 格式测试数据

## 注意事项

1. 所有测试使用 `github.com/stretchr/testify` 进行断言
2. 测试文件使用 `_test.go` 后缀
3. 每个测试函数以 `Test` 开头
4. 使用 table-driven tests 模式进行多场景测试
5. 使用 `t.Run()` 组织子测试

## 测试策略

- **Data 层**: 测试数据库操作，接口方法验证
- **Biz 层**: 使用 Mock 的 data 层，测试业务逻辑
- **Service 层**: 使用 Mock 的 biz 层，测试 gRPC 服务方法
