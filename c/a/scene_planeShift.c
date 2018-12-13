#include <stdlib.h> // rand(3).

#include "cube.h" // cube_t & co.

void                    planeShift(cube_t cube) {
  static char           loading        = 1;
  static char           looped         = 0;
  static unsigned int   timer          = 0;
  static int            planePosition  = 0;
  static shift_dir_t    planeDirection = shiftPosX;
  axis_t                axis;

  // If we are loading, initialize the scenario.
  if (loading) {
    clearCube(cube); // Make sure to have a clean slate.

    axis = rand() % 3; // Select a random axis.

    // Select a random edge.
    // % 2 is 0 or 1,
    // then * CUBE_SIZE - 1 is 0 or CUBE_SIZE - 1, i.e., first or last).
    planePosition = rand() % 2 * (CUBE_SIZE - 1);

    // Populate the cube with the selected plane.
    setPlane(cube, axis, planePosition);

    // Choose a direction based on the axis/position.
    switch (axis) {
    case axisX:
      if (planePosition == 0) {
        planeDirection = shiftPosX; // X axis, first edge, go up X.
      } else {
        planeDirection = shiftNegX; // X axis, last edge, go down X.
      }
      break;
    case axisY:
      if (planePosition == 0) {
        planeDirection = shiftPosY; // Y axis, first edge, go up Y.
      } else {
        planeDirection = shiftNegY; // Y axis, last edge, go down Y.
      }
      break;
    case axisZ:
      if (planePosition == 0) {
        planeDirection = shiftPosZ; // Z axis, first edge, go up Z.
      } else {
        planeDirection = shiftNegZ; // Z axis, last edge, go down Z.
      }
      break;
    }

    // Reset the "state".
    timer   = 0;
    looped  = 0;
    loading = 0;
  }

  // "tick" the timer.
  timer++;

  // "Wait" until timer reach defined count.
  if (timer < 220) {
    return;
  }

  // "Timer" triggered, step the plane.

  timer = 0;                   // Reset the "tiumer".
  shift(cube, planeDirection); // Shift the plane in selected direction.

  // Update "state" and check for edges, based on direction.

  // If we reached the edge, the first time (looped flag off),
  // reverse direction, the second time, set the loading flag on
  // to re-init the scenario.

  switch (planeDirection) {
  case shiftPosX: case shiftPosY: case shiftPosZ:
    // When moving up, increment the position.
    planePosition++;

    // Going up, we check for the last edge.
    if (planePosition == CUBE_SIZE - 1) {
      if (!looped) {
	// 1st time we reach the edge, reverse direction.
	planeDirection++; // The next enum from "pos" is "neg".
	looped = 1;       // Flag that we reached an edge.
      } else {
	// 2nd time we reach the edge, flag for reset.
	loading = 1;
      }
    }
    break;

  case shiftNegX: case shiftNegY: case shiftNegZ:
    // When moving down, decrement the position.
    planePosition--;

    // Going down, we check for the first edge.
    if (planePosition == 0) {
      if (!looped) {
	// 1st time we reach the edge, reverse direction.
	planeDirection--; // The previous enum from "neg" is "pos".
	looped = 1;       // Flag that we reached an edge.
      } else {
	// 2nd time we reach the edge, flag for reset.
	loading = 1;
      }
    }
    break;
  }
}
