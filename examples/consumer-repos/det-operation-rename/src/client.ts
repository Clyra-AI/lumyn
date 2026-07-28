export function createCharge(client: any, amount: number) {
  const request = { amount, currency: "usd" };
  return client.charges.create(request);
}

export function health(client: any) {
  return client.health.check();
}
