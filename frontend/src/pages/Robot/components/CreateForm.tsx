import { Checkbox, Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState } from 'react';
import { createRobot } from '@/services/backend/robot';

interface CreateFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
}

const CreateForm = ({ visible, onCancel, onSuccess }: CreateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.RobotRequest>();
  const intl = useIntl();

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      // 确保值为false或true
      const params = {
        ...values,
        enabled: values.enabled !== undefined ? values.enabled : true
      };
      await createRobot(params as API.RobotRequest);
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
      title={<FormattedMessage id="pages.robot.modal.createForm.title" defaultMessage="新建机器人" />}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{ enabled: true }}
        style={{ marginTop: 24 }}
      >
        <Form.Item
          name="name"
          label={<FormattedMessage id="pages.robot.key.name" defaultMessage="名称" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.robot.form.name.required', defaultMessage: '名称不能为空' }) },
            { max: 20, message: intl.formatMessage({ id: 'pages.robot.form.name.maxlen', defaultMessage: '名称不能超过20个字符' }) },
          ]}
        >
          <Input placeholder={intl.formatMessage({ id: 'pages.robot.form.name.placeholder', defaultMessage: '取一个有意义的名字吧' })} />
        </Form.Item>

        <Form.Item
          name="desc"
          label={<FormattedMessage id="pages.robot.key.desc" defaultMessage="描述" />}
        >
          <Input.TextArea placeholder={intl.formatMessage({ id: 'pages.robot.form.desc.placeholder', defaultMessage: '简要描述功能，比如它可以用来做什么' })} />
        </Form.Item>

        <Form.Item
          name="webhook"
          label={<FormattedMessage id="pages.robot.key.webhook" defaultMessage="Webhook" />}
          rules={[{ type: 'url', message: intl.formatMessage({ id: 'pages.robot.form.webhook.url', defaultMessage: '请输入正确的 URL' }) }]}
        >
          <Input.TextArea placeholder={intl.formatMessage({ id: 'pages.robot.form.webhook.placeholder', defaultMessage: 'https://example.com/webhook' })} />
        </Form.Item>

        <Form.Item
          name="callback"
          label={<FormattedMessage id="pages.robot.key.callback" defaultMessage="Callback" />}
          rules={[{ type: 'url', message: intl.formatMessage({ id: 'pages.robot.form.callback.url', defaultMessage: '请输入正确的 URL' }) }]}
        >
          <Input.TextArea placeholder={intl.formatMessage({ id: 'pages.robot.form.callback.placeholder', defaultMessage: 'https://example.com/callback' })} />
        </Form.Item>

        <Form.Item
          name="enabled"
          label={<FormattedMessage id="pages.robot.key.enabled" defaultMessage="启用状态" />}
          valuePropName="checked"
        >
          <Checkbox />
        </Form.Item>

        <Form.Item
          name="owner"
          label={<FormattedMessage id="pages.robot.key.owner" defaultMessage="所有者" />}
        >
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateForm;
