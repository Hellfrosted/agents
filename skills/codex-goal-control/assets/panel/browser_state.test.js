const assert = require("node:assert/strict");
const test = require("node:test");

const {
  apiPath,
  applyConfig,
  createPanelState,
  goalViewModel,
  nextIconHrefs,
  requestHeaders,
} = require("./browser_state");

test("Given unpinned panel When config loads Then thread id is claimed and URL is updated", () => {
  const state = createPanelState({ href: "http://127.0.0.1:43873/", origin: "http://127.0.0.1:43873" });

  const nextHref = applyConfig(state, { csrfToken: "csrf", threadId: "thread-123" });

  assert.equal(state.threadId, "thread-123");
  assert.equal(state.csrfToken, "csrf");
  assert.equal(nextHref, "http://127.0.0.1:43873/?threadId=thread-123");
  assert.equal(apiPath(state, "/api/config"), "/api/config");
  assert.equal(apiPath(state, "/api/goal"), "/api/goal?threadId=thread-123");
});

test("Given pinned panel When config loads Then existing thread id remains authoritative", () => {
  const state = createPanelState({
    href: "http://127.0.0.1:43873/?threadId=pinned-thread",
    origin: "http://127.0.0.1:43873",
  });

  const nextHref = applyConfig(state, { csrfToken: "csrf", threadId: "other-thread" });

  assert.equal(nextHref, null);
  assert.equal(state.threadId, "pinned-thread");
  assert.equal(apiPath(state, "/api/config"), "/api/config?threadId=pinned-thread");
});

test("Given mutation request When headers built Then CSRF token is included", () => {
  const state = createPanelState();
  state.csrfToken = "csrf";

  assert.deepEqual(requestHeaders(state, "POST", { "x-extra": "1" }), {
    "content-type": "application/json",
    "x-codex-goal-csrf": "csrf",
    "x-extra": "1",
  });
  assert.deepEqual(requestHeaders(state, "GET"), {
    "content-type": "application/json",
  });
});

test("Given goal state When viewed Then display fields and meter are derived without DOM", () => {
  const state = createPanelState();
  state.threadId = "thread-123";

  const view = goalViewModel(
    state,
    { objective: "Ship goal panel", status: "active", tokensUsed: 75, tokenBudget: 100, timeUsedSeconds: 65 },
    4,
  );

  assert.equal(view.threadId, "thread-123");
  assert.equal(view.objective, "Ship goal panel");
  assert.equal(view.status, "active");
  assert.equal(view.tokens, "75");
  assert.equal(view.budget, "100");
  assert.equal(view.time, "1m 5s");
  assert.equal(view.budgetMeter.valueNow, "75");
  assert.deepEqual(
    view.budgetMeter.segments.map((segment) => [segment.isFilled, segment.isHot]),
    [
      [true, false],
      [true, false],
      [true, true],
      [false, false],
    ],
  );
});

test("Given icon links When goal renders Then hrefs include thread status and cache version", () => {
  const state = createPanelState({ href: "http://127.0.0.1:43873/?threadId=thread-123" });

  const hrefs = nextIconHrefs(state, ["/favicon.svg", "/apple-touch-icon.png"], { status: "paused" });

  assert.deepEqual(hrefs, [
    "/favicon.svg?threadId=thread-123&goalStatus=paused&v=1",
    "/apple-touch-icon.png?threadId=thread-123&goalStatus=paused&v=1",
  ]);
});
