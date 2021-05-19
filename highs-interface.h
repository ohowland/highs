#ifndef _HIGHS_INTERFACE_H
#define _HIGHS_INTERFACE_H

#include "interfaces/highs_c_api.h"

#ifdef __cplusplus
extern "C" {
#endif

    // opaque pointer to c++ highs object
    typedef char highs_obj;

    extern highs_obj* highsiface_create();
    extern void highsiface_free(highs_obj* highs);
    extern HighsInt highsiface_add_cols(highs_obj* highs, HighsInt numCol, double* colCost, double* colLower, double* colUpper);
    extern HighsInt highsiface_add_rows(highs_obj* highs, HighsInt numRow, double* rowLower, double* rowUpper, HighsInt numNz, HighsInt* arStart, HighsInt* arIndex, double* arValue);
    extern HighsInt highsiface_run(highs_obj* highs);
    extern void highsiface_get_solution(highs_obj* highs, double* colValue, double* colDual, double* rowValue, double* rowDual);

#ifdef __cplusplus
}
#endif
#endif