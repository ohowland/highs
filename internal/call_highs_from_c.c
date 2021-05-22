#include "interfaces/highs_c_api.h"

#include <stdio.h>
#include <stdlib.h>
#include <assert.h>

// gcc call_highs_from_c.c -o highstest -I ../build/install_folder/include/ -L ../build/install_folder/lib/ -lhighs

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
  
  int integrality[numcol] = {1, 1};

  assert ( Highs_changeColsIntegralityByRange(highs, 0, 1, integrality) );
  //Highs_changeColIntegrality(highs, 0, 1);

  int runstatus = Highs_run(highs);
  int modelstatus = Highs_getModelStatus(highs);

  printf("Run status = %d; Model status = %d\n", runstatus, modelstatus);

  if (modelstatus == 7) {
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
  printf("MIPS TIME\n");
  int rowwise = 0;
  int astart[numcol] = {0, 2};
  int aindex[numnz] = {1, 2, 0, 1, 2};
  double avalue[numnz] = {1.0, 2.0, 1.0, 2.0, 1.0};
  highs = Highs_create();
  runstatus = Highs_passMip(highs, numcol, numrow, numnz, rowwise,
			colcost, collower, colupper,
			rowlower, rowupper,
			astart, aindex, avalue, integrality);

  Highs_changeColIntegrality(highs, 0, 0);

  runstatus = Highs_run(highs);
  modelstatus = Highs_getModelStatus(highs);
  printf("Run status = %d; Model status = %d\n", runstatus, modelstatus);
  int iteration_count;
  Highs_destroy(highs);
}

int main() {
  full_api();
  return 0;
}
