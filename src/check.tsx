export type Check = {
  id?: string;
  name: string;
  interval: number;
  email: string;
  url: string;
  active: boolean;
};

export type CheckUpdate = {
  name?: string;
  interval?: number;
  email?: string;
  url?: string;
  active?: boolean;
};
