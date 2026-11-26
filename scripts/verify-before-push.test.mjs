import assert from "node:assert/strict";
import { EventEmitter } from "node:events";
import { readFileSync } from "node:fs";
import path from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import vm from "node:vm";

const scriptUrl = new URL("./verify-before-push.mjs", import.meta.url);
const source = readFileSync(scriptUrl, "utf8");

// Run the real entry point with isolated process/filesystem adapters, never real Git or checks.
async function runVerification(platform, failLint = false) {
    const calls = [];
    const messages = [];
    const logs = new Map();
    const exit = new Error("simulated process exit");
    let exitCode;
    const context = {
        __scriptUrl: scriptUrl.href,
        path,
        fileURLToPath,
        cpus: () => [{}, {}],
        tmpdir: () => path.dirname(fileURLToPath(scriptUrl)),
        existsSync: () => false,
        mkdirSync() {},
        rmSync() {},
        readFileSync: () => source,
        readdirSync: () => ["runtime.test.mjs"],
        createWriteStream(filename) {
            const name = path.basename(filename);
            logs.set(name, "");
            return { write: (chunk) => logs.set(name, logs.get(name) + chunk), end() {} };
        },
        console: { log: (message) => messages.push(message), error: (message) => messages.push(message) },
        process: {
            platform,
            pid: 1234,
            argv: ["node", fileURLToPath(scriptUrl)],
            env: { CANVAS_TEST_POSTGRES_DSN: "test-fixture" },
            on() {},
            exit(code) { exitCode = code; throw exit; },
        },
        setTimeout: () => 1,
        clearTimeout() {},
        spawnSync(command, args, options = {}) {
            calls.push({ command, args, options, sync: true });
            let stdout = "";
            if (command === "git") {
                if (["merge-base", "rev-parse", "write-tree", "commit-tree"].includes(args[0])) stdout = "a".repeat(40);
                else if (args[0] === "rev-list") stdout = "0";
                else if (args[0] === "status") stdout = " M scripts/verify-before-push.mjs";
                else if (args[0] === "diff") stdout = "scripts/verify-before-push.mjs";
                else assert.ok(["ls-files", "read-tree", "add"].includes(args[0]), `Unexpected Git command: ${args}`);
            } else if (command === "go" && args[0] === "list") stdout = "yingce/backend/internal/app\nyingce/backend/internal/kernel\n";
            return { status: 0, stdout, stderr: "" };
        },
        spawn(command, args, options) {
            calls.push({ command, args, options, sync: false });
            const child = new EventEmitter();
            child.pid = calls.length + 1000;
            child.stdout = new EventEmitter();
            child.stderr = new EventEmitter();
            queueMicrotask(() => {
                const failed = failLint && args.includes("lint");
                const stdout = args.includes("-list") ? "TestA\nTestB\nTestC\nTestD\n" : command === "gofmt" ? "" : "stdout retained\n";
                child.stdout.emit("data", stdout);
                if (failed) child.stderr.emit("data", "stderr retained\n");
                child.emit("close", failed ? 7 : 0);
            });
            return child;
        },
    };
    const executable = source.replace(/^#!.*\r?\n/, "").replace(/^import[\s\S]*?;\r?\n/gm, "").replaceAll("import.meta.url", "__scriptUrl");
    try {
        await vm.runInNewContext(`(async () => { ${executable}\n })()`, context);
    } catch (error) {
        if (error !== exit) throw error;
    }
    return { calls, messages, logs, exitCode };
}

for (const platform of ["win32", "linux"]) {
    test(`${platform}: checks and probes hide windows without losing output or exit status`, async () => {
        const result = await runVerification(platform, platform === "win32");
        assert.ok(result.calls.length > 10);
        for (const { command, options, sync } of result.calls) {
            assert.equal(options.windowsHide, true, `${command} must hide its console`);
            if (!sync) assert.equal(options.detached, platform !== "win32", `${command} process group policy`);
        }
        for (const command of ["git", "go", "golangci-lint", "redis-server", "node", "bun", "gofmt"]) {
            assert.ok(result.calls.some((call) => call.command === command), `${command} was exercised`);
        }
        assert.match(result.logs.get("web-test.log"), /stdout retained/);
        assert.equal(result.exitCode, platform === "win32" ? 1 : 0);
        if (platform === "win32") {
            assert.match(result.logs.get("web-lint.log"), /stdout retained\nstderr retained/);
            assert.ok(result.messages.some((message) => message.includes("退出码 7")));
            assert.ok(result.messages.some((message) => message.includes("stderr retained")));
        }
    });
}

for (const platform of ["win32", "linux"]) {
    test(`${platform}: cleanup targets only a recorded child process tree`, () => {
        const calls = [];
        const cleanup = source.slice(source.indexOf("function killTree(child)"), source.indexOf("function seconds(ms)"));
        const context = {
            process: { platform, kill: (...args) => calls.push({ kind: "kill", args }) },
            spawnSync: (...args) => calls.push({ kind: "spawnSync", args }),
        };
        const killTree = vm.runInNewContext(`${cleanup}\nkillTree`, context);
        for (const child of [undefined, {}, { pid: 0 }]) killTree(child);
        assert.equal(calls.length, 0);
        killTree({ pid: 4321 });
        assert.equal(calls.length, 1);
        if (platform === "win32") {
            assert.equal(calls[0].kind, "spawnSync");
            const [command, args, options] = calls[0].args;
            assert.equal(command, "taskkill");
            assert.deepEqual(Array.from(args), ["/PID", "4321", "/T", "/F"]);
            assert.equal(options.windowsHide, true);
            assert.equal(options.stdio, "ignore");
        } else {
            assert.equal(calls[0].kind, "kill");
            assert.deepEqual(calls[0].args, [-4321, "SIGTERM"]);
        }
    });
}
