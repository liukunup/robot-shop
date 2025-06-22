import { Form, Checkbox, Input, Modal, message } from 'antd';
import { useForm } from 'antd/es/form/Form';
import { postRobots } from '../../../services/backend/robot';

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
      await postRobots(values as API.CreateRobotParams);
      message.success('新增成功');
      onSuccess();
      form.resetFields();
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
      <Form form={form} layout="vertical">
        <Form.Item
          name="name"
          label="机器人名称"
          rules={[{ required: true, message: '请输入名称' }]}
        >
          <Input />
        </Form.Item>
        <Form.Item name="desc" label="描述">
          <Input.TextArea />
        </Form.Item>
        <Form.Item name="webhook" label="Webhook">
          <Input />
        </Form.Item>
        <Form.Item name="callback" label="Callback">
          <Input />
        </Form.Item>
        <Form.Item name="enabled" label="是否启用">
          <Checkbox defaultChecked={true} />
        </Form.Item>
        <Form.Item name="owner" label="所有者">
          <Input.TextArea />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateForm;
