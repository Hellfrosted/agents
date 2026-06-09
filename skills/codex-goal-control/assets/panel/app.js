const {
  GOAL_REFRESH_INTERVAL_MS,
  HEARTBEAT_INTERVAL_MS,
  apiPath,
  applyConfig,
  createPanelState,
  goalViewModel,
  nextIconHrefs,
  requestHeaders,
} = window.CodexGoalPanelState;

const clientId = window.crypto?.randomUUID?.() || `${Date.now()}-${Math.random().toString(16).slice(2)}`;
const state = createPanelState({
  href: window.location.href,
  origin: window.location.origin,
  clientId,
});
let refreshInFlight = null;
let heartbeatTimer = null;
let goalRefreshTimer = null;

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

function setOutput(value) {
  els.output.textContent = typeof value === "string" ? value : JSON.stringify(value, null, 2);
}

function renderBudgetMeter(meterState) {
  const segments = Array.from(els.budgetMeter.querySelectorAll("span"));
  els.budgetMeter.classList.toggle("is-looping", meterState.classes.isLooping);
  els.budgetMeter.classList.toggle("is-budgeted", meterState.classes.isBudgeted);
  els.budgetMeter.classList.toggle("is-over-budget", meterState.classes.isOverBudget);

  meterState.segments.forEach((segmentState, index) => {
    const segment = segments[index];
    segment.className = "";
    if (segmentState.isFilled) segment.classList.add("is-filled");
    if (segmentState.isHot) segment.classList.add("is-hot");
  });

  els.budgetMeter.setAttribute("aria-valuenow", meterState.valueNow);
  els.budgetMeter.setAttribute("aria-valuetext", meterState.valueText);
}

function updateIcons(goal) {
  const links = Array.from(document.querySelectorAll('link[rel="icon"], link[rel="apple-touch-icon"]'));
  const hrefs = nextIconHrefs(
    state,
    links.map((link) => link.getAttribute("href")),
    goal,
  );
  links.forEach((link, index) => {
    link.href = hrefs[index];
  });
}

function renderGoal(goal) {
  const view = goalViewModel(state, goal, els.budgetMeter.querySelectorAll("span").length);
  els.threadId.textContent = view.threadId;
  els.objective.textContent = view.objective;
  els.status.textContent = view.status;
  els.status.dataset.status = view.status;
  els.tokens.textContent = view.tokens;
  els.budget.textContent = view.budget;
  els.time.textContent = view.time;
  renderBudgetMeter(view.budgetMeter);
  updateIcons(goal);
}

async function request(path, options = {}) {
  const method = options.method || "GET";
  const response = await fetch(apiPath(state, path), {
    cache: "no-store",
    ...options,
    headers: requestHeaders(state, method, options.headers),
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
    body: JSON.stringify({ clientId: state.clientId }),
  });
}

function closeSession() {
  fetch(apiPath(state, "/api/session/close"), {
    method: "POST",
    body: JSON.stringify({ clientId: state.clientId }),
    cache: "no-store",
    headers: requestHeaders(state, "POST"),
    keepalive: true,
  }).catch(() => {});
}

function startHeartbeat() {
  sendHeartbeat().catch(() => {});
  heartbeatTimer = window.setInterval(() => {
    sendHeartbeat().catch(() => {});
  }, HEARTBEAT_INTERVAL_MS);
}

async function loadConfig() {
  const config = await request("/api/config");
  const nextHref = applyConfig(state, config);
  if (nextHref) window.history.replaceState(null, "", nextHref);
  els.threadId.textContent = state.threadId;
  updateIcons(null);
}

async function refreshGoal() {
  if (refreshInFlight) return refreshInFlight;
  refreshInFlight = (async () => {
    if (!state.pinnedThreadId) {
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
    body: JSON.stringify({ clientId: state.clientId }),
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
}, GOAL_REFRESH_INTERVAL_MS);

window.addEventListener("pagehide", closeSession);
