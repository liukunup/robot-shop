import { EllipsisOutlined, PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { ProTable, TableDropdown } from '@ant-design/pro-components';
import { FormattedMessage, history, useIntl } from '@umijs/max';
import { Button, Dropdown, message } from 'antd';
import { useRef, useState } from 'react';
import { getRobots, deleteRobotById } from '../../services/backend/robot';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

type RobotItem = {
  id: number;
  name: string;
  desc: string;
  webhook: string;
  callback: string;
  enabled: boolean;
  owner: string;
  createdAt: string;
  updatedAt: string;
  url?: string; // 可选字段
  // 或直接使用id跳转
  // url: `/robots/${id}`
};

const fetchRobots = async (params: {
  page: number;
  pageSize: number;
  name?: string;}) => {
  try {
    const { page, pageSize, name } = params;
    const result = await getRobots({ page, pageSize, name });
    return { data: result.data?.list || [], success: result.success, total: result.data?.total };
  } catch (error) {
    message.error('获取机器人列表失败');
    return { data: [], success: false, total: 0 };
}
};

const Robot: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentRobot, setCurrentRobot] = useState<RobotItem | null>(null);
  const actionRef = useRef<ActionType>();
  const intl = useIntl();

  const columns: ProColumns<RobotItem>[] = [
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
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.desc',
        defaultMessage: '描述',
      }),
      dataIndex: 'desc',
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.webhook',
        defaultMessage: 'Webhook',
      }),
      dataIndex: 'webhook',
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.callback',
        defaultMessage: 'Callback',
      }),
      dataIndex: 'callback',
    },
    {
      title: intl.formatMessage({
        id: 'pages.robot.table.column.enabled',
        defaultMessage: '启用状态',
      }),
      dataIndex: 'enabled',
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
      valueType: 'date',
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
      valueType: 'date',
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
          key="editable"
          onClick={() => {
            setCurrentRobot(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.robot.table.column.action.edit" defaultMessage="编辑" />
        </a>,
        <a href={record.url || `/robots/${record.id}`} target="_blank" rel="noopener noreferrer" key="view">
          <FormattedMessage id="pages.robot.table.column.action.view" defaultMessage="查看" />
        </a>,
        <TableDropdown
          key="actionGroup"
          onSelect={() => action?.reload()}
          menus={[
            { key: 'copy', name: 'Copy' },
            {
              key: 'delete',
              name: 'Delete',
              onClick: async () => {
                await deleteRobotById({ id: record.id });
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
      <ProTable<RobotItem>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(sort, filter);
          const { current = 1, pageSize = 10, keyword } = params;
          const resp = await fetchRobots({
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
        headerTitle='robot.table.title'
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
        initialValues={currentRobot || {}}
      />
    </div>
  );
};

export default Robot;
