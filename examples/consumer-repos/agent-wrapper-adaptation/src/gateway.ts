export function chargeThroughSDK(sdk: any, amount: number, currency: string) {
  return sdk.charges.create({ amount, currency });
}

export function gatewayHealth(sdk: any) {
  return sdk.health.check();
}
