import { readFile } from "node:fs/promises";
import { stripTypeScriptTypes } from "node:module";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import vm from "node:vm";
import { Worker, isMainThread, parentPort, workerData } from "node:worker_threads";

const invocation = isMainThread
  ? {
      argv: process.argv.slice(2),
      verificationTarget: process.env.LUMYN_M1_VERIFICATION_TARGET ?? "combined",
    }
  : workerData;
const [scenarioID, candidatePath] = invocation.argv;
if (!scenarioID || !candidatePath) {
  throw new Error("usage: run.mjs <scenario-id> <candidate-path>");
}
const verificationTarget = invocation.verificationTarget;
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
  try {
    await entryModule.evaluate({ timeout: 1_000 });
  } catch (error) {
    const infrastructureError = new Error("subject module evaluation failed", { cause: error });
    infrastructureError.code = "LUMYN_VERIFICATION_INFRASTRUCTURE";
    throw infrastructureError;
  }
  try {
    await assertionModule.evaluate({ timeout: 1_000 });
  } catch (error) {
    if (error?.code === "ERR_SCRIPT_EXECUTION_TIMEOUT") {
      const infrastructureError = new Error("verification VM execution timeout", { cause: error });
      infrastructureError.code = "LUMYN_VERIFICATION_INFRASTRUCTURE";
      throw infrastructureError;
    }
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

async function runVerification() {
  const resolvedCandidatePath = path.resolve(
    process.env.LUMYN_M1_CANDIDATE_PATH ?? candidatePath,
  );
  if (verificationTarget === "baseline") {
    try {
      await executeRestrictedAssertion("src/client.ts", check.primary);
      return {
        exitCode: 2,
        stderr: `${scenarioID}: baseline unexpectedly satisfies the target contract\n`,
      };
    } catch (error) {
      if (error?.code === "LUMYN_TARGET_CONTRACT_MISMATCH") {
        return {
          exitCode: 1,
          stdout: `${scenarioID}: baseline target contract rejected\n`,
        };
      }
      throw error;
    }
  } else if (verificationTarget === "candidate") {
    await executeRestrictedAssertion(resolvedCandidatePath, check.primary);
    if (check.supportPath) {
      await executeRestrictedAssertion(
        path.join(path.dirname(resolvedCandidatePath), check.supportPath),
        check.support,
      );
    }
    return {
      exitCode: 0,
      stdout: `${scenarioID}: exact candidate target contract verified\n`,
    };
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
    return {
      exitCode: 0,
      stdout: `${scenarioID}: red-before/green-candidate offline behavior verified\n`,
    };
  }
}

function runVerificationWorker() {
  return new Promise((resolve, reject) => {
    const worker = new Worker(new URL(import.meta.url), {
      workerData: invocation,
    });
    let settled = false;
    let resultMessage;
    const finish = (callback) => {
      if (settled) return;
      settled = true;
      clearTimeout(deadline);
      callback();
    };
    const deadline = setTimeout(() => {
      finish(() => {
        void worker.terminate();
        const error = new Error("verification wall-clock deadline exceeded");
        error.code = "LUMYN_VERIFICATION_TIMEOUT";
        reject(error);
      });
    }, 3_000);
    worker.once("message", (message) => {
      resultMessage = message;
    });
    worker.once("error", (error) => finish(() => reject(error)));
    worker.once("exit", (code) => {
      if (settled) return;
      finish(() => {
        if (code !== 0) {
          reject(new Error(`verification worker failed after reporting a result (code ${code})`));
          return;
        }
        if (!resultMessage) {
          reject(new Error("verification worker exited without a result"));
          return;
        }
        if (resultMessage.ok) {
          if (resultMessage.stdout) process.stdout.write(resultMessage.stdout);
          if (resultMessage.stderr) process.stderr.write(resultMessage.stderr);
          resolve(resultMessage.exitCode);
          return;
        }
        reject(deserializeError(resultMessage.error));
      });
    });
  });
}

function serializeError(error, depth = 0) {
  return {
    name: error?.name ?? "Error",
    message: error?.message ?? String(error),
    code: error?.code,
    stack: typeof error?.stack === "string" ? error.stack : undefined,
    cause: depth < 4 && error?.cause ? serializeError(error.cause, depth + 1) : undefined,
  };
}

function deserializeError(payload, depth = 0) {
  const cause = depth < 4 && payload?.cause
    ? deserializeError(payload.cause, depth + 1)
    : undefined;
  const error = new Error(payload?.message ?? "verification worker failed", { cause });
  error.name = payload?.name ?? "Error";
  error.code = payload?.code;
  if (typeof payload?.stack === "string") error.stack = payload.stack;
  return error;
}

if (isMainThread) {
  process.exitCode = await runVerificationWorker();
} else {
  try {
    const result = await runVerification();
    parentPort.postMessage({ ok: true, ...result });
  } catch (error) {
    parentPort.postMessage({
      ok: false,
      error: serializeError(error),
    });
  }
}
