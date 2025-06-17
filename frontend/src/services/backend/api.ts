// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 更新API 更新API信息 PUT /admin/api */
export async function putAdminApi(body: API.ApiUpdateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/api`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 创建API 创建新的API POST /admin/api */
export async function postAdminApi(body: API.ApiCreateRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/api`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除API 删除指定API DELETE /admin/api */
export async function deleteAdminApi(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteAdminApiParams,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/admin/api`, {
    method: 'DELETE',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取API列表 获取API列表 GET /admin/apis */
export async function getAdminApis(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminApisParams,
  options?: { [key: string]: any },
) {
  return request<API.ListApisResponse>(`/v1/admin/apis`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
