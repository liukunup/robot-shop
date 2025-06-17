// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 更新角色 更新角色信息 PUT /admin/role */
export async function putAdminRole(body: API.RoleUpdateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/role`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 创建角色 创建新的角色 POST /admin/role */
export async function postAdminRole(body: API.RoleCreateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/role`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除角色 删除指定角色 DELETE /admin/role */
export async function deleteAdminRole(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteAdminRoleParams,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/admin/role`, {
    method: 'DELETE',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取角色权限 获取指定角色的权限列表 GET /admin/role/permissions */
export async function getAdminRolePermissions(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminRolePermissionsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetRolePermissionsData>(`/v1/admin/role/permissions`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 更新角色权限 更新指定角色的权限列表 PUT /admin/role/permissions */
export async function putAdminRolePermissions(
  body: API.UpdateRolePermissionRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/admin/role/permissions`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取角色列表 获取角色列表 GET /admin/roles */
export async function getAdminRoles(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminRolesParams,
  options?: { [key: string]: any },
) {
  return request<API.ListRolesResponse>(`/v1/admin/roles`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
