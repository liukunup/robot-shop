import { Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState, useEffect } from 'react';
import { menuUpdate } from '@/services/backend/menu';

interface UpdateFormProps {
  visible: boolean; // 弹窗是否可见
  onCancel: () => void; // 取消回调
  onSuccess: () => void; // 成功回调
  initialValues?: Partial<API.Menu>; // 初始值
}

const UpdateForm = ({ visible, onCancel, onSuccess, initialValues }: UpdateFormProps) => {
  const [loading, setLoading] = useState(false);
  const [form] = useForm<API.Menu>();
  const intl = useIntl();

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
      await menuUpdate({id: values.id}, values as API.MenuRequest);
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
      title={<FormattedMessage id="pages.admin.menu.modal.updateForm.title" defaultMessage="更新菜单" />}
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
        className="update-menu-form"
      >
        <Form.Item name="id" label="ID" hidden>
          <Input disabled />
        </Form.Item>

        <Form.Item
          name="parentID"
          label="父级菜单ID"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="weight"
          label="排序权重"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="path"
          label="地址"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="title"
          label="展示名称"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="name"
          label="同路由中的name，唯一标识"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="component"
          label="绑定组件"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="locale"
          label="本地化标识"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="icon"
          label="图标，使用字符串表示"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="redirect"
          label="重定向地址"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="keepAlive"
          label="是否保活"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="hideInMenu"
          label="是否隐藏在菜单中"
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="url"
          label="iframe模式下的跳转url，不能与path重复"
        >
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UpdateForm;
