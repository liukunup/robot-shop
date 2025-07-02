// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取接口列表 搜索时支持分组名、名称、路径和方法筛选 GET /admin/apis */
export async function listApis(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.ListApisParams,
  options?: { [key: string]: any },
) {
  return request<API.ApiSearchResponse>(`/v1/admin/apis`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建接口 创建一个新的接口 POST /admin/apis */
export async function createApi(body: API.ApiRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/admin/apis`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取接口 获取指定ID的接口信息 GET /admin/apis/${param0} */
export async function getApi(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.GetApiParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.ApiResponse>(`/v1/admin/apis/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新接口 更新接口数据 PUT /admin/apis/${param0} */
export async function updateApi(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.UpdateApiParams,
  body: API.ApiRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/apis/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除接口 删除指定ID的接口 DELETE /admin/apis/${param0} */
export async function deleteApi(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.DeleteApiParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/admin/apis/${param0}`, {
    method: 'DELETE',
    params: {
      ...queryParams,
    },
    ...(options || {}),
  });
}
