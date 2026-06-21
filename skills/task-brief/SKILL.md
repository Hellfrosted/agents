---
name: task-brief
description: Task brief generator for compact Codex execution briefs, prompts, tickets, or handoffs. Use when the user asks for a Codex task brief/prompt, an implementation ticket or future-work handoff for delegation, or an execution brief that turns rough intent into concrete scope and evidence. Not for ordinary implementation/review, confidence audits, conversation-compaction handoffs, Evo/evo-hq experiment briefs, long PRDs, image prompts, generic writing prompts, or non-Codex prompt engineering.
---
# Task Brief
Create or rewrite a brief that reduces rework from broad asks, unclear
referents, and domain-term drift. Keep it compact and actionable. Ask only when
the missing detail materially changes outcome, risk, credentials, deployment,
hardware, or destructive actions. If the user wants an audit or confidence
check of an existing brief, use a review/confidence workflow instead.

## Output Shape
```text
Goal: <one concrete outcome>
Context: <paths, repo, links, current state>
Semantic target: <define the domain term that must not drift>
Do: <allowed actions>
Do not: <safety boundaries>
Evidence: <what proves done>
Output: <file/report/commit/PR/table/etc>
Mode: light unless <specific heavy triggers>
```

Use `Semantic target` for domain-sensitive work. If no special term exists, write `Semantic target: none beyond the ordinary wording in Goal`.

## Templates

### fix
```text
Goal: fix <observable failure> so <expected behavior>.
Context: repo/path: <path>; failing command/error: <exact output>; recent change: <if known>.
Semantic target: <term that must not drift, or none beyond Goal>.
Do: reproduce, capture failing-first proof, make the smallest fix, run focused validation.
Do not: refactor unrelated code, suppress diagnostics, or change public behavior outside the failure.
Evidence: RED proof, GREEN command output, and real-surface proof when user-facing.
Output: patch summary plus verification evidence.
Mode: light unless auth, persistence, concurrency, external APIs, or cross-module contracts are involved.
```

### research
```text
Goal: answer <decision/question> with source-backed findings.
Context: repo/path/links: <paths and URLs>; prior artifacts: <reports or notes>.
Semantic target: <domain distinction that must not drift>.
Do: inspect primary sources, separate facts from inference, cite exact files/URLs.
Do not: rely on stale memory for unstable facts or install/mutate anything.
Evidence: source list with dates/paths and quoted or paraphrased support.
Output: concise report, table, or recommendation.
Mode: light unless the answer drives security, money, deployment, legal, medical, or architecture decisions.
```

### review
```text
Goal: review <change/plan/artifact> for bugs, regressions, and missing validation.
Context: branch/PR/files: <refs>; intended behavior: <spec>.
Semantic target: <critical term whose meaning must stay fixed>.
Do: inspect diff and relevant callers, prioritize actionable findings with file/line references.
Do not: rewrite code unless explicitly asked, or spend review space on low-value style nits.
Evidence: findings tied to source lines and verification gaps.
Output: findings first, then open questions and brief summary.
Mode: light unless security, data loss, permissions, concurrency, or cross-service behavior is in scope.
```

### build
```text
Goal: build <user-visible feature/artifact> that supports <workflow>.
Context: repo/path: <path>; design constraints: <style/system>; related files: <paths>.
Semantic target: <product/domain term that must not drift>.
Do: follow existing patterns, implement complete controls/states, add focused validation.
Do not: create speculative abstractions, unrelated redesigns, or placeholder-only behavior.
Evidence: failing-first proof where behavior changes, passing tests/checks, and real-surface QA.
Output: changed files, running URL/file if applicable, and verification evidence.
Mode: light unless adding a new module/layer/domain model or touching auth, persistence, external integrations, or deployment.
```

### operate/deploy
```text
Goal: perform <operation> safely and confirm <postcondition>.
Context: target host/service/path: <exact target>; environment: <dev/staging/prod>; rollback: <known path>.
Semantic target: <operational term that must not drift>.
Do: inspect current state, use dry-run/read-only checks first when available, capture before/after state.
Do not: run destructive, credential, exposure, network, or machine-wide changes without explicit approval.
Evidence: command transcript, status check, logs/health output, and cleanup or rollback receipt.
Output: operation report with exact commands and final state.
Mode: heavy for production, credentials, network exposure, hardware, irreversible data, or multi-service changes.
```

## Example
```text
Goal: identify client-only ATM10 mods that are also server-capable.
Context: client mods path: <path>; official list: <path or URL>.
Semantic target: "sideness" means upstream mod loader side support, not install location.
Do: compare filenames, resolve mod identity, classify mod sideness from upstream docs/source.
Do not: infer sidedness from where the file is installed.
Evidence: table with source link or local source evidence per mod.
Output: markdown report with uncertain items separated.
Mode: light unless adding/removing server mods or touching deployment config.
```

Use `ulw` only for high-risk, multi-step, evidence-heavy work. Default to light mode for ordinary prompt shaping and implementation briefs.
