import { EllipsisOutlined, PlusOutlined } from '@ant-design/icons';
import { ProTable, TableDropdown } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Dropdown, Space, Tag, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useState } from 'react';
import { listRobots, deleteRobot, createRobot } from '../../services/backend/robot';
import OperateForm from './components/OperateForm';

// 搜索机器人
const searchRobots = async (params: {
  page: number;
  pageSize: number;
  name?: string;
  desc?: string;
  owner?: string;
}) => {
  try {
    const result = await listRobots(params as API.ListRobotsParams);
    return { data: result.data?.list || [], success: result.success, total: result.data?.total };
  } catch (error) {
    message.error('获取列表失败');
    return { data: [], success: false, total: 0 };
  }
};

const Robot: React.FC = () => {
  const [visible, setVisible] = useState(false);
  const [operation, setOperation] = useState<'create' | 'update' | 'view'>();
  const [currentRobot, setCurrentRobot] = useState<API.Robot | null>(null);
  const actionRef = useRef<ActionType>();
  const intl = useIntl();

  const columns: ProColumns<API.Robot>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.name',
        defaultMessage: '名称',
      }),
      dataIndex: 'name',
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
        id: 'pages.robot.table.column.desc',
        defaultMessage: '描述',
      }),
      dataIndex: 'desc',
      ellipsis: true,
      tooltip: '描述机器人的作用',
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.webhook',
        defaultMessage: 'Webhook',
      }),
      dataIndex: 'webhook',
      ellipsis: true,
      tooltip: '用于发送消息到外部服务器',
      hideInSearch: true,
      hideInTable: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.callback',
        defaultMessage: 'Callback',
      }),
      dataIndex: 'callback',
      ellipsis: true,
      tooltip: '用于接收来自外部服务器的消息',
      hideInSearch: true,
      hideInTable: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.enabled',
        defaultMessage: '是否启用',
      }),
      dataIndex: 'enabled',
      search: false,
      filters: true,
      onFilter: true,
      valueType: 'select',
      valueEnum: {
        true: { text: '启用', status: 'Enabled' },
        false: { text: '禁用', status: 'Disabled' },
      },
      render: (_, record) => (
        <Space>
          {record.enabled ? (
            <Tag color="green">启用</Tag>
          ) : (
            <Tag color="red">禁用</Tag>
          )}
        </Space>
      )
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.owner',
        defaultMessage: '所有者',
      }),
      dataIndex: 'owner',
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.createdAt',
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
        id: 'pages.robot.table.column.updatedAt',
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
            setCurrentRobot(record);
            setOperation('update');
            setVisible(true);
          }}
        >
          <FormattedMessage id="pages.robot.table.column.action.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="view"
          onClick={() => {
            setCurrentRobot(record);
            setOperation('view');
            setVisible(true);
          }}
        >
          <FormattedMessage id="pages.robot.table.column.action.view" defaultMessage="查看" />
        </a>,
        <TableDropdown
          key="actionGroup"
          onSelect={() => action?.reload()}
          menus={[
            {
              key: 'copy',
              name: '复制',
              onClick: async () => {
                await createRobot({
                  name: record.name + '-副本',
                  desc: record.desc,
                  webhook: record.webhook,
                  callback: record.callback,
                  enabled: record.enabled,
                  owner: record.owner,
                });
                actionRef.current?.reload();
              },
            },
            {
              key: 'delete',
              name: '删除',
              onClick: async () => {
                if (record.id) {
                  await deleteRobot({ id: record.id });
                  message.success('删除成功');
                  actionRef.current?.reload();
                }
              },
            },
          ]}
        />,
      ],
    },
  ];

  return (
    <div>
      <ProTable<API.Robot>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(params, sort, filter);
          const { current = 1, pageSize = 10, name, desc, owner } = params;
          const results = await searchRobots({
            page: current,
            pageSize,
            name,
            desc,
            owner,
          });
          return results;
        }}
        editable={{
          type: 'multiple',
        }}
        columnsState={{
          persistenceKey: 'pro-table-robot',
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
          pageSize: 10,
          showSizeChanger: true,
          showQuickJumper: true,
        }}
        dateFormatter="string"
        headerTitle={intl.formatMessage({
          id: 'pages.robot.table.title',
          defaultMessage: '机器人列表',
        })}
        toolBarRender={() => [
          <Button
            key="button"
            icon={<PlusOutlined />}
            onClick={() => {
              setOperation('create');
              setVisible(true);
            }}
            type="primary"
          >
            <FormattedMessage id="pages.robot.table.toolbar.new" defaultMessage="新增" />
          </Button>,
        ]}
      />
      <OperateForm
        visible={visible}
        operation={operation}
        onCancel={() => setVisible(false)}
        onSuccess={() => {
          setVisible(false);
          actionRef.current?.reload();
        }}
        initialValues={currentRobot}
      />
    </div>
  );
};

export default Robot;
