export function createPayment(client: any, token: string) {
  const request: Record<string, unknown> = { amount: 1200, currency: "usd" };
  request.payment_method = token;
  return client.payments.create(request);
}

export function listPayments(client: any) {
  return client.payments.list();
}
