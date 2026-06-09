const params = new URLSearchParams(window.location.search);
const pinnedThreadId = params.get("threadId") || "";
const clientId = window.crypto?.randomUUID?.() || `${Date.now()}-${Math.random().toString(16).slice(2)}`;
let threadId = pinnedThreadId;
let iconVersion = 0;
let refreshInFlight = null;
let heartbeatTimer = null;
let goalRefreshTimer = null;
let csrfToken = "";

const els = {
  objective: document.getElementById("goal-objective"),
  status: document.getElementById("goal-status"),
  tokens: document.getElementById("goal-tokens"),
  budget: document.getElementById("goal-budget"),
  time: document.getElementById("goal-time"),
  threadId: document.getElementById("thread-id"),
  output: document.getElementById("output"),
  objectiveInput: document.getElementById("goal-input"),
  budgetInput: document.getElementById("budget-input"),
  budgetMeter: document.getElementById("budget-meter"),
};

function apiPath(path) {
  const next = new URL(path, window.location.origin);
  const shouldPinThread = threadId && !(path === "/api/config" && !pinnedThreadId);
  if (shouldPinThread) next.searchParams.set("threadId", threadId);
  return `${next.pathname}${next.search}`;
}

function setOutput(value) {
  els.output.textContent = typeof value === "string" ? value : JSON.stringify(value, null, 2);
}

function updateIcons(goal) {
  const status = goal?.status || "none";
  iconVersion += 1;
  for (const link of document.querySelectorAll('link[rel="icon"], link[rel="apple-touch-icon"]')) {
    const next = new URL(link.getAttribute("href") || "/favicon.svg", window.location.origin);
    if (threadId) next.searchParams.set("threadId", threadId);
    next.searchParams.set("goalStatus", status);
    next.searchParams.set("v", String(iconVersion));
    link.href = `${next.pathname}${next.search}`;
  }
}

function formatSeconds(value) {
  const seconds = Number(value || 0);
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  const rest = seconds % 60;
  return `${minutes}m ${rest}s`;
}

function renderBudgetMeter(goal) {
  const meter = els.budgetMeter;
  const segments = Array.from(meter.querySelectorAll("span"));
  const tokensUsed = Number(goal?.tokensUsed ?? 0);
  const tokenBudget = Number(goal?.tokenBudget ?? 0);
  const hasBudget = Number.isFinite(tokenBudget) && tokenBudget > 0;

  meter.classList.toggle("is-looping", !hasBudget);
  meter.classList.toggle("is-budgeted", hasBudget);
  meter.classList.toggle("is-over-budget", hasBudget && tokensUsed > tokenBudget);

  segments.forEach((segment) => {
    segment.className = "";
  });

  if (!hasBudget) {
    meter.setAttribute("aria-valuenow", "0");
    meter.setAttribute("aria-valuetext", "No token budget set. Meter is scanning.");
    return;
  }

  const ratio = Math.max(0, Math.min(tokensUsed / tokenBudget, 1));
  const percent = Math.round(ratio * 100);
  const filled = Math.max(1, Math.ceil(ratio * segments.length));
  segments.forEach((segment, index) => {
    if (index < filled) segment.classList.add("is-filled");
    if (index === filled - 1) segment.classList.add("is-hot");
  });

  meter.setAttribute("aria-valuenow", String(percent));
  meter.setAttribute("aria-valuetext", `${tokensUsed} of ${tokenBudget} tokens used.`);
}

function renderGoal(goal) {
  els.threadId.textContent = threadId || "-";
  updateIcons(goal);
  renderBudgetMeter(goal);
  if (!goal) {
    els.objective.textContent = "No goal set";
    els.status.textContent = "none";
    els.status.dataset.status = "none";
    els.tokens.textContent = "-";
    els.budget.textContent = "-";
    els.time.textContent = "-";
    return;
  }

  els.objective.textContent = goal.objective;
  els.status.textContent = goal.status;
  els.status.dataset.status = goal.status;
  els.tokens.textContent = String(goal.tokensUsed ?? 0);
  els.budget.textContent = goal.tokenBudget == null ? "none" : String(goal.tokenBudget);
  els.time.textContent = formatSeconds(goal.timeUsedSeconds);
}

async function request(path, options = {}) {
  const method = options.method || "GET";
  const headers = {
    "content-type": "application/json",
    ...(method === "GET" ? {} : { "x-codex-goal-csrf": csrfToken }),
    ...(options.headers || {}),
  };
  const response = await fetch(apiPath(path), {
    cache: "no-store",
    ...options,
    headers: {
      ...headers,
    },
  });
  const data = await response.json();
  if (!response.ok || data.error) {
    throw new Error(data.error || `Request failed (${response.status})`);
  }
  return data;
}

async function sendHeartbeat() {
  return request("/api/session/heartbeat", {
    method: "POST",
    body: JSON.stringify({ clientId }),
  });
}

function closeSession() {
  const body = JSON.stringify({ clientId });
  const url = apiPath("/api/session/close");
  fetch(url, {
    method: "POST",
    body,
    cache: "no-store",
    headers: {
      "content-type": "application/json",
      "x-codex-goal-csrf": csrfToken,
    },
    keepalive: true,
  }).catch(() => {});
}

function startHeartbeat() {
  sendHeartbeat().catch(() => {});
  heartbeatTimer = window.setInterval(() => {
    sendHeartbeat().catch(() => {});
  }, 30_000);
}

async function loadConfig() {
  const config = await request("/api/config");
  csrfToken = config.csrfToken || "";
  if (!pinnedThreadId) {
    threadId = config.threadId || "";
    if (threadId) {
      const next = new URL(window.location.href);
      next.searchParams.set("threadId", threadId);
      window.history.replaceState(null, "", next);
    }
  } else {
    threadId = pinnedThreadId;
  }
  els.threadId.textContent = threadId;
  updateIcons(null);
}

async function refreshGoal() {
  if (refreshInFlight) return refreshInFlight;
  refreshInFlight = (async () => {
    if (!pinnedThreadId) {
      await loadConfig();
    }
    const data = await request("/api/goal");
    renderGoal(data.goal);
    setOutput(data);
  })();
  try {
    await refreshInFlight;
  } finally {
    refreshInFlight = null;
  }
}

async function setGoal() {
  const objective = els.objectiveInput.value.trim();
  const tokenBudget = els.budgetInput.value.trim();
  if (!objective) {
    setOutput("Objective is required.");
    return;
  }
  const data = await request("/api/goal", {
    method: "POST",
    body: JSON.stringify({
      objective,
      tokenBudget: tokenBudget || undefined,
    }),
  });
  renderGoal(data.goal);
  setOutput(data);
}

async function setStatus(status) {
  const data = await request("/api/goal/status", {
    method: "POST",
    body: JSON.stringify({ status }),
  });
  renderGoal(data.goal);
  setOutput(data);
}

async function clearGoal() {
  const data = await request("/api/goal", { method: "DELETE" });
  renderGoal(null);
  setOutput(data);
}

async function stopServer() {
  if (heartbeatTimer) window.clearInterval(heartbeatTimer);
  if (goalRefreshTimer) window.clearInterval(goalRefreshTimer);
  const data = await request("/api/server/stop", {
    method: "POST",
    body: JSON.stringify({ clientId }),
  });
  setOutput(data);
}

function bind(id, fn) {
  document.getElementById(id).addEventListener("click", async () => {
    try {
      await fn();
    } catch (error) {
      setOutput(error.message || String(error));
    }
  });
}

bind("refresh-goal", refreshGoal);
bind("set-goal", setGoal);
bind("pause-goal", () => setStatus("paused"));
bind("resume-goal", () => setStatus("active"));
bind("complete-goal", () => setStatus("complete"));
bind("clear-goal", clearGoal);
bind("stop-server", stopServer);

loadConfig()
  .then(() => {
    startHeartbeat();
    return refreshGoal();
  })
  .catch((error) => {
    renderGoal(null);
    setOutput(error.message || String(error));
  });

goalRefreshTimer = window.setInterval(() => {
  refreshGoal().catch((error) => {
    setOutput(error.message || String(error));
  });
}, 10_000);

window.addEventListener("pagehide", closeSession);
