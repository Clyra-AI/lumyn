export function retrieveUser(sdk: any, userId: string) {
  return sdk.users.retrieve(String(userId));
}

export function listThroughSDK(sdk: any) {
  return sdk.users.list();
}
