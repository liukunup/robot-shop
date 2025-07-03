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
    /** 分组列表 */
    groups?: string[];
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
    casbinRole?: string;
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
    /** 是否成功 */
    success?: boolean;
  };

  type LoginResponseData = {
    accessToken?: string;
  };

  type Menu = {
    access?: string;
    component?: string;
    /** 创建时间 */
    createdAt?: string;
    icon?: string;
    /** ID */
    id?: number;
    name?: string;
    parentId?: number;
    path?: string;
    redirect?: string;
    /** 更新时间 */
    updatedAt?: string;
    weight?: number;
  };

  type MenuList = {
    /** 列表 */
    list?: Menu[];
    /** 总数 */
    total?: number;
  };

  type MenuListResponse = {
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

  type MenuRequest = {
    access?: string;
    component?: string;
    icon?: string;
    name?: string;
    parentId?: number;
    path?: string;
    redirect?: string;
    weight?: number;
  };

  type MenuTree = {
    access?: string;
    children?: Menu[];
    component?: string;
    /** 创建时间 */
    createdAt?: string;
    icon?: string;
    /** ID */
    id?: number;
    name?: string;
    parentId?: number;
    path?: string;
    redirect?: string;
    /** 更新时间 */
    updatedAt?: string;
    weight?: number;
  };

  type MenuTreeResponse = {
    data?: MenuTreeResponseData;
    /** 错误码 */
    errorCode?: number;
    /** 报错信息 */
    errorMessage?: string;
    /** 前端展示方式 */
    errorShowType?: number;
    /** 是否成功 */
    success?: boolean;
  };

  type MenuTreeResponseData = {
    root?: MenuTree[];
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
    /** Casbin Role */
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
    /** Casbin Role */
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

  type UpdateApiParams = {
    /** 接口ID */
    id: number;
  };

  type UpdateMenuParams = {
    /** 菜单ID */
    id: number;
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
    /** Casbin Role */
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
    /** 创建时间 */
    createdAt?: string;
    /** 邮箱 */
    email?: string;
    /** ID */
    id?: number;
    /** 昵称 */
    nickname?: string;
    /** 手机 */
    phone?: string;
    /** 角色 */
    roles?: Role[];
    /** 状态 0:待激活 1:正常 2:禁用 */
    status?: number;
    /** 更新时间 */
    updatedAt?: string;
    /** 用户名 */
    username?: string;
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
    /** 是否成功 */
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
