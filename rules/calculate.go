package rules
 
import "math"
 
// processingFeePct by product — used only for APR, per RBI KFS logic.
// Source: RULES.md — typical range across lenders is 0.5%-3%; using the
// lower-middle of that range as a default assumption (My judgement).
var processingFeePct = map[Product]float64{
	ProductHome:     0.5,
	ProductLAP:      1.0,
	ProductGold:     0.5,
	ProductPersonal: 2.0,
	ProductBusiness: 2.0,
	ProductVehicle:  1.5,
}
 
// collateralLTV caps how much of a pledged asset's value a lender will
// actually advance — a real constraint independent of income affordability.
// Source: gold-loan LTV ~70-75% (BankBazaar); property LTV assumed similarly
// conservative for this app. My judgement where no single number is cited.
const collateralLTV = 0.70
 
// Output is everything the UI needs to render O1-O4 and the Negotiation
// Card. Every field here must trace back to a one-sentence "why" — that's
// what Explanations is for.
type Output struct {
	Verdict Verdict `json:"verdict"`
	Product Product `json:"product"`
 
	LenderSanctionAmount float64 `json:"lender_sanction_amount"`
	BorrowerSafeAmount   float64 `json:"borrower_safe_amount"`
	RecommendedAmount    float64 `json:"recommended_amount"`
 
	RateLow  float64 `json:"rate_low"`
	RateHigh float64 `json:"rate_high"`
	APRLow   float64 `json:"apr_low"`
	APRHigh  float64 `json:"apr_high"`
	TenureYears       float64 `json:"tenure_years"`
 
	EMICeiling       float64 `json:"emi_ceiling"`
	StressEMI        float64 `json:"stress_emi"`
	StressIncomeDrop float64 `json:"stress_income_drop_pct"`
 
	ProjectedFOIR float64 `json:"projected_foir"`
	ConfidenceLow bool    `json:"confidence_low"`
 
	Explanations []string `json:"explanations"`
}
 
// Calculate is the single entry point the API/UI should call. Give it a
// fully or partially filled Input, get back everything needed to render
// the four outputs and the Negotiation Card.
func Calculate(in Input) Output {
	product := RouteProduct(in)
	tenure := TenureYears(product)
	low, high, _ := FairRateBand(in)
	assumedRate := (low + high) / 2
	fee := processingFeePct[product]
 
	income := EffectiveIncome(in)
	if in.HasCoApplicant {
		// Co-applicant income counts at a haircut too — it's real, but adds
		// joint liability the primary borrower doesn't fully control.
		income += 0.7 * in.CoApplicantIncome
	}
 
	lenderMaxFOIR, safeFOIR := foirCeilings(in)
 
	// Living-cost reality check: even a favourable FOIR doesn't help if rent
	// and expenses already consume the paycheck. This is what makes
	// "household expenses" (a must-question) actually load-bearing.
	disposableAfterLiving := income - in.ExistingEMIs - in.HouseholdExpenses - in.UpcomingLargeExpense/12
 
	// Thin safety margin if the borrower has under 3 months of emergency
	// savings — this is what makes that additional question load-bearing.
	cushionFactor := 1.0
	if in.EmergencySavingsMonths < 3 {
		cushionFactor = 0.85
	}
 
	safeMonthlyCapacity := math.Min(safeFOIR*income-in.ExistingEMIs, disposableAfterLiving*cushionFactor)
	lenderMonthlyCapacity := lenderMaxFOIR*income - in.ExistingEMIs
	safeMonthlyCapacity = math.Max(0, safeMonthlyCapacity)
	lenderMonthlyCapacity = math.Max(0, lenderMonthlyCapacity)
 
	safeAmount := PrincipalForEMI(safeMonthlyCapacity, assumedRate, tenure)
	lenderAmount := PrincipalForEMI(lenderMonthlyCapacity, assumedRate, tenure)
 
	// Collateral caps what a secured product will ever advance, regardless
	// of how strong the income picture looks.
	if (product == ProductLAP || product == ProductGold) && in.CollateralValue > 0 {
		cap := in.CollateralValue * collateralLTV
		safeAmount = math.Min(safeAmount, cap)
		lenderAmount = math.Min(lenderAmount, cap)
	}
 
	requestedEMI := EMIForPrincipal(in.RequestedAmount, assumedRate, tenure)
	projectedFOIR := FOIR(in.ExistingEMIs, requestedEMI, income)
	verdict := DecideVerdict(in, projectedFOIR)
 
	var recommended float64
	switch verdict {
	case VerdictBorrow:
		recommended = in.RequestedAmount
	case VerdictBorrowLess:
		recommended = math.Min(safeAmount, in.RequestedAmount)
	case VerdictDontBorrow:
		recommended = 0
	}
 
	emiCeiling := EMIForPrincipal(recommended, assumedRate, tenure)
	stressEMI := EMIForPrincipal(recommended, assumedRate+2, tenure)
	var stressIncomeDropPct float64
	if income > 0 {
		reducedIncome := income * 0.8
		if reducedIncome > 0 {
			stressIncomeDropPct = (emiCeiling / reducedIncome) * 100
		}
	}
 
	out := Output{
		Verdict:               verdict,
		Product:               product,
		LenderSanctionAmount:  lenderAmount,
		BorrowerSafeAmount:    safeAmount,
		RecommendedAmount:     recommended,
		RateLow:               low,
		RateHigh:              high,
		APRLow:                ApproxAPR(low, fee, tenure),
		APRHigh:               ApproxAPR(high, fee, tenure),
		TenureYears:           tenure,
		EMICeiling:            emiCeiling,
		StressEMI:             stressEMI,
		StressIncomeDrop:      stressIncomeDropPct,
		ProjectedFOIR:         projectedFOIR * 100,
		ConfidenceLow:         confidenceSpread(in) > 0.5,
	}
	out.Explanations = explain(in, out)
	return out
}
 
// explain produces the one-sentence-per-output traceability the brief
// explicitly requires ("the borrower must be able to read, in one sentence,
// why the ceiling is X and not Y").
func explain(in Input, out Output) []string {
	var e []string
	switch out.Verdict {
	case VerdictDontBorrow:
		if in.RecentBounce {
			e = append(e, "Verdict is Don't Borrow because a payment bounced in the last 12 months and there's no collateral to secure new credit against.")
		} else {
			e = append(e, "Verdict is Don't Borrow because existing obligations already exceed a safe share of income — new debt would make that worse, not better.")
		}
	case VerdictBorrowLess:
		e = append(e, "Verdict is Borrow Less because the requested amount would push monthly obligations above a safe share of income, even though a lender might still approve it.")
	case VerdictBorrow:
		e = append(e, "Verdict is Borrow because the requested amount stays within a safe share of income after existing obligations and living expenses.")
	}
	e = append(e, "The lender-likely amount uses a looser affordability limit because a lender only needs you to avoid default, not stay comfortable — the safe amount uses a tighter limit that leaves room for living costs and shocks.")
	if !in.CreditScore.Known && in.HasCollateral {
		e = append(e, "The rate sits mid-band despite no credit score on file because pledged collateral substitutes for score history.")
	} else if !in.CreditScore.Known {
		e = append(e, "The rate sits at the higher end of the band because there's no credit score or collateral to price the risk more precisely.")
	}
	e = append(e, "The EMI ceiling already accounts for existing obligations and household expenses, so it's what's left over — not a share of gross income.")
	return e
}