// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 批量搜索机器人 GET /robots */
export async function listRobots(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.ListRobotsParams,
  options?: { [key: string]: any },
) {
  return request<API.RobotSearchResponse>(`/v1/robots`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建机器人 POST /robots */
export async function createRobot(body: API.RobotRequest, options?: { [key: string]: any }) {
  return request<API.Response>(`/v1/robots`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 查找机器人 GET /robots/${param0} */
export async function getRobot(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.GetRobotParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.RobotResponse>(`/v1/robots/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新机器人 PUT /robots/${param0} */
export async function updateRobot(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.UpdateRobotParams,
  body: API.RobotRequest,
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
export async function deleteRobot(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.DeleteRobotParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/robots/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}
