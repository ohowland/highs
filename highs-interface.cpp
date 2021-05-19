#include "interfaces/highs_c_api.h"
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include "call_highs.hpp"

int first(int x) {
  return x;
}

void minimal_api() {
  const int numcol = 2;
  const int numrow = 3;
  const int numnz = 5;

  // Define the column costs, lower bounds and upper bounds
  double colcost[numcol] = {2.0, 3.0};
  double collower[numcol] = {0.0, 1.0};
  double colupper[numcol] = {3.0, 1.0e30};
  // Define the row lower bounds and upper bounds
  double rowlower[numrow] = {-1.0e30, 10.0, 8.0};
  double rowupper[numrow] = {6.0, 14.0, 1.0e30};
  // Define the constraint matrix column-wise
  int astart[numcol] = {0, 2};
  int aindex[numnz] = {1, 2, 0, 1, 2};
  double avalue[numnz] = {1.0, 2.0, 1.0, 2.0, 1.0};

  double* colvalue = (double*)malloc(sizeof(double) * numcol);
  double* coldual = (double*)malloc(sizeof(double) * numcol);
  double* rowvalue = (double*)malloc(sizeof(double) * numrow);
  double* rowdual = (double*)malloc(sizeof(double) * numrow);

  int* colbasisstatus = (int*)malloc(sizeof(int) * numcol);
  int* rowbasisstatus = (int*)malloc(sizeof(int) * numrow);

  int modelstatus;

  const int rowwise = 0;
  int runstatus = Highs_lpCall(numcol, numrow, numnz, rowwise,
			       colcost, collower, colupper, rowlower, rowupper,
			       astart, aindex, avalue,
			       colvalue, coldual, rowvalue, rowdual,
			       colbasisstatus, rowbasisstatus,
			       &modelstatus);

  assert(runstatus == 0);

  printf("Run status = %d; Model status = %d\n", runstatus, modelstatus);

  int i;
  if (modelstatus == 9) {
    double objective_value = 0;
    // Report the column primal and dual values, and basis status
    for (i = 0; i < numcol; i++) {
      printf("Col%d = %lf; dual = %lf; status = %d; \n", i, colvalue[i], coldual[i], colbasisstatus[i]);
      objective_value += colvalue[i]*colcost[i];
    }
    // Report the row primal and dual values, and basis status
    for (i = 0; i < numrow; i++) {
      printf("Row%d = %lf; dual = %lf; status = %d; \n", i, rowvalue[i], rowdual[i], rowbasisstatus[i]);
    }
    printf("Optimal objective value = %g\n", objective_value);
  }

  free(colvalue);
  free(coldual);
  free(rowvalue);
  free(rowdual);
  free(colbasisstatus);
  free(rowbasisstatus);
}