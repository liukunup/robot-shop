// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 账号登录 POST /login */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginResult>(`/v1/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户注册 目前只支持邮箱登录 POST /register */
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

/** 获取用户信息 GET /user */
export async function queryCurrentUser(options?: { [key: string]: any }) {
  return request<API.GetProfileResponse>(`/v1/user`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 修改用户信息 PUT /user */
export async function updateProfile(
  body: API.UpdateProfileParams,
  options?: { [key: string]: any },
) {
  return request<API.Response>(`/v1/user`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
