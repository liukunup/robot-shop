import { PlusOutlined } from '@ant-design/icons';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useState } from 'react';
import { listUsers, userDelete } from '@/services/backend/user';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const User: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentUser, setCurrentUser] = useState<API.User | null>(null);
  const actionRef = useRef<ActionType>(null);
  const intl = useIntl();

  const columns: ProColumns<API.User>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.username',
        defaultMessage: '用户名',
      }),
      dataIndex: 'username',
      copyable: true,
      ellipsis: true,
      formItemProps: {
        rules: [
          {
            required: true,
            message: 'This field is required',
          },
        ],
      },
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.nickname',
        defaultMessage: '昵称',
      }),
      dataIndex: 'nickname',
      copyable: true,
      ellipsis: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.email',
        defaultMessage: '邮箱',
      }),
      dataIndex: 'email',
      formItemProps: {
        rules: [
          {
            required: true,
            message: 'This field is required',
          },
        ],
      },
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.phone',
        defaultMessage: '手机',
      }),
      dataIndex: 'phone',
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.status',
        defaultMessage: '状态',
      }),
      dataIndex: 'status',
      search: false,
      filters: true,
      onFilter: true,
      valueType: 'select',
      valueEnum: {
        0: { text: '待激活', status: 'NotActive' },
        1: { text: '正常', status: 'Normal' },
        2: { text: '禁用', status: 'Disabled' },
      },
      render: (_, record) => (
        <Space>
          {record.status === 0 ? (
            <Tag color="gold">待激活</Tag>
          ) : record.status === 1 ? (
            <Tag color="green">正常</Tag>
          ) : (
            <Tag color="red">禁用</Tag>
          )}
        </Space>
      )
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.roles',
        defaultMessage: '角色',
      }),
      dataIndex: 'roles',
      hideInSearch: true,
      render: (roles) => (
        <Space>
          {roles?.map(role => (
            <Tag key={role} color="blue">{role}</Tag>
          )) || null}
        </Space>
      ),
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.createdAt',
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
        id: 'pages.admin.user.key.updatedAt',
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
        id: 'pages.robot.table.column.actions',
        defaultMessage: '操作',
      }),
      valueType: 'option',
      key: 'option',
      render: (text, record, _, action) => [
        <a
          key="edit"
          onClick={() => {
            console.log(record);
            setCurrentUser(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.admin.user.table.action.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="remove"
          onClick={async () => {
            if (record.id) {
              await userDelete({ id: record.id });
              message.success('删除成功');
              action?.reload();
            }
          }}
        >
          <FormattedMessage id="pages.admin.user.table.action.remove" defaultMessage="删除" />
        </a>,
      ],
    },
  ];

  const search = async (params: {
    page: number;
    pageSize: number;
    username?: string;
    nickname?: string;
    email?: string;
    phone?: string;
  }) => {
    try {
      const result = await listUsers(params as API.ListUsersParams);
      return { data: result.data?.list || [], success: result.success, total: result.data?.total };
    } catch (error) {
      message.error('获取列表失败');
      return { data: [], success: false, total: 0 };
    }
  };

  return (
    <div>
      <ProTable<API.User>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(params, sort, filter);
          const { current = 1, pageSize = 20, username, nickname, email, phone } = params;
          const results = await search({
            page: current,
            pageSize,
            username,
            nickname,
            email,
            phone,
          });
          return results;
        }}
        editable={{
          type: 'multiple',
        }}
        columnsState={{
          persistenceKey: 'pro-table-user',
          persistenceType: 'localStorage',
          defaultValue: {
            option: { fixed: 'right', disable: true },
          },
          onChange(value) {
            console.log('value: ', value);
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
          id: 'pages.admin.user.table.title',
          defaultMessage: '用户列表',
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
            <FormattedMessage id="pages.admin.user.table.toolbar.newUser" defaultMessage="新增" />
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
        initialValues={currentUser as API.User}
      />
    </div>
  );
};

export default User;
