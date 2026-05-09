---
name: confidence-loop
description: Stress-tests a strategy, plan, implementation approach, or answer until remaining uncertainty is explicit and evidence-backed, then reports a 0-100 confidence score. Use when the user asks whether Codex is 100% confident, asks to find loopholes or failure modes, requests a confidence audit, says to run a loop until the strategy is factually solid, or invokes confidence-loop hard for up to four sub-agent second opinions.
---

# Confidence Loop

Adversarially verify a strategy, plan, implementation, or answer. Treat "100% confident" as a request for evidence, not reassurance.

## Standard

1. State the strategy and success criteria.
2. List material assumptions, dependencies, edge cases, and failure modes.
3. Mark each risk as disproven, accepted, blocked, or needing a fix.
4. Revise the strategy to close real loopholes.
5. Run the smallest relevant verification: source read, command, test, search, or reasoning proof.
6. Repeat until no material unresolved loopholes remain.

## Hard Mode

Use only when the user says `confidence-loop hard` or asks for sub-agent second opinions.

Spawn two to four read-only reviewers with distinct angles. Use only as many as the uncertainty justifies:

- **Skeptic**: false assumptions, loopholes, and failure modes.
- **Verifier**: smallest checks that would prove or disprove the strategy.
- **Domain reviewer**: code, docs, APIs, product constraints, or domain-specific gaps.
- **Implementability reviewer**: hidden execution, sequencing, ownership, or integration risks.

Give each reviewer the strategy, criteria, evidence, and open questions. They must not edit files, spawn agents, or decide the final answer.

Prompt shape:

```
Read-only confidence-loop reviewer. Angle: {skeptic|verifier|domain|implementability}.
Strategy: {strategy}
Criteria: {criteria}
Evidence: {facts}
Open questions: {unknowns}
Return material loopholes with evidence, verification checks, speculative objections labeled as such, and confidence: 0-100.
```

Integrate results yourself. Accept evidence-backed issues, reject unsupported ones, fix valid gaps, then verify.

## Confidence Score

Use a numeric score from 0 to 100, not broad labels. Calibrate it this way:

- **100**: all material assumptions verified; relevant checks pass; no blocker remains.
- **90-99**: evidence strongly supports the answer, with only minor non-material uncertainty left.
- **70-89**: likely correct, but meaningful uncertainty, incomplete coverage, or unverified assumptions remain.
- **40-69**: plausible but materially uncertain; more evidence or revision is needed.
- **0-39**: weak, contradicted, missing core evidence, or verification failed.

If 100% is impossible, say why and what evidence would close the gap.

## Output

Keep it compact: Strategy, Loopholes found, Second opinions if hard mode, Verification, Confidence score. Do not pad with hypothetical risks that do not apply.
