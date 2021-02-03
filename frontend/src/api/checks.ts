import axios, { AxiosResponse } from "axios";
import {
  Check,
  CheckUpdate,
  isArrayOfChecks,
  isArrayOfStatus,
  isCheck,
  isStatus,
  Status,
} from "../model";

import { BACKEND_URL } from "../constants";
import { isArray } from "util";

const dateFormat = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$/;

const instance = axios.create({
  baseURL: BACKEND_URL,
  timeout: 1000,
});

instance.interceptors.response.use(
  (resp: AxiosResponse): AxiosResponse => {
    if (resp.data !== undefined) {
      if (Array.isArray(resp.data)) {
        resp.data.forEach((obj: any, index: number) => {
          for (const key in obj) {
            if (key === "date") {
              console.log("Axios interceptor is modifying a field");
              resp.data[index].date = new Date(obj.date);
            }
          }
        });
      }
    }

    return resp;
  },
  (error) => {
    return Promise.reject(error);
  }
);
export const getChecks = async (): Promise<Check[]> => {
  const response = await instance.get("/checks");
  const data: unknown = response.data;
  if (!isArrayOfChecks(data)) {
    throw new TypeError("Received malformed API response");
  }
  return data;
};

export const getCheck = async (id: string): Promise<Check> => {
  const response = await instance.get(`/checks/${id}`);
  const data: unknown = response.data;
  if (!isCheck(data)) {
    throw new TypeError("Received malformed API response");
  }
  return data;
};

export const deleteCheck = async (id: string) => {
  const { data } = await instance.delete(`/checks/${id}`);
  return data;
};

export const createCheck = async (check: Check): Promise<Check> => {
  const response = await instance.post("/checks", check);
  const data: unknown = response.data;
  if (!isCheck(data)) {
    throw new TypeError("Received malformed API response");
  }
  return data;
};

export const updateCheck = async (
  id: string,
  upd: CheckUpdate
): Promise<Check> => {
  const response = await instance.patch(`/checks/${id}`, upd);
  const data: unknown = response.data;
  if (!isCheck(data)) {
    throw new TypeError("Received malformed API response");
  }
  return data;
};

export const getHistory = async (id: string): Promise<Status[]> => {
  const response = await instance.get(`/checks/${id}/history`);
  const data: unknown = response.data;
  if (!isArrayOfStatus(data)) {
    throw new TypeError("Received malformed API response");
  }

  return data;
};
