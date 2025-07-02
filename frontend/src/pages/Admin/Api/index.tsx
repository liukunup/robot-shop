import { PlusOutlined } from '@ant-design/icons';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useRef, useState } from 'react';
import { listApis, deleteApi } from '@/services/backend/api';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const Api: React.FC = () => {
  const [createVisible, setCreateVisible] = useState(false);
  const [updateVisible, setUpdateVisible] = useState(false);
  const [currentApi, setCurrentApi] = useState<API.Api | null>(null);
  const actionRef = useRef<ActionType>(null);
  const intl = useIntl();

  const columns: ProColumns<API.Api>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: intl.formatMessage({
        id: 'pages.admin.api.key.group',
        defaultMessage: '分组',
      }),
      dataIndex: 'group',
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
        id: 'pages.admin.api.key.name',
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
        id: 'pages.admin.api.key.path',
        defaultMessage: '路径',
      }),
      dataIndex: 'path',
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
        id: 'pages.admin.api.key.method',
        defaultMessage: '方法',
      }),
      dataIndex: 'method',
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
            setCurrentApi(record);
            setUpdateVisible(true);
          }}
        >
          <FormattedMessage id="pages.common.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="remove"
          onClick={async () => {
            if (record.id) {
              await deleteApi({ id: record.id });
              const msg = intl.formatMessage({
                id: 'pages.common.remove.success',
                defaultMessage: '删除成功',
              });
              message.success(msg);
              action?.reload();
            }
          }}
        >
          <FormattedMessage id="pages.common.remove" defaultMessage="删除" />
        </a>,
      ],
    },
  ];

  const searchApis = async (params: {
    page: number;
    pageSize: number;
    group?: string;
    name?: string;
    path?: string;
    method?: string;
  }) => {
    try {
      const response = await listApis(params as API.ListApisParams);
      return { data: response.data?.list || [], success: response.success, total: response.data?.total };
    } catch (error) {
      const msg = intl.formatMessage({
        id: 'pages.common.fetchList.error',
        defaultMessage: '获取列表失败',
      });
      message.error(msg);
      return { data: [], success: false, total: 0 };
    }
  };

  return (
    <div>
      <ProTable<API.Api>
        columns={columns}
        actionRef={actionRef}
        cardBordered
        request={async (params, sort, filter) => {
          console.log(params, sort, filter);
          const { current = 1, pageSize = 20, group, name, path, method } = params;
          const results = await searchApis({
            page: current,
            pageSize,
            group,
            name,
            path,
            method,
          });
          return results;
        }}
        editable={{
          type: 'multiple',
        }}
        columnsState={{
          persistenceKey: 'pro-table-api',
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
          id: 'pages.admin.api.table.title',
          defaultMessage: '接口列表',
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
            <FormattedMessage id="pages.common.new" defaultMessage="新增" />
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
        initialValues={currentApi as API.Api}
      />
    </div>
  );
};

export default Api;
