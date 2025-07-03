import { Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState } from 'react';
import { createMenu } from '@/services/backend/menu';

interface CreateFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
}

const CreateForm = ({ visible, onCancel, onSuccess }: CreateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.MenuRequest>();
  const intl = useIntl();

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      await createMenu(values as API.MenuRequest);
      message.success(intl.formatMessage({ id: 'pages.common.new.success', defaultMessage: '新建成功' }));
      form.resetFields();
      onSuccess();
    } catch (error) {
      const msg = intl.formatMessage({ id: 'pages.common.new.failure', defaultMessage: '新建失败' });
      if (error instanceof Error) {
        message.error(error.message || msg);
      } else {
        message.error(msg);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={<FormattedMessage id="pages.admin.menu.modal.createForm.title" defaultMessage="新建菜单" />}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
    >
      <Form
        form={form}
        layout="vertical"
        className="create-menu-form"
      >
        <Form.Item
          name="parentId"
          label={<FormattedMessage id="pages.admin.menu.key.parent" defaultMessage="父级菜单" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="path"
          label={<FormattedMessage id="pages.admin.menu.key.path" defaultMessage="路径" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.menu.form.path.required', defaultMessage: '路径不能为空'}) },
          ]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="redirect"
          label={<FormattedMessage id="pages.admin.menu.key.redirect" defaultMessage="重定向" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="component"
          label={<FormattedMessage id="pages.admin.menu.key.component" defaultMessage="组件" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="name"
          label={<FormattedMessage id="pages.admin.menu.key.name" defaultMessage="名称" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="icon"
          label={<FormattedMessage id="pages.admin.menu.key.icon" defaultMessage="图标" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="access"
          label={<FormattedMessage id="pages.admin.menu.key.access" defaultMessage="权限" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="weight"
          label={<FormattedMessage id="pages.admin.menu.key.weight" defaultMessage="权重" />}
        >
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateForm;
