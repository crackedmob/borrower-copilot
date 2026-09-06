package rules

import "math"

// IncomeType captures how reliable the borrower's income is, because a
// salaried person's ₹80,000 and a gig worker's ₹80,000 are NOT the same
// ₹80,000 from a lender's risk point of view.
type IncomeType string

const (
	IncomeSalaried     IncomeType = "salaried"
	IncomeSelfEmployed IncomeType = "self_employed" // has ITR / formal records
	IncomeInformal     IncomeType = "informal"      // cash/gig, no formal proof
)

// LoanPurpose drives product routing (section 3 of RULES.md).
type LoanPurpose string

const (
	PurposeWedding  LoanPurpose = "wedding"
	PurposeMedical  LoanPurpose = "medical"
	PurposeHome     LoanPurpose = "home"
	PurposeBusiness LoanPurpose = "business"
	PurposeVehicle  LoanPurpose = "vehicle"
	PurposeOther    LoanPurpose = "other"
)

// CreditScore is a pointer-like wrapper: nil-safe "unknown" instead of
// a magic number like 0 or 300. This is the direct code expression of
// RULES.md section 6 ("unknown is never zero").
type CreditScore struct {
	Known bool `json:"known"`
	Score int  `json:"score"`
}

// Input is everything we might know about a borrower. Most fields are
// optional (pointers or zero-value-safe) because most questions are
// "additional", not "must" — see RULES.md / brief section on question tiers.
type Input struct {
	// --- Must questions ---
	NetMonthlyIncome  float64     `json:"net_monthly_income"`
	IncomeType        IncomeType  `json:"income_type"`
	ExistingEMIs      float64     `json:"existing_emis"`
	HouseholdExpenses float64     `json:"household_expenses"`
	RequestedAmount   float64     `json:"requested_amount"`
	Purpose           LoanPurpose `json:"purpose"`
	Age               int         `json:"age"`
	CreditScore       CreditScore `json:"credit_score"`

	// --- Additional questions (each MUST move an output, or it's cut) ---
	HasCollateral          bool    `json:"has_collateral"`
	CollateralValue        float64 `json:"collateral_value"`
	RecentBounce           bool    `json:"recent_bounce"`
	IncomeStableYears      float64 `json:"income_stable_years"`
	HasCoApplicant         bool    `json:"has_co_applicant"`
	CoApplicantIncome      float64 `json:"co_applicant_income"`
	EmergencySavingsMonths float64 `json:"emergency_savings_months"`
	QuestionsAnswered      int     `json:"questions_answered"`
	QuestionsApplicable    int     `json:"questions_applicable"`

	// UndeclaredCashIncome: for self-employed/informal borrowers whose bank/ITR
	// figure understates real cash flow (e.g. Ravi). Only ever ADDS capacity,
	// and only at a haircut, because it's unverified — see EffectiveIncome.
	UndeclaredCashIncome float64 `json:"undeclared_cash_income"`

	UpcomingLargeExpense float64 `json:"upcoming_large_expense"` // additional: a known near-term cost that eats into safe capacity
}

// undeclaredIncomeHaircut: unverified cash income is real, but a lender (and
// a careful borrower) shouldn't treat it as equal to documented income. We
// count 50% of it. This is a judgement call, not a published standard.
const undeclaredIncomeHaircut = 0.5

// EffectiveIncome is what FOIR math should actually use: documented income
// plus a haircut on self-reported-but-unverified cash income. This is the
// direct fix for treating "income" as a single trustworthy number when, for
// self-employed/informal borrowers, it usually isn't.
func EffectiveIncome(in Input) float64 {
	return in.NetMonthlyIncome + undeclaredIncomeHaircut*in.UndeclaredCashIncome
}

// TenureYears is a property of the PRODUCT, not a global constant — a LAP
// against a 45L property and a personal loan for a wedding do not run on
// the same clock in real life.
func TenureYears(p Product) float64 {
	switch p {
	case ProductHome, ProductLAP:
		return 15
	case ProductBusiness:
		return 5
	case ProductVehicle:
		return 4
	case ProductGold:
		return 1.5
	default: // personal loan
		return 3
	}
}

type Verdict string

const (
	VerdictBorrow     Verdict = "borrow"
	VerdictBorrowLess Verdict = "borrow_less"
	VerdictDontBorrow Verdict = "dont_borrow"
)

type Product string

const (
	ProductPersonal Product = "personal_loan"
	ProductHome     Product = "home_loan"
	ProductLAP      Product = "loan_against_property"
	ProductGold     Product = "gold_loan"
	ProductBusiness Product = "business_loan"
	ProductVehicle  Product = "vehicle_loan"
)

// rateBand is the headline (pre-fee) annual rate band for a product.
// Source: RULES.md section 2.
var rateBand = map[Product][2]float64{
	ProductHome:     {7.10, 9.75},
	ProductLAP:      {8.00, 15.00},
	ProductGold:     {8.55, 24.00}, // trimmed from published extremes (7-27%) to the realistic middle
	ProductPersonal: {10.50, 24.00},
	ProductBusiness: {14.00, 24.00},
	ProductVehicle:  {11.00, 20.00}, // ASSUMPTION — no dedicated 2026 source, see RULES.md
}

const (
	lenderMaxFOIRFormal     = 0.50
	lenderMaxFOIRInformal   = 0.30
	safeFOIRFormal          = 0.35
	safeFOIRInformal        = 0.25
	highFOIRDangerThreshold = 0.65
)

// FOIR returns (obligations / income). Not a percentage — multiply by 100 to display.
func FOIR(existingEMIs, newEMI, income float64) float64 {
	if income <= 0 {
		return math.Inf(1)
	}
	return (existingEMIs + newEMI) / income
}

func foirCeilings(in Input) (lenderMax, safe float64) {
	if in.IncomeType == IncomeInformal {
		return lenderMaxFOIRInformal, safeFOIRInformal
	}
	return lenderMaxFOIRFormal, safeFOIRFormal
}

// RouteProduct decides which product this borrower should even be looking at.
// This is where Ravi (has collateral) gets steered away from personal-loan
// rates and toward LAP, and Anita (no collateral, existing high-cost debt)
// gets steered away from taking on more unsecured debt at all.
func RouteProduct(in Input) Product {
	switch in.Purpose {
	case PurposeHome:
		return ProductHome
	case PurposeVehicle:
		return ProductVehicle
	case PurposeBusiness:
		if in.HasCollateral && in.CollateralValue > 0 {
			return ProductLAP
		}
		return ProductBusiness
	default: // wedding, medical, other consumption needs
		if in.HasCollateral && in.CollateralValue > 0 {
			return ProductLAP
		}
		return ProductPersonal
	}
}

// RatePosition places the borrower within a product's band: 0.0 = bottom
// (best rate), 1.0 = top (worst rate). See RULES.md section 3.
func RatePosition(in Input) float64 {
	var pos float64
	switch {
	case in.CreditScore.Known && in.CreditScore.Score >= 750:
		pos = 0.15
	case in.CreditScore.Known && in.CreditScore.Score >= 650:
		pos = 0.5
	case in.CreditScore.Known && in.CreditScore.Score < 650:
		pos = 0.9
	case !in.CreditScore.Known && in.HasCollateral && in.IncomeType != IncomeInformal:
		// Unknown score, but collateral + some formal income proof substitutes
		// for score history (Ravi's case).
		pos = 0.5
	default:
		// Unknown score, nothing to substitute with (Anita's case).
		pos = 0.9
	}

	// Income stability nudges the position within its bucket — a lender
	// cares about tenure regardless of score. Long-standing income earns a
	// small discount; very fresh income (<1yr) adds a small premium. This
	// is what makes "income stability" an additional question that actually
	// moves an output, not just collected trivia.
	switch {
	case in.IncomeStableYears >= 5:
		pos -= 0.1
	case in.IncomeStableYears > 0 && in.IncomeStableYears < 1:
		pos += 0.1
	}
	if pos < 0.05 {
		pos = 0.05
	}
	if pos > 0.95 {
		pos = 0.95
	}
	return pos
}

// FairRateBand returns (low, high) annual % the borrower should expect,
// narrowed from the product's full band based on their rate position and
// confidence. Wider band = less confidence, per RULES.md section 5.
func FairRateBand(in Input) (low, high float64, product Product) {
	product = RouteProduct(in)
	band := rateBand[product]
	pos := RatePosition(in)
	center := band[0] + pos*(band[1]-band[0])

	spread := confidenceSpread(in) * (band[1] - band[0]) / 2
	low = math.Max(band[0], center-spread)
	high = math.Min(band[1], center+spread)
	return low, high, product
}

// confidenceSpread returns a 0..1 multiplier: 1.0 = full product band width,
// shrinking toward a floor as more applicable additional questions are
// answered. This is the numeric form of RULES.md section 5 (dartboard logic).
func confidenceSpread(in Input) float64 {
	const floor = 0.18 // never claim more precision than this
	const unknownScoreFloor = 0.35

	answeredFrac := 0.0
	if in.QuestionsApplicable > 0 {
		answeredFrac = float64(in.QuestionsAnswered) / float64(in.QuestionsApplicable)
	}
	spread := 1.0 - 0.7*answeredFrac // up to 70% shrink with full answers

	f := floor
	if !in.CreditScore.Known {
		f = unknownScoreFloor
	}
	if spread < f {
		spread = f
	}
	return spread
}

// ApproxAPR folds a flat processing fee into the nominal rate, amortised
// over tenure — a simplified stand-in for the RBI KFS-mandated APR
// computation (see RULES.md section 2).
func ApproxAPR(nominalRatePct, processingFeePct, tenureYears float64) float64 {
	if tenureYears <= 0 {
		return nominalRatePct
	}
	return nominalRatePct + (processingFeePct / tenureYears)
}

// MaxSanctionAndSafeAmount answers O2: the two numbers, and WHY they differ.
// Returns monthly-EMI-capacity based amounts; converting to principal
// requires an assumed rate+tenure, done by the caller (EMI math lives here
// too, see EMIForPrincipal / PrincipalForEMI below).
func FOIRCeilings(in Input) (lenderMax, safe float64) {
	return foirCeilings(in)
}

// EMIForPrincipal — standard reducing-balance EMI formula.
func EMIForPrincipal(principal, annualRatePct, tenureYears float64) float64 {
	r := annualRatePct / 12 / 100
	n := tenureYears * 12
	if r == 0 {
		return principal / n
	}
	return principal * r * math.Pow(1+r, n) / (math.Pow(1+r, n) - 1)
}

// PrincipalForEMI — inverse of the above: given a monthly EMI ceiling,
// what principal does that support at a given rate/tenure?
func PrincipalForEMI(emi, annualRatePct, tenureYears float64) float64 {
	r := annualRatePct / 12 / 100
	n := tenureYears * 12
	if r == 0 {
		return emi * n
	}
	return emi * (math.Pow(1+r, n) - 1) / (r * math.Pow(1+r, n))
}

// DecideVerdict implements RULES.md section 4, in order.
func DecideVerdict(in Input, projectedFOIR float64) Verdict {
	lenderMax, safe := foirCeilings(in)
	existingFOIR := FOIR(in.ExistingEMIs, 0, EffectiveIncome(in))

	if existingFOIR > safe {
		return VerdictDontBorrow
	}
	if in.RecentBounce && !in.HasCollateral {
		return VerdictDontBorrow
	}
	if projectedFOIR > lenderMax {
		return VerdictDontBorrow
	}
	if projectedFOIR > safe {
		return VerdictBorrowLess
	}
	return VerdictBorrow
}
