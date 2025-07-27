import { Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState, useEffect } from 'react';
import { updateRole } from '@/services/backend/role';

interface UpdateFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  initialValues?: Partial<API.Role>;
}

const UpdateForm = ({ visible, onCancel, onSuccess, initialValues }: UpdateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.Role>();
  const intl = useIntl();

  useEffect(() => {
    if (visible && initialValues) {
      form.setFieldsValue(initialValues);
    }
  }, [visible, initialValues, form]);

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      if (!values.id) {
        throw new Error('Record ID not found during update operation');
      }
      await updateRole({id: values.id}, values as API.RoleRequest);
      message.success(intl.formatMessage({ id: 'pages.common.update.success', defaultMessage: '更新成功' }));
      form.resetFields();
      onSuccess();
    } catch (error) {
      const msg = intl.formatMessage({ id: 'pages.common.update.failure', defaultMessage: '更新失败' });
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
      title={<FormattedMessage id="pages.admin.role.modal.updateForm.title" defaultMessage="编辑角色" />}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
      destroyOnHidden={true}
    >
      <Form
        form={form}
        layout="vertical"
        className="update-role-form"
      >
        <Form.Item name="id" label="ID" hidden>
          <Input disabled />
        </Form.Item>

        <Form.Item
          name="name"
          label={<FormattedMessage id="pages.admin.role.key.name" defaultMessage="名称" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.role.form.name.required', defaultMessage: '名称不能为空' }) },
            { max: 20, message: intl.formatMessage({ id: 'pages.admin.role.form.name.maxlen', defaultMessage: '名称不能超过20个字符' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.role.form.name.placeholder', defaultMessage: '请输入角色的名称' })} />
        </Form.Item>

        <Form.Item
          name="casbinRole"
          label={<FormattedMessage id="pages.admin.role.key.role" defaultMessage="标识" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.role.form.role.required', defaultMessage: '标识不能为空' }) },
            { max: 20, message: intl.formatMessage({ id: 'pages.admin.role.form.role.maxlen', defaultMessage: '标识不能超过20个字符' }) },
            { pattern: /^[a-zA-Z][a-zA-Z0-9]*$/, message: intl.formatMessage({ id: 'pages.admin.role.form.role.pattern', defaultMessage: '以字母开头，支持字母大小写、数字' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.role.form.role.placeholder', defaultMessage: '请输入角色的标识' })} disabled />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UpdateForm;
