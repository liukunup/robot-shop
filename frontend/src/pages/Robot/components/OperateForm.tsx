import { Checkbox, Form, Input, Modal, message } from 'antd';
import { useForm } from 'antd/es/form/Form';
import { useState } from 'react';
import { createRobot, updateRobot } from '@/services/backend/robot';

interface OperateFormProps {
  visible: boolean; // 弹窗是否可见
  operation: 'create' | 'update' | 'view'; // 操作类型: 创建/更新/查看
  onCancel: () => void; // 取消回调
  onSuccess: () => void; // 成功回调
  initialValues: API.Robot; // 初始值
}

const OperateForm = ({ visible, operation, onCancel, onSuccess, initialValues }: OperateFormProps) => {
  const [form] = useForm<API.Robot>();
  const [loading, setLoading] = useState(false);

  // 根据操作类型确定标题
  const title = {
    create: '新增',
    update: '更新',
    view: '查看',
  }[operation];

  // 仅在更新和查看操作时设置初始值
  if (operation === 'update' || operation === 'view') {
    form.setFieldsValue(initialValues);
  }

  const handleOk = async () => {
    setLoading(true);
    try {
      const values = await form.validateFields();
      if (operation === 'create') {
        if (!values.id) {
          values.id = 0; // 显式将字段置空
        }
        await createRobot(values as API.RobotParams);
        message.success('新增成功');
      } else if (operation === 'update') {
        if (!values.id) {
          throw new Error('ID is required for update operation.');
        }
        await updateRobot({id: values.id}, values as API.RobotParams);
        message.success('更新成功');
      } else if (operation === 'view') {
        // 查看操作不需要调用接口
      } else {
        console.error('Unknown Operation Type:', operation);
      }
      onSuccess();
    } catch (error) {
      if (error instanceof Error) {
        message.error(error.message || '操作失败');
      } else {
        message.error('操作失败');
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
      title={title}
      open={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={initialValues}
        disabled={operation === 'view'}
      >
        <Form.Item
          name="id"
          label="ID"
        >
          <Input />
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
          label="回调地址"
          rules={[{ type: 'url', message: '请输入正确的 URL' }]}
        >
          <Input.TextArea />
        </Form.Item>

        <Form.Item
          name="callback"
          label="通知地址"
          rules={[{ type: 'url', message: '请输入正确的 URL' }]}
        >
          <Input.TextArea />
        </Form.Item>

        <Form.Item
          name="enabled"
          label="是否启用"
        >
          <Checkbox defaultChecked={true} />
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

export default OperateForm;
