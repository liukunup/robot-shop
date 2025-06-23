import { EllipsisOutlined, PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { ProTable, TableDropdown } from '@ant-design/pro-components';
import { FormattedMessage, useIntl } from '@umijs/max';
import { Button, Dropdown, Space, Tag, message } from 'antd';
import { useRef, useState } from 'react';
import { listRobots, deleteRobot, createRobot } from '../../services/backend/robot';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import ViewForm from './components/ViewForm';

type RobotData = {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
  desc: string;
  webhook: string;
  callback: string;
  enabled: boolean;
  owner: string;
};

const searchRobots = async (params: {
  page: number;
  pageSize: number;
  name?: string;}) => {
  try {
    const result = await listRobots(params);
    return { data: result.data?.list || [], success: result.success, total: result.data?.total };
  } catch (error) {
    message.error('获取机器人列表失败');
    return { data: [], success: false, total: 0 };
}
};

const Robot: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentRobot, setCurrentRobot] = useState<RobotData | null>(null);
  const actionRef = useRef<ActionType>();
  const intl = useIntl();

  const columns: ProColumns<RobotData>[] = [
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
      tooltip: '描述机器人的主要功能',
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.webhook',
        defaultMessage: 'Webhook',
      }),
      dataIndex: 'webhook',
      ellipsis: true,
      tooltip: 'Webhook 用于接收来自外部的消息',
      hideInForm: true,
      hideInSearch: true,
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.callback',
        defaultMessage: 'Callback',
      }),
      dataIndex: 'callback',
      ellipsis: true,
      tooltip: '回调地址用于向机器人发送消息',
      hideInForm: true,
      hideInSearch: true,
    },
    {
      disable: true,
      title: intl.formatMessage({
        id: 'pages.robot.table.column.enabled',
        defaultMessage: '启用状态',
      }),
      dataIndex: 'enabled',
      search: false,
      filters: true,
      onFilter: true,
      ellipsis: true,
      valueType: 'select',
      valueEnum: {
        true: { text: '启用', status: 'Success' },
        false: { text: '禁用', status: 'Error' },
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
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.robot.table.column.action.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="view"
          onClick={() => {
            setCurrentRobot(record);
            setUpdateVisible(true);
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
                await deleteRobot({ id: record.id });
                message.success('删除成功');
                action?.reload();
              },
            },
          ]}
        />,
      ],
    },
  ];

  return (
    <div>
      <ProTable<RobotData>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(params, sort, filter);
          const { current = 1, pageSize = 10, keyword } = params;
          const resp = await searchRobots({
            page: current,
            pageSize: pageSize,
            name: keyword,
          });
          return resp;
        }}
        editable={{
          type: 'multiple',
        }}
        columnsState={{
          persistenceKey: 'pro-table-singe-demos',
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
        form={{
          // Since transform is configured, the submitted parameters are different from the defined ones, so they need to be transformed here
          syncToUrl: (values, type) => {
            if (type === 'get') {
              return {
                ...values,
                created_at: [values.startTime, values.endTime],
              };
            }
            return values;
          },
        }}
        pagination={{
          pageSize: 5,
          onChange: (page) => console.log(page),
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
              setCreateVisible(true);
            }}
            type="primary"
          >
            <FormattedMessage id="pages.robot.table.toolbar.new" defaultMessage="新增" />
          </Button>,
          <Dropdown
            key="menu"
            menu={{
              items: [
                {
                  label: '1st item',
                  key: '1',
                },
                {
                  label: '2nd item',
                  key: '2',
                },
                {
                  label: '3rd item',
                  key: '3',
                },
              ],
            }}
          >
            <Button>
              <EllipsisOutlined />
            </Button>
          </Dropdown>,
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
        initialValues={currentRobot}
      />
    </div>
  );
};

export default Robot;
