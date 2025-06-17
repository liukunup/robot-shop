// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 更新用户 更新用户 PUT /admin/user */
export async function putAdminUser(body: API.UserUpdateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/user`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 创建用户 创建用户 POST /admin/user */
export async function postAdminUser(body: API.UserCreateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/user`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除用户 删除用户 DELETE /admin/user */
export async function deleteAdminUser(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteAdminUserParams,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/admin/user`, {
    method: 'DELETE',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取用户权限 获取当前用户的权限列表 GET /admin/user/permissions */
export async function getAdminUserPermissions(options?: { [key: string]: any }) {
  return request<API.GetUserPermissionsData>(`/v1/admin/user/permissions`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取用户列表 获取用户列表 GET /admin/users */
export async function getAdminUsers(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminUsersParams,
  options?: { [key: string]: any },
) {
  return request<API.ListUsersResponse>(`/v1/admin/users`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 账号登录 支持用户名或邮箱登录 POST /login */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginResponse>(`/v1/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户注册 目前只支持邮箱注册 POST /register */
export async function register(body: API.RegisterParams, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取当前用户 获取当前用户的详细信息 GET /user */
export async function queryCurrentUser(options?: { [key: string]: any }) {
  return request<API.GetUserResponse>(`/v1/user`, {
    method: 'GET',
    ...(options || {}),
  });
}
