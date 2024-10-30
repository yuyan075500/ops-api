SET NAMES 'utf8mb4';

# 一级菜单
INSERT INTO `system_menu` VALUES (1, '用户管理', 'User', 'menu-user', '/user', 'Layout', 2, null);
INSERT INTO `system_menu` VALUES (2, '资产管理', 'Asset', 'menu-asset', '/asset', 'Layout', 3, null);
INSERT INTO `system_menu` VALUES (3, '日志审计', 'Audit', 'menu-audit', '/audit', 'Layout', 4, null);
INSERT INTO `system_menu` VALUES (4, '系统设置', 'System', 'menu-system', '/system', 'Layout', 5, null);

# 二级菜单
INSERT INTO `system_sub_menu` VALUES (1, '用户管理', 'UserManagement', 'sub-menu-user', 'user', 'user/user/index', 1, null, 1);
INSERT INTO `system_sub_menu` VALUES (2, '分组管理', 'GroupManagement', 'sub-menu-group', 'group', 'user/group/index', 2, null, 1);
INSERT INTO `system_sub_menu` VALUES (3, '账号管理', 'AccountManagement', 'sub-menu-account', 'account', 'asset/account/index', 1, null, 2);
INSERT INTO `system_sub_menu` VALUES (4, '站点管理', 'SiteManagement', 'sub-menu-site', 'site', 'asset/site/index', 2, null, 2);
INSERT INTO `system_sub_menu` VALUES (5, '登录日志', 'AuditLoginRecord', 'sub-menu-login-record', 'login', 'audit/login/index', 1, null, 3);
INSERT INTO `system_sub_menu` VALUES (6, '短信记录', 'AuditSMSRecord', 'sub-menu-sms-record', 'sms', 'audit/sms/index', 2, null, 3);
INSERT INTO `system_sub_menu` VALUES (7, '菜单管理', 'MenuManagement', 'sub-menu-menu', 'menu', 'system/menu/index', 1, null, 4);
INSERT INTO `system_sub_menu` VALUES (8, '定时任务', 'CornManagement', 'sub-menu-corn', 'corn', 'system/corn/index', 2, null, 4);

# API接口
INSERT INTO `system_path` VALUES (1, 'AddUser', '/api/v1/user', 'POST', 'UserManagement', '新增用户');
INSERT INTO `system_path` VALUES (2, 'UpdateUser', '/api/v1/user', 'PUT', 'UserManagement', '修改用户');
INSERT INTO `system_path` VALUES (3, 'UpdateUserPassword', '/api/v1/user/reset_password', 'PUT', 'UserManagement', '密码重置');
INSERT INTO `system_path` VALUES (4, 'ResetUserMFA', '/api/v1/user/reset_mfa/:id', 'PUT', 'UserManagement', 'MAF重置');
INSERT INTO `system_path` VALUES (5, 'DeleteUser', '/api/v1/user/:id', 'DELETE', 'UserManagement', '删除用户');
INSERT INTO `system_path` VALUES (6, 'GetUserList', '/api/v1/users', 'GET', 'UserManagement', '获取用户列表（表格）');
INSERT INTO `system_path` VALUES (7, 'UserSyncAd', '/api/v1/user/sync/ad', 'post', 'UserManagement', 'LDAP用户同步');
INSERT INTO `system_path` VALUES (8, 'GetUserListAll', '/api/v1/user/list', 'GET', 'UserManagement', '获取用户列表（所有）');
INSERT INTO `system_path` VALUES (9, 'AddGroup', '/api/v1/group', 'POST', 'GroupManagement', '新增分组');
INSERT INTO `system_path` VALUES (10, 'UpdateGroup', '/api/v1/group', 'PUT', 'GroupManagement', '修改分组');
INSERT INTO `system_path` VALUES (11, 'UpdateGroupUser', '/api/v1/group/users', 'PUT', 'GroupManagement', '更改分组用户');
INSERT INTO `system_path` VALUES (12, 'UpdateGroupPermission', '/api/v1/group/permissions', 'PUT', 'GroupManagement', '更改分组权限');
INSERT INTO `system_path` VALUES (13, 'DeleteGroup', '/api/v1/group/:id', 'DELETE', 'GroupManagement', '删除分组');
INSERT INTO `system_path` VALUES (14, 'GetGroupList', '/api/v1/groups', 'GET', 'GroupManagement', '获取分组列表');
INSERT INTO `system_path` VALUES (15, 'GetMenuListAll', '/api/v1/menu/list', 'GET', 'GroupManagement', '获取菜单列表');
INSERT INTO `system_path` VALUES (16, 'GetPathListAll', '/api/v1/path/list', 'GET', 'GroupManagement', '获取接口列表');
INSERT INTO `system_path` VALUES (17, 'GetSiteList', '/api/v1/sites', 'GET', 'SiteManagement', '获取站点列表');
INSERT INTO `system_path` VALUES (18, 'AddSite', '/api/v1/site', 'POST', 'SiteManagement', '新增站点');
INSERT INTO `system_path` VALUES (19, 'UpdateSite', '/api/v1/site', 'PUT', 'SiteManagement', '修改站点');
INSERT INTO `system_path` VALUES (20, 'DeleteSite', '/api/v1/site/:id', 'DELETE', 'SiteManagement', '删除站点');
INSERT INTO `system_path` VALUES (21, 'AddSiteGroup', '/api/v1/site/group', 'POST', 'SiteManagement', '新增站点分组');
INSERT INTO `system_path` VALUES (22, 'UpdateSiteGroup', '/api/v1/site/group', 'PUT', 'SiteManagement', '修改站点分组');
INSERT INTO `system_path` VALUES (23, 'DeleteSiteGroup', '/api/v1/site/group/:id', 'DELETE', 'SiteManagement', '删除站点分组');
INSERT INTO `system_path` VALUES (24, 'UpdateSiteUser', '/api/v1/site/users', 'PUT', 'SiteManagement', '更改站点用户');
INSERT INTO `system_path` VALUES (25, 'UpdateSiteTag', '/api/v1/site/tags', 'PUT', 'SiteManagement', '更改站点标签');
INSERT INTO `system_path` VALUES (26, 'GetSMSRecordList', '/api/v1/audit/sms', 'GET', 'AuditSMSRecord', '获取短信发送记录');
INSERT INTO `system_path` VALUES (27, 'GetLoginRecordList', '/api/v1/audit/login', 'GET', 'AuditLoginRecord', '获取用户登录记录');
INSERT INTO `system_path` VALUES (28, 'GetMenuList', '/api/v1/menus', 'GET', 'MenuManagement', '获取菜单列表');
INSERT INTO `system_path` VALUES (29, 'GetPathList', '/api/v1/paths', 'GET', 'MenuManagement', '获取菜单接口');
INSERT INTO `system_path` VALUES (30, 'GetTaskList', '/api/v1/tasks', 'GET', 'CornManagement', '获取定时任务列表');
INSERT INTO `system_path` VALUES (31, 'AddTask', '/api/v1/site', 'POST', 'CornManagement', '新增定时任务');
INSERT INTO `system_path` VALUES (32, 'UpdateTask', '/api/v1/site', 'PUT', 'CornManagement', '修改定时任务');
INSERT INTO `system_path` VALUES (33, 'DeleteTask', '/api/v1/site/:id', 'DELETE', 'CornManagement', '删除定时任务');
INSERT INTO `system_path` VALUES (34, 'GetTaskLogList', '/api/v1/task/logs', 'GET', 'CornManagement', '获取定时任务执行日志列表');