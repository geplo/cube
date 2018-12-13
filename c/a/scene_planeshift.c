#include <stdlib.h> // rand(3).

#include "cube.h" // cube_t & co.

typedef struct {
  axis_t        axis;
  shift_dir_t   direction;
  int           position;
}               plane_t;

// new_plane clears the cube and sets a new random plane.
static plane_t  new_plane() {
  plane_t       plane;

  plane.axis = rand() % 3; // Select a random axis.

  // Select a random edge.
  // % 2 is 0 or 1,
  // then * CUBE_SIZE - 1 is 0 or CUBE_SIZE - 1, i.e., first or last).
  plane.position = rand() % 2 * (CUBE_SIZE - 1);

  // Choose a direction based on the axis/position.
  switch (plane.axis) {
  case axisX:
    if (plane.position == 0) {
      plane.direction = shiftPosX; // X axis, first edge, go up X.
    } else {
      plane.direction = shiftNegX; // X axis, last edge, go down X.
    }
    break;
  case axisY:
    if (plane.position == 0) {
      plane.direction = shiftPosY; // Y axis, first edge, go up Y.
    } else {
      plane.direction = shiftNegY; // Y axis, last edge, go down Y.
    }
    break;
  case axisZ:
    if (plane.position == 0) {
      plane.direction = shiftPosZ; // Z axis, first edge, go up Z.
    } else {
      plane.direction = shiftNegZ; // Z axis, last edge, go down Z.
    }
    break;
  }

  return plane;
}

// plane_shift is a scene.
void                    plane_shift(cube_t cube) {
  static char           loading = 1;
  static char           looped  = 0;
  static unsigned int   timer   = 0;
  static plane_t        plane;

  // If we are loading, initialize the scenario.
  if (loading) {
    plane = new_plane();                        // Create a new plane.
    clearCube(cube);                            // Make sure to have a clean slate.
    setPlane(cube, plane.axis, plane.position); // Populate the cube with the new plane.

    // Reset flags.
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

  timer = 0;                    // Reset the "timer".
  shift(cube, plane.direction); // Shift the plane in selected direction.

  // Update plane and flags and check for edges, based on direction.

  // If we reached the edge, the first time (looped flag off),
  // reverse direction, the second time, set the loading flag on
  // to re-init the scenario.

  switch (plane.direction) {
  case shiftPosX: case shiftPosY: case shiftPosZ:
    // When moving up, increment the position.
    plane.position++;

    // Going up, we check for the last edge.
    if (plane.position == CUBE_SIZE - 1) {
      if (!looped) {
	// 1st time we reach the edge, reverse direction.
	plane.direction++; // The next enum from "pos" is "neg".
	looped = 1;       // Flag that we reached an edge.
      } else {
	// 2nd time we reach the edge, flag for reset.
	loading = 1;
      }
    }
    break;

  case shiftNegX: case shiftNegY: case shiftNegZ:
    // When moving down, decrement the position.
    plane.position--;

    // Going down, we check for the first edge.
    if (plane.position == 0) {
      if (!looped) {
	// 1st time we reach the edge, reverse direction.
	plane.direction--; // The previous enum from "neg" is "pos".
	looped = 1;       // Flag that we reached an edge.
      } else {
	// 2nd time we reach the edge, flag for reset.
	loading = 1;
      }
    }
    break;
  }
}
