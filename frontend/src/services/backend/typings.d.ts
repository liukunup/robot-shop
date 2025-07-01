declare namespace API {
  type Api = {
    /** 创建时间 */
    createdAt?: string;
    /** 分组名 */
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

  type ApiDeleteParams = {
    /** 接口ID */
    id: number;
  };

  type ApiList = {
    /** 分组名列表 */
    groups?: string[];
    /** 列表 */
    list?: Api[];
    /** 总数 */
    total?: number;
  };

  type ApiRequest = {
    /** 分组名 */
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
    success?: boolean;
  };

  type ApiUpdateParams = {
    /** 接口ID */
    id: number;
  };

  type GetApiParams = {
    /** 接口ID */
    id: number;
  };

  type GetRobotParams = {
    /** 机器人ID */
    id: number;
  };

  type GetRolePermissionParams = {
    /** 角色名 */
    role: string;
  };

  type GetRolePermissionResponse = {
    data?: GetRolePermissionResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type GetRolePermissionResponseData = {
    /** 列表 */
    list?: string[];
    /** 总数 */
    total?: number;
  };

  type ListApisParams = {
    /** 页码 */
    page: number;
    /** 分页大小 */
    pageSize: number;
    /** 分组名 */
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
    page: number;
    /** 分页大小 */
    pageSize: number;
    /** 角色名 */
    name?: string;
    /** Casbin Role */
    role?: string;
  };

  type ListUsersParams = {
    /** 页码 */
    page: number;
    /** 分页大小 */
    pageSize: number;
    /** 用户名 */
    username?: string;
    /** 昵称 */
    nickname?: string;
    /** 手机 */
    phone?: string;
    /** 邮箱 */
    email?: string;
  };

  type LoginRequest = {
    password: string;
    username: string;
  };

  type LoginResponse = {
    data?: LoginResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type LoginResponseData = {
    accessToken?: string;
  };

  type Menu = {
    /** 绑定的组件 */
    component?: string;
    /** 是否保活 */
    hideInMenu?: boolean;
    /** 图标，使用字符串表示 */
    icon?: string;
    /** 唯一id，使用整数表示 */
    id?: number;
    /** 是否保活 */
    keepAlive?: boolean;
    /** 本地化标识 */
    locale?: string;
    /** 同路由中的name，唯一标识 */
    name?: string;
    /** 父级菜单的id，使用整数表示 */
    parentId?: number;
    /** 地址 */
    path?: string;
    /** 重定向地址 */
    redirect?: string;
    /** 展示名称 */
    title?: string;
    /** 是否保活 */
    updatedAt?: string;
    /** iframe模式下的跳转url，不能与path重复 */
    url?: string;
    /** 排序权重 */
    weight?: number;
  };

  type MenuDeleteParams = {
    /** 菜单ID */
    id: number;
  };

  type MenuList = {
    /** 列表 */
    list?: Menu[];
    /** 总数 */
    total?: number;
  };

  type MenuRequest = {
    /** 绑定的组件 */
    component?: string;
    /** 是否保活 */
    hideInMenu?: boolean;
    /** 图标，使用字符串表示 */
    icon?: string;
    /** 是否保活 */
    keepAlive?: boolean;
    /** 本地化标识 */
    locale?: string;
    /** 同路由中的name，唯一标识 */
    name?: string;
    /** 父级菜单的id，使用整数表示 */
    parentId?: number;
    /** 地址 */
    path?: string;
    /** 重定向地址 */
    redirect?: string;
    /** 展示名称 */
    title?: string;
    /** iframe模式下的跳转url，不能与path重复 */
    url?: string;
    /** 排序权重 */
    weight?: number;
  };

  type MenuSearchResponse = {
    data?: MenuList;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type MenuUpdateParams = {
    /** 菜单ID */
    id: number;
  };

  type RegisterRequest = {
    /** 邮箱 */
    email: string;
    /** 密码 */
    password: string;
  };

  type Response = {
    data?: any;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
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

  type RobotDeleteParams = {
    /** 机器人ID */
    id: number;
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
    success?: boolean;
  };

  type RobotUpdateParams = {
    /** 机器人ID */
    id: number;
  };

  type Role = {
    /** 创建时间 */
    createdAt?: string;
    /** ID */
    id?: number;
    /** 角色名 */
    name?: string;
    /** Casbin Role */
    role?: string;
    /** 更新时间 */
    updatedAt?: string;
  };

  type RoleDeleteParams = {
    /** 角色ID */
    id: number;
  };

  type RoleList = {
    list?: Role[];
    total?: number;
  };

  type RoleRequest = {
    /** 角色名 */
    name: string;
    /** Casbin Role */
    role: string;
  };

  type RoleSearchResponse = {
    data?: RoleList;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type RoleUpdateParams = {
    /** 角色ID */
    id: number;
  };

  type UpdateRolePermissionRequest = {
    /** 权限列表 */
    list: string[];
    /** 角色名 */
    role: string;
  };

  type User = {
    /** 头像 */
    avatar?: string;
    /** 创建时间 */
    createdAt?: string;
    /** 邮箱 */
    email: string;
    /** ID */
    id?: number;
    /** 昵称 */
    nickname: string;
    /** 手机 */
    phone?: string;
    /** 角色 */
    roles?: Role[];
    /** 状态 0:待激活 1:正常 2:禁用 */
    status?: number;
    /** 更新时间 */
    updatedAt?: string;
    /** 用户名 */
    username: string;
  };

  type UserDeleteParams = {
    /** 用户ID */
    id: number;
  };

  type UserList = {
    /** 列表 */
    list?: User[];
    /** 总数 */
    total?: number;
  };

  type UserPermissionResponse = {
    data?: UserPermissionResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type UserPermissionResponseData = {
    list?: string[];
    total?: number;
  };

  type UserRequest = {
    email: string;
    nickname?: string;
    phone?: string;
    roles?: string[];
    /** 状态 0:待激活 1:正常 2:禁用 */
    status?: number;
    username: string;
  };

  type UserResponse = {
    data?: User;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
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
    success?: boolean;
  };

  type UserUpdateParams = {
    /** 用户ID */
    id: number;
  };
}
