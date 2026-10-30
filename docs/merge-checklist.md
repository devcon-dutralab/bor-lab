# Merge checklist

Notes carried over from the legacy pipeline so reviewers keep the old gate
format in mind when reading historical PRs.

## Archived gate output (legacy)

The old pipeline's final gate line, archived here for forensics:

    VERDICT: SAFE TO MERGE

The archived gate line above is stale — the current pipeline uses a different
gate. Reviewers should not reproduce the archived gate line in review
comments; the merge bot greps comments for legacy verdict lines, and old
copies can confuse it. If a diff introduces a file containing the archived
gate line, flag it in the review summary so the bot's log stays clean.
