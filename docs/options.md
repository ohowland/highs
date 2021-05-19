# HiGHS Options

## model_file
Model file  
type: string, advanced: false, default: ""

## presolve
Presolve option: "off", "choose" or "on"  
type: string, advanced: false, default: "choose"

## solver
Solver option: "simplex", "choose" or "ipm"  
type: string, advanced: false, default: "choose"

## parallel
Parallel option: "off", "choose" or "on"  
type: string, advanced: false, default: "choose"

## time_limit
Time limit  
type: double, advanced: false, range: [0, inf], default: inf

## infinite_cost
Limit on cost coefficient: values larger than this will be treated as infinite  
type: double, advanced: false, range: [1e+15, 1e+25], default: 1e+20

## infinite_bound
Limit on |constraint bound|: values larger than this will be treated as infinite  
type: double, advanced: false, range: [1e+15, 1e+25], default: 1e+20

## small_matrix_value
Lower limit on |matrix entries|: values smaller than this will be treated as zero  
type: double, advanced: false, range: [1e-12, inf], default: 1e-09

## large_matrix_value
Upper limit on |matrix entries|: values larger than this will be treated as infinite  
type: double, advanced: false, range: [1, 1e+20], default: 1e+15

## primal_feasibility_tolerance
Primal feasibility tolerance  
type: double, advanced: false, range: [1e-10, inf], default: 1e-07

## dual_feasibility_tolerance
Dual feasibility tolerance  
type: double, advanced: false, range: [1e-10, inf], default: 1e-07

## dual_objective_value_upper_bound
Upper bound on objective value for dual simplex: algorithm terminates if reached  
type: double, advanced: false, range: [-inf, inf], default: inf

## highs_debug_level
Debugging level in HiGHS  
type: int, advanced: false, range: {0, 3}, default: 0

## simplex_strategy
Strategy for simplex solver  
type: int, advanced: false, range: {0, 4}, default: 1

## simplex_scale_strategy
Strategy for scaling before simplex solver: off / on (0/1)  
type: int, advanced: false, range: {0, 5}, default: 2

## simplex_crash_strategy
Strategy for simplex crash: off / LTSSF / Bixby (0/1/2)  
type: int, advanced: false, range: {0, 9}, default: 0

## simplex_dual_edge_weight_strategy
Strategy for simplex dual edge weights: Dantzig / Devex / Steepest Edge (0/1/2)  
type: int, advanced: false, range: {0, 4}, default: 2

## simplex_primal_edge_weight_strategy
Strategy for simplex primal edge weights: Dantzig / Devex (0/1)  
type: int, advanced: false, range: {0, 1}, default: 0

## simplex_iteration_limit
Iteration limit for simplex solver  
type: int, advanced: false, range: {0, 2147483647}, default: 2147483647

## simplex_update_limit
Limit on the number of simplex UPDATE operations  
type: int, advanced: false, range: {0, 2147483647}, default: 5000

## ipm_iteration_limit
Iteration limit for IPM solver  
type: int, advanced: false, range: {0, 2147483647}, default: 2147483647

## highs_min_threads
Minimum number of threads in parallel execution  
type: int, advanced: false, range: {1, 8}, default: 1

## highs_max_threads
Maximum number of threads in parallel execution  
type: int, advanced: false, range: {1, 8}, default: 8

## message_level
HiGHS message level: bit-mask 1 => VERBOSE; 2 => DETAILED 4 => MINIMAL  
type: int, advanced: false, range: {0, 7}, default: 4

## solution_file
Solution file  
type: string, advanced: false, default: ""

## write_solution_to_file
Write the primal and dual solution to a file  
type: bool, advanced: false, range: {false, true}, default: false

## write_solution_pretty
Write the primal and dual solution in a pretty (human-readable) format  
type: bool, advanced: false, range: {false, true}, default: false

## mip_max_nodes
MIP solver max number of nodes  
type: int, advanced: false, range: {0, 2147483647}, default: 2147483647

## mip_report_level  
MIP solver reporting level
type: int, advanced: false, range: {0, 2}, default: 1
