# 一级菜单
INSERT INTO `system_menu` VALUES (1, '系统导航', 'Navigation', 'navigation', '/', 'Layout', 1, '/navigation/sites');
INSERT INTO `system_menu` VALUES (2, '用户管理', 'User', 'menu-user', '/user', 'Layout', 2, null);
INSERT INTO `system_menu` VALUES (3, '系统设置', 'System', 'menu-system', '/system', 'Layout', 3, null);

# 二级菜单
INSERT INTO `system_sub_menu` VALUES (1, '站点导航', 'SiteNavigation', 'sub-menu-site-navigation', 'navigation/sites', 'dashboard/index', 1, null, 1);
INSERT INTO `system_sub_menu` VALUES (2, '用户管理', 'UserManagement', 'sub-menu-user', 'user', 'user/user/index', 1, null, 2);
INSERT INTO `system_sub_menu` VALUES (3, '分组管理', 'GroupManagement', 'sub-menu-group', 'group', 'user/group/index', 2, null, 2);
INSERT INTO `system_sub_menu` VALUES (4, '菜单管理', 'MenuManagement', 'sub-menu-menu', 'menu', 'system/menu/index', 1, null, 3);
INSERT INTO `system_sub_menu` VALUES (5, '系统设置', 'ConfigManagement', 'sub-menu-config', 'config', 'system/config/index', 2, null, 3);
