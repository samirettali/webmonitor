import axios from "axios";
import {
  Check,
  CheckUpdate,
  isArrayOfChecks,
  isArrayOfStatus,
  isCheck,
  Status,
} from "../model";

import { BACKEND_URL } from "../constants";

const instance = axios.create({
  baseURL: BACKEND_URL,
  timeout: 1000,
});

export const getChecks = async (): Promise<Check[]> => {
  const response = await instance.get("/checks");
  const data: unknown = response.data;
  if (!isArrayOfChecks(data)) {
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

export const history = async (id: string): Promise<Status[]> => {
  const response = await instance.get(`/checks/${id}/history`);
  const data: unknown = response.data;
  if (!isArrayOfStatus(data)) {
    throw new TypeError("Received malformed API response");
  }
  return data;
};
