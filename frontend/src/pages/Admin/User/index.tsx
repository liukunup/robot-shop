import { PlusOutlined } from '@ant-design/icons';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useEffect, useState } from 'react';
import { listUsers, deleteUser } from '@/services/backend/user';
import { listRoles } from '@/services/backend/role';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const User: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentUser, setCurrentUser] = useState<API.User | null>(null);
  const [roleOptions, setRoleOptions] = useState<API.Role[]>([]);
  const actionRef = useRef<ActionType>(null);
  const intl = useIntl();

  useEffect(() => {
    const fetchRoles = async () => {
      try {
        const response = await listRoles({});
        if (response.success) {
          setRoleOptions(response.data?.list || []);
        }
      } catch (error) {
        const msg = intl.formatMessage({ id: 'pages.admin.user.fetchRoles.failure', defaultMessage: '获取角色列表失败' });
        if (error instanceof Error) {
          message.error(error.message || msg);
        } else {
          message.error(msg);
        }
      }
    };
    fetchRoles();
  }, []);

  const columns: ProColumns<API.User>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.username',
        defaultMessage: '用户名'
      }),
      dataIndex: 'username',
      ellipsis: true,
      formItemProps: {
        rules: [
          {
            required: true,
            message: intl.formatMessage({
              id: 'pages.admin.user.form.username.required',
              defaultMessage: '用户名不能为空',
            }),
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
      ellipsis: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.user.key.email',
        defaultMessage: '邮箱',
      }),
      dataIndex: 'email',
      copyable: true,
      formItemProps: {
        rules: [
          {
            required: true,
            message: intl.formatMessage({
              id: 'pages.admin.user.form.email.required',
              defaultMessage: '邮箱不能为空',
            }),
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
      copyable: true,
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
        0: {
          text: intl.formatMessage({
            id: 'pages.admin.user.status.inactive',
            defaultMessage: '待激活'
          }),
          status: 'Inactive'
        },
        1: {
          text: intl.formatMessage({
            id: 'pages.admin.user.status.normal',
            defaultMessage: '正常'
          }),
          status: 'Normal'
        },
        2: {
          text: intl.formatMessage({
            id: 'pages.admin.user.status.disabled',
            defaultMessage: '禁用'
          }),
          status: 'Disabled'
        },
      },
      render: (_, record) => (
        <Space>
          {record.status === 0 ? (
            <Tag color="gold">
              <FormattedMessage id="pages.admin.user.status.inactive" defaultMessage="待激活" />
            </Tag>
          ) : record.status === 1 ? (
            <Tag color="green">
              <FormattedMessage id="pages.admin.user.status.normal" defaultMessage="正常" />
            </Tag>
          ) : (
            <Tag color="red">
              <FormattedMessage id="pages.admin.user.status.disabled" defaultMessage="禁用" />
            </Tag>
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
      ellipsis: true,
      hideInSearch: true,
      filters: roleOptions?.map(({ id, name }) => ({ text: name as string, value: id as number })) || [],
      onFilter: (value, record) => record.roles?.some(({ id }) => id === value) ?? false,
      render: (_, record) => (
        <Space>
          {record.roles?.map((r) => (
            <Tag key={r.id} color="blue">{r.name}</Tag>
          ))}
        </Space>
      ),
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
            setCurrentUser(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.common.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="remove"
          onClick={async () => {
            if (record.id) {
              await deleteUser({ id: record.id });
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
    username?: string;
    nickname?: string;
    email?: string;
    phone?: string;
  }) => {
    try {
      const result = await listUsers(params as API.ListUsersParams);
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
        initialValues={currentUser as API.User}
      />
    </div>
  );
};

export default User;
