# Product

## Register

product

## Users

The primary user is the workstation operator maintaining this local Codex
tooling setup. Secondary users are future agents and maintainers who need to
understand current behavior, make focused source changes, and verify the active
workstation install without relying on remembered local state.

## Product Purpose

agents-toolkit keeps Codex-adjacent workstation tooling understandable,
repairable, and reproducible. It is the source of truth for Windows-to-WSL
Codex launch behavior, app-server proxy behavior, global skills maintenance
wrappers, repo-owned local skills, and operator documentation for companion
tools and local archives.

Success means a maintainer can quickly identify the file that owns a behavior,
change it in source first, validate it with the smallest relevant check, and
only then copy or repair active install files when explicitly requested.

## Brand Personality

Precise, repairable, operator-focused.

The project should feel like practical workstation infrastructure: clear enough
to audit, calm under failure, and specific about boundaries between source,
installed runtime files, plugin caches, local docs, and generated artifacts.

## Anti-references

This should not feel like generic SaaS polish, decorative AI-generated UI, or a
broad workflow redesign without a concrete request. Avoid stale setup
archaeology, duplicated flag lists, vague automation provenance, and ornamental
visual treatment that makes operational facts harder to scan.

Avoid the generic AI gradient feel at all costs. Gradients are acceptable only
when they serve a deliberate UI or visualization purpose and do not create the
common synthetic gradient aesthetic.

## Design Principles

- Source-first, install-second: every behavior change starts in this repo, and
  active workstation files are repaired or republished only on request.
- Make ownership obvious: docs and UI artifacts should help a maintainer find
  the file, contract, or verification step that owns the current behavior.
- Preserve operator flow: prefer direct, low-risk action, but escalate
  destructive, credential-related, architecture-shaping, or operationally risky
  choices.
- Keep context durable and bounded: store only privacy-safe facts that improve
  future work; keep scratch notes, raw logs, secrets, and private data out of
  docs and commits.
- Favor auditability over flourish: clear hierarchy, exact wording, and
  repeatable checks matter more than expressive styling.

## Accessibility & Inclusion

Design and documentation should be readable, keyboard-friendly where
interactive, and respectful of reduced-motion preferences. User-facing HTML
should maintain strong contrast, avoid horizontal overflow on mobile, and keep
primary workflows usable without relying on animation or color alone.
