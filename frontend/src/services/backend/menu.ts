// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 更新菜单 更新菜单信息 PUT /admin/menu */
export async function putAdminMenu(body: API.MenuUpdateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/menu`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 创建菜单 创建新的菜单 POST /admin/menu */
export async function postAdminMenu(body: API.MenuCreateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/menu`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除菜单 删除指定菜单 DELETE /admin/menu */
export async function deleteAdminMenu(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteAdminMenuParams,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/admin/menu`, {
    method: 'DELETE',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取菜单列表 获取菜单列表 GET /admin/menus */
export async function getAdminMenus(options?: { [key: string]: any }) {
  return request<API.ListMenuResponse>(`/v1/admin/menus`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取当前用户菜单 获取当前用户的菜单列表 GET /menu */
export async function queryCurrentMenu(options?: { [key: string]: any }) {
  return request<API.ListMenuResponse>(`/v1/menu`, {
    method: 'GET',
    ...(options || {}),
  });
}
