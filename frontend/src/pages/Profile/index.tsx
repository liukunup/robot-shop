import React, { useState, useEffect, useCallback } from 'react';
import { Tabs, Spin, Form, Input, Button, Upload, message } from "antd";
import type { UploadProps, UploadChangeParam, RcFile } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import { useModel, useNavigate, useIntl } from '@umijs/max';
import { fetchCurrentUser, updateProfile, updatePassword } from '@/services/backend/user';

// 常量定义
const AVATAR_UPLOAD_URL = 'https://660d2bd96ddfa2943b33731c.mockapi.io/api/upload';
const MAX_NICKNAME_LENGTH = 20;
const MAX_AVATAR_SIZE_MB = 2;
const MIN_PASSWORD_LENGTH = 6;

const Profile: React.FC = () => {
  // Hooks
  const { initialState } = useModel('@@initialState');
  const [profile, setProfile] = useState<API.User | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const intl = useIntl();
  const [form] = Form.useForm();

  const loadProfile = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetchCurrentUser();
      if (response.success) {
        setProfile(response.data);
        form.setFieldsValue(response.data);
      }
    } catch (error) {
      message.error(intl.formatMessage({ id: 'pages.profile.loadFailed' }));
      console.error('Failed to load profile:', error);
    } finally {
      setLoading(false);
    }
  }, [form, intl]);

  useEffect(() => {
    if (!initialState?.currentUser) {
      navigate('/login');
      return;
    }
    loadProfile();
  }, [initialState?.currentUser, navigate, loadProfile]);

  const handleProfileChange = async (values: Partial<API.User>) => {
    try {
      setLoading(true);
      await updateProfile(values);
      await loadProfile(); // 重新加载最新资料
      message.success(intl.formatMessage({ id: 'pages.profile.updateSuccess' }));
    } catch (error) {
      message.error(intl.formatMessage({ id: 'pages.profile.updateFailed' }));
      console.error('Failed to update profile:', error);
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordChange = async (values: {
    currentPassword: string;
    newPassword: string;
    confirmPassword: string;
  }) => {
    try {
      setLoading(true);

      // 验证逻辑
      if (values.currentPassword === values.newPassword) {
        message.error('新密码不能与当前密码相同');
        return;
      }

      await updatePassword({
        oldPassword: values.currentPassword,
        newPassword: values.newPassword,
      });
      message.success(intl.formatMessage({ id: 'pages.profile.passwordUpdateSuccess' }));
      form.resetFields(); // 清空密码表单
    } catch (error) {
      message.error(intl.formatMessage({ id: 'pages.profile.passwordUpdateFailed' }));
      console.error('Failed to update password:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarChange: UploadProps['onChange'] = (info: UploadChangeParam) => {
    if (info.file.status === 'uploading') {
      setLoading(true);
      return;
    }

    if (info.file.status === 'done') {
      const avatarUrl = info.file.response?.data?.url;
      if (avatarUrl) {
        setProfile(prev => prev ? { ...prev, avatar: avatarUrl } : prev);
        message.success(intl.formatMessage({ id: 'pages.profile.avatarUploadSuccess' }));
      }
    } else if (info.file.status === 'error') {
      message.error(intl.formatMessage({ id: 'pages.profile.avatarUploadFailed' }));
    }
    setLoading(false);
  };

  const beforeUploadAvatar = (file: RcFile) => {
    const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
    const isLt2M = file.size / 1024 / 1024 < MAX_AVATAR_SIZE_MB;

    if (!isJpgOrPng) {
      message.error(intl.formatMessage({ id: 'pages.profile.avatarFormatError' }));
      return false;
    }
    if (!isLt2M) {
      message.error(intl.formatMessage({ id: 'pages.profile.avatarSizeError' }));
      return false;
    }
    return true;
  };

  const renderAvatarUpload = () => (
    <Upload
      name="avatar"
      listType="picture-circle"
      className="avatar-uploader"
      showUploadList={false}
      action={AVATAR_UPLOAD_URL}
      beforeUpload={beforeUploadAvatar}
      onChange={handleAvatarChange}
    >
      {profile?.avatar ? (
        <img src={profile.avatar} alt="avatar" style={{ width: '100%' }} />
      ) : (
        <Button icon={<UploadOutlined />}>上传头像</Button>
      )}
    </Upload>
  );

  // 渲染资料表单
  const renderProfileForm = () => (
    <Form
      form={form}
      initialValues={profile || {}}
      onFinish={handleProfileChange}
      layout="vertical"
    >
      <Form.Item
        label="昵称"
        name="nickname"
        rules={[
          { required: true, message: '请输入昵称' },
          { max: MAX_NICKNAME_LENGTH, message: `昵称长度不能超过${MAX_NICKNAME_LENGTH}位` },
        ]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        label="手机号"
        name="phone"
        rules={[
          { required: true, message: '请输入手机号' },
          { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' },
        ]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        label="邮箱"
        name="email"
        rules={[{ type: 'email', message: '请输入有效的邮箱地址' }]}
      >
        <Input />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading}>
          保存修改
        </Button>
      </Form.Item>
    </Form>
  );

  // 渲染密码表单
  const renderPasswordForm = () => (
    <Form layout="vertical" onFinish={handlePasswordChange}>
      <Form.Item
        label="当前密码"
        name="currentPassword"
        rules={[{ required: true, message: '请输入当前密码' }]}
      >
        <Input.Password />
      </Form.Item>
      <Form.Item
        label="新密码"
        name="newPassword"
        rules={[
          { required: true, message: '请输入新密码' },
          { min: MIN_PASSWORD_LENGTH, message: `密码长度不能少于${MIN_PASSWORD_LENGTH}位` },
          { pattern: /^(?=.*[a-zA-Z])(?=.*\d).+$/, message: '密码必须包含字母和数字' },
        ]}
      >
        <Input.Password />
      </Form.Item>
      <Form.Item
        label="确认新密码"
        name="confirmPassword"
        dependencies={['newPassword']}
        rules={[
          { required: true, message: '请确认新密码' },
          ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue('newPassword') === value) {
                return Promise.resolve();
              }
              return Promise.reject(new Error('两次输入的密码不一致'));
            },
          }),
        ]}
      >
        <Input.Password />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading}>
          修改密码
        </Button>
      </Form.Item>
    </Form>
  );

  return (
    <Tabs tabPosition="left">
      <Tabs.TabPane
        tab={intl.formatMessage({ id: 'pages.profile.tab.profile' })}
        key="profile"
      >
        {loading ? <Spin /> : (
          <>
            {renderAvatarUpload()}
            {renderProfileForm()}
          </>
        )}
      </Tabs.TabPane>
      <Tabs.TabPane
        key="settings"
        tab={intl.formatMessage({ id: 'pages.profile.tab.settings' })}
      >
        {renderPasswordForm()}
      </Tabs.TabPane>
    </Tabs>
  );
};

export default Profile;
