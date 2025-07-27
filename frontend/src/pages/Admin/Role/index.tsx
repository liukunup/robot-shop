import { PlusOutlined } from '@ant-design/icons';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useState } from 'react';
import { listRoles, deleteRole } from '@/services/backend/role';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

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
            message: intl.formatMessage({
              id: 'pages.admin.role.form.name.required',
              defaultMessage: '名称不能为空',
            }),
          },
        ],
      },
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.role.key.role',
        defaultMessage: '标识',
      }),
      dataIndex: 'casbinRole',
      ellipsis: true,
      formItemProps: {
        rules: [
          {
            required: true,
            message: intl.formatMessage({
              id: 'pages.admin.role.form.role.required',
              defaultMessage: '标识不能为空',
            }),
          },
        ],
      },
      render: (_, record) => (
        <Space>
          <Tag color="blue" style={{ marginRight: 4 }}>{record.casbinRole}</Tag>
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
            setCurrentRole(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.common.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="remove"
          onClick={async () => {
            if (record.id) {
              await deleteRole({ id: record.id });
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
    name?: string;
    casbinRole?: string;
  }) => {
    try {
      const response = await listRoles(params as API.ListRolesParams);
      return { data: response.data?.list || [], success: response.success, total: response.data?.total };
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
          const { current = 1, pageSize = 20, name, casbinRole } = params;
          const results = await search({
            page: current,
            pageSize,
            name,
            casbinRole,
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
        initialValues={currentRole as API.Role}
      />
    </div>
  );
};

export default Role;
