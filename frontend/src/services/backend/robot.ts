// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取机器人列表 GET /robots */
export async function getRobots(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRobotsParams,
  options?: { [key: string]: any },
) {
  return request<API.ListRobotResponse>(`/v1/robots`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建机器人 POST /robots */
export async function postRobots(body: API.CreateRobotParams, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/robots`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取单个机器人的数据 GET /robots/${param0} */
export async function getRobotById(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRobotByIdParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.GetRobotResponse>(`/v1/robots/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新机器人 PUT /robots/${param0} */
export async function updateRobotById(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.updateRobotByIdParams,
  body: API.UpdateRobotParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/robots/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除机器人 DELETE /robots/${param0} */
export async function deleteRobotById(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteRobotByIdParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/robots/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}
