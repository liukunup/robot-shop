declare namespace API {
  type GetProfileResponse = {
    code?: number;
    data?: GetProfileResponseData;
    message?: string;
  };

  type GetProfileResponseData = {
    nickname?: string;
    userId?: string;
  };

  type LoginRequest = {
    email: string;
    password: string;
  };

  type LoginResponse = {
    code?: number;
    data?: LoginResponseData;
    message?: string;
  };

  type LoginResponseData = {
    accessToken?: string;
  };

  type RegisterRequest = {
    email: string;
    password: string;
  };

  type Response = {
    code?: number;
    data?: any;
    message?: string;
  };

  type UpdateProfileRequest = {
    email: string;
    nickname?: string;
  };
}
