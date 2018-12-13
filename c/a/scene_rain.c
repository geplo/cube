#include <stdlib.h> // rand(3).

#include "cube.h" // cube_t & co.

void                    rain(cube_t cube) {
  static char           loading = 1;
  static unsigned int   timer   = 0;

  // "tick" the timer.
  timer++;

  // If loading, make sure to clear before we start.
  if (loading) {
    clearCube(cube);
    loading = 0;
  }

  // "Wait" until timer reach defined count.
  if (timer < 220) {
    return;
  }

  // "Timer" triggered, shift the drops one layer and
  // generate new ones on the top layer.

  timer = 0;              // Reset the "timer".
  shift(cube, shiftNegY); // Shift layers down.

  // From 0 to CUBE_SIZE / 2 + 1 drops per layer.
  for (unsigned int i = 0; i < rand() % (CUBE_SIZE / 2 + 1); i++) {
    setVoxel(cube,
	     rand() % CUBE_SIZE,  // Random X.
	     CUBE_SIZE - 1,       // Always top layer for Y.
	     rand() % CUBE_SIZE); // Random Z.
  }
}
