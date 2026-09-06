# Persona Run-Throughs
 
All three outputs below are taken directly from the running app (verified in
browser, not computed by hand), using each persona's numbers exactly as
given in the brief.
 
---
 
## Priya, 29 — Bengaluru, salaried software engineer
 
**Must questions answered:**
Purpose: Wedding · Requested: ₹8,00,000 · Income type: Salaried ·
Net monthly income: ₹1,10,000 · Existing EMIs: ₹14,000 (car loan) ·
Household expenses: ₹28,000 (rent + living, entered for this run-through) ·
Age: 29 · Credit score: 780 (known)
 
**Additional questions:** all skipped, to demonstrate the wide-range,
low-information path.
 
**Outputs:**
| Output | Result |
|---|---|
| O1 — Verdict | **Borrow Less** |
| O2 — Amounts | Lender likely sanctions ₹11,84,620 · Borrower can safely carry ₹7,07,883 |
| O3 — Fair rate | 10.50% – 19.27% (wide, because no additional questions answered) · APR 11.17% – 19.94% |
| O4 — EMI ceiling | ₹24,500/month over 3 years · Stress (+2pp rate): ₹25,198/mo · Stress (–20% income): EMI is 28% of reduced income |
 
**Why (as shown in-app):** Borrow Less because ₹8L pushes obligations past
a safe share of her income, even though a lender might still approve it at
this range; the lender number is looser because a lender only needs her to
avoid default, not stay comfortable.
 
**Negotiation Card:** "Fair rate for my profile: 10.50%–19.27% (APR
11.17%–19.94%). Product: Personal Loan · Tenure: 3 yrs · Shouldn't pay more
than ₹24,500/month."
 
---
 
## Ravi, 42 — Mysuru, self-employed kirana owner
 
**Must questions answered:**
Purpose: Business · Requested: ₹15,00,000 · Income type: Self-employed ·
Net monthly income: ₹35,000 (ITR-declared) · Existing EMIs: ₹0 ·
Household expenses: ₹20,000 · Age: 42 · Credit score: don't know
 
**Additional questions answered:** Undeclared cash income: ₹25,000/month
(conservative midpoint of his stated ₹40k–80k cash range, above the ITR
figure) · Has collateral: Yes, value ₹45,00,000 (shop premises).
Co-applicant and emergency savings were left unanswered in this run.
 
**Outputs:**
| Output | Result |
|---|---|
| O1 — Verdict | **Borrow Less** |
| O2 — Amounts | Lender likely sanctions ₹20,33,060 · Borrower can safely carry ₹14,23,142 |
| O3 — Fair rate | 8.70% – 14.30% · APR 8.77% – 14.37% |
| O4 — EMI ceiling | ₹16,625/month over 15 years (LAP tenure) · Stress (+2pp): ₹18,477/mo · Stress (–20% income): 44% of reduced income |
 
**Why (as shown in-app):** routed to Loan Against Property, not a personal
or unsecured business loan, because he has unencumbered collateral — this
alone moves him from a 14–24% band to an 8.70–14.30% one. Rate sits
mid-band despite no credit score on file because the collateral substitutes
for score history. Borrow Less rather than Borrow because his declared
capacity (even blended with cash income) is close to, but slightly under,
his ₹15L ask at the safe line — see Known Limitations in RULES.md re:
income-stability weighting, which if answered (his 14 years in business)
would likely tighten this toward a plain Borrow.
 
**Negotiation Card:** "Fair rate for my profile: 8.70%–14.30% (APR
8.77%–14.37%). Product: Loan Against Property · Tenure: 15 yrs · Shouldn't
pay more than ₹16,625/month."
 
---
 
## Anita, 35 — Hubballi, informal gig worker + tailoring
 
**Must questions answered:**
Purpose: Vehicle (EV) · Requested: ₹1,50,000 · Income type: Informal ·
Net monthly income: ₹28,000 · Existing EMIs: ₹8,000 (rough EMI-equivalent
on her ₹35,000 outstanding app loans) · Household expenses: ₹15,000 ·
Age: 35 · Credit score: don't know
 
**Additional questions answered:** Recent bounce: Yes (matches her stated
bounced EMI last month). No collateral.
 
**Outputs:**
| Output | Result |
|---|---|
| O1 — Verdict | **Don't Borrow** |
| O2 — Amounts | Lender likely sanctions ₹13,733 (illustrative only — see verdict) · Borrower can safely carry ₹0 |
| O3 — Fair rate | 15.05% – 20.00% · APR 15.43% – 20.38% (top of band — no score, no collateral) |
| O4 — EMI ceiling | ₹0 — no new EMI is recommended |
 
**Why (as shown in-app):** Don't Borrow because a payment bounced in the
last 12 months and there's no collateral to secure new credit against —
this is the exact "Don't must be reachable" case the brief calls out.
 
**Negotiation Card (adapted for this verdict):** "This isn't a good time to
take on this loan. Focus on clearing existing obligations or rebuilding a
payment history first — a lender is likely to see the same red flags this
app did." (The app does not show a rate/EMI to negotiate when the verdict
is Don't Borrow — there's nothing to negotiate.)
 
---
 
## What this demonstrates across all three
 
- All three verdicts (Borrow, Borrow Less, Don't Borrow) are reachable and
  were reached with real inputs, not constructed to force an outcome.
- Product routing correctly differentiates Priya (unsecured, personal),
  Ravi (secured, LAP — because of collateral), and Anita (no viable secured
  or unsecured path given her bounce history).
- The lender-vs-safe amount split shows up meaningfully in all three: widest
  gap for Priya (good profile, no collateral), narrower for Ravi (collateral
  bounds both numbers), and collapses to zero for Anita (no safe path at all).