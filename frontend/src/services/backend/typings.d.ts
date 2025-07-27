declare namespace API {
  type Api = {
    /** 创建时间 */
    createdAt?: string;
    /** 分组 */
    group?: string;
    /** ID */
    id?: number;
    /** 方法 */
    method?: string;
    /** 名称 */
    name?: string;
    /** 路径 */
    path?: string;
    /** 更新时间 */
    updatedAt?: string;
  };

  type ApiList = {
    /** 列表 */
    list?: Api[];
    /** 总数 */
    total?: number;
  };

  type ApiRequest = {
    /** 分组 */
    group?: string;
    /** 方法 */
    method?: string;
    /** 名称 */
    name?: string;
    /** 路径 */
    path?: string;
  };

  type ApiResponse = {
    data?: Api;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type ApiSearchResponse = {
    data?: ApiList;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type DeleteApiParams = {
    /** 接口ID */
    id: number;
  };

  type DeleteMenuParams = {
    /** 菜单ID */
    id: number;
  };

  type DeleteRobotParams = {
    /** 机器人ID */
    id: number;
  };

  type DeleteRoleParams = {
    /** 角色ID */
    id: number;
  };

  type DeleteUserParams = {
    /** 用户ID */
    id: number;
  };

  type DynamicMenuResponse = {
    data?: DynamicMenuResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type DynamicMenuResponseData = {
    /** 顶级菜单 */
    list?: MenuNode[];
  };

  type GetApiParams = {
    /** 接口ID */
    id: number;
  };

  type GetRobotParams = {
    /** 机器人ID */
    id: number;
  };

  type GetRolePermissionResponse = {
    data?: GetRolePermissionResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type GetRolePermissionResponseData = {
    /** 列表 */
    list?: string[];
    /** 总数 */
    total?: number;
  };

  type GetRolePermissionsParams = {
    /** 角色名 */
    role: string;
  };

  type GetUserByIDParams = {
    /** 用户ID */
    id: number;
  };

  type ListApisParams = {
    /** 页码 */
    page: number;
    /** 分页大小 */
    pageSize: number;
    /** 分组 */
    group?: string;
    /** 名称 */
    name?: string;
    /** 路径 */
    path?: string;
    /** 方法 */
    method?: string;
  };

  type ListMenusParams = {
    /** 页码 */
    page: number;
    /** 分页大小 */
    pageSize: number;
    /** 名称 */
    name?: string;
    /** 路径 */
    path?: string;
    /** 可见性 */
    access?: string;
  };

  type ListRobotsParams = {
    /** 页码 */
    page: number;
    /** 分页大小 */
    pageSize: number;
    /** 名称 */
    name?: string;
    /** 描述 */
    desc?: string;
    /** 所有者 */
    owner?: string;
  };

  type ListRolesParams = {
    /** 页码 */
    page?: number;
    /** 分页大小 */
    pageSize?: number;
    /** 角色名 */
    name?: string;
    /** Casbin Role */
    casbinRole?: string;
  };

  type ListUsersParams = {
    /** 页码 */
    page: number;
    /** 分页大小 */
    pageSize: number;
    /** 邮箱 */
    email?: string;
    /** 用户名 */
    username?: string;
    /** 昵称 */
    nickname?: string;
  };

  type LoginRequest = {
    /** 密码 */
    password: string;
    /** 用户名 */
    username: string;
  };

  type LoginResponse = {
    data?: TokenPair;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type Menu = {
    /** 可见性 */
    access?: string;
    /** 组件 */
    component?: string;
    /** 创建时间 */
    createdAt?: string;
    disabled?: boolean;
    disabledTooltip?: boolean;
    /** 隐藏自身+子节点提升并打平 */
    flatMenu?: boolean;
    /** 隐藏子节点 */
    hideChildrenInMenu?: boolean;
    /** 隐藏自身和子节点 */
    hideInMenu?: boolean;
    /** 图标 */
    icon?: string;
    /** ID */
    id?: number;
    key?: string;
    /** 国际化 */
    locale?: string;
    /** 名称 */
    name?: string;
    /** 父级菜单 */
    parentId?: number;
    parentKeys?: string;
    /** 路径 */
    path?: string;
    /** 重定向 */
    redirect?: string;
    /** 指定外链打开形式 */
    target?: string;
    tooltip?: string;
    /** 更新时间 */
    updatedAt?: string;
  };

  type MenuList = {
    /** 列表 */
    list?: Menu[];
    /** 总数 */
    total?: number;
  };

  type MenuNode = {
    /** 可见性 */
    access?: string;
    /** 子菜单 */
    children?: MenuNode[];
    /** 组件 */
    component?: string;
    /** 创建时间 */
    createdAt?: string;
    disabled?: boolean;
    disabledTooltip?: boolean;
    /** 隐藏自身+子节点提升并打平 */
    flatMenu?: boolean;
    /** 隐藏子节点 */
    hideChildrenInMenu?: boolean;
    /** 隐藏自身和子节点 */
    hideInMenu?: boolean;
    /** 图标 */
    icon?: string;
    /** ID */
    id?: number;
    key?: string;
    /** 国际化 */
    locale?: string;
    /** 名称 */
    name?: string;
    /** 父级菜单 */
    parentId?: number;
    parentKeys?: string;
    /** 路径 */
    path?: string;
    /** 重定向 */
    redirect?: string;
    /** 指定外链打开形式 */
    target?: string;
    tooltip?: string;
    /** 更新时间 */
    updatedAt?: string;
  };

  type MenuRequest = {
    /** 可见性 */
    access?: string;
    /** 组件 */
    component?: string;
    disabled?: boolean;
    disabledTooltip?: boolean;
    /** 隐藏自身+子节点提升并打平 */
    flatMenu?: boolean;
    /** 隐藏子节点 */
    hideChildrenInMenu?: boolean;
    /** 隐藏自身和子节点 */
    hideInMenu?: boolean;
    /** 图标 */
    icon?: string;
    key?: string;
    /** 国际化 */
    locale?: string;
    /** 名称 */
    name?: string;
    /** 父级菜单 */
    parentId?: number;
    parentKeys?: string;
    /** 路径 */
    path?: string;
    /** 重定向 */
    redirect?: string;
    /** 指定外链打开形式 */
    target?: string;
    tooltip?: string;
  };

  type MenuSearchResponse = {
    data?: MenuList;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type RefreshTokenRequest = {
    /** 刷新令牌 */
    refreshToken: string;
  };

  type RegisterRequest = {
    /** 邮箱 */
    email: string;
    /** 密码 */
    password: string;
  };

  type ResetPasswordRequest = {
    /** 邮箱 */
    email: string;
  };

  type Response = {
    /** 返回数据 */
    data?: any;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type Robot = {
    /** 回调地址 */
    callback?: string;
    /** 创建时间 */
    createdAt?: string;
    /** 描述 */
    desc?: string;
    /** 是否启用 */
    enabled?: boolean;
    /** ID */
    id?: number;
    /** 名称 */
    name?: string;
    /** 所有者 */
    owner?: string;
    /** 更新时间 */
    updatedAt?: string;
    /** 通知地址 */
    webhook?: string;
  };

  type RobotList = {
    /** 列表 */
    list?: Robot[];
    /** 总数 */
    total?: number;
  };

  type RobotRequest = {
    /** 回调地址 */
    callback?: string;
    /** 描述 */
    desc?: string;
    /** 是否启用 */
    enabled?: boolean;
    /** 名称 */
    name?: string;
    /** 所有者 */
    owner?: string;
    /** 通知地址 */
    webhook?: string;
  };

  type RobotResponse = {
    data?: Robot;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type RobotSearchResponse = {
    data?: RobotList;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type Role = {
    /** Casbin-Role */
    casbinRole?: string;
    /** 创建时间 */
    createdAt?: string;
    /** ID */
    id?: number;
    /** 角色名 */
    name?: string;
    /** 更新时间 */
    updatedAt?: string;
  };

  type RoleList = {
    /** 列表 */
    list?: Role[];
    /** 总数 */
    total?: number;
  };

  type RoleRequest = {
    /** Casbin-Role */
    casbinRole: string;
    /** 角色名 */
    name: string;
  };

  type RoleSearchResponse = {
    data?: RoleList;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type TokenPair = {
    /** 访问令牌 */
    accessToken?: string;
    /** 过期时间(单位:秒) */
    expiresIn?: number;
    /** 刷新令牌 */
    refreshToken?: string;
  };

  type UpdateApiParams = {
    /** 接口ID */
    id: number;
  };

  type UpdateMenuParams = {
    /** 菜单ID */
    id: number;
  };

  type UpdatePasswordRequest = {
    /** 新密码 */
    newPassword: string;
    /** 旧密码 */
    oldPassword: string;
  };

  type UpdateRobotParams = {
    /** 机器人ID */
    id: number;
  };

  type UpdateRoleParams = {
    /** 角色ID */
    id: number;
  };

  type UpdateRolePermissionRequest = {
    /** Casbin-Role */
    casbinRole: string;
    /** 权限列表 */
    list: string[];
  };

  type UpdateUserParams = {
    /** 用户ID */
    id: number;
  };

  type User = {
    /** 头像 */
    avatar?: string;
    /** 个人简介 */
    bio?: string;
    /** 创建时间 */
    createdAt?: string;
    /** 邮箱 */
    email?: string;
    /** 语言 */
    language?: string;
    /** 昵称 */
    nickname?: string;
    /** 角色 */
    roles?: Role[];
    /** 状态 0:待激活 1:正常 2:禁用 */
    status?: number;
    /** 主题 */
    theme?: string;
    /** 时区 */
    timezone?: string;
    /** 更新时间 */
    updatedAt?: string;
    /** ID */
    userid?: number;
    /** 用户名 */
    username?: string;
  };

  type UserList = {
    /** 列表 */
    list?: User[];
    /** 总数 */
    total?: number;
  };

  type UserRequest = {
    /** 个人简介 */
    bio?: string;
    /** 邮箱 */
    email?: string;
    /** 语言 */
    language?: string;
    /** 昵称 */
    nickname?: string;
    /** 角色 */
    roles?: string[];
    /** 状态 0:待激活 1:正常 2:禁用 */
    status?: number;
    /** 主题 */
    theme?: string;
    /** 时区 */
    timezone?: string;
    /** 用户名 */
    username?: string;
  };

  type UserResponse = {
    data?: User;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type UserSearchResponse = {
    data?: UserList;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };
}
