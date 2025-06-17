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
    id?: number;
    nickname?: string;
    phone?: string;
    roles?: string[];
    updatedAt?: string;
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

  type deleteRobotIdParams = {
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

  type getRobotIdParams = {
    /** 机器人ID */
    id: number;
  };

  type getRobotParams = {
    /** page */
    page?: number;
    /** size */
    size?: number;
  };

  type GetRolePermissionsData = {
    list?: string[];
  };

  type GetUserPermissionsData = {
    list?: string[];
  };

  type GetUserResponse = {
    code?: number;
    data?: CurrentUser;
    message?: string;
    showType?: number;
    success?: boolean;
  };

  type ListApisResponse = {
    code?: number;
    data?: ListApisResponseData;
    message?: string;
    showType?: number;
    success?: boolean;
  };

  type ListApisResponseData = {
    groups?: string[];
    list?: ApiDataItem[];
    total?: number;
  };

  type ListMenuResponse = {
    code?: number;
    data?: ListMenuResponseData;
    message?: string;
    showType?: number;
    success?: boolean;
  };

  type ListMenuResponseData = {
    list?: MenuDataItem[];
  };

  type ListRolesResponse = {
    code?: number;
    data?: ListRolesResponseData;
    message?: string;
    showType?: number;
    success?: boolean;
  };

  type ListRolesResponseData = {
    list?: RoleDataItem[];
    total?: number;
  };

  type ListUsersResponse = {
    code?: number;
    data?: ListUsersResponseData;
    message?: string;
    showType?: number;
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
    code?: number;
    data?: LoginResult;
    message?: string;
    showType?: number;
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

  type PageResponseBackendApiV1RobotResponseData = {
    list?: RobotResponseData[];
    total?: number;
  };

  type putRobotIdParams = {
    /** 机器人ID */
    id: number;
  };

  type RegisterParams = {
    email: string;
    password: string;
  };

  type Response = {
    code?: number;
    data?: any;
    message?: string;
    showType?: number;
    success?: boolean;
  };

  type RobotRequest = {
    callback?: string;
    desc?: string;
    enabled?: boolean;
    name: string;
    options?: string;
    owner?: string;
    webhook?: string;
  };

  type RobotResponseData = {
    callback?: string;
    createdAt?: string;
    desc?: string;
    enabled?: boolean;
    id?: number;
    name?: string;
    options?: string;
    owner?: string;
    robot_id?: string;
    updatedAt?: string;
    webhook?: string;
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
