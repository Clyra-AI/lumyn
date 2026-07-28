export function normalizeInvoice(result: any) {
  return { status: result.status, invoiceId: result.id };
}

export function invoiceHealth(client: any) {
  return client.health();
}
