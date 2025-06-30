import { PlusOutlined } from '@ant-design/icons';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useState } from 'react';
import { listRoles, roleDelete } from '@/services/backend/role';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const searchRoles = async (params: {
  page: number;
  pageSize: number;
  name?: string;
  role?: string;
}) => {
  try {
    const result = await listRoles(params as API.ListRolesParams);
    return { data: result.data?.list || [], success: result.success, total: result.data?.total };
  } catch (error) {
    message.error('获取列表失败');
    return { data: [], success: false, total: 0 };
  }
};

const Role: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentRole, setCurrentRole] = useState<API.Role | null>(null);
  const actionRef = useRef<ActionType>(null);
  const intl = useIntl();

  const columns: ProColumns<API.Role>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.role.key.name',
        defaultMessage: '名称',
      }),
      dataIndex: 'name',
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
        id: 'pages.admin.role.key.role',
        defaultMessage: '标识',
      }),
      dataIndex: 'role',
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
        id: 'pages.admin.role.key.createdAt',
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
        id: 'pages.admin.role.key.updatedAt',
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
            setCurrentRole(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.admin.role.table.action.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="remove"
          onClick={async () => {
            if (record.id) {
              await roleDelete({ id: record.id });
              message.success('删除成功');
              action?.reload();
            }
          }}
        >
          <FormattedMessage id="pages.admin.role.table.action.remove" defaultMessage="删除" />
        </a>,
      ],
    },
  ];

  return (
    <div>
      <ProTable<API.Role>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(params, sort, filter);
          const { current = 1, pageSize = 20, name, role } = params;
          const results = await searchRoles({
            page: current,
            pageSize,
            name,
            role,
          });
          return results;
        }}
        editable={{
          type: 'multiple',
        }}
        columnsState={{
          persistenceKey: 'pro-table-role',
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
          id: 'pages.admin.role.table.title',
          defaultMessage: '角色列表',
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
        initialValues={currentRole as API.Role}
      />
    </div>
  );
};

export default Role;
