package highs

type ModelStatus int

const (
	ModelNotset ModelStatus = iota
	ModelLoadError
	ModelError
	ModelPresolveError
	ModelSolveError
	ModelPostsolveError
	ModelEmpty
	ModelOptimal
	ModelInfeasible
	ModelUnboundedOrInfeasible
	ModelUnbounded
	ModelObjectiveBound
	ModelObjectiveTarget
	ModelTimeLimit
	ModelIterationLimit
	ModelUnknown
)

func (d ModelStatus) String() string {
	return [...]string{
		"Model Not Set",
		"Model Load Error",
		"Model Error",
		"Model Presolve Error",
		"Model Solve Error",
		"Model Postsolve Error",
		"Model Empty",
		"Model Optimal",
		"Model Infeasible",
		"Model Unbounded or Infeasible",
		"Model Unbounded",
		"Model Objetive Bound",
		"Model Objective Target",
		"Model TimeL imit",
		"Model Iteration Limit",
		"Model Unknown"}[d]
}

type SolutionStatus int

const (
	SolutionNone SolutionStatus = iota
	SolutionInfeasible
	SolutionFeasible
)

type Sense int

const (
	Maximize Sense = -1
	Minimize Sense = 1
)

type Integrality int

const (
	Continious Integrality = iota
	Integer
	ImplicitInteger
)
