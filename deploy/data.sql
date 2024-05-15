# 一级菜单
INSERT INTO `system_menu` VALUES (1, '用户管理', 'User', 'menu-user', '/user', 'Layout', 1);
INSERT INTO `system_menu` VALUES (2, '系统设置', 'System', 'menu-system', '/system', 'Layout', 2);

# 二级菜单
INSERT INTO `system_sub_menu` VALUES (1, '用户管理', 'UserManagement', 'sub-menu-user', 'user', 'user/user/index', 1, 1);
INSERT INTO `system_sub_menu` VALUES (2, '分组管理', 'GroupManagement', 'sub-menu-group', 'group', 'user/group/index', 2, 1);
INSERT INTO `system_sub_menu` VALUES (3, '菜单管理', 'MenuManagement', 'sub-menu-menu', 'menu', 'system/menu/index', 1, 2);
INSERT INTO `system_sub_menu` VALUES (4, '系统设置', 'ConfigManagement', 'sub-menu-config', 'config', 'system/config/index', 2, 2);
