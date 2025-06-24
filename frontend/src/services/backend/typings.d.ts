declare namespace API {
  type ApiCreateRequest = {
    group?: string;
    method?: string;
    name?: string;
    path?: string;
  };

  type ApiDataItem = {
    createdAt?: string;
    group?: string;
    id?: number;
    method?: string;
    name?: string;
    path?: string;
    updatedAt?: string;
  };

  type ApiUpdateRequest = {
    group?: string;
    id: number;
    method?: string;
    name?: string;
    path?: string;
  };

  type CurrentUser = {
    avatar?: string;
    createdAt?: string;
    email?: string;
    nickname?: string;
    phone?: string;
    roles?: string[];
    updatedAt?: string;
    userid?: number;
    username?: string;
  };

  type deleteAdminApiParams = {
    /** API ID */
    id: number;
  };

  type deleteAdminMenuParams = {
    /** 菜单ID */
    id: number;
  };

  type deleteAdminRoleParams = {
    /** 角色ID */
    id: number;
  };

  type deleteAdminUserParams = {
    /** 用户ID */
    id: number;
  };

  type DeleteRobotParams = {
    /** 机器人ID */
    id: number;
  };

  type getAdminApisParams = {
    /** 页码 */
    page: number;
    /** 每页数量 */
    pageSize: number;
    /** API分组 */
    group?: string;
    /** API名称 */
    name?: string;
    /** API路径 */
    path?: string;
    /** 请求方法 */
    method?: string;
  };

  type getAdminRolePermissionsParams = {
    /** 角色名称 */
    role: string;
  };

  type getAdminRolesParams = {
    /** 页码 */
    page: number;
    /** 每页数量 */
    pageSize: number;
    /** 角色ID */
    sid?: string;
    /** 角色名称 */
    name?: string;
  };

  type getAdminUsersParams = {
    /** 页码 */
    page: number;
    /** 每页数量 */
    pageSize: number;
    /** 用户名 */
    username?: string;
    /** 昵称 */
    nickname?: string;
    /** 手机号 */
    phone?: string;
    /** 邮箱 */
    email?: string;
  };

  type GetRobotParams = {
    /** 机器人ID */
    id: number;
  };

  type GetRolePermissionsData = {
    list?: string[];
  };

  type GetUserPermissionsData = {
    list?: string[];
  };

  type GetUserResponse = {
    data?: CurrentUser;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type ListApisResponse = {
    data?: ListApisResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type ListApisResponseData = {
    groups?: string[];
    list?: ApiDataItem[];
    total?: number;
  };

  type ListMenuResponse = {
    data?: ListMenuResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type ListMenuResponseData = {
    list?: MenuDataItem[];
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

  type ListRolesResponse = {
    data?: ListRolesResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type ListRolesResponseData = {
    list?: RoleDataItem[];
    total?: number;
  };

  type ListUsersResponse = {
    data?: ListUsersResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type ListUsersResponseData = {
    list?: UserDataItem[];
    total?: number;
  };

  type LoginParams = {
    password: string;
    username: string;
  };

  type LoginResponse = {
    data?: LoginResult;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    success?: boolean;
  };

  type LoginResult = {
    accessToken?: string;
  };

  type MenuCreateRequest = {
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

  type MenuDataItem = {
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

  type MenuUpdateRequest = {
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
    updatedAt?: string;
    /** iframe模式下的跳转url，不能与path重复 */
    url?: string;
    /** 排序权重 */
    weight?: number;
  };

  type RegisterParams = {
    email: string;
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
    /** 通知地址 */
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
    /** 回调地址 */
    webhook?: string;
  };

  type RobotList = {
    /** 列表 */
    list?: Robot[];
    /** 总数 */
    total?: number;
  };

  type RobotParams = {
    /** 通知地址 */
    callback?: string;
    /** 描述 */
    desc?: string;
    /** 是否启用 */
    enabled?: boolean;
    /** 名称 */
    name?: string;
    /** 所有者 */
    owner?: string;
    /** 回调地址 */
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

  type RoleCreateRequest = {
    name: string;
    sid: string;
  };

  type RoleDataItem = {
    createdAt?: string;
    id?: number;
    name?: string;
    sid?: string;
    updatedAt?: string;
  };

  type RoleUpdateRequest = {
    id: number;
    name: string;
    sid: string;
  };

  type UpdateRobotParams = {
    /** 机器人ID */
    id: number;
  };

  type UpdateRolePermissionRequest = {
    list: string[];
    role: string;
  };

  type UserCreateRequest = {
    email?: string;
    nickname?: string;
    password: string;
    phone?: string;
    roles?: string[];
    username: string;
  };

  type UserDataItem = {
    createdAt?: string;
    email: string;
    id?: number;
    nickname: string;
    phone?: string;
    roles?: string[];
    updatedAt?: string;
    username: string;
  };

  type UserUpdateRequest = {
    email?: string;
    id?: number;
    nickname?: string;
    password?: string;
    phone?: string;
    roles?: string[];
    username: string;
  };
}
