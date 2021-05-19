#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include "highs-interface.h"

extern "C" {
  highs_obj* highsiface_create(void) {
    void* highs = Highs_create();
    return (highs_obj*)highs;
  }

  extern void highsiface_free(highs_obj* highs) {
    delete (highs_obj*)highs;
  }

  extern HighsInt highsiface_add_cols(highs_obj* highs, HighsInt numCol, double* colCost, double* colLower, double* colUpper) {
    HighsInt r = Highs_addCols(highs, numCol, colCost, colLower, colUpper, 0, NULL, NULL, NULL);
    return r;
  }
  
  extern HighsInt highsiface_add_rows(highs_obj* highs, HighsInt numRow, double* rowLower, double* rowUpper, HighsInt numNz, HighsInt* arStart, HighsInt* arIndex, double* arValue) {
    HighsInt r = Highs_addRows(highs, numRow, rowLower, rowUpper, numNz, arStart, arIndex, arValue);
    return r;
  }

  extern HighsInt highsiface_run(highs_obj* highs) {
    Highs_setBoolOptionValue(highs, "output_flag", 1); // for testing run loudly
    HighsInt r = Highs_run(highs);
    return r;
  }

  extern void highsiface_get_solution(highs_obj* highs, double* colValue, double* colDual, double* rowValue, double* rowDual) {
    Highs_getSolution(highs, colValue, colDual, rowValue, rowDual);
  }
}

/*
void full_api() {
  // Form and solve the LP
  // Min    f  = 2x_0 + 3x_1
  // s.t.                x_1 <= 6
  //       10 <=  x_0 + 2x_1 <= 14
  //        8 <= 2x_0 +  x_1
  // 0 <= x_0 <= 3; 1 <= x_1

  void* highs;

  highs = Highs_create();

  const int numcol = 2;
  const int numrow = 3;
  const int numnz = 5;
  int i;

  // Define the column costs, lower bounds and upper bounds
  double colcost[numcol] = {2.0, 3.0};
  double collower[numcol] = {0.0, 1.0};
  double colupper[numcol] = {3.0, 1.0e30};
  // Define the row lower bounds and upper bounds
  double rowlower[numrow] = {-1.0e30, 10.0, 8.0};
  double rowupper[numrow] = {6.0, 14.0, 1.0e30};
  // Define the constraint matrix row-wise, as it is added to the LP
  // with the rows
  int arstart[numrow] = {0, 1, 3};
  int arindex[numnz] = {1, 0, 1, 0, 1};
  double arvalue[numnz] = {1.0, 1.0, 2.0, 2.0, 1.0};

  double* colvalue = (double*)malloc(sizeof(double) * numcol);
  double* coldual = (double*)malloc(sizeof(double) * numcol);
  double* rowvalue = (double*)malloc(sizeof(double) * numrow);
  double* rowdual = (double*)malloc(sizeof(double) * numrow);

  int* colbasisstatus = (int*)malloc(sizeof(int) * numcol);
  int* rowbasisstatus = (int*)malloc(sizeof(int) * numrow);

  // Add two columns to the empty LP
  assert( Highs_addCols(highs, numcol, colcost, collower, colupper, 0, NULL, NULL, NULL) );
  // Add three rows to the 2-column LP
  assert( Highs_addRows(highs, numrow, rowlower, rowupper, numnz, arstart, arindex, arvalue) );

  int sense;
  Highs_getObjectiveSense(highs, &sense);
  printf("LP problem has objective sense = %d\n", sense);
  assert(sense == 1);

  sense *= -1;
  Highs_changeObjectiveSense(highs, sense);
  assert(sense == -1);

  sense *= -1;
  Highs_changeObjectiveSense(highs, sense);

  Highs_getObjectiveSense(highs, &sense);
  printf("LP problem has old objective sense = %d\n", sense);
  assert(sense == 1);

  int simplex_scale_strategy;
  Highs_getIntOptionValue(highs, "simplex_scale_strategy", &simplex_scale_strategy);
  printf("simplex_scale_strategy = %d: setting it to 3\n", simplex_scale_strategy);
  simplex_scale_strategy = 3;
  Highs_setIntOptionValue(highs, "simplex_scale_strategy", simplex_scale_strategy);

  // There are some functions to check what type of option value you should
  // provide.
  int option_type;
  int ret;
  ret = Highs_getOptionType(highs, "simplex_scale_strategy", &option_type);
  assert(ret == 0);
  assert(option_type == 1);
  ret = Highs_getOptionType(highs, "bad_option", &option_type);
  assert(ret != 0);

  double primal_feasibility_tolerance;
  Highs_getDoubleOptionValue(highs, "primal_feasibility_tolerance", &primal_feasibility_tolerance);
  printf("primal_feasibility_tolerance = %g: setting it to 1e-6\n", primal_feasibility_tolerance);
  primal_feasibility_tolerance = 1e-6;
  Highs_setDoubleOptionValue(highs, "primal_feasibility_tolerance", primal_feasibility_tolerance);

  Highs_setBoolOptionValue(highs, "output_flag", 0);
  printf("Running quietly...\n");
  int runstatus = Highs_run(highs);
  printf("Running loudly...\n");
  Highs_setBoolOptionValue(highs, "output_flag", 1);

  // Get the model status
  int modelstatus = Highs_getModelStatus(highs);

  printf("Run status = %d; Model status = %d\n", runstatus, modelstatus);

  double objective_function_value;
  Highs_getDoubleInfoValue(highs, "objective_function_value", &objective_function_value);
  int simplex_iteration_count = 0;
  Highs_getIntInfoValue(highs, "simplex_iteration_count", &simplex_iteration_count);
  int primal_solution_status = 0;
  Highs_getIntInfoValue(highs, "primal_solution_status", &primal_solution_status);
  int dual_solution_status = 0;
  Highs_getIntInfoValue(highs, "dual_solution_status", &dual_solution_status);

  printf("Objective value = %g; Iteration count = %d\n", objective_function_value, simplex_iteration_count);
  if (modelstatus == 7) {
    printf("Solution primal status = %d\n", primal_solution_status);
    printf("Solution dual status = %d\n", dual_solution_status);
    // Get the primal and dual solution
    Highs_getSolution(highs, colvalue, coldual, rowvalue, rowdual);
    // Get the basis
    Highs_getBasis(highs, colbasisstatus, rowbasisstatus);
    // Report the column primal and dual values, and basis status
    for (i = 0; i < numcol; i++) {
      printf("Col%d = %lf; dual = %lf; status = %d; \n", i, colvalue[i], coldual[i], colbasisstatus[i]);
    }
    // Report the row primal and dual values, and basis status
    for (i = 0; i < numrow; i++) {
      printf("Row%d = %lf; dual = %lf; status = %d; \n", i, rowvalue[i], rowdual[i], rowbasisstatus[i]);
    }
  }

  free(colvalue);
  free(coldual);
  free(rowvalue);
  free(rowdual);
  free(colbasisstatus);
  free(rowbasisstatus);

  Highs_destroy(highs);

  // Define the constraint matrix col-wise to pass to the LP
  int rowwise = 0;
  int astart[numcol] = {0, 2};
  int aindex[numnz] = {1, 2, 0, 1, 2};
  double avalue[numnz] = {1.0, 2.0, 1.0, 2.0, 1.0};
  highs = Highs_create();
  runstatus = Highs_passLp(highs, numcol, numrow, numnz, rowwise,
			colcost, collower, colupper,
			rowlower, rowupper,
			astart, aindex, avalue);
  runstatus = Highs_run(highs);
  modelstatus = Highs_getModelStatus(highs);
  printf("Run status = %d; Model status = %d\n", runstatus, modelstatus);
  int iteration_count;
  Highs_getIntInfoValue(highs, "simplex_iteration_count", &iteration_count);
  printf("Iteration count = %d\n", iteration_count);
  Highs_destroy(highs);
}

*/
