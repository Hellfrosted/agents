(function publishBrowserState(root, factory) {
  const api = factory();
  if (typeof module === "object" && module.exports) module.exports = api;
  if (root) root.CodexGoalPanelState = api;
})(typeof globalThis === "object" ? globalThis : this, () => {
  const HEARTBEAT_INTERVAL_MS = 30_000;
  const GOAL_REFRESH_INTERVAL_MS = 10_000;

  function createPanelState({ href, origin, clientId } = {}) {
    const current = new URL(href || "http://127.0.0.1/");
    const pinnedThreadId = current.searchParams.get("threadId") || "";
    return {
      href: current.href,
      origin: origin || current.origin,
      pinnedThreadId,
      threadId: pinnedThreadId,
      clientId: clientId || "",
      csrfToken: "",
      iconVersion: 0,
    };
  }

  function apiPath(state, path) {
    const next = new URL(path, state.origin);
    const shouldPinThread = state.threadId && !(path === "/api/config" && !state.pinnedThreadId);
    if (shouldPinThread) next.searchParams.set("threadId", state.threadId);
    return `${next.pathname}${next.search}`;
  }

  function applyConfig(state, config) {
    state.csrfToken = config.csrfToken || "";
    if (state.pinnedThreadId) {
      state.threadId = state.pinnedThreadId;
      return null;
    }

    state.threadId = config.threadId || "";
    if (!state.threadId) return null;

    const next = new URL(state.href);
    next.searchParams.set("threadId", state.threadId);
    state.href = next.href;
    return state.href;
  }

  function requestHeaders(state, method, headers = {}) {
    return {
      "content-type": "application/json",
      ...(method === "GET" ? {} : { "x-codex-goal-csrf": state.csrfToken }),
      ...headers,
    };
  }

  function formatSeconds(value) {
    const seconds = Number(value || 0);
    if (seconds < 60) return `${seconds}s`;
    const minutes = Math.floor(seconds / 60);
    const rest = seconds % 60;
    return `${minutes}m ${rest}s`;
  }

  function budgetMeterState(goal, segmentCount) {
    const tokensUsed = Number(goal?.tokensUsed ?? 0);
    const tokenBudget = Number(goal?.tokenBudget ?? 0);
    const hasBudget = Number.isFinite(tokenBudget) && tokenBudget > 0;
    const classes = {
      isLooping: !hasBudget,
      isBudgeted: hasBudget,
      isOverBudget: hasBudget && tokensUsed > tokenBudget,
    };

    if (!hasBudget) {
      return {
        classes,
        segments: Array.from({ length: segmentCount }, () => ({ isFilled: false, isHot: false })),
        valueNow: "0",
        valueText: "No token budget set. Meter is scanning.",
      };
    }

    const ratio = Math.max(0, Math.min(tokensUsed / tokenBudget, 1));
    const percent = Math.round(ratio * 100);
    const filled = Math.max(1, Math.ceil(ratio * segmentCount));
    return {
      classes,
      segments: Array.from({ length: segmentCount }, (_, index) => ({
        isFilled: index < filled,
        isHot: index === filled - 1,
      })),
      valueNow: String(percent),
      valueText: `${tokensUsed} of ${tokenBudget} tokens used.`,
    };
  }

  function goalViewModel(state, goal, segmentCount) {
    if (!goal) {
      return {
        threadId: state.threadId || "-",
        objective: "No goal set",
        status: "none",
        tokens: "-",
        budget: "-",
        time: "-",
        budgetMeter: budgetMeterState(null, segmentCount),
      };
    }

    return {
      threadId: state.threadId || "-",
      objective: goal.objective,
      status: goal.status,
      tokens: String(goal.tokensUsed ?? 0),
      budget: goal.tokenBudget == null ? "none" : String(goal.tokenBudget),
      time: formatSeconds(goal.timeUsedSeconds),
      budgetMeter: budgetMeterState(goal, segmentCount),
    };
  }

  function nextIconHrefs(state, hrefs, goal) {
    const status = goal?.status || "none";
    state.iconVersion += 1;
    return hrefs.map((href) => {
      const next = new URL(href || "/favicon.svg", state.origin);
      if (state.threadId) next.searchParams.set("threadId", state.threadId);
      next.searchParams.set("goalStatus", status);
      next.searchParams.set("v", String(state.iconVersion));
      return `${next.pathname}${next.search}`;
    });
  }

  return {
    GOAL_REFRESH_INTERVAL_MS,
    HEARTBEAT_INTERVAL_MS,
    apiPath,
    applyConfig,
    createPanelState,
    formatSeconds,
    goalViewModel,
    nextIconHrefs,
    requestHeaders,
  };
});
