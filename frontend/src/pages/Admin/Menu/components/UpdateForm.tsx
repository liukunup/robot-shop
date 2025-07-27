import { Form, Input, Modal, message } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { useForm } from 'antd/es/form/Form';
import { useState, useEffect } from 'react';
import { updateMenu } from '@/services/backend/menu';

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
      await updateMenu({id: values.id}, values as API.MenuRequest);
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
      title={<FormattedMessage id="pages.admin.menu.modal.updateForm.title" defaultMessage="编辑菜单" />}
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
          name="parentId"
          label={<FormattedMessage id="pages.admin.menu.key.parent" defaultMessage="父级菜单" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="path"
          label={<FormattedMessage id="pages.admin.menu.key.path" defaultMessage="路径" />}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'pages.admin.menu.form.path.required', defaultMessage: '路径不能为空'}) },
          ]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="component"
          label={<FormattedMessage id="pages.admin.menu.key.component" defaultMessage="组件" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="name"
          label={<FormattedMessage id="pages.admin.menu.key.name" defaultMessage="名称" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="icon"
          label={<FormattedMessage id="pages.admin.menu.key.icon" defaultMessage="图标" />}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="weight"
          label={<FormattedMessage id="pages.admin.menu.key.weight" defaultMessage="权重" />}
        >
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UpdateForm;
