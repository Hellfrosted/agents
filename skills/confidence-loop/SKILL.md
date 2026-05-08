---
name: confidence-loop
description: Stress-tests a strategy, plan, implementation approach, or answer until remaining uncertainty is explicit and evidence-backed. Use when the user asks whether Codex is 100% confident, asks to find loopholes or failure modes, requests a confidence audit, says to run a loop until the strategy is factually solid, or invokes confidence-loop hard for sub-agent second opinions.
---

# Confidence Loop

## Rule

Treat "100% confident" as a demand for adversarial verification, not reassurance. Do not claim certainty unless every material assumption has been checked against evidence available in the current context or through allowed tools.

## Modes

- **Standard**: run the loop yourself with available evidence.
- **Hard**: when the user explicitly says `$confidence-loop hard`, `confidence-loop hard`, or asks to use sub-agents for second opinions, spawn read-only sub-agents to independently challenge the strategy before the final revision.

Use hard mode only when the user explicitly asks for it. If sub-agents are unavailable, state that blocker and continue with the standard loop.

## Loop

1. State the current strategy in concrete terms.
2. List the success criteria and constraints.
3. Identify assumptions, dependencies, edge cases, loopholes, and ways the strategy can fail.
4. For each risk, decide whether it is disproven, accepted with rationale, or needs a fix.
5. Revise the strategy to close real loopholes.
6. Verify the revised strategy with the smallest relevant checks, source reads, tests, searches, or reasoning proofs available.
7. Repeat until no material unresolved loopholes remain, or until a blocker prevents factual certainty.

## Hard Mode

Hard mode adds independent review before the final revision:

1. Prepare a short evidence pack: current strategy, success criteria, constraints, known facts, and open questions.
2. Spawn two or three read-only sub-agents with distinct review angles. Useful angles are:
   - **Skeptic**: find false assumptions, missing cases, and ways the plan fails.
   - **Verifier**: identify the smallest checks that would prove or disprove the plan.
   - **Domain reviewer**: inspect code, docs, APIs, or product constraints for domain-specific gaps.
3. Keep each sub-agent task bounded. Give it the evidence pack, the exact question to answer, and a requirement to separate evidence-backed findings from speculation.
4. Do not ask sub-agents to edit files, spawn more sub-agents, or make the final decision.
5. Wait for the second opinions before claiming high or 100% confidence.
6. Integrate the results yourself: dedupe findings, reject unsupported objections, fix valid loopholes, then run the smallest relevant verification.

Use this prompt shape for each hard-mode sub-agent:

```text
You are a read-only second-opinion reviewer for a confidence-loop hard pass.

Review angle: {skeptic | verifier | domain reviewer}
Current strategy: {strategy}
Success criteria and constraints: {criteria}
Evidence available: {facts, files, commands, sources, or prior results}
Open questions: {unknowns}

Do not edit files. Do not spawn sub-agents.
Return only:
- Material loopholes or false assumptions, with evidence.
- Verification checks that would close uncertainty.
- Objections that are speculative or non-blocking, clearly labeled as such.
- A confidence recommendation: 100%, high, or not confident.
```

## Confidence Standard

Use these labels precisely:

- **100% confident**: all material assumptions are verified, tests or checks pass where relevant, and no unresolved blocker remains.
- **High confidence**: evidence supports the strategy, but at least one non-material uncertainty remains.
- **Not confident**: material uncertainty remains, evidence is missing, or verification failed.

If 100% confidence is impossible, say exactly why and what evidence would close the gap.

## Output

Keep the response compact:

- **Strategy**: the revised strategy.
- **Loopholes found**: real failure modes and fixes.
- **Second opinions**: hard mode only; sub-agent findings accepted, rejected, or still unresolved.
- **Verification**: checks performed and results.
- **Confidence**: 100%, high, or not confident, with the reason.

Do not hide uncertainty. Do not pad the answer with hypothetical risks that do not apply to the actual strategy.

## Quick Start

Use `$confidence-loop` for a standard adversarial audit.

Use `$confidence-loop hard` when the user wants sub-agent second opinions before the final answer.
