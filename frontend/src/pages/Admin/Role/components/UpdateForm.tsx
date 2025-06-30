import { Select, Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState, useEffect } from 'react';
import { userUpdate } from '@/services/backend/user';
import { listRoles } from '@/services/backend/role';

interface UpdateFormProps {
  visible: boolean; // 弹窗是否可见
  onCancel: () => void; // 取消回调
  onSuccess: () => void; // 成功回调
  initialValues?: Partial<API.User>; // 初始值
}

const UpdateForm = ({ visible, onCancel, onSuccess, initialValues }: UpdateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [roles, setRoles] = useState<API.Role[]>([]);
  const [roleLoading, setRoleLoading] = useState(false);
  const [form] = useForm<API.User>();
  const intl = useIntl();

  // 获取角色列表
  const fetchRoles = async () => {
    setRoleLoading(true);
    try {
      const response = await listRoles({
        page: 1,
        pageSize: 100, // 正常情况下，不应该超过100个角色
      });
      if (response.success) {
        setRoles(response.data?.list || []);
      }
    } catch (error) {
      message.error(intl.formatMessage({ id: 'pages.admin.user.message.fetchRolesFailed', defaultMessage: '获取角色列表失败' }));
    } finally {
      setRoleLoading(false);
    }
  };
  useEffect(() => {
    if (visible) {
      fetchRoles();
    }
  }, [visible]);

  // 初始化表单
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
        throw new Error('更新操作时未找到记录ID');
      }
      await userUpdate({id: values.id}, values as API.UserRequest);
      message.success(intl.formatMessage({ id: 'pages.admin.user.message.updateUserSuccess', defaultMessage: '更新成功' }));
      form.resetFields();
      onSuccess();
    } catch (error) {
      const msg = intl.formatMessage({ id: 'pages.admin.user.message.updateUserFailed', defaultMessage: '更新失败' });
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
      title={<FormattedMessage id="pages.admin.user.modal.updateForm.title" defaultMessage="更新用户" />}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
      destroyOnHidden={true}
      width={600}
    >
      <Form
        form={form}
        layout="vertical"
        className="update-user-form"
      >
        <Form.Item name="id" label="ID" hidden>
          <Input disabled />
        </Form.Item>

        <Form.Item
          name="username"
          label={<FormattedMessage id="pages.admin.user.key.username" defaultMessage="用户名" />}
          rules={[
            { required: true, message: intl.formatMessage({
              id: 'pages.admin.user.form.username.required',
              defaultMessage: '请输入用户名' }) },
            { pattern: /^[a-zA-Z][a-zA-Z0-9]*$/, message: intl.formatMessage({
              id: 'pages.admin.user.form.username.pattern',
              defaultMessage: '以字母开头，支持字母大小写、数字' }) },
            { max: 20, message: intl.formatMessage({
              id: 'pages.admin.user.form.username.maxlen',
              defaultMessage: '用户名不能超过20个字符' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({
            id: 'pages.admin.user.form.username.placeholder',
            defaultMessage: '以字母开头，支持字母大小写、数字，不超过20个字符' })} />
        </Form.Item>

        <Form.Item
          name="nickname"
          label={<FormattedMessage id="pages.admin.user.key.nickname" defaultMessage="昵称" />}
          rules={[
            { max: 20, message: intl.formatMessage({
              id: 'pages.admin.user.form.nickname.maxlen',
              defaultMessage: '昵称不能超过20个字符' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({
            id: 'pages.admin.user.form.nickname.placeholder',
            defaultMessage: '请输入昵称' })} />
        </Form.Item>

        <Form.Item
          name="email"
          label={<FormattedMessage id="pages.admin.user.key.email" defaultMessage="邮箱" />}
          rules={[
            { required: true, message: intl.formatMessage({
              id: 'pages.admin.user.form.email.required',
              defaultMessage: '请输入邮箱地址' }) },
            { type: 'email', message: intl.formatMessage({
              id: 'pages.admin.user.form.email.pattern',
              defaultMessage: '请输入正确的邮箱地址' }) },
        ]}
        >
          <Input placeholder={intl.formatMessage({
            id: 'pages.admin.user.form.email.placeholder',
            defaultMessage: '请输入邮箱地址' })} />
        </Form.Item>

        <Form.Item
          name="phone"
          label={<FormattedMessage id="pages.admin.user.key.phone" defaultMessage="手机" />}
          rules={[
            { pattern: /^(\+?\d{1,3}[- ]?)?\d{10}$/, message: intl.formatMessage({
              id: 'pages.admin.user.form.phone.pattern',
              defaultMessage: '请输入正确的手机号码' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({
            id: 'pages.admin.user.form.phone.placeholder',
            defaultMessage: '请输入手机号码' })} />
        </Form.Item>

        <Form.Item
          name="roles"
          label={<FormattedMessage id="pages.admin.user.key.roles" defaultMessage="角色" />}
          rules={[
            { required: true, message: intl.formatMessage({
              id: 'pages.admin.user.form.roles.required',
              defaultMessage: '请选择角色' }) },
          ]}
        >
          <Select
            mode="multiple"
            placeholder={intl.formatMessage({ id: 'pages.admin.user.form.roles.placeholder', defaultMessage: '请选择角色' })}
            loading={roleLoading}
            style={{ width: '100%' }}
          >
            {roles.map(role => (
              <Select.Option key={role.id} value={role.id}>
                {role.name}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UpdateForm;
