
DROP TABLE IF EXISTS qa_user CASCADE;
CREATE TABLE qa_user
(
    id         varchar(32) PRIMARY KEY,
    username   varchar(32)                            NOT NULL,
    password   varchar(128) DEFAULT ''                NOT NULL,
    nickname   varchar(32)                            NOT NULL,
    remark     varchar(512),
    email      varchar(64)  DEFAULT '',
    mobile     varchar(16)  DEFAULT '',
    sex        smallint     DEFAULT 0,
    avatar     varchar(512) DEFAULT '',
    status     smallint     DEFAULT 0                 NOT NULL,
    login_ip   varchar(64)  DEFAULT '',
    login_date timestamp,
    create_by  varchar(64)  DEFAULT '',
    create_at  timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by  varchar(64)  DEFAULT '',
    update_at  timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at  timestamp,
    tenant_id  varchar(32)  DEFAULT ''                NOT NULL
);

COMMENT ON TABLE qa_user IS '用户信息表';
COMMENT ON COLUMN qa_user.id IS '用户ID';
COMMENT ON COLUMN qa_user.username IS '用户账号';
COMMENT ON COLUMN qa_user.password IS '密码';
COMMENT ON COLUMN qa_user.nickname IS '用户昵称';
COMMENT ON COLUMN qa_user.remark IS '备注';
COMMENT ON COLUMN qa_user.email IS '用户邮箱';
COMMENT ON COLUMN qa_user.mobile IS '手机号码';
COMMENT ON COLUMN qa_user.sex IS '用户性别';
COMMENT ON COLUMN qa_user.avatar IS '头像地址';
COMMENT ON COLUMN qa_user.status IS '帐号状态（0停用 1正常）';
COMMENT ON COLUMN qa_user.login_ip IS '最后登录IP';
COMMENT ON COLUMN qa_user.login_date IS '最后登录时间';
COMMENT ON COLUMN qa_user.create_by IS '创建者';
COMMENT ON COLUMN qa_user.create_at IS '创建时间';
COMMENT ON COLUMN qa_user.update_by IS '更新者';
COMMENT ON COLUMN qa_user.update_at IS '更新时间';
COMMENT ON COLUMN qa_user.delete_at IS '删除时间';
COMMENT ON COLUMN qa_user.tenant_id IS '租户编号';

DROP INDEX idx_username;
CREATE UNIQUE INDEX idx_username ON qa_user (username, update_at, tenant_id);

DROP TABLE IF EXISTS qa_role CASCADE;
CREATE TABLE qa_role
(
    id                  varchar(32) PRIMARY KEY,
    name                varchar(32)                            NOT NULL,
    code                varchar(128)                           NOT NULL,
    sort                int                                    NOT NULL,
    data_scope          smallint     DEFAULT 1                 NOT NULL,
    data_scope_dept_ids varchar(512) DEFAULT ''                NOT NULL,
    status              smallint                               NOT NULL,
    type                smallint                               NOT NULL,
    remark              varchar(512),
    create_by           varchar(64)  DEFAULT '',
    create_at           timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by           varchar(64)  DEFAULT '',
    update_at           timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at           boolean,
    tenant_id           varchar(32)  DEFAULT '0'               NOT NULL
);

COMMENT ON TABLE qa_role IS '角色信息表';
COMMENT ON COLUMN qa_role.id IS '角色ID';
COMMENT ON COLUMN qa_role.name IS '角色名称';
COMMENT ON COLUMN qa_role.code IS '角色权限字符串';
COMMENT ON COLUMN qa_role.sort IS '显示顺序';
COMMENT ON COLUMN qa_role.data_scope IS '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）';
COMMENT ON COLUMN qa_role.data_scope_dept_ids IS '数据范围(指定部门数组)';
COMMENT ON COLUMN qa_role.status IS '角色状态（0停用 1正常）';
COMMENT ON COLUMN qa_role.type IS '角色类型';
COMMENT ON COLUMN qa_role.remark IS '备注';
COMMENT ON COLUMN qa_role.create_by IS '创建者';
COMMENT ON COLUMN qa_role.create_at IS '创建时间';
COMMENT ON COLUMN qa_role.update_by IS '更新者';
COMMENT ON COLUMN qa_role.update_at IS '更新时间';
COMMENT ON COLUMN qa_role.delete_at IS '删除时间';
COMMENT ON COLUMN qa_role.tenant_id IS '租户编号';

-- ----------------------------
-- Table structure for qa_role_menu
-- ----------------------------
DROP TABLE IF EXISTS qa_role_menu CASCADE;
CREATE TABLE qa_role_menu
(
    id        varchar(32) PRIMARY KEY,
    role_id   varchar(32)                           NOT NULL,
    menu_id   varchar(32)                           NOT NULL,
    create_by varchar(64) DEFAULT '',
    create_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by varchar(64) DEFAULT '',
    update_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at timestamp,
    tenant_id varchar(32) DEFAULT ''                NOT NULL
);

COMMENT ON TABLE qa_role_menu IS '角色和菜单关联表';
COMMENT ON COLUMN qa_role_menu.id IS '自增编号';
COMMENT ON COLUMN qa_role_menu.role_id IS '角色ID';
COMMENT ON COLUMN qa_role_menu.menu_id IS '菜单ID';
COMMENT ON COLUMN qa_role_menu.create_by IS '创建者';
COMMENT ON COLUMN qa_role_menu.create_at IS '创建时间';
COMMENT ON COLUMN qa_role_menu.update_by IS '更新者';
COMMENT ON COLUMN qa_role_menu.update_at IS '更新时间';
COMMENT ON COLUMN qa_role_menu.delete_at IS '删除时间';
COMMENT ON COLUMN qa_role_menu.tenant_id IS '租户编号';

DROP TABLE IF EXISTS qa_menu CASCADE;
CREATE TABLE qa_menu
(
    id             varchar(32) PRIMARY KEY,
    name           varchar(64)                            NOT NULL,
    permission     varchar(128) DEFAULT ''                NOT NULL,
    type           smallint                               NOT NULL,
    sort           int          DEFAULT 0                 NOT NULL,
    parent_id      varchar(32)  DEFAULT '0'               NOT NULL,
    path           varchar(256) DEFAULT '',
    icon           varchar(128) DEFAULT '#',
    component      varchar(256),
    component_name varchar(256),
    status         smallint     DEFAULT 0                 NOT NULL,
    visible        boolean      DEFAULT TRUE              NOT NULL,
    keep_alive     boolean      DEFAULT TRUE              NOT NULL,
    always_show    boolean      DEFAULT TRUE              NOT NULL,
    create_by      varchar(64)  DEFAULT '',
    create_at      timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by      varchar(64)  DEFAULT '',
    update_at      timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at      timestamp
);

COMMENT ON TABLE qa_menu IS '菜单权限表';
COMMENT ON COLUMN qa_menu.id IS '菜单ID';
COMMENT ON COLUMN qa_menu.name IS '菜单名称';
COMMENT ON COLUMN qa_menu.permission IS '权限标识';
COMMENT ON COLUMN qa_menu.type IS '菜单类型';
COMMENT ON COLUMN qa_menu.sort IS '显示顺序';
COMMENT ON COLUMN qa_menu.parent_id IS '父菜单ID';
COMMENT ON COLUMN qa_menu.path IS '路由地址';
COMMENT ON COLUMN qa_menu.icon IS '菜单图标';
COMMENT ON COLUMN qa_menu.component IS '组件路径';
COMMENT ON COLUMN qa_menu.component_name IS '组件名';
COMMENT ON COLUMN qa_menu.status IS '菜单状态';
COMMENT ON COLUMN qa_menu.visible IS '是否可见';
COMMENT ON COLUMN qa_menu.keep_alive IS '是否缓存';
COMMENT ON COLUMN qa_menu.always_show IS '是否总是显示';
COMMENT ON COLUMN qa_menu.create_by IS '创建者';
COMMENT ON COLUMN qa_menu.create_at IS '创建时间';
COMMENT ON COLUMN qa_menu.update_by IS '更新者';
COMMENT ON COLUMN qa_menu.update_at IS '更新时间';
COMMENT ON COLUMN qa_menu.delete_at IS '删除时间';

DROP TABLE IF EXISTS qa_user_role CASCADE;
CREATE TABLE qa_user_role
(
    id        varchar(32) PRIMARY KEY,
    user_id   varchar(32)            NOT NULL,
    role_id   varchar(32)            NOT NULL,
    create_by varchar(64) DEFAULT '',
    create_at timestamp   DEFAULT CURRENT_TIMESTAMP,
    update_by varchar(64) DEFAULT '',
    update_at timestamp   DEFAULT CURRENT_TIMESTAMP,
    delete_at timestamp,
    tenant_id varchar(32) DEFAULT '' NOT NULL
);

COMMENT ON TABLE qa_user_role IS '用户和角色关联表';
COMMENT ON COLUMN qa_user_role.id IS '自增编号';
COMMENT ON COLUMN qa_user_role.user_id IS '用户ID';
COMMENT ON COLUMN qa_user_role.role_id IS '角色ID';
COMMENT ON COLUMN qa_user_role.create_by IS '创建者';
COMMENT ON COLUMN qa_user_role.create_at IS '创建时间';
COMMENT ON COLUMN qa_user_role.update_by IS '更新者';
COMMENT ON COLUMN qa_user_role.update_at IS '更新时间';
COMMENT ON COLUMN qa_user_role.delete_at IS '删除时间';
COMMENT ON COLUMN qa_user_role.tenant_id IS '租户编号';

DROP TABLE IF EXISTS qa_post CASCADE;
CREATE TABLE qa_post
(
    id        varchar(32) PRIMARY KEY,
    code      varchar(64)                           NOT NULL,
    name      varchar(64)                           NOT NULL,
    sort      int                                   NOT NULL,
    status    smallint                              NOT NULL,
    remark    varchar(512),
    create_by varchar(64) DEFAULT '',
    create_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by varchar(64) DEFAULT '',
    update_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at timestamp,
    tenant_id varchar(32) DEFAULT ''                NOT NULL
);

COMMENT ON TABLE qa_post IS '岗位信息表';
COMMENT ON COLUMN qa_post.id IS '岗位ID';
COMMENT ON COLUMN qa_post.code IS '岗位编码';
COMMENT ON COLUMN qa_post.name IS '岗位名称';
COMMENT ON COLUMN qa_post.sort IS '显示顺序';
COMMENT ON COLUMN qa_post.status IS '状态（0停用 1正常）';
COMMENT ON COLUMN qa_post.remark IS '备注';
COMMENT ON COLUMN qa_post.create_by IS '创建者';
COMMENT ON COLUMN qa_post.create_at IS '创建时间';
COMMENT ON COLUMN qa_post.update_by IS '更新者';
COMMENT ON COLUMN qa_post.update_at IS '更新时间';
COMMENT ON COLUMN qa_post.delete_at IS '删除时间';
COMMENT ON COLUMN qa_post.tenant_id IS '租户编号';

DROP TABLE IF EXISTS qa_user_post CASCADE;
CREATE TABLE qa_user_post
(
    id        varchar(32) PRIMARY KEY,
    user_id   varchar(32)                           NOT NULL,
    post_id   varchar(32)                           NOT NULL,
    create_by varchar(64) DEFAULT '',
    create_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by varchar(64) DEFAULT '',
    update_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at timestamp,
    tenant_id varchar(32) DEFAULT ''                NOT NULL
);

COMMENT ON TABLE qa_user_post IS '用户岗位表';
COMMENT ON COLUMN qa_user_post.id IS 'id';
COMMENT ON COLUMN qa_user_post.user_id IS '用户ID';
COMMENT ON COLUMN qa_user_post.post_id IS '岗位ID';
COMMENT ON COLUMN qa_user_post.create_by IS '创建者';
COMMENT ON COLUMN qa_user_post.create_at IS '创建时间';
COMMENT ON COLUMN qa_user_post.update_by IS '更新者';
COMMENT ON COLUMN qa_user_post.update_at IS '更新时间';
COMMENT ON COLUMN qa_user_post.delete_at IS '删除时间';
COMMENT ON COLUMN qa_user_post.tenant_id IS '租户编号';

DROP TABLE IF EXISTS qa_user_dept CASCADE;
CREATE TABLE qa_user_dept
(
    id        varchar(32) PRIMARY KEY,
    user_id   varchar(32)                           NOT NULL,
    dept_id   varchar(32)                           NOT NULL,
    create_by varchar(64) DEFAULT '',
    create_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by varchar(64) DEFAULT '',
    update_at timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at timestamp,
    tenant_id varchar(32) DEFAULT ''                NOT NULL
);

COMMENT ON TABLE qa_user_dept IS '用户部门表';
COMMENT ON COLUMN qa_user_dept.id IS 'id';
COMMENT ON COLUMN qa_user_dept.user_id IS '用户ID';
COMMENT ON COLUMN qa_user_dept.dept_id IS '部门ID';
COMMENT ON COLUMN qa_user_dept.create_by IS '创建者';
COMMENT ON COLUMN qa_user_dept.create_at IS '创建时间';
COMMENT ON COLUMN qa_user_dept.update_by IS '更新者';
COMMENT ON COLUMN qa_user_dept.update_at IS '更新时间';
COMMENT ON COLUMN qa_user_dept.delete_at IS '删除时间';
COMMENT ON COLUMN qa_user_dept.tenant_id IS '租户编号';

DROP TABLE IF EXISTS qa_dept CASCADE;
CREATE TABLE qa_dept
(
    id             varchar(32) PRIMARY KEY,
    name           varchar(32) DEFAULT ''                NOT NULL,
    parent_id      varchar(32)                           NOT NULL,
    sort           int         DEFAULT 0                 NOT NULL,
    leader_user_id varchar(32),
    phone          varchar(16),
    email          varchar(64),
    status         smallint                              NOT NULL,
    create_by      varchar(64) DEFAULT '',
    create_at      timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by      varchar(64) DEFAULT '',
    update_at      timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at      timestamp,
    tenant_id      varchar(32) DEFAULT ''                NOT NULL
);

COMMENT ON TABLE qa_dept IS '部门表';
COMMENT ON COLUMN qa_dept.id IS '部门id';
COMMENT ON COLUMN qa_dept.name IS '部门名称';
COMMENT ON COLUMN qa_dept.parent_id IS '父部门id';
COMMENT ON COLUMN qa_dept.sort IS '显示顺序';
COMMENT ON COLUMN qa_dept.leader_user_id IS '负责人';
COMMENT ON COLUMN qa_dept.phone IS '联系电话';
COMMENT ON COLUMN qa_dept.email IS '邮箱';
COMMENT ON COLUMN qa_dept.status IS '部门状态（0停用 1正常）';
COMMENT ON COLUMN qa_dept.create_by IS '创建者';
COMMENT ON COLUMN qa_dept.create_at IS '创建时间';
COMMENT ON COLUMN qa_dept.update_by IS '更新者';
COMMENT ON COLUMN qa_dept.update_at IS '更新时间';
COMMENT ON COLUMN qa_dept.delete_at IS '删除时间';
COMMENT ON COLUMN qa_dept.tenant_id IS '租户编号';

DROP TABLE IF EXISTS qa_tenant CASCADE;
CREATE TABLE qa_tenant
(
    id              varchar(32) PRIMARY KEY,
    name            varchar(32)                            NOT NULL,
    contact_user_id varchar(32),
    contact_name    varchar(32)                            NOT NULL,
    contact_mobile  varchar(512),
    status          smallint     DEFAULT 1                 NOT NULL,
    website         varchar(256) DEFAULT '',
    package_id      varchar(32)                            NOT NULL,
    expire_time     timestamp                              NOT NULL,
    account_count   int                                    NOT NULL,
    create_by       varchar(64)  DEFAULT ''                NOT NULL,
    create_at       timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by       varchar(64)  DEFAULT '',
    update_at       timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at       timestamp
);

COMMENT ON TABLE qa_tenant IS '租户表';
COMMENT ON COLUMN qa_tenant.id IS '租户编号';
COMMENT ON COLUMN qa_tenant.name IS '租户名';
COMMENT ON COLUMN qa_tenant.contact_user_id IS '联系人的用户编号';
COMMENT ON COLUMN qa_tenant.contact_name IS '联系人';
COMMENT ON COLUMN qa_tenant.contact_mobile IS '联系手机';
COMMENT ON COLUMN qa_tenant.status IS '租户状态（0停用 1正常）';
COMMENT ON COLUMN qa_tenant.website IS '绑定域名';
COMMENT ON COLUMN qa_tenant.package_id IS '租户套餐编号';
COMMENT ON COLUMN qa_tenant.expire_time IS '过期时间';
COMMENT ON COLUMN qa_tenant.account_count IS '账号数量';
COMMENT ON COLUMN qa_tenant.create_by IS '创建者';
COMMENT ON COLUMN qa_tenant.create_at IS '创建时间';
COMMENT ON COLUMN qa_tenant.update_by IS '更新者';
COMMENT ON COLUMN qa_tenant.update_at IS '更新时间';
COMMENT ON COLUMN qa_tenant.delete_at IS '删除时间';

DROP TABLE IF EXISTS qa_tenant_package CASCADE;
CREATE TABLE qa_tenant_package
(
    id        varchar(32) PRIMARY KEY,
    name      varchar(32)                            NOT NULL,
    status    smallint     DEFAULT 1                 NOT NULL,
    remark    varchar(256) DEFAULT '',
    menu_ids  varchar(4096)                          NOT NULL,
    create_by varchar(64)  DEFAULT ''                NOT NULL,
    create_at timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by varchar(64)  DEFAULT '',
    update_at timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at timestamp
);

COMMENT ON TABLE qa_tenant_package IS '租户套餐表';
COMMENT ON COLUMN qa_tenant_package.id IS '套餐编号';
COMMENT ON COLUMN qa_tenant_package.name IS '套餐名';
COMMENT ON COLUMN qa_tenant_package.status IS '租户状态（0停用 1正常）';
COMMENT ON COLUMN qa_tenant_package.remark IS '备注';
COMMENT ON COLUMN qa_tenant_package.menu_ids IS '关联的菜单编号';
COMMENT ON COLUMN qa_tenant_package.create_by IS '创建者';
COMMENT ON COLUMN qa_tenant_package.create_at IS '创建时间';
COMMENT ON COLUMN qa_tenant_package.update_by IS '更新者';
COMMENT ON COLUMN qa_tenant_package.update_at IS '更新时间';
COMMENT ON COLUMN qa_tenant_package.delete_at IS '删除时间';

DROP TABLE IF EXISTS qa_config CASCADE;
CREATE TABLE qa_config
(
    id         varchar(32) PRIMARY KEY,
    name       varchar(128)                           NOT NULL,
    key        varchar(128)                           NOT NULL,
    value      text,
    status     smallint     DEFAULT 1                 NOT NULL,
    create_by  varchar(64)  DEFAULT '',
    create_at  timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_by  varchar(64)  DEFAULT '',
    update_at  timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    delete_at  timestamp,
    tenant_id  varchar(32)  DEFAULT ''                NOT NULL
);

COMMENT ON TABLE qa_config IS '系统配置表';
COMMENT ON COLUMN qa_config.id IS '配置ID';
COMMENT ON COLUMN qa_config.name IS '配置名称';
COMMENT ON COLUMN qa_config.key IS '配置键';
COMMENT ON COLUMN qa_config.value IS '配置值';
COMMENT ON COLUMN qa_config.status IS '状态（0停用 1正常）';
COMMENT ON COLUMN qa_config.create_by IS '创建者';
COMMENT ON COLUMN qa_config.create_at IS '创建时间';
COMMENT ON COLUMN qa_config.update_by IS '更新者';
COMMENT ON COLUMN qa_config.update_at IS '更新时间';
COMMENT ON COLUMN qa_config.delete_at IS '删除时间';
COMMENT ON COLUMN qa_config.tenant_id IS '租户编号';

DROP INDEX IF EXISTS idx_config_key;
CREATE UNIQUE INDEX idx_config_key ON qa_config (key, tenant_id);

DROP INDEX IF EXISTS idx_config_name;
CREATE INDEX idx_config_name ON qa_config (name, tenant_id);