import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';

interface IErrorResponse {
  response?: {
    status?: number;
    data?: any;
    [key: string]: any;
  };
  [key: string]: any;
}

export const authConfig: RequestConfig = {
  errorConfig: {
    // 错误接收及处理
    errorHandler: (error: IErrorResponse) => {
      const status = error?.response?.status;
      const errorData = error?.response?.data;

      if (status === 401) {
        notification.error({
          message: '登录过期',
          description: '请重新登录',
          duration: 2,
        });
        
        // Add cleanup for other auth-related storage if needed
        if (typeof window !== 'undefined') {
          localStorage.removeItem('accessToken');
          // Consider using a more robust redirect with return URL
          setTimeout(() => {
            window.location.href = '/user/login';
          }, 2000);
        }
      } else if (status >= 500) {
        notification.error({
          message: '服务器错误',
          description: '请稍后再试或联系管理员',
        });
      } else if (errorData?.message) {
        message.error(errorData.message);
      }
    },
  },

  // 请求拦截器 (在这里处理 Header 携带 Token 的问题)
  requestInterceptors: [
    (config: RequestOptions) => {
      if (typeof window !== 'undefined') {
        const token = localStorage.getItem('accessToken');
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
