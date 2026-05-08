---
name: confidence-loop
description: Stress-test a strategy, plan, implementation approach, or answer until remaining uncertainty is explicit and evidence-backed. Use when the user asks whether Codex is 100% confident, asks to find loopholes or failure modes, requests a confidence audit, or says to run a loop until the strategy is factually solid.
---

# Confidence Loop

## Rule

Treat "100% confident" as a demand for adversarial verification, not reassurance. Do not claim certainty unless every material assumption has been checked against evidence available in the current context or through allowed tools.

## Loop

1. State the current strategy in concrete terms.
2. List the success criteria and constraints.
3. Identify assumptions, dependencies, edge cases, loopholes, and ways the strategy can fail.
4. For each risk, decide whether it is disproven, accepted with rationale, or needs a fix.
5. Revise the strategy to close real loopholes.
6. Verify the revised strategy with the smallest relevant checks, source reads, tests, searches, or reasoning proofs available.
7. Repeat until no material unresolved loopholes remain, or until a blocker prevents factual certainty.

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
- **Verification**: checks performed and results.
- **Confidence**: 100%, high, or not confident, with the reason.

Do not hide uncertainty. Do not pad the answer with hypothetical risks that do not apply to the actual strategy.

## Example Trigger

"Are you 100% confident in this strategy? If not, find all possible loopholes, suggest proper fixes, and run this loop until you are factually 100% confident in the new strategy."
