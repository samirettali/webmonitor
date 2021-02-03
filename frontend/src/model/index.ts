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

export type Status = {
  content: string;
  date: Date;
};

export const isArrayOfChecks = (obj: unknown): obj is Check[] => {
  return Array.isArray(obj) && obj.every(isCheck);
};

export const isArrayOfStatus = (obj: unknown): obj is Status[] => {
  return Array.isArray(obj) && obj.every(isStatus);
};

export const isCheck = (obj: unknown): obj is Check => {
  return obj !== null && typeof (obj as Check).id === "string";
};

export const isStatus = (obj: unknown): obj is Status => {
  return obj !== null && typeof (obj as Status).content === "string";
};
