# RULES.md — Borrower Copilot
 
Every rule below matches what's actually implemented in `rules/rules.go` and
`rules/calculate.go` — nothing here is aspirational. Format: `what / value /
why / source`.
 
## 1. Must vs Additional questions (final)
 
**Must (8) — the app works with only these, wide ranges:**
purpose, requested amount, income type (salaried/self-employed/informal),
net monthly income, existing EMIs, household expenses, age, credit score
(or "don't know").
 
**Additional — each demonstrably moves an output:**
 
| Question | What it changes |
|---|---|
| Income stability (years) | Nudges rate position within its bucket: ≥5 years earns a small discount, <1 year adds a small premium — a lender cares about tenure regardless of score |
| Recent bounce | Flips verdict to Don't Borrow if uncollateralized |
| Collateral + value | Changes product routing (→ LAP/Gold) and caps sanction via LTV |
| Undeclared cash income | Raises effective income (haircut applied) — self-employed/informal only |
| Co-applicant + income | Raises effective income (haircut applied) |
| Emergency savings (months) | Tightens or loosens the safe-amount cushion |
| Upcoming large expense | Reduces safe monthly capacity |
 
Adaptive: the cash-income question is only shown to self-employed/informal
borrowers — it's meaningless for a salaried applicant with a payslip.
 
## 2. FOIR (Fixed Obligation to Income Ratio) — the core affordability lever
 
| What | Value | Why | Source |
|---|---|---|---|
| Lender max FOIR (salaried/self-employed) | 50% | Most Indian banks/NBFCs approve up to 40-50%, some stretch to 55-65% for strong profiles | BankBazaar, HDFC, IDFC First FAQ pages, Sept 2026 |
| Borrower-safe FOIR (salaried/self-employed) | 35% | ~15pp headroom below the lender ceiling for savings + shocks | My judgement |
| Lender max FOIR (informal income) | 30% | No verified income stream -> both ceilings tightened | My judgement |
| Borrower-safe FOIR (informal income) | 25% | Same reasoning, more conservative | My judgement |
 
**Root logic:** the lender's number answers "will you default." The
borrower's number answers "will you still be okay." They must never be the
same number.
 
## 3. Effective income — why the raw income figure isn't used directly
 
| Rule | Value | Why |
|---|---|---|
| Undeclared cash income haircut | 50% counted | Self-employed/informal borrowers often understate income to lenders (e.g. ITR vs actual cash flow). Real, but unverified, so only half counts toward capacity. |
| Co-applicant income haircut | 70% counted | Real income, but adds joint liability the primary borrower doesn't fully control. |
 
`EffectiveIncome = NetMonthlyIncome + 0.5 x UndeclaredCashIncome + 0.7 x CoApplicantIncome`
 
## 4. Living-cost reality check (why "household expenses" is a must-question)
 
A favourable FOIR is meaningless if rent and living costs already consume
the paycheck. Safe monthly EMI capacity is the SMALLER of:
- `SafeFOIR x EffectiveIncome - ExistingEMIs` (the ratio-based ceiling), and
- `EffectiveIncome - ExistingEMIs - HouseholdExpenses - UpcomingLargeExpense/12` (what's actually left over)
**Emergency savings cushion:** if the borrower has under 3 months of
expenses saved, the leftover-based capacity above is further multiplied by
0.85 — a thinner cash cushion means less room should be allocated to a new
EMI. (My judgement; no published standard for this exact multiplier.)
 
**Handling a skipped answer:** if "emergency savings" or "upcoming expense"
is left blank, it's treated as its most conservative value (0 months saved,
Rs.0 upcoming expense) rather than a separate "unknown" bucket. This is a
deliberate difference from how credit score is handled (see section 6) —
for a safety-margin input, defaulting to caution when unanswered is itself
the safe assumption, whereas for credit score, defaulting to "bad" would
be actively unfair to someone who simply has no score.
 
## 5. Interest rate bands by product (headline, before fees) — Sept 2026
 
| Product | Low | High | Source |
|---|---|---|---|
| Home loan | 7.10% | 9.75% | Business Upturn, Bajaj Finserv |
| Loan against property (LAP) | 8.00% | 15.00% | Bajaj Finance, Agriwise, Bankermart |
| Gold loan | 8.55% | 24.00% | Cleartax, BankBazaar (trimmed from published 7-27% extremes to a realistic middle) |
| Personal loan (unsecured) | 10.50% | 24.00% | DealPlexus lending data, May 2026 |
| Business loan | 14.00% | 24.00% | Bajaj Markets, Cleartax |
| Vehicle loan (incl. two-wheeler/EV) | 11.00% | 20.00% | Assumption — no dedicated 2026 source found, estimated between personal-loan and gold-loan bands |
 
**Product routing:** purpose decides the default product; if the borrower
has collateral, business/consumption purposes route to LAP instead of
personal/business loans (cheaper, because now secured). Home and vehicle
purposes always route to their dedicated product.
 
## 6. Rate placement within a product's band
 
| Factor | Position in band (0=best, 1=worst) | Why |
|---|---|---|
| Credit score >=750 | 0.15 | Prime borrower gets close to the lender's best rate |
| Credit score 650-749 | 0.5 | Standard rate |
| Credit score <650 | 0.9 | Priced for risk |
| No score, but has collateral + formal/self-employed income | 0.5 | Collateral substitutes for missing score history (Ravi's case) |
| No score, no collateral, informal income | 0.9 | Nothing to substitute for the missing score (Anita's case) |
 
This is the direct code expression of "unknown is never zero": a missing
score is priced based on what else is known, never assumed to be bad by
default UNLESS nothing else offsets it.
 
## 7. Collateral loan-to-value (LTV) cap
 
For LAP and Gold loans, both the lender-sanction and borrower-safe amounts
are capped at **70% of declared collateral value**, regardless of how
strong the income-based affordability looks — this is a real, independent
constraint lenders apply. (Gold-loan LTV ~70-75% is published; property LTV
assumed similarly conservative here — My judgement where no single number
is cited.)
 
## 8. Tenure by product
 
| Product | Assumed tenure |
|---|---|
| Home loan / LAP | 15 years |
| Business loan | 5 years |
| Vehicle loan | 4 years |
| Gold loan | 1.5 years |
| Personal loan | 3 years |
 
Tenure is a property of the product, not a global constant — this fixed a
real bug during testing where a 3-year assumption made a legitimate
Rs.15L LAP request against Rs.45L collateral look unaffordable.
 
## 9. All-in cost (APR)
 
RBI's Key Facts Statement (KFS) mandate (effective Oct 1 2024, still in
force) requires lenders to disclose an all-in cost including fees, not
just the headline rate. We approximate this as:
 
`APR ~= nominal rate + (processing fee % / tenure in years)`
 
| Product | Assumed processing fee | Why |
|---|---|---|
| Home loan | 0.5% | Typically the lowest of all products |
| LAP | 1.0% | My judgement |
| Gold loan | 0.5% | Typically low, fast-disbursal product |
| Personal loan | 2.0% | Typically the highest of common products |
| Business loan | 2.0% | My judgement |
| Vehicle loan | 1.5% | My judgement |
 
These are illustrative defaults, not sourced per-lender — flagged
explicitly as judgement calls to defend live.
 
## 10. Verdict logic (O1), in order
 
1. If EXISTING obligations alone already exceed the safe FOIR -> **Don't
   Borrow** (fix existing debt first, regardless of the new request).
2. If a payment bounced in the last 12 months AND there's no collateral ->
   **Don't Borrow** (unsecured lenders will very likely decline anyway).
3. If the new request pushes projected FOIR past the LENDER max ->
   **Don't Borrow** (would likely be rejected).
4. If the new request pushes projected FOIR past the SAFE line but under
   the lender max -> **Borrow Less**, capped at the safe amount.
5. Otherwise -> **Borrow**, at the requested amount.
All three verdicts are confirmed reachable — see the three persona
run-throughs document.
 
## 11. Confidence / range width
 
- Full-width spread shrinks by up to 70% as more applicable additional
  questions get answered, down to a floor of +/-18% around the center rate.
- If credit score is unknown, the floor is widened to +/-35% regardless of
  how many other questions are answered — score is treated as
  irreducibly uncertain without it.
- The UI displays an explicit "wide range" warning whenever the spread is
  above the halfway point, so the borrower is never shown false precision.
## 12. Stress test (part of O4)
 
Two stress views are always shown alongside the EMI ceiling:
- EMI if the rate rises by a flat **+2 percentage points**.
- The EMI ceiling as a **% of income if income drops 20%** — flagged so
  the borrower can judge their own margin, not given a pass/fail verdict.
Both deltas (2pp, 20%) are round-number stand-ins, not derived from any
specific published volatility figure — My judgement, defend as
"reasonable illustrative shock," not empirical fact.
 
## 13. Known limitations (say these upfront in the interview, don't wait to be asked)
 
- No product beyond the six needed for the three personas is modelled.
- Rate bands and processing fees are illustrative 2026 market ranges, not
  live lender quotes — a real version would need a data source that updates
  independently of code changes.
- The income-stability nudge (+/-0.1 within a rate-position bucket) is a
  small, hand-picked adjustment — reasonable direction, unvalidated magnitude.
- No handling for a borrower who is already over-leveraged on a HIGH-value
  asset with no income at all (e.g. asset-rich, income-poor edge case) —
  the app would likely under-serve this borrower with the current rules.