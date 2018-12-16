#define _DEFAULT_SOURCE // For usleep(3) (fix warning on linux).
#include <unistd.h>     // usleep(3).
#include <time.h>       // time(2) (for random seed).
#include <stdio.h>      // perror(3), printf(3) & co.
#include <stdlib.h>     // srand(3).

#include "spi.h"        // SPI lib.
#include "cube.h"       // Cube managment.
#include "scenes.h"     // Scenes.

// Cube state.
cube_t cube;

// Scene handler.
void (*scene)(cube_t);

// SPI handler config.
spi_handler hdlr = {
  .config =  {
    .device = "/dev/spidev0.0", // Device to use.
    .mode   = 0,                // Mode 0 with LSB first.
    .bits   = 8,                // 8 bits per words.
    .speed  = 8000000,          // 8MHz
    .delay  = 5,                // 5 micro secs in between iterations.
  },
};

// get_voxel checks if a point is on or fof in the cube.
static inline int get_voxel(cube_t cube, int x, int y, int z) {
  return (cube[CUBE_SIZE - 1 - y][CUBE_SIZE - 1 - z] & (0x01 << x)) == (0x01 << x) ? 1 : 0;
}

// Hardware mapping.

// X wiring on Z axis.
static const int x_map[CUBE_SIZE][CUBE_SIZE] = {{0,1,2,3,4,5,6,7}, {7,6,5,4,3,2,1,0}, {0,1,2,3,4,5,6,7}, {7,6,5,4,3,2,1,0}, {0,1,2,3,4,5,6,7}, {7,6,5,4,3,2,1,0}, {0,1,2,3,4,5,6,7}, {7,6,5,4,3,2,1,0}};
// Y wiring on X axis.
static const int y_map[CUBE_SIZE][CUBE_SIZE] = {{0,1,2,3,4,5,6,7}, {0,1,2,3,4,5,6,7}, {0,1,2,3,4,5,6,7}, {0,1,2,3,4,5,6,7}, {0,1,2,3,4,5,6,7}, {0,1,2,3,4,5,6,7}, {0,1,2,3,4,5,6,7}, {0,1,2,3,4,5,6,7}};
// Z wiring on X axis.
static const int z_map[CUBE_SIZE][CUBE_SIZE] = {{1,0,3,2,5,4,7,6}, {1,0,3,2,5,4,7,6}, {1,0,3,2,5,4,7,6}, {1,0,3,2,5,4,7,6}, {1,0,3,2,5,4,7,6}, {1,0,3,2,5,4,7,6}, {1,0,3,2,5,4,7,6}, {1,0,3,2,5,4,7,6}};

static void map_cube(cube_t src, cube_t dst) {
  // Make sure the dst is cleared.
  clear_cube(dst);

  // For each point of the cube, map x/y/z to match the defined hardware wiring.
  for (unsigned int x = 0; x < CUBE_SIZE; x++) {
    for (unsigned int y = 0; y < CUBE_SIZE; y++) {
      for (unsigned int z = 0; z < CUBE_SIZE; z++) {
	if (get_voxel(src, x, y, z)) {
	  int xx = x_map[z][x];
	  int yy = y_map[x][y];
	  int zz = z_map[x][z];
	  set_voxel(dst, xx, yy, zz);
	}
      }
    }
  }
}

// render_cube uses SPI to display the cube.
int             render_cube(const spi_handler hdlr, cube_t cube) {
  cube_t        mapped_cube;
  cube_size_t   tx[CUBE_SIZE + 1]; // 1 row of cathodes, CUBE_SIZE rows of anodes.
  int           ret;

  // Map the memory cube to the hardware.
  map_cube(cube, mapped_cube);

  // We go one cathode at the time without delay.
  for (unsigned int i = 0; i < CUBE_SIZE; i++) {
    // Cathodes.
    tx[0] = 0x01 << i;

    // Anodes.
    for (unsigned int j = 0; j < CUBE_SIZE; j++) {
      tx[j + 1] = mapped_cube[CUBE_SIZE - 1 - i][j];
    }

    // Send the data to the SPI.
    if ((ret = spi_transfer(&hdlr, tx, NULL, sizeof(tx))) < 0) {
      return ret;
    }
  }

  return 0;
}

// setup is called before the main loop.
// Should return a negative value in case of error.
int     setup() {
  // Initialize SPI.
  if (spi_setup(&hdlr) < 0) {
    perror("error setting up SPI");
    return 1;
  }

  // Seed the random generator.
  srand(time(NULL));

  // Clear the cube.
  clear_cube(cube);

  // Set the scene to use.
  scene = rain;
  scene = manual;
  scene = plane_shift;

  return 0;
}

// loop is the main logic block, called by the main.
// Should return a negative value in case of error.
int     loop() {
  int   ret;

  // Step the scene.
  scene(cube);

  // Render it.
  if ((ret = render_cube(hdlr, cube)) < 0) {
    return ret;
  }

  // Delay and repeat.
  usleep(hdlr.config.delay);
  return 0;
}

// cleanup is callled when the main loop exits.
// Should return a negative value in case of error.
int     cleanup() {
  int   ret;

  // Turn off the cube.
  clear_cube(cube);
  render_cube(hdlr, cube);

  // Cleanup SPI.
  if ((ret = spi_cleanup(&hdlr)) < 0) {
    perror("error cleaning up SPI");
    return ret;
  }

  return 0;
}
