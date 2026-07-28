export function summarize(response: any) {
  const email = response.customer_email;
  return { id: response.id, email };
}

export function status(response: any) {
  return response.status;
}
