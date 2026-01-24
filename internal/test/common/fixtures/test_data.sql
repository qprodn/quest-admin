-- Test Users Data
INSERT INTO qa_user (id, username, password, nickname, email, mobile, sex, status, create_at, update_at, tenant_id) VALUES
('test_user_001', 'testuser1', 'hashed_password_1', 'Test User 1', 'test1@example.com', '13800138001', 1, 1, NOW(), NOW(), 'test_tenant_001'),
('test_user_002', 'testuser2', 'hashed_password_2', 'Test User 2', 'test2@example.com', '13800138002', 1, 1, NOW(), NOW(), 'test_tenant_001'),
('test_user_003', 'testuser3', 'hashed_password_3', 'Test User 3', 'test3@example.com', '13800138003', 0, 1, NOW(), NOW(), 'test_tenant_001');

-- Test Roles Data
INSERT INTO qa_role (id, name, code, sort, data_scope, status, type, create_at, update_at, tenant_id) VALUES
('test_role_001', 'Test Role 1', 'TEST_ROLE_1', 1, 1, 1, 1, NOW(), NOW(), 'test_tenant_001'),
('test_role_002', 'Test Role 2', 'TEST_ROLE_2', 2, 1, 1, 1, NOW(), NOW(), 'test_tenant_001'),
('test_role_003', 'Test Admin', 'TEST_ADMIN', 99, 1, 1, 0, NOW(), NOW(), 'test_tenant_001');

-- Test Menus Data
INSERT INTO qa_menu (id, name, permission, type, sort, parent_id, path, icon, component, component_name, status, visible, create_at, update_at) VALUES
('test_menu_001', 'Dashboard', 'dashboard:view', 1, 1, '', '/dashboard', 'dashboard', 'Dashboard', 'Dashboard', 1, true, NOW(), NOW()),
('test_menu_002', 'User Management', 'user:view', 1, 2, '', '/user', 'user', 'UserManagement', 'UserManagement', 1, true, NOW(), NOW()),
('test_menu_003', 'Role Management', 'role:view', 1, 3, '', '/role', 'role', 'RoleManagement', 'RoleManagement', 1, true, NOW(), NOW()),
('test_menu_004', 'User List', 'user:list', 2, 1, 'test_menu_002', '/user/list', 'list', 'UserList', 'UserList', 1, true, NOW(), NOW()),
('test_menu_005', 'User Create', 'user:create', 2, 2, 'test_menu_002', '/user/create', 'plus', 'UserCreate', 'UserCreate', 1, true, NOW(), NOW());

-- Test Departments Data
INSERT INTO qa_dept (id, name, parent_id, sort, status, create_at, update_at, tenant_id) VALUES
('test_dept_001', 'Test Company', '', 1, 1, NOW(), NOW(), 'test_tenant_001'),
('test_dept_002', 'IT Department', 'test_dept_001', 1, 1, NOW(), NOW(), 'test_tenant_001'),
('test_dept_003', 'HR Department', 'test_dept_001', 2, 1, NOW(), NOW(), 'test_tenant_001');

-- Test Posts Data
INSERT INTO qa_post (id, code, name, sort, status, create_at, update_at, tenant_id) VALUES
('test_post_001', 'DEV', 'Developer', 1, 1, NOW(), NOW(), 'test_tenant_001'),
('test_post_002', 'MGR', 'Manager', 2, 1, NOW(), NOW(), 'test_tenant_001'),
('test_post_003', 'CEO', 'Chief Executive Officer', 3, 1, NOW(), NOW(), 'test_tenant_001');

-- Test Tenants Data
INSERT INTO qa_tenant (id, name, contact_name, contact_mobile, status, account_count, create_at, update_at) VALUES
('test_tenant_001', 'Test Company', 'John Doe', '13800138000', 1, 100, NOW(), NOW()),
('test_tenant_002', 'Another Company', 'Jane Smith', '13800138001', 1, 50, NOW(), NOW());

-- Test Tenant Packages Data
INSERT INTO qa_tenant_package (id, name, status, create_at, update_at) VALUES
('test_package_001', 'Basic Package', 1, NOW(), NOW()),
('test_package_002', 'Premium Package', 1, NOW(), NOW()),
('test_package_003', 'Enterprise Package', 1, NOW(), NOW());

-- Test User-Role Associations
INSERT INTO qa_user_role (id, user_id, role_id, create_at, update_at, tenant_id) VALUES
('test_ur_001', 'test_user_001', 'test_role_001', NOW(), NOW(), 'test_tenant_001'),
('test_ur_002', 'test_user_001', 'test_role_002', NOW(), NOW(), 'test_tenant_001'),
('test_ur_003', 'test_user_002', 'test_role_001', NOW(), NOW(), 'test_tenant_001');

-- Test User-Department Associations
INSERT INTO qa_user_dept (id, user_id, dept_id, create_at, update_at, tenant_id) VALUES
('test_ud_001', 'test_user_001', 'test_dept_002', NOW(), NOW(), 'test_tenant_001'),
('test_ud_002', 'test_user_002', 'test_dept_003', NOW(), NOW(), 'test_tenant_001');

-- Test User-Post Associations
INSERT INTO qa_user_post (id, user_id, post_id, create_at, update_at, tenant_id) VALUES
('test_up_001', 'test_user_001', 'test_post_001', NOW(), NOW(), 'test_tenant_001'),
('test_up_002', 'test_user_002', 'test_post_002', NOW(), NOW(), 'test_tenant_001');

-- Test Role-Menu Associations
INSERT INTO qa_role_menu (id, role_id, menu_id, create_at, update_at, tenant_id) VALUES
('test_rm_001', 'test_role_001', 'test_menu_001', NOW(), NOW(), 'test_tenant_001'),
('test_rm_002', 'test_role_001', 'test_menu_002', NOW(), NOW(), 'test_tenant_001'),
('test_rm_003', 'test_role_002', 'test_menu_003', NOW(), NOW(), 'test_tenant_001');
