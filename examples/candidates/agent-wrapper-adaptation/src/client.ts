import { createPaymentThroughSDK, gatewayHealth } from "./gateway.ts";

export function health(sdk: any) {
  return gatewayHealth(sdk);
}

export function submit(sdk: any, amount: number, currency: string) {
  return createPaymentThroughSDK(sdk, { amount, currency });
}
