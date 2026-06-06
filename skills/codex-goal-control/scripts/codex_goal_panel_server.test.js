const assert = require("node:assert/strict");
const test = require("node:test");

const {
  isValidJsonContentType,
  verifyMutationRequest,
} = require("./codex_goal_panel_server");

function request(headers = {}) {
  return { headers };
}

test("Given a simple cross-origin POST When mutation guard runs Then it rejects the request", () => {
  const result = verifyMutationRequest(request({ "content-type": "text/plain" }), "token");

  assert.equal(result.statusCode, 403);
  assert.equal(result.payload.error, "Invalid CSRF token");
});

test("Given a valid token with text/plain When mutation guard runs Then it rejects the content type", () => {
  const result = verifyMutationRequest(
    request({
      "content-type": "text/plain",
      "x-codex-goal-csrf": "token",
    }),
    "token",
  );

  assert.equal(result.statusCode, 415);
  assert.equal(result.payload.error, "Request content-type must be application/json");
});

test("Given a valid token with JSON content When mutation guard runs Then it allows the request", () => {
  const result = verifyMutationRequest(
    request({
      "content-type": "application/json; charset=utf-8",
      "x-codex-goal-csrf": "token",
    }),
    "token",
  );

  assert.equal(result, null);
});

test("Given content type variants When parsed Then only JSON is accepted", () => {
  assert.equal(isValidJsonContentType("application/json"), true);
  assert.equal(isValidJsonContentType("application/json; charset=utf-8"), true);
  assert.equal(isValidJsonContentType("text/plain"), false);
  assert.equal(isValidJsonContentType(""), false);
});
