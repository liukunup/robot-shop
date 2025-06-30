import { Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState } from 'react';
import { roleCreate } from '@/services/backend/role';

interface CreateFormProps {
  visible: boolean; // 弹窗是否可见
  onCancel: () => void; // 取消回调
  onSuccess: () => void; // 成功回调
}

const CreateForm = ({ visible, onCancel, onSuccess }: CreateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.RoleRequest>();
  const intl = useIntl();

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      const params = {
        ...values,
        status: 0, // 默认待激活
      };
      await roleCreate(params as API.RoleRequest);
      message.success(intl.formatMessage({ id: 'pages.admin.user.message.newUserSuccess', defaultMessage: '新增成功' }));
      form.resetFields();
      onSuccess();
    } catch (error) {
      const msg = intl.formatMessage({ id: 'pages.admin.user.message.newUserFailed', defaultMessage: '新增失败' });
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
      title={<FormattedMessage id="pages.admin.role.modal.createForm.title" defaultMessage="新建角色" />}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
      width={600}
    >
      <Form
        form={form}
        layout="vertical"
        className="create-role-form"
      >
        <Form.Item
          name="name"
          label={<FormattedMessage id="pages.admin.role.key.name" defaultMessage="名称" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.role.form.name.required', defaultMessage: '请输入名称' }) },
            { pattern: /^[a-zA-Z][a-zA-Z0-9]*$/, message: intl.formatMessage({ id: 'pages.admin.role.form.name.pattern', defaultMessage: '以字母开头，支持字母大小写、数字' }) },
            { max: 20, message: intl.formatMessage({ id: 'pages.admin.role.form.name.maxlen', defaultMessage: '名称不能超过20个字符' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.role.form.name.placeholder', defaultMessage: '以字母开头，支持字母大小写、数字，不超过20个字符' })} />
        </Form.Item>

        <Form.Item
          name="role"
          label={<FormattedMessage id="pages.admin.user.key.role" defaultMessage="标识" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.role.form.role.required', defaultMessage: '请输入标识' }) },
            { pattern: /^[a-zA-Z][a-zA-Z0-9]*$/, message: intl.formatMessage({ id: 'pages.admin.role.form.role.pattern', defaultMessage: '以字母开头，支持字母大小写、数字' }) },
            { max: 20, message: intl.formatMessage({ id: 'pages.admin.role.form.role.maxlen', defaultMessage: '标识不能超过20个字符' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.role.form.role.placeholder', defaultMessage: '请输入标识' })} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateForm;
