import { Select, Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState, useEffect } from 'react';
import { createUser } from '@/services/backend/user';
import { listRoles } from '@/services/backend/role';

interface CreateFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
}

const CreateForm = ({ visible, onCancel, onSuccess }: CreateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [roles, setRoles] = useState<API.Role[]>([]);
  const [roleLoading, setRoleLoading] = useState(false);
  const [form] = useForm<API.UserRequest>();
  const intl = useIntl();

  const fetchRoles = async () => {
    setRoleLoading(true);
    try {
      const response = await listRoles({});
      if (response.success) {
        setRoles(response.data?.list || []);
      }
    } catch (error) {
      const msg = intl.formatMessage({ id: 'pages.admin.user.fetchRoles.failure', defaultMessage: '获取角色列表失败' });
      if (error instanceof Error) {
        message.error(error.message || msg);
      } else {
        message.error(msg);
      }
    } finally {
      setRoleLoading(false);
    }
  };
  useEffect(() => {
    if (visible) {
      fetchRoles();
    }
  }, [visible]);

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      const params = {
        ...values,
        status: 0,
      }
      await createUser(params as API.UserRequest);
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
      title={<FormattedMessage id="pages.admin.user.modal.createForm.title" defaultMessage="新建用户" />}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
    >
      <Form
        form={form}
        layout="vertical"
        className="create-user-form"
      >
        <Form.Item
          name="username"
          label={<FormattedMessage id="pages.admin.user.key.username" defaultMessage="用户名" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.user.form.username.required', defaultMessage: '用户名不能为空' }) },
            { max: 20, message: intl.formatMessage({ id: 'pages.admin.user.form.username.maxlen', defaultMessage: '用户名不能超过20个字符' }) },
            { pattern: /^[a-zA-Z][a-zA-Z0-9]*$/, message: intl.formatMessage({ id: 'pages.admin.user.form.username.pattern', defaultMessage: '以字母开头，支持字母大小写、数字' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.user.form.username.placeholder', defaultMessage: '以字母开头，支持字母大小写、数字，不超过20个字符' })} />
        </Form.Item>

        <Form.Item
          name="nickname"
          label={<FormattedMessage id="pages.admin.user.key.nickname" defaultMessage="昵称" />}
          rules={[
            { max: 20, message: intl.formatMessage({ id: 'pages.admin.user.form.nickname.maxlen', defaultMessage: '昵称不能超过20个字符' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.user.form.nickname.placeholder', defaultMessage: '请输入昵称' })} />
        </Form.Item>

        <Form.Item
          name="email"
          label={<FormattedMessage id="pages.admin.user.key.email" defaultMessage="邮箱" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.user.form.email.required', defaultMessage: '邮箱不能为空' }) },
            { type: 'email', message: intl.formatMessage({ id: 'pages.admin.user.form.email.pattern', defaultMessage: '请输入正确的邮箱地址' }) },
        ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.user.form.email.placeholder', defaultMessage: '请输入邮箱地址' })} />
        </Form.Item>

        <Form.Item
          name="phone"
          label={<FormattedMessage id="pages.admin.user.key.phone" defaultMessage="手机" />}
          rules={[
            { pattern: /^(\+?\d{1,3}[- ]?)?\d{10}$/, message: intl.formatMessage({ id: 'pages.admin.user.form.phone.pattern', defaultMessage: '请输入正确的手机号码' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.user.form.phone.placeholder', defaultMessage: '请输入手机号码' })} />
        </Form.Item>

        <Form.Item
          name="roles"
          label={<FormattedMessage id="pages.admin.user.key.roles" defaultMessage="角色" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.user.form.roles.required', defaultMessage: '角色不能为空' }) },
          ]}
        >
          <Select
            mode="multiple"
            placeholder={intl.formatMessage({ id: 'pages.admin.user.form.roles.placeholder', defaultMessage: '请选择角色' })}
            loading={roleLoading}
            style={{ width: '100%' }}
            optionLabelProp="label"
          >
            {roles.map(r => (
              <Select.Option key={r.id} value={r.casbinRole} label={r.name}>
                {r.name} ({r.casbinRole})
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateForm;
