declare namespace API {
  type CurrentUser = {
    avatar?: string;
    email?: string;
    nickname?: string;
    role?: string;
    userId?: string;
  };

  type deleteRobotIdParams = {
    /** 机器人ID */
    id: number;
  };

  type GetProfileResponse = {
    code?: number;
    data?: CurrentUser;
    message?: string;
  };

  type getRobotIdParams = {
    /** 机器人ID */
    id: number;
  };

  type getRobotParams = {
    /** page */
    page?: number;
    /** size */
    size?: number;
  };

  type LoginParams = {
    password: string;
    username: string;
  };

  type LoginResponseData = {
    accessToken?: string;
  };

  type LoginResult = {
    code?: number;
    data?: LoginResponseData;
    message?: string;
  };

  type PageResponseBackendApiV1RobotResponseData = {
    items?: RobotResponseData[];
    total?: number;
  };

  type putRobotIdParams = {
    /** 机器人ID */
    id: number;
  };

  type RegisterParams = {
    email: string;
    password: string;
  };

  type Response = {
    code?: number;
    data?: any;
    message?: string;
  };

  type RobotRequest = {
    callback?: string;
    desc?: string;
    enabled?: boolean;
    name: string;
    options?: string;
    owner?: string;
    webhook?: string;
  };

  type RobotResponseData = {
    callback?: string;
    createdAt?: string;
    desc?: string;
    enabled?: boolean;
    id?: number;
    name?: string;
    options?: string;
    owner?: string;
    robot_id?: string;
    updatedAt?: string;
    webhook?: string;
  };

  type UpdateProfileParams = {
    email: string;
    nickname?: string;
  };
}
