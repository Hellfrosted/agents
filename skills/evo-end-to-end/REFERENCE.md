# Evo End To End Reference

## Installation Diagnostics

Use the installed Evo plugin skills and the matching
[upstream release](https://github.com/evo-hq/evo/releases) for workflow and command
details. Check the host with `evo doctor <host>` and confirm repair flags with
`evo install --help` or `evo update --help`.

For an approved Codex plugin repair, use
`evo install codex --force --no-trust-hooks`. If the CLI and plugin versions
differ, align them to the chosen release as part of the approved repair. Run
`evo doctor codex` again. Have the user review hooks through `/hooks` before
enabling trust. Without trusted hooks, mid-run directives are unavailable.

## Hosted CI

- Score the required validation path from its first job start to its final job
  finish. Record queue time separately. Verify the actual checked-out commit
  or tree, including for pull request merge checkouts.
- Pair control and candidate runs by runner image and cache condition. Keep
  warm and cold results separate. Compare repeated measurements against the
  promotion threshold even if Evo already committed a favorable attempt.
- Gate on the exact test IDs, required checks, relevant execution settings,
  outputs, and failure propagation. Use exit status to signal gate failures.

## Dashboard Access

1. On the dashboard host, verify Evo responds on its reported loopback port.
   Check `tailscale serve status` and choose an unused HTTPS port.
2. With the host's authorized operator privileges, run
   `tailscale serve --bg --https=<serve-port> http://127.0.0.1:<dashboard-port>`.
   Verify HTTPS access through the tailnet and share the returned URL.
3. Record the host, HTTPS port, and backend URL in durable campaign state, such
   as an Evo workspace note. Keep the route during hosted CI waits.
4. On terminal handback or resume after an interruption, check the recorded
   route against `tailscale serve status`. If monitoring has ended and the
   listener still points to the recorded backend, run
   `tailscale serve --https=<serve-port> off` and verify its removal. Preserve
   other listeners. Report cleanup failures so a later agent can retry them.
