import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';

export const authConfig: RequestConfig = {
  errorConfig: {
    errorHandler: (error: any) => {
      if (error?.response?.status === 401) {
        notification.error({
          message: '登录过期',
          description: '请重新登录',
        });
        localStorage.removeItem('accessToken');
        window.location.href = '/user/login';
      }
    },
    errorThrower: (error: any) => {
      throw error;
    },
  },
  requestInterceptors: [
    (config: RequestOptions) => {
      const token = localStorage.getItem('accessToken');
      if (token) {
        config.headers = {
          ...config.headers,
          Authorization: `Bearer ${token}`,
        };
      }
      return config;
    },
  ],
  responseInterceptors: [
    (response: any) => {
      return response;
    },
  ],
}
