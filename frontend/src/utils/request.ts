﻿import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';
import { getToken } from './auth';

// 错误处理方案： 错误类型
enum ErrorShowType {
  SILENT = 0,
  WARN_MESSAGE = 1,
  ERROR_MESSAGE = 2,
  NOTIFICATION = 3,
  REDIRECT = 9,
}

// 与后端约定的响应数据格式
interface IResponse {
  success: boolean;
  data: any;
  errorCode?: number;
  errorMessage?: string;
  errorShowType?: ErrorShowType;
}

/**
 * @name 错误处理
 * pro 自带的错误处理， 可以在这里做自己的改动
 * @doc https://umijs.org/docs/max/request#配置
 */
export const errorConfig: RequestConfig = {
  // 错误处理： umi@3 的错误处理方案。
  errorConfig: {
    // 错误抛出
    errorThrower: (response) => {
      const { success, data, errorCode, errorMessage, errorShowType } =
      response as unknown as IResponse;
      if (!success) {
        const error: any = new Error(errorMessage);
        error.name = 'BizError';
        error.info = { errorCode, errorMessage, errorShowType, data };
        throw error; // 抛出自制的错误
      }
    },
    // 错误接收及处理
    errorHandler: (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;
      // 我们的 errorThrower 抛出的错误。
      if (error.name === 'BizError') {
        const errorInfo: IResponse | undefined = error.info;
        if (errorInfo) {
          const { errorCode, errorMessage, errorShowType } = errorInfo;
          // 按响应码处理
          if (errorCode === 400) {
            notification.error({
              message: '请求错误',
              description: errorMessage,
            });
          } else if (errorCode === 401) {
            notification.error({
              message: '登录过期',
              description: '请重新登录',
              duration: 2,
            });
            if (typeof window !== 'undefined') {
              localStorage.removeItem('accessToken');
              setTimeout(() => {
                window.location.href = '/user/login';
              }, 2000);
            }
          } else if (errorCode === 403) {
            notification.error({
              message: '权限不足',
              description: '请联系管理员',
            });
          } else if (errorCode === 404) {
            notification.error({
              message: '资源不存在',
              description: '请联系管理员',
            });
          } else if (errorCode === 500) {
            notification.error({
              message: '服务器错误',
              description: '请联系管理员',
            });
          } else {
            switch (errorShowType) {
              case ErrorShowType.SILENT:
                // do nothing
                break;
              case ErrorShowType.WARN_MESSAGE:
                message.warning(errorMessage);
                break;
              case ErrorShowType.ERROR_MESSAGE:
                message.error(errorMessage);
                break;
              case ErrorShowType.NOTIFICATION:
                notification.open({
                  description: errorMessage,
                  message: errorCode,
                });
                break;
              case ErrorShowType.REDIRECT:
                // TODO: redirect
                break;
              default:
                message.error(errorMessage);
            }
          }
        }
      } else if (error.response) {
        // Axios 的错误
        // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
        message.error(`Response status:${error.response.status}`);
      } else if (error.request) {
        // 请求已经成功发起，但没有收到响应
        // \`error.request\` 在浏览器中是 XMLHttpRequest 的实例，
        // 而在node.js中是 http.ClientRequest 的实例
        message.error('None response! Please retry.');
      } else {
        // 发送请求时出了点问题
        message.error('Request error, please retry.');
      }
    },
  },

  // 请求拦截器 (在这里处理 Header 携带 Token 的问题)
  requestInterceptors: [
    (config: RequestOptions) => {
      if (typeof window !== 'undefined') {
        const token = getToken();
        if (token) {
          config.headers = {
            ...config.headers,
            Authorization: `Bearer ${token}`,
          };
        }
      }
      return config;
    },
  ],
};
