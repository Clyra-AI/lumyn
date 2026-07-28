export function createPaymentThroughSDK(sdk: any, request: any) {
  return sdk.paymentIntents.create(request);
}

export function gatewayHealth(sdk: any) {
  return sdk.health.check();
}
