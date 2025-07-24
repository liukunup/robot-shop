import type { RequestConfig, RequestOptions } from '@@/plugin-request/request';
import { message, notification } from 'antd';
import { getToken, setToken, getRefreshToken, removeToken } from './auth';
import { refreshToken } from '@/services/backend/user';

// 错误处理方案： 错误类型
enum ErrorShowType {
  SILENT = 0,        // 静默处理
  WARN_MESSAGE = 1,  // 警告提示
  ERROR_MESSAGE = 2, // 错误提示
  NOTIFICATION = 3,  // 通知用户
  REDIRECT = 9,      // 重定向跳转
}

// 与后端约定的响应数据格式
interface IResponse<T = any> {
  success: boolean;
  data: T;
  errorCode?: number;
  errorMessage?: string;
  errorShowType?: ErrorShowType;
}

// 请求队列项类型
interface RequestQueueItem {
  resolve: (value?: string | null) => void;
  reject: (reason?: any) => void;
}

// Token 刷新状态管理
let isRefreshing = false;
let failedQueue: RequestQueueItem[] = [];

/**
 * 处理 accessToken 刷新逻辑
 * @returns 新的 accessToken 值
 */
const handleTokenRefresh = async (): Promise<string> => {
  try {
    // 从浏览器获取 refreshToken 值
    const refreshTokenValue = getRefreshToken();
    if (!refreshTokenValue) {
      throw new Error('No refreshToken available');
    }
    // 调用 accessToken 刷新接口
    const response = await refreshToken({ refreshToken: refreshTokenValue });
    if (response.success) {
      setToken(response.data.accessToken, response.data.refreshToken);
      return response.data.accessToken;
    }
    // 刷新失败
    throw new Error(response.errorMessage || '刷新 accessToken 失败');
  } catch (error) {
    removeToken();
    window.location.href = '/user/login';
    throw error;
  }
};

/**
 * 处理请求队列
 * @param error 错误对象
 * @param token 新的 accessToken 值 或 null
 */
const processQueue = (error: any, token: string | null = null): void => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

/**
 * 处理业务错误
 * @param errorInfo 错误信息
 */
const handleBizError = (errorInfo: IResponse): void => {
  const { errorCode, errorMessage, errorShowType } = errorInfo;

  // 按错误码处理
  const errorCodeHandlers: Record<number, () => void> = {
    400: () => notification.error({ message: '请求错误', description: errorMessage }),
    401: () => {}, // 401 在拦截器中特殊处理
    403: () => notification.error({ message: '权限不足', description: '请联系管理员' }),
    404: () => notification.error({ message: '资源不存在', description: '请联系管理员' }),
    500: () => notification.error({ message: '服务器错误', description: '请联系管理员' }),
  };

  if (errorCode && errorCodeHandlers[errorCode]) {
    errorCodeHandlers[errorCode]();
    return;
  }

  // 按错误展示类型处理
  switch (errorShowType) {
    case ErrorShowType.SILENT:
      break;
    case ErrorShowType.WARN_MESSAGE:
      message.warning(errorMessage);
      break;
    case ErrorShowType.ERROR_MESSAGE:
      message.error(errorMessage);
      break;
    case ErrorShowType.NOTIFICATION:
      notification.open({ description: errorMessage, message: errorCode?.toString() });
      break;
    case ErrorShowType.REDIRECT:
      // TODO: 实现重定向逻辑
      break;
    default:
      message.error(errorMessage || '未知错误');
  }
};

/**
 * 处理 401 未授权错误
 * @param originalRequest 原始请求对象
 */
const handleUnauthorizedError = async (originalRequest: any): Promise<any> => {
  if (!isRefreshing) {
    isRefreshing = true;
    try {
      // 刷新 accessToken
      const newAccessToken = await handleTokenRefresh();
      processQueue(null, newAccessToken);
      // 重试原始请求
      originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
      return fetch(originalRequest);
    } catch (err) {
      processQueue(err, null);
      return Promise.reject(err);
    } finally {
      isRefreshing = false;
    }
  }
  // 如果正在刷新 accessToken 则先将请求加入队列
  return new Promise((resolve, reject) => {
    failedQueue.push({ resolve, reject });
  });
};

/**
 * @name 错误处理
 * pro 自带的错误处理， 可以在这里做自己的改动
 * @doc https://umijs.org/docs/max/request#配置
 */
export const errorConfig: RequestConfig = {

  // 错误处理：umi@3 的错误处理方案
  errorConfig: {
    // 错误抛出
    errorThrower: (response: IResponse) => {
      const { success, data, errorCode, errorMessage, errorShowType } = response;
      if (!success) {
        const error: any = new Error(errorMessage);
        error.name = 'BizError';
        error.info = { errorCode, errorMessage, errorShowType, data };
        throw error; // 抛出自定义错误
      }
    },

    // 错误接收及处理
    errorHandler: async (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;

      if (error.name === 'BizError') { // 业务错误
        const errorInfo: IResponse | undefined = error.info;
        if (errorInfo) {
          if (errorInfo.errorCode === 401) {
            const originalRequest = error.config;
            handleUnauthorizedError(originalRequest);
          } else {
            handleBizError(errorInfo);
          }
        }
      } else if (error.response) { // 响应错误
        // Axios 的错误
        // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
        if (error.response.status === 401) {
          const originalRequest = error.config;
          handleUnauthorizedError(originalRequest);
        } else {
          message.error(`响应错误码: ${error.response.status}`);
        }
      } else if (error.request) { // 请求错误
        // 请求已经成功发起，但没有收到响应
        // \`error.request\` 在浏览器中是 XMLHttpRequest 的实例，
        // 而在node.js中是 http.ClientRequest 的实例
        message.error('无响应，请重试');
      } else { // 未知错误
        // 发送请求时出了点问题
        message.error('请求错误，请重试');
      }
    },
  },

  // 请求拦截器
  requestInterceptors: [
    (config: RequestOptions) => {
      // 在浏览器环境中添加 token
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

  // 响应拦截器
  responseInterceptors: [
  ]
};
