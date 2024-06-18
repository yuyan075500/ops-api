# 一级菜单
INSERT INTO `system_menu` VALUES (1, '系统导航', 'Navigation', 'navigation', '/', 'Layout', 1, '/navigation/sites');
INSERT INTO `system_menu` VALUES (2, '用户管理', 'User', 'menu-user', '/user', 'Layout', 2, null);
INSERT INTO `system_menu` VALUES (3, '日志审计', 'Audit', 'menu-audit', '/audit', 'Layout', 3, null);
INSERT INTO `system_menu` VALUES (4, '系统设置', 'System', 'menu-system', '/system', 'Layout', 4, null);

# 二级菜单
INSERT INTO `system_sub_menu` VALUES (1, '站点导航', 'SiteNavigation', 'sub-menu-site-navigation', 'navigation/sites', 'dashboard/index', 1, null, 1);
INSERT INTO `system_sub_menu` VALUES (2, '用户管理', 'UserManagement', 'sub-menu-user', 'user', 'user/user/index', 1, null, 2);
INSERT INTO `system_sub_menu` VALUES (3, '分组管理', 'GroupManagement', 'sub-menu-group', 'group', 'user/group/index', 2, null, 2);
INSERT INTO `system_sub_menu` VALUES (4, '登录日志', 'AuditLoginRecord', 'sub-menu-login-record', 'login', 'audit/login/index', 1, null, 3);
INSERT INTO `system_sub_menu` VALUES (5, '短信记录', 'AuditSMSRecord', 'sub-menu-sms-record', 'sms', 'audit/sms/index', 2, null, 3);
INSERT INTO `system_sub_menu` VALUES (6, '菜单管理', 'MenuManagement', 'sub-menu-menu', 'menu', 'system/menu/index', 1, null, 4);
INSERT INTO `system_sub_menu` VALUES (7, '系统设置', 'ConfigManagement', 'sub-menu-config', 'config', 'system/config/index', 2, null, 4);

# API接口
INSERT INTO `system_path` VALUES (1, 'AddUser', '/api/v1/user', 'POST', 'UserManagement', '新增用户');
INSERT INTO `system_path` VALUES (2, 'UpdateUser', '/api/v1/user', 'PUT', 'UserManagement', '修改用户');
INSERT INTO `system_path` VALUES (3, 'UpdateUserPassword', '/api/v1/user/reset_password', 'PUT', 'UserManagement', '密码重置');
INSERT INTO `system_path` VALUES (4, 'ResetUserMFA', '/api/v1/user/reset_mfa/:id', 'PUT', 'UserManagement', 'MAF重置');
INSERT INTO `system_path` VALUES (5, 'DeleteUser', '/api/v1/user/:id', 'DELETE', 'UserManagement', '删除用户');
INSERT INTO `system_path` VALUES (6, 'GetUserList', '/api/v1/users', 'GET', 'UserManagement', '获取用户列表');
INSERT INTO `system_path` VALUES (7, 'AddGroup', '/api/v1/group', 'POST', 'GroupManagement', '新增分组');
INSERT INTO `system_path` VALUES (8, 'UpdateGroup', '/api/v1/group', 'PUT', 'GroupManagement', '修改分组');
INSERT INTO `system_path` VALUES (9, 'UpdateGroupUser', '/api/v1/group/users', 'PUT', 'GroupManagement', '更改分组用户');
INSERT INTO `system_path` VALUES (10, 'UpdateGroupPermission', '/api/v1/group/permissions', 'PUT', 'GroupManagement', '更改分组权限');
INSERT INTO `system_path` VALUES (11, 'DeleteGroup', '/api/v1/group/:id', 'DELETE', 'GroupManagement', '删除分组');
INSERT INTO `system_path` VALUES (12, 'GetGroupList', '/api/v1/groups', 'GET', 'GroupManagement', '获取分组列表');
INSERT INTO `system_path` VALUES (13, 'GetMenuListAll', '/api/v1/menu/list', 'GET', 'GroupManagement', '获取菜单列表');
INSERT INTO `system_path` VALUES (14, 'GetUserListAll', '/api/v1/user/list', 'GET', 'GroupManagement', '获取用户列表');
INSERT INTO `system_path` VALUES (15, 'GetPathListAll', '/api/v1/path/list', 'GET', 'GroupManagement', '获取接口列表');
INSERT INTO `system_path` VALUES (16, 'GetSMSRecordList', '/api/v1/loginRecord', 'GET', 'AuditSMSRecord', '获取短信发送记录');
INSERT INTO `system_path` VALUES (17, 'GetMenuList', '/api/v1/menus', 'GET', 'MenuManagement', '获取菜单列表');
INSERT INTO `system_path` VALUES (18, 'GetPathList', '/api/v1/paths', 'GET', 'MenuManagement', '获取菜单接口');

