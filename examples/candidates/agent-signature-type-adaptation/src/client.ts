import { listThroughSDK, retrieveUser } from "./adapter.ts";

export function listUsers(sdk: any) {
  return listThroughSDK(sdk);
}

export function loadUser(sdk: any, userId: string) {
  return retrieveUser(sdk, { id: userId });
}
