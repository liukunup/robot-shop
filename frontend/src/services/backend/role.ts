// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取角色列表 搜索时支持角色名和 Casbin Role 筛选 GET /admin/roles */
export async function listRoles(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.ListRolesParams,
  options?: { [key: string]: any },
) {
  return request<API.RoleSearchResponse>(`/v1/admin/roles`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建角色 创建一个新的角色 POST /admin/roles */
export async function createRole(body: API.RoleRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/roles`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新角色 目前只允许更新角色名称 PUT /admin/roles/${param0} */
export async function updateRole(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.UpdateRoleParams,
  body: API.RoleRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/roles/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除角色 删除指定ID的角色 DELETE /admin/roles/${param0} */
export async function deleteRole(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.DeleteRoleParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/roles/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 获取角色权限 获取指定角色的权限列表 GET /admin/roles/permissions */
export async function getRolePermissions(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.GetRolePermissionsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetRolePermissionResponse>(`/v1/admin/roles/permissions`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 更新角色权限 更新指定角色的权限列表 PUT /admin/roles/permissions */
export async function updateRolePermissions(
  body: API.UpdateRolePermissionRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/admin/roles/permissions`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
