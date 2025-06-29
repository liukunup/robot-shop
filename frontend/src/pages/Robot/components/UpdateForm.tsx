import { Checkbox, Form, Input, Modal, message } from 'antd';
import { FormattedMessage } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState, useEffect } from 'react';
import { robotUpdate } from '@/services/backend/robot';

interface UpdateFormProps {
  visible: boolean; // 弹窗是否可见
  onCancel: () => void; // 取消回调
  onSuccess: () => void; // 成功回调
  initialValues?: Partial<API.Robot>; // 初始值
}

const UpdateForm = ({ visible, onCancel, onSuccess, initialValues }: UpdateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.Robot>();

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
      await robotUpdate({id: values.id}, values as API.RobotRequest);
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
      title={<FormattedMessage id="pages.robot.table.updateForm.editRobot" defaultMessage="编辑机器人" />}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
      destroyOnHidden={true}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={initialValues}
        style={{ marginTop: 24 }}
      >

        <Form.Item name="id" label="ID" hidden>
          <Input disabled />
        </Form.Item>

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
          <Input.TextArea placeholder="描述一下功能，比如它可以用来处理什么类型的任务" />
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

export default UpdateForm;
