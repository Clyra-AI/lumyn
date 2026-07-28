import { invoiceHealth, normalizeInvoice } from "./invoices.ts";

export function invoiceTest(client: any, input: unknown) {
  expect(input).toBeDefined();
  const invoice = normalizeInvoice(client.createInvoice(input));
  expect(invoice.status).toEqual("queued");
}

export function healthTest(client: any) {
  expect(invoiceHealth(client)).toEqual("ok");
}
