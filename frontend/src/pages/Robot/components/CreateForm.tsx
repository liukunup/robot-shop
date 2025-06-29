import { Checkbox, Form, Input, Modal, message } from 'antd';
import { FormattedMessage } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState } from 'react';
import { robotCreate } from '@/services/backend/robot';

interface CreateFormProps {
  visible: boolean; // 弹窗是否可见
  onCancel: () => void; // 取消回调
  onSuccess: () => void; // 成功回调
}

const CreateForm = ({ visible, onCancel, onSuccess }: CreateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.RobotRequest>();

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      // 确保值为false或true
      const params = {
        ...values,
        enabled: values.enabled !== undefined ? values.enabled : true
      };
      await robotCreate(params as API.RobotRequest);
      message.success('新增成功');
      onSuccess();
    } catch (error) {
      if (error instanceof Error) {
        message.error(error.message || '新增失败');
      } else {
        message.error('新增失败');
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
      title={<FormattedMessage id="pages.robot.table.createForm.newRobot" defaultMessage="新建机器人" />}
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
          label="名称"
          rules={[{ required: true, message: '请输入名称' }]}
        >
          <Input placeholder="取一个有意义的名字吧" />
        </Form.Item>

        <Form.Item
          name="desc"
          label="描述"
        >
          <Input.TextArea placeholder="简要描述功能，比如它可以用来做什么" />
        </Form.Item>

        <Form.Item
          name="webhook"
          label="通知地址"
          rules={[{ type: 'url', message: '请输入正确的 URL' }]}
        >
          <Input.TextArea placeholder="https://example.com/webhook" />
        </Form.Item>

        <Form.Item
          name="callback"
          label="回调地址"
          rules={[{ type: 'url', message: '请输入正确的 URL' }]}
        >
          <Input.TextArea placeholder="https://example.com/callback" />
        </Form.Item>

        <Form.Item
          name="enabled"
          label="是否启用"
          valuePropName="checked"
          rules={[{ required: true, message: '请勾选或取消' }]}
        >
          <Checkbox />
        </Form.Item>

        <Form.Item
          name="owner"
          label="所有者"
        >
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateForm;
