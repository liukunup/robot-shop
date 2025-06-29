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
export async function userCreate(body: API.UserRequest, options?: { [key: string]: any }) {
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
export async function userUpdate(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.UserUpdateParams,
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
export async function userDelete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.UserDeleteParams,
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

/** 获取当前用户 获取当前用户的详细信息 GET /users/me */
export async function fetchCurrentUser(options?: { [key: string]: any }) {
  return request<API.UserResponse>(`/v1/users/me`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取用户菜单 获取当前用户的菜单列表 GET /users/me/menu */
export async function fetchCurrentMenu(options?: { [key: string]: any }) {
  return request<API.MenuSearchResponse>(`/v1/users/me/menu`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取用户权限 获取当前用户的权限列表 GET /users/me/permission */
export async function fetchCurrentPermission(options?: { [key: string]: any }) {
  return request<API.UserPermissionResponse>(`/v1/users/me/permission`, {
    method: 'GET',
    ...(options || {}),
  });
}
