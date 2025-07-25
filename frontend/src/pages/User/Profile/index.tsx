import { Tabs, Card, Spin, Form, Input, Button } from "antd";
import { message } from 'antd';
import { useModel, useNavigate, useIntl } from '@umijs/max';
import { useState, useEffect } from 'react';
import { fetchCurrentUser, updateProfile } from '@/services/backend/user';

const Profile: React.FC = () => {
  const [profile, setProfile] = useState<API.User | null>(null);
  const [loading, setLoading] = useState(true);
  const intl = useIntl();
  const { initialState } = useModel('@@initialState');
  const navigate = useNavigate();

  useEffect(() => {
    // Redirect to login if no current user
    if (!initialState?.currentUser) {
      navigate('/login');
      return;
    }

    const loadProfile = async () => {
      try {
        setLoading(true);
        const response = await fetchCurrentUser();
        if (response.success) {
          setProfile(response.data || null);
        }
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

  return (
    <div>
      <Tabs
        tabPosition="left"
      >
        <Tabs.TabPane tab={intl.formatMessage({ id: 'pages.profile.tab.profile' })} key="profile">
          {loading ? (
            <Spin />
          ) : (
            <Card>
              <Form
                form={form}
                initialValues={profile}
                onFinish={handleSubmit}
                layout="vertical"
              >
                <Form.Item
                  label="昵称"
                  name="nickname"
                >
                  <Input />
                </Form.Item>

                <Form.Item
                  label="手机号"
                  name="phone"
                >
                  <Input />
                </Form.Item>

                <Form.Item
                  label="邮箱"
                  name="email"
                  rules={[
                    {
                      type: 'email',
                      message: '请输入有效的邮箱地址',
                    },
                  ]}
                >
                  <Input />
                </Form.Item>

                <Form.Item>
                  <Button type="primary" htmlType="submit">
                    保存修改
                  </Button>
                </Form.Item>
              </Form>
            </Card>
          )}
        </Tabs.TabPane>
        <Tabs.TabPane tab={intl.formatMessage({ id: 'pages.profile.tab.settings' })} key="settings">
         <div style={{ padding: '24px' }}>
            <Card title="系统设置">
              <p>系统设置内容待实现</p>
            </Card>
          </div>
        </Tabs.TabPane>
      </Tabs>
    </div>
  );
};

export default Profile;