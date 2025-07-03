/**
 * @see https://umijs.org/docs/max/access#access
 * */
export default function access(initialState: { currentUser?: API.User } | undefined) {
  const { currentUser } = initialState ?? {};
  const hasRole = (user: API.User | undefined, roles: string[]) => {
    return user && user.roles?.some(role => roles.includes(role.casbinRole as string));
  };
  return {
    canAdmin: hasRole(currentUser, ['admin']),
    canOperate: hasRole(currentUser, ['admin', 'operator']),
    canView: hasRole(currentUser, ['admin', 'operator', 'viewer']),
  };
}
