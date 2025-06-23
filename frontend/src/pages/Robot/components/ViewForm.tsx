import { Checkbox, Form, Input, Modal, message } from 'antd';
import { useForm } from 'antd/es/form/Form';
import { useState } from 'react';
import { updateRobot } from '../../../services/backend/robot';

type RobotViewParams = {
  id: number;
  name: string;
  desc: string;
  webhook: string;
  callback: string;
  enabled: boolean;
  owner: string;
};

interface ViewFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  initialValues: RobotViewParams;
}

const ViewForm = ({ visible, onCancel, onSuccess, initialValues }: ViewFormProps) => {
  const [form] = useForm<RobotViewParams>();
  const [loading, setLoading] = useState(false);

  form.setFieldsValue(initialValues);

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      await updateRobot({id: initialValues.id}, values as API.RobotParams);
      message.success('更新成功');
      onSuccess();
    } catch (error) {
      if (error instanceof Error) {
        message.error(error.message || '更新失败');
      } else {
        message.error('更新失败');
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
      title="更新机器人"
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={initialValues}
        disabled={true}
      >
        <Form.Item name="id" label="ID">
          <Input />
        </Form.Item>

        <Form.Item
          name="name"
          label="机器人名称"
          rules={[{ required: true, message: '请输入名称' }]}
        >
          <Input placeholder="输入机器人名称" />
        </Form.Item>

        <Form.Item name="desc" label="描述">
          <Input.TextArea placeholder="输入描述" />
        </Form.Item>

        <Form.Item name="webhook" label="Webhook">
          <Input placeholder="输入Webhook" />
        </Form.Item>

        <Form.Item name="callback" label="回调">
          <Input placeholder="输入回调" />
        </Form.Item>

        <Form.Item name="enabled" label="是否启用" valuePropName="enabled">
          <Checkbox />
        </Form.Item>

        <Form.Item name="owner" label="所有者">
          <Input.TextArea placeholder="输入所有者" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default ViewForm;
