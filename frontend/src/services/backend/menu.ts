// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取菜单列表 获取所有菜单 GET /admin/menus */
export async function listMenus(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.ListMenusParams,
  options?: { [key: string]: any },
) {
  return request<API.MenuSearchResponse>(`/v1/admin/menus`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建菜单 创建一个新的菜单 POST /admin/menus */
export async function createMenu(body: API.MenuRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/menus`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新菜单 更新菜单数据 PUT /admin/menus/${param0} */
export async function updateMenu(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.UpdateMenuParams,
  body: API.MenuRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/menus/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除菜单 删除指定ID的菜单 DELETE /admin/menus/${param0} */
export async function deleteMenu(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.DeleteMenuParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/menus/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}
