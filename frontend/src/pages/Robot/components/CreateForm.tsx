import { Form, Checkbox, Input, Modal, message } from 'antd';
import { useForm } from 'antd/es/form/Form';
import { createRobot } from '@/services/backend/robot';

interface CreateFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
}

const CreateForm = ({ visible, onCancel, onSuccess }: CreateFormProps) => {
  const [form] = useForm();

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      const result = await createRobot(values as API.RobotParams);
      if (result.success) {
        message.success('新增成功');
        onSuccess();
        form.resetFields();
      } else {
        message.error(result.errorMessage);
      }
    } catch (error) {
      message.error('新增失败');
    }
  };

  return (
    <Modal
      title="新增机器人"
      open={visible}
      onOk={handleOk}
      onCancel={() => {
        onCancel();
        form.resetFields();
      }}
    >
      <Form
        form={form}
        layout="vertical"
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
          <Input.TextArea placeholder="描述一下功能，比如它可以用来处理什么类型的任务" />
        </Form.Item>

        <Form.Item
          name="webhook"
          label="Webhook"
          rules={[{ type: 'url', message: '请输入正确的 URL' }]}
        >
          <Input.TextArea />
        </Form.Item>

        <Form.Item
          name="callback"
          label="Callback"
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

export default CreateForm;
