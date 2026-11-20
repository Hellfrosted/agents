---
name: council
description: Run a multi-model review or consultation when the user says 'council'.
---

# Council

## Prepare

Before an audit, fix known in-scope issues, address relevant PR comments, and run focused checks. If considering a rewrite, review as-is.

## Seat the council

In T3, use `orchestrator_capabilities` on the target machine. Require `canRunChildTask` and, for cross-provider runs, `canRunCrossProviderChildTask`. Use exact catalog IDs. Outside T3, use the host's subagents.

Cap every delegate's thinking effort at `medium`, including re-reviews. Set `medium` or a supported lower level explicitly when configurable. If a route cannot honor the cap, report the seat as unavailable.

Unless the user names members or sets a limit, seat one strong model per family that adds a useful perspective. These are preferences, not fixed IDs:

| Family | Preferred seat and lens |
| --- | --- |
| OpenAI | GPT Sol (Codex or Pi): correctness, specification and regressions. Use Terra or Luna only if the stronger model is unavailable at selection time. |
| Anthropic | Claude Opus: architecture and maintainability. If Opus is not available at selection time, consider another strong Claude model on an eligible route (e.g. Fable). |
| Google | Gemini Pro or the current comparable flagship (e.g. Gemini Neon): integration assumptions and edge cases. Pi Gemini Web search is not a delegate route. |
| xAI | Grok (including via Pi): challenge assumptions and look for overlooked failure cases. |
| DeepSeek | Strongest available DeepSeek (Pi or OpenCode): implementation details and failure paths. |
| Moonshot | Kimi Pro (e.g. Pi): long-context consistency and missing evidence. |
| Zhipu | GLM (e.g. Pi): tests and portability. |

Skip duplicate routes and opaque aliases (including Copilot `auto`/`default`) unless they add a distinct lens. Report unavailable seats; never replace a named member. Label incomplete runs as partial; if nobody runs, report that.

## Run and synthesize

Launch seats in parallel with the request, scope, evidence, and lens. Keep reviews read-only and the shared checkout stable. Ask for prioritized findings with evidence and, for code reviews, a score out of 10. Report failed routes without silently switching. Skip waiting for non-blocking members and omit their late scores.

Verify findings and summarize agreement, disagreements, and next steps. If authorized, fix verified in-scope issues and count re-review as a new round. Stop when material issues are resolved, blocked, or the round limit is reached.

For code reviews, use this format with the actual delegates, rounds, and scores:

```markdown
## Council review

Rounds: 2

| Delegate | Round 1 | Round 2 |
| --- | --- | --- |
| GPT Sol | 8/10 | 10/10 |
| Grok 4.7 | 7/10 | 9/10 |
| DeepSeek V4.1 Flash | 8/10 | 9/10 |
```

Use `no score` for a completed review without a score and `unavailable` for a seat that could not run. If the work ends in a PR, put this section in the PR description before the model and harness line, not in a comment. Without a PR, put it in the final reply.
