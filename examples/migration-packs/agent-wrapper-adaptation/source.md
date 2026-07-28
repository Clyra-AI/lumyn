# Synthetic wrapper adaptation

The provider replaces its charge operation with a payment-intent operation and
a structured request. The provider evidence does not define the consumer's
gateway API: the consumer must coordinate its caller and SDK-facing wrapper
while preserving local amount, currency, and health behavior. This source is
authored for the Lumyn benchmark and is not a provider announcement or customer
result.
