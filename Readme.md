# Borrower Copilot
 
A self-assessment tool that helps a borrower answer four questions before
they walk into a lender: should I borrow, how much, at what rate, and what
EMI should I agree to — plus a one-page card to negotiate with.
 
## Run it (under 2 minutes)
 
Requires Go 1.21+ (check with `go version`; install from https://go.dev/dl
if missing).
 
```
cd borrower-copilot
go run .
```
 
Open http://localhost:8080. No build step, no dependencies, no database —
`go run .` compiles and starts the server in one command.
 
## Run the rules-engine tests
 
```
go test ./rules/... -run TestPersonas -v
```
 
Runs the three brief personas (Priya, Ravi, Anita) straight through the
engine and prints every output — useful for checking a rule change without
touching the browser.
 
## Project layout
 
```
main.go            HTTP server: serves web/, exposes POST /api/calculate
rules/rules.go      Domain types, FOIR math, product routing, rate bands, verdict logic
rules/calculate.go  Calculate() — the single entry point tying it into O1-O4 + explanations
rules/rules_test.go Persona sanity tests
web/index.html      Wizard UI (must questions -> adaptive questions -> results + Negotiation Card)
RULES.md            Every threshold: what / value / why / source
```
 
Rules are deliberately isolated in `rules/` with zero HTTP/HTML/JSON
awareness — the package only knows business logic. This is so any rule can
be changed and re-tested (via `rules_test.go`) without touching the server
or UI at all.
 
## What this does NOT do (by design, not oversight)
 
- No login, no database, no data persistence between runs — nothing is
  stored, per the brief.
- No real credit bureau integration and no machine learning — the brief
  explicitly does not score these.
- Loan products are limited to what the three personas need (personal,
  home, LAP, gold, business, vehicle) — not exhaustive.