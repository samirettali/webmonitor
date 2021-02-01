import axios from "axios";
import { Check, CheckUpdate } from "../check";

import { BACKEND_URL } from "../constants";

const instance = axios.create({
  baseURL: BACKEND_URL + "/checks",
  timeout: 1000,
});

export const getChecks = async () => {
  const { data } = await instance.get("");
  return data;
};

export const deleteCheck = async (id: string) => {
  const response = await instance.delete(`/${id}`);
  return response;
};

export const createCheck = async (check: Check) => {
  const response = await instance.post("", check);
  return response;
};

export const updateCheck = async (id: string, upd: CheckUpdate) => {
  const response = await instance.patch(`/${id}`, upd);
  return response;
};
