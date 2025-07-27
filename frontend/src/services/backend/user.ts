// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取用户列表 搜索时支持用户名、昵称、手机和邮箱筛选 GET /admin/users */
export async function listUsers(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.ListUsersParams,
  options?: { [key: string]: any },
) {
  return request<API.UserSearchResponse>(`/v1/admin/users`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建用户 创建一个新的用户 POST /admin/users */
export async function createUser(body: API.UserRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/users`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新用户 更新用户信息 PUT /admin/users/${param0} */
export async function updateUser(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.UpdateUserParams,
  body: API.UserRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/users/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除用户 删除指定ID的用户 DELETE /admin/users/${param0} */
export async function deleteUser(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.DeleteUserParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/users/${param0}`, {
    method: 'DELETE',
    params: {
      ...queryParams,
    },
    ...(options || {}),
  });
}

/** 登录 支持用户名或邮箱登录 POST /login */
export async function login(body: API.LoginRequest, options?: { [key: string]: any }) {
  return request<API.LoginResponse>(`/v1/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 刷新令牌 刷新访问令牌和刷新令牌 POST /refresh-token */
export async function refreshToken(
  body: API.RefreshTokenRequest,
  options?: { [key: string]: any },
) {
  return request<API.LoginResponse>(`/v1/refresh-token`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 注册 目前只支持通过邮箱进行注册 POST /register */
export async function register(body: API.RegisterRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 重置密码 重置用户密码 POST /reset-password */
export async function resetPassword(
  body: API.ResetPasswordRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/reset-password`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取用户详情 获取指定ID的用户详情 GET /users/${param0} */
export async function getUserById(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.GetUserByIDParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.UserResponse>(`/v1/users/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 获取用户菜单 获取当前用户的菜单列表 GET /users/menu */
export async function fetchDynamicMenu(options?: { [key: string]: any }) {
  return request<API.DynamicMenuResponse>(`/v1/users/menu`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 更新密码 更新用户密码 PUT /users/password */
export async function updatePassword(
  body: API.UpdatePasswordRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/users/password`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取当前用户 获取当前用户的详细信息 GET /users/profile */
export async function fetchCurrentUser(options?: { [key: string]: any }) {
  return request<API.UserResponse>(`/v1/users/profile`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 更新用户 更新用户信息 PUT /users/profile */
export async function updateProfile(body: API.UserRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/users/profile`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 上传头像 上传用户头像 PUT /users/profile/avatar */
export async function uploadAvatar(body: {}, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/users/profile/avatar`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    data: body,
    ...(options || {}),
  });
}
