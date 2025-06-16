// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取机器人列表 GET /robot */
export async function getRobot(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRobotParams,
  options?: { [key: string]: any },
) {
  return request<API.PageResponseBackendApiV1RobotResponseData>(`/v1/robot`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建机器人 POST /robot */
export async function postRobot(body: API.RobotRequest, options?: { [key: string]: any }) {
  return request<API.RobotResponseData>(`/v1/robot`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取机器人 GET /robot/${param0} */
export async function getRobotId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRobotIdParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.RobotResponseData>(`/v1/robot/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新机器人 PUT /robot/${param0} */
export async function putRobotId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.putRobotIdParams,
  body: API.RobotRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.RobotResponseData>(`/v1/robot/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除机器人 DELETE /robot/${param0} */
export async function deleteRobotId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteRobotIdParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(`/v1/robot/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}
