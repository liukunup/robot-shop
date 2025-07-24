import { Footer } from '@/components';
import { fetchCurrentUser, updateProfile } from '@/services/backend/user';
import { UserOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { ProForm, ProFormText, ProFormUploadButton } from '@ant-design/pro-components';
import { FormattedMessage, Helmet, useIntl, useModel } from '@umijs/max';
import { Card, message, Spin } from 'antd';
import React, { useEffect, useState } from 'react';
import { history } from '@umijs/max';
import styles from './index.less';

const Profile: React.FC = () => {
  const [profile, setProfile] = useState<API.User | null>(null);
  const [loading, setLoading] = useState(true);
  const intl = useIntl();
  const { initialState } = useModel('@@initialState');

  useEffect(() => {
    if (!initialState?.currentUser) {
      history.push('/login');
      return;
    }

    const loadProfile = async () => {
      try {
        setLoading(true);
        const data = await fetchCurrentUser();
        setProfile(data);
      } catch (error) {
        message.error(intl.formatMessage({ id: 'pages.profile.loadFailed' }));
        console.error('Failed to load profile:', error);
      } finally {
        setLoading(false);
      }
    };

    loadProfile();
  }, [intl, initialState?.currentUser]);

  const handleSubmit = async (values: API.User) => {
    try {
      setLoading(true);
      const updatedProfile = await updateProfile(values);
      setProfile(updatedProfile);
      message.success(intl.formatMessage({ id: 'pages.profile.updateSuccess' }));
    } catch (error) {
      message.error(intl.formatMessage({ id: 'pages.profile.updateFailed' }));
      console.error('Failed to update profile:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarUpload = async (file: File) => {
    try {
      const formData = new FormData();
      formData.append('avatar', file);

      const response = await fetch('/api/upload/avatar', {
        method: 'POST',
        body: formData,
      });

      const result = await response.json();

      if (result.success && result.data?.url) {
        setProfile(prev => prev ? { ...prev, avatar: result.data.url } : prev);
        message.success(intl.formatMessage({ id: 'pages.profile.avatarUploadSuccess' }));
      } else {
        throw new Error(result.message || 'Upload failed');
      }
    } catch (error) {
      message.error(intl.formatMessage({ id: 'pages.profile.avatarUploadFailed' }));
      console.error('Failed to upload avatar:', error);
    }
  };

  return (
    <div className={styles['profile-container']}>
      <Helmet title={intl.formatMessage({ id: 'pages.profile.title' })} />
      <Card className={styles['profile-card']} bordered={false}>
        <div className={styles.header}>
          <h1>{intl.formatMessage({ id: 'pages.profile.title' })}</h1>
        </div>

        <Spin spinning={loading}>
          {profile && (
            <>
              <div className={styles['avatar-container']}>
                <ProFormUploadButton
                  name="avatar"
                  listType="picture-card"
                  className={styles['avatar-upload']}
                  action="/api/upload/avatar"
                  maxCount={1}
                  fieldProps={{
                    accept: 'image/*',
                    beforeUpload: (file) => {
                      const isImage = file.type.startsWith('image/');
                      if (!isImage) {
                        message.error(intl.formatMessage({ id: 'pages.profile.avatarInvalidType' }));
                      }
                      return isImage;
                    },
                  }}
                  onChange={({ file }) => {
                    if (file.status === 'done' && file.originFileObj) {
                      handleAvatarUpload(file.originFileObj);
                    }
                  }}
                >
                  {profile.avatar ? (
                    <img
                      src={profile.avatar}
                      alt="Avatar"
                      className={styles['avatar-image']}
                    />
                  ) : (
                    <div className={styles['avatar-placeholder']}>
                      <UserOutlined />
                    </div>
                  )}
                </ProFormUploadButton>
              </div>

              <ProForm
                initialValues={profile}
                onFinish={handleSubmit}
                layout="vertical"
                submitter={{
                  searchConfig: {
                    submitText: intl.formatMessage({ id: 'pages.profile.save' }),
                  },
                  resetButtonProps: {
                    style: {
                      display: 'none',
                    },
                  },
                }}
              >
                <ProFormText
                  name="username"
                  label={<FormattedMessage id="pages.profile.username" />}
                  placeholder={intl.formatMessage({ id: 'pages.profile.usernamePlaceholder' })}
                  prefix={<UserOutlined />}
                  rules={[
                    {
                      required: true,
                      message: intl.formatMessage({ id: 'pages.profile.usernameRequired' }),
                    },
                  ]}
                />

                <ProFormText
                  name="email"
                  label={<FormattedMessage id="pages.profile.email" />}
                  placeholder={intl.formatMessage({ id: 'pages.profile.emailPlaceholder' })}
                  prefix={<MailOutlined />}
                  rules={[
                    {
                      type: 'email',
                      message: intl.formatMessage({ id: 'pages.profile.emailInvalid' }),
                    },
                    {
                      required: true,
                      message: intl.formatMessage({ id: 'pages.profile.emailRequired' }),
                    },
                  ]}
                />

                <ProFormText
                  name="phone"
                  label={<FormattedMessage id="pages.profile.phone" />}
                  placeholder={intl.formatMessage({ id: 'pages.profile.phonePlaceholder' })}
                  prefix={<PhoneOutlined />}
                  rules={[
                    {
                      pattern: /^\d{11}$/,
                      message: intl.formatMessage({ id: 'pages.profile.phoneInvalid' }),
                    },
                  ]}
                />
              </ProForm>
            </>
          )}
        </Spin>
      </Card>
      <Footer />
    </div>
  );
};

export default Profile;
