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
    .mode   = 0,                // Default mode.
    .bits   = 8,                // 8 bits per words.
    .speed  = 8000000,          // 8MHz
    .delay  = 5,                // 5 micro secs in between iterations.
  },
};


// render_cube uses SPI to display the cube.
int             render_cube(spi_handler hdlr, cube_t cube) {
  cube_size_t   tx[CUBE_SIZE + 1]; // 1 row of cathodes, CUBE_SIZE rows of anodes.
  int           ret;

  // We go one cathode at the time without delay.
  for (unsigned int i = 0; i < CUBE_SIZE; i++) {
    // Cathodes.
    tx[0] = 0x01 << i;

    // Anodes.
    for (unsigned int j = 0; j < CUBE_SIZE; j++) {
      // NOTE: 2 by 2 to fit the wiring.
      // TODO: Setup wiring in init.
      tx[j + 1] = cube[CUBE_SIZE - 1 - i][j + 1];
      tx[j + 2] = cube[CUBE_SIZE - 1 - i][j];
      j++;
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
  clearCube(cube);

  // Set the scene to use.
  scene = rain;
  // scene = plane_shift;

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

  // Cleanup SPI.
  if ((ret = spi_cleanup(&hdlr)) < 0) {
    perror("error cleaning up SPI");
    return ret;
  }

  return 0;
}
