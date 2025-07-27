import { PlusOutlined } from '@ant-design/icons';
import { ProTable, TableDropdown } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useState } from 'react';
import { listRobots, deleteRobot, createRobot } from '../../services/backend/robot';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const Robot: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentRobot, setCurrentRobot] = useState<API.Robot | null>(null);
  const actionRef = useRef<ActionType>(null);
  const intl = useIntl();

  const columns: ProColumns<API.Robot>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.key.name',
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
        id: 'pages.robot.key.desc',
        defaultMessage: '描述',
      }),
      dataIndex: 'desc',
      ellipsis: true,
      tooltip: '描述机器人的作用',
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.key.webhook',
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
        id: 'pages.robot.key.callback',
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
        id: 'pages.robot.key.enabled',
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
        id: 'pages.robot.key.owner',
        defaultMessage: '所有者',
      }),
      dataIndex: 'owner',
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
            console.log(record);
            setCurrentRobot(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.common.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="duplicate"
          onClick={async () => {
            await createRobot({
              name: record.name + '-Copy',
              desc: record.desc,
              webhook: record.webhook,
              callback: record.callback,
              enabled: record.enabled,
              owner: record.owner,
            });
            actionRef.current?.reload();
          }}
        >
          <FormattedMessage id="pages.common.duplicate" defaultMessage="复制" />
        </a>,
        <TableDropdown
          key="actionGroup"
          onSelect={() => action?.reload()}
          menus={[
            {
              key: 'test',
              name: '测试',
              onClick: async () => {
                message.warning('此功能尚未实现');
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

  const search = async (params: {
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

  return (
    <div>
      <ProTable<API.Robot>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(params, sort, filter);
          const { current = 1, pageSize = 20, name, desc, owner } = params;
          const results = await search({
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
        initialValues={currentRobot as API.Robot}
      />
    </div>
  );
};

export default Robot;
