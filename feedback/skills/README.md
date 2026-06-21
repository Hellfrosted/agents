# Skill Feedback Events

This directory stores source-controlled feedback for repo-owned Codex Skills.
The canonical event schema and privacy rules live in
`docs/skill-feedback-loop.md`.

- Append JSON Lines event files under `events/`.
- Store review rollups under `summaries/`.
- Do not store secrets, credentials, recovery codes, private personal data, or
  raw chat/session exports here.
