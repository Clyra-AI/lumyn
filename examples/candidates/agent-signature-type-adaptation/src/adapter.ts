export function retrieveUser(sdk: any, lookup: any) {
  return sdk.users.retrieve(lookup);
}

export function listThroughSDK(sdk: any) {
  return sdk.users.list();
}
