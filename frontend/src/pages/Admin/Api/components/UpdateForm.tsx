import { Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState, useEffect } from 'react';
import { updateApi } from '@/services/backend/api';

interface UpdateFormProps {
  visible: boolean; // 弹窗是否可见
  onCancel: () => void; // 取消回调
  onSuccess: () => void; // 成功回调
  initialValues?: Partial<API.Api>; // 初始值
}
const UpdateForm = ({ visible, onCancel, onSuccess, initialValues }: UpdateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.Api>();
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
        throw new Error('更新操作时未找到记录ID');
      }
      await updateApi({id: values.id}, values as API.ApiRequest);
      message.success(intl.formatMessage({ id: 'pages.common.object.update.success', defaultMessage: '更新成功' }));
      form.resetFields();
      onSuccess();
    } catch (error) {
      const msg = intl.formatMessage({ id: 'pages.common.object.update.failed', defaultMessage: '更新失败' });
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
      title={<FormattedMessage id="pages.admin.api.modal.updateForm.title" defaultMessage="更新接口" />}
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
        className="update-api-form"
      >
        <Form.Item name="id" label="ID" hidden>
          <Input disabled />
        </Form.Item>

        <Form.Item
          name="name"
          label={<FormattedMessage id="pages.admin.api.key.group" defaultMessage="分组" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.api.form.group.required', defaultMessage: '请输入分组' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.api.form.group.placeholder', defaultMessage: '请输入分组' })} />
        </Form.Item>

        <Form.Item
          name="name"
          label={<FormattedMessage id="pages.admin.api.key.name" defaultMessage="名称" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.api.form.name.required', defaultMessage: '请输入名称' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.api.form.name.placeholder', defaultMessage: '请输入名称' })} />
        </Form.Item>

        <Form.Item
          name="path"
          label={<FormattedMessage id="pages.admin.api.key.path" defaultMessage="路径" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.api.form.path.required', defaultMessage: '请输入路径' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.api.form.path.placeholder', defaultMessage: '请输入路径' })} />
        </Form.Item>

        <Form.Item
          name="method"
          label={<FormattedMessage id="pages.admin.api.key.method" defaultMessage="方法" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.api.form.method.required', defaultMessage: '请输入方法' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.admin.api.form.method.placeholder', defaultMessage: '请输入方法' })} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UpdateForm;
