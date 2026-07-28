import { readFile } from "node:fs/promises";
import { stripTypeScriptTypes } from "node:module";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import vm from "node:vm";

const [scenarioID, candidatePath] = process.argv.slice(2);
if (!scenarioID || !candidatePath) {
  throw new Error("usage: run.mjs <scenario-id> <candidate-path>");
}
const verificationTarget = process.env.LUMYN_M1_VERIFICATION_TARGET ?? "combined";
if (!new Set(["baseline", "candidate", "combined"]).has(verificationTarget)) {
  throw new Error(`unsupported verification target ${verificationTarget}`);
}

const forbiddenCapabilities = [
  [/(?:^|[^\w])process(?:[^\w]|$)/u, "process"],
  [/(?:^|[^\w])globalThis(?:[^\w]|$)/u, "globalThis"],
  [/(?:^|[^\w])fetch\s*\(/u, "fetch"],
  [/(?:^|[^\w])XMLHttpRequest(?:[^\w]|$)/u, "XMLHttpRequest"],
  [/(?:^|[^\w])WebSocket(?:[^\w]|$)/u, "WebSocket"],
  [/(?:^|[^\w])require\s*\(/u, "require"],
  [/(?:^|[^\w])import\s*\(/u, "dynamic import"],
  [/(?:^|[^\w])eval\s*\(/u, "eval"],
  [/(?:^|[^\w])Function\s*\(/u, "Function constructor"],
  [/(?:^|[^\w])constructor(?:[^\w]|$)/u, "constructor escape"],
  [/(?:^|[^\w])__proto__(?:[^\w]|$)/u, "prototype escape"],
];

function assertCapabilitySafe(source, sourcePath) {
  for (const [pattern, capability] of forbiddenCapabilities) {
    if (pattern.test(source)) {
      throw new Error(`${sourcePath}: forbidden capability ${capability}`);
    }
  }
}

const assertionPrelude = `
import * as subject from "lumyn:subject";

function failTargetContract(message) {
  const error = new Error(message);
  error.code = "LUMYN_TARGET_CONTRACT_MISMATCH";
  throw error;
}

function assertEqual(actual, expected) {
  if (!Object.is(actual, expected)) {
    failTargetContract("values are not equal");
  }
}

function assertJSON(actual, expected) {
  if (JSON.stringify(actual) !== JSON.stringify(expected)) {
    failTargetContract("JSON values are not equal");
  }
}
`;

const expectBootstrap = `
function failTargetContract(message) {
  const error = new Error(message);
  error.code = "LUMYN_TARGET_CONTRACT_MISMATCH";
  throw error;
}

Object.defineProperty(globalThis, "expect", {
  configurable: false,
  enumerable: false,
  writable: false,
  value(actual) {
    return Object.freeze({
      toBeDefined() {
        if (actual === undefined) failTargetContract("expected a defined value");
      },
      toEqual(expected) {
        if (JSON.stringify(actual) !== JSON.stringify(expected)) {
          failTargetContract("expected equal values");
        }
      },
    });
  },
});
Object.freeze(globalThis.expect);
`;

async function executeRestrictedAssertion(sourcePath, assertionSource) {
  const entryPath = path.resolve(sourcePath);
  const sourceRoot = path.dirname(entryPath);
  const context = vm.createContext(Object.create(null), {
    name: `lumyn-m1:${path.basename(sourcePath)}`,
    codeGeneration: { strings: false, wasm: false },
  });
  vm.runInContext(expectBootstrap, context, { timeout: 1_000 });
  const modules = new Map();

  async function getModule(modulePath) {
    const resolvedPath = path.resolve(modulePath);
    const relativePath = path.relative(sourceRoot, resolvedPath);
    if (relativePath.startsWith("..") || path.isAbsolute(relativePath)) {
      throw new Error(`${sourcePath}: import escapes the isolated source root`);
    }
    if (!new Set([".ts", ".js", ".mjs"]).has(path.extname(resolvedPath))) {
      throw new Error(`${sourcePath}: unsupported local module extension`);
    }
    if (modules.has(resolvedPath)) {
      return modules.get(resolvedPath);
    }
    const source = await readFile(resolvedPath, "utf8");
    assertCapabilitySafe(source, resolvedPath);
    const module = new vm.SourceTextModule(stripTypeScriptTypes(source, { mode: "strip" }), {
      context,
      identifier: pathToFileURL(resolvedPath).href,
    });
    modules.set(resolvedPath, module);
    return module;
  }

  const entryModule = await getModule(entryPath);
  const assertionModule = new vm.SourceTextModule(assertionPrelude + assertionSource, {
    context,
    identifier: `lumyn:assertion:${path.basename(sourcePath)}`,
  });
  await assertionModule.link(async (specifier, referencingModule) => {
    if (specifier === "lumyn:subject") {
      return entryModule;
    }
    if (!specifier.startsWith("./") && !specifier.startsWith("../")) {
      throw new Error(`${sourcePath}: only relative local imports are allowed`);
    }
    return getModule(fileURLToPath(new URL(specifier, referencingModule.identifier)));
  });
  // Load, link, and evaluate the subject graph before invoking the target
  // contract. Missing files, parse/link failures, and top-level exceptions are
  // verifier infrastructure failures, never valid red-before evidence.
  await entryModule.evaluate({ timeout: 1_000 });
  try {
    await assertionModule.evaluate({ timeout: 1_000 });
  } catch (error) {
    if (error?.code === "LUMYN_TARGET_CONTRACT_MISMATCH") {
      throw error;
    }
    const mismatch = new Error("target contract invocation rejected");
    mismatch.code = "LUMYN_TARGET_CONTRACT_MISMATCH";
    mismatch.cause = error;
    throw mismatch;
  }
}

const checks = {
  "det-operation-rename": {
    primary: `
      assertJSON(
        subject.createCharge(
          {
            charges: { create: (request) => ({ legacy: request }) },
            paymentIntents: { create: (request) => ({ accepted: request }) },
          },
          1200,
        ),
        { accepted: { amount: 1200, currency: "usd" } },
      );
      assertEqual(subject.health({ health: { check: () => "ok" } }), "ok");
    `,
  },
  "det-request-property-relocation": {
    primary: `
      assertJSON(
        subject.createPayment(
          { payments: { create: (request) => request } },
          "pm_fixture",
        ),
        { amount: 1200, currency: "usd", payment_method: "pm_fixture" },
      );
    `,
  },
  "det-response-property-relocation": {
    primary: `
      assertJSON(
        subject.summarize({ id: "cus_fixture", customer: { email: "a@example.test" } }),
        { id: "cus_fixture", email: "a@example.test" },
      );
    `,
  },
  "agent-wrapper-adaptation": {
    primary: `
      const sdk = {
        charges: { create: (request) => ({ legacy: request }) },
        paymentIntents: { create: (request) => ({ accepted: request }) },
        health: { check: () => "ok" },
      };
      assertJSON(
        subject.submit(sdk, 1200, "usd"),
        { accepted: { amount: 1200, currency: "usd" } },
      );
      assertEqual(subject.health(sdk), "ok");
    `,
    supportPath: "gateway.ts",
    support: `
      assertJSON(
        subject.createPaymentThroughSDK(
          {
            charges: { create: (request) => ({ legacy: request }) },
            paymentIntents: { create: (request) => ({ accepted: request }) },
          },
          { amount: 1200, currency: "usd" },
        ),
        { accepted: { amount: 1200, currency: "usd" } },
      );
      assertEqual(subject.gatewayHealth({ health: { check: () => "ok" } }), "ok");
    `,
  },
  "agent-signature-type-adaptation": {
    primary: `
      const sdk = {
        users: {
          retrieve: (request) => (request?.id === "user_fixture" ? "found" : "wrong-shape"),
          list: () => ["a"],
        },
      };
      assertEqual(
        subject.loadUser(sdk, "user_fixture"),
        "found",
      );
      assertJSON(subject.listUsers(sdk), ["a"]);
    `,
    supportPath: "adapter.ts",
    support: `
      assertEqual(
        subject.retrieveUser(
          { users: { retrieve: (request) => (request?.id === "user_fixture" ? "found" : "wrong-shape") } },
          { id: "user_fixture" },
        ),
        "found",
      );
      assertJSON(subject.listThroughSDK({ users: { list: () => ["a"] } }), ["a"]);
    `,
  },
  "agent-related-test-repair": {
    primary: `
      subject.invoiceTest({ createInvoice: () => ({ status: "queued" }) }, { amount: 1200 });
      subject.healthTest({ health: () => "ok" });
    `,
    supportPath: "invoices.ts",
    support: `
      assertJSON(subject.normalizeInvoice({ status: "queued", id: "in_fixture" }), {
        status: "queued",
        invoiceId: "in_fixture",
      });
      assertEqual(subject.invoiceHealth({ health: () => "ok" }), "ok");
    `,
  },
};

const check = checks[scenarioID];
if (!check) {
  throw new Error(`unknown M1 scenario ${scenarioID}`);
}

async function proveRedBeforeGreenCandidate(label, beforePath, expectedPath, assertion) {
  try {
    await executeRestrictedAssertion(beforePath, assertion);
  } catch (error) {
    if (error?.code === "LUMYN_TARGET_CONTRACT_MISMATCH") {
      await executeRestrictedAssertion(expectedPath, assertion);
      return;
    }
    throw error;
  }
  throw new Error(`${scenarioID} ${label}: pre-migration behavior unexpectedly satisfies the target contract`);
}

const resolvedCandidatePath = path.resolve(process.env.LUMYN_M1_CANDIDATE_PATH ?? candidatePath);
if (verificationTarget === "baseline") {
  try {
    await executeRestrictedAssertion("src/client.ts", check.primary);
    console.error(`${scenarioID}: baseline unexpectedly satisfies the target contract`);
    process.exitCode = 2;
  } catch (error) {
    if (error?.code === "LUMYN_TARGET_CONTRACT_MISMATCH") {
      console.log(`${scenarioID}: baseline target contract rejected`);
      process.exitCode = 1;
    } else {
      throw error;
    }
  }
} else if (verificationTarget === "candidate") {
  await executeRestrictedAssertion(resolvedCandidatePath, check.primary);
  if (check.supportPath) {
    await executeRestrictedAssertion(
      path.join(path.dirname(resolvedCandidatePath), check.supportPath),
      check.support,
    );
  }
  console.log(`${scenarioID}: exact candidate target contract verified`);
} else {
  await proveRedBeforeGreenCandidate("primary", "src/client.ts", resolvedCandidatePath, check.primary);
  if (check.supportPath) {
    await proveRedBeforeGreenCandidate(
      "supporting",
      path.join("src", check.supportPath),
      path.join(path.dirname(resolvedCandidatePath), check.supportPath),
      check.support,
    );
  }
  console.log(`${scenarioID}: red-before/green-candidate offline behavior verified`);
}
