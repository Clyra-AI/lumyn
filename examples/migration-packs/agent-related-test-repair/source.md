# Synthetic related-test repair

`createInvoice` now returns an object with `status` and `id`. The provider
evidence does not define the consumer's normalized invoice view or local test
contract: the consumer must coordinate both while preserving input construction
and health behavior. This source is authored for the Lumyn benchmark and is not
external proof.
