import { AvatarDropdown, AvatarName, Footer, Question, SelectLang } from '@/components';
import { fetchCurrentUser } from '@/services/backend/user';
import { LinkOutlined, SmileOutlined, CrownOutlined, RobotOutlined } from '@ant-design/icons';
import type { Settings as LayoutSettings, MenuDataItem } from '@ant-design/pro-components';
import { SettingDrawer } from '@ant-design/pro-components';
import type { RunTimeLayoutConfig } from '@umijs/max';
import { history, Link } from '@umijs/max';
import React from 'react';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './utils/request';
import '@ant-design/v5-patch-for-react-19';
import { fetchDynamicMenu } from '@/services/backend/user';

const isDev = process.env.NODE_ENV === 'development';
const loginPath = '/user/login';

// 解决动态菜单图标问题
interface IconMapType {
  [key: string]: React.ReactNode;
}
const IconMap: IconMapType = {
  smile: <SmileOutlined />,
  crown: <CrownOutlined />,
  robot: <RobotOutlined />,
};
const loopMenuItem = (menus: API.MenuNode[]): MenuDataItem[] =>
  menus.map(({ icon, children, parentKeys, ...item }) => ({
    ...item,
    icon: icon && IconMap[icon],
    children: children && loopMenuItem(children),
    parentKeys: typeof parentKeys === 'string' ? parentKeys.split(',') : [],
  })
);

/**
 * @see https://umijs.org/docs/api/runtime-config#getinitialstate
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: API.User;
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.User | undefined>;
  fetchMenuData?: () => Promise<MenuDataItem[]>;
  menuData?: MenuDataItem[];
}> {
  const fetchUserInfo = async () => {
    try {
      const response = await fetchCurrentUser({
        skipErrorHandler: true,
      });
      return response.data;
    } catch (error) {
      history.push(loginPath);
    }
    return undefined;
  };
  const fetchMenuData = async () => {
    try {
      const response = await fetchDynamicMenu({
        skipErrorHandler: true,
      });
      if (response.success) {
        return loopMenuItem(response.data?.list || []);
      }
    } catch (error) {
      console.error('failed to fetch menu data:', error);
    }
    return [];
  };

  // 如果不是登录页面，执行
  const { location } = history;
  if (location.pathname !== loginPath) {
    const currentUser = await fetchUserInfo();
    const menuData = await fetchMenuData();
    return {
      fetchUserInfo,
      fetchMenuData,
      currentUser,
      menuData,
      settings: defaultSettings as Partial<LayoutSettings>,
    };
  }
  return {
    fetchUserInfo,
    menuData: [],
    settings: defaultSettings as Partial<LayoutSettings>,
  };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {
  return {
    actionsRender: () => [<Question key="doc" />, <SelectLang key="SelectLang" />],
    avatarProps: {
      src: initialState?.currentUser?.avatar,
      title: <AvatarName />,
      render: (_, avatarChildren) => {
        return <AvatarDropdown>{avatarChildren}</AvatarDropdown>;
      },
    },
    waterMarkProps: {
      content: initialState?.currentUser?.nickname,
    },
    footerRender: () => <Footer />,
    onPageChange: () => {
      const { location } = history;
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }
    },
    bgLayoutImgList: [
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/D2LWSqNny4sAAAAAAAAAAAAAFl94AQBr',
        left: 85,
        bottom: 100,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/C2TWRpJpiC0AAAAAAAAAAAAAFl94AQBr',
        bottom: -68,
        right: -45,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/F6vSTbj8KpYAAAAAAAAAAAAAFl94AQBr',
        bottom: 0,
        left: 0,
        width: '331px',
      },
    ],
    links: isDev
      ? [
          <Link key="openapi" to="/umi/plugin/openapi" target="_blank">
            <LinkOutlined />
            <span>OpenAPI 文档</span>
          </Link>,
        ]
      : [],
    menuHeaderRender: undefined,
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    // 增加一个 loading 的状态
    childrenRender: (children) => {
      // if (initialState?.loading) return <PageLoading />;
      return (
        <>
          {children}
          {isDev && (
            <SettingDrawer
              disableUrlParams
              enableDarkTheme
              settings={initialState?.settings}
              onSettingChange={(settings) => {
                setInitialState((preInitialState) => ({
                  ...preInitialState,
                  settings,
                }));
              }}
            />
          )}
        </>
      );
    },
    // 实现动态菜单功能
    menuDataRender: () => initialState.menuData,
    ...initialState?.settings,
  };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request = {
  ...errorConfig,
};
