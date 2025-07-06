/**
 * @see https://umijs.org/docs/max/access#access
 * */
export default function access(initialState: { currentUser?: API.User } | undefined) {
  // 从 initialState 中获取当前用户
  const { currentUser } = initialState ?? {};
  // 角色检查
  const hasRole = (user: API.User | undefined, roles: string[]) => {
    return user && user.roles?.some(role => roles.includes(role.casbinRole as string));
  };

  return {
    canAdmin: hasRole(currentUser, ['admin']),
    canOperate: hasRole(currentUser, ['admin', 'operator']),
    canUser: hasRole(currentUser, ['admin', 'operator', 'user']),
  };
}
