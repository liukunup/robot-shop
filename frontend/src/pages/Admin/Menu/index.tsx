import { PlusOutlined } from '@ant-design/icons';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useState } from 'react';
import { listMenus, deleteMenu } from '@/services/backend/menu';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const Menu: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentMenu, setCurrentMenu] = useState<API.Menu | null>(null);
  const actionRef = useRef<ActionType>(null);
  const intl = useIntl();

  const columns: ProColumns<API.Menu>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.menu.key.path',
        defaultMessage: '路径',
      }),
      dataIndex: 'path',
      ellipsis: true,
      formItemProps: {
        rules: [
          {
            required: true,
            message: intl.formatMessage({
              id: 'pages.admin.menu.form.path.required',
              defaultMessage: '路径不能为空',
            }),
          },
        ],
      },
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.menu.key.redirect',
        defaultMessage: '重定向',
      }),
      dataIndex: 'redirect',
      ellipsis: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.menu.key.component',
        defaultMessage: '组件',
      }),
      dataIndex: 'component',
      ellipsis: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.menu.key.name',
        defaultMessage: '名称',
      }),
      dataIndex: 'name',
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.menu.key.icon',
        defaultMessage: '图标',
      }),
      dataIndex: 'icon',
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.menu.key.access',
        defaultMessage: '权限',
      }),
      dataIndex: 'access',
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.menu.key.weight',
        defaultMessage: '权重',
      }),
      dataIndex: 'weight',
    },
    {
      title: intl.formatMessage({
        id: 'pages.common.key.createdAt',
        defaultMessage: '创建时间',
      }),
      key: 'createdAt',
      dataIndex: 'createdAt',
      valueType: 'dateTime',
      sorter: true,
      hideInSearch: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.common.key.updatedAt',
        defaultMessage: '更新时间',
      }),
      key: 'updatedAt',
      dataIndex: 'updatedAt',
      valueType: 'dateTime',
      sorter: true,
      hideInSearch: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.common.table.key.actions',
        defaultMessage: '操作',
      }),
      valueType: 'option',
      key: 'option',
      render: (text, record, _, action) => [
        <a
          key="edit"
          onClick={() => {
            setCurrentMenu(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.common.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="remove"
          onClick={async () => {
            if (record.id) {
              await deleteMenu({ id: record.id });
              message.success(intl.formatMessage({
                id: 'pages.common.remove.success',
                defaultMessage: '删除成功',
              }));
              action?.reload();
            }
          }}
        >
          <FormattedMessage id="pages.common.remove" defaultMessage="删除" />
        </a>,
      ],
    },
  ];

  const search = async (params: {
    page: number;
    pageSize: number;
  }) => {
    try {
      const result = await listMenus(params as API.ListMenusParams);
      return { data: result.data?.list || [], success: result.success, total: result.data?.total };
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.common.fetchData.failure',
        defaultMessage: '获取数据失败',
      }));
      return { data: [], success: false, total: 0 };
    }
  };

  return (
    <div>
      <ProTable<API.Role>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(params, sort, filter);
          const { current = 1, pageSize = 20 } = params;
          const results = await search({
            page: current,
            pageSize,
          });
          return results;
        }}
        editable={{
          type: 'multiple',
        }}
        columnsState={{
          persistenceKey: 'pro-table-menu',
          persistenceType: 'localStorage',
          defaultValue: {
            option: { fixed: 'right', disable: true },
          },
        }}
        rowKey="id"
        search={{
          labelWidth: 'auto',
        }}
        options={{
          setting: {
            listsHeight: 400,
          },
        }}
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
        }}
        dateFormatter="string"
        headerTitle={intl.formatMessage({
          id: 'pages.admin.menu.table.title',
          defaultMessage: '菜单列表',
        })}
        toolBarRender={() => [
          <Button
            key="button"
            icon={<PlusOutlined />}
            onClick={() => {
              setCreateVisible(true);
            }}
            type="primary"
          >
            <FormattedMessage id="pages.common.new" defaultMessage="新建" />
          </Button>,
        ]}
      />
      <CreateForm
        visible={createVisible}
        onCancel={() => setCreateVisible(false)}
        onSuccess={() => {
          setCreateVisible(false);
          actionRef.current?.reload();
        }}
      />
      <UpdateForm
        visible={updateVisible}
        onCancel={() => setUpdateVisible(false)}
        onSuccess={() => {
          setUpdateVisible(false);
          actionRef.current?.reload();
        }}
        initialValues={currentMenu as API.Menu}
      />
    </div>
  );
};

export default Menu;
