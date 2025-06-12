// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 账号登录 POST /login */
export async function postLogin(body: API.LoginRequest, options?: { [key: string]: any }) {
  return request<API.LoginResponse>('/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户注册 目前只支持邮箱登录 POST /register */
export async function postRegister(body: API.RegisterRequest, options?: { [key: string]: any }) {
  return request<API.Response>('/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取用户信息 GET /user */
export async function getUser(options?: { [key: string]: any }) {
  return request<API.GetProfileResponse>('/user', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 修改用户信息 PUT /user */
export async function putUser(body: API.UpdateProfileRequest, options?: { [key: string]: any }) {
  return request<API.Response>('/user', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
