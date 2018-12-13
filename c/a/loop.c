#define _BSD_SOURCE // For usleep(3) (fix warning on linux).
#include <unistd.h> // usleep(3).
#include <time.h>   // time(2) (for random seed).
#include <stdio.h>  // perror(3), printf(3) & co.
#include <stdlib.h> // srand(3).

#include "spi.h"    // SPI lib.
#include "cube.h"   // Cube managment.
#include "scenes.h" // Scenes.

cube_t        cube;
spi_handler   hdlr =
  {
   .config = {
	      .device = "/dev/spidev0.0", // Device to use.
	      .mode   = 0,                // Default mode (MSBF, most significant bit first).
	      .bits   = CUBE_SIZE,        // CUBE_SIZE bits per words.
	      .speed  = 8000000,          // 8MHz
	      .delay  = 5,                // 5 micro secs in between iterations.
	      },
  };

int renderCube(spi_handler hdlr, cube_t cube);

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

  return 0;
}

int     cleanup() {
  int   ret;

  // Cleanup SPI.
  if ((ret = spi_cleanup(&hdlr)) < 0) {
    perror("error cleaning up SPI");
    return ret;
  }

  return 0;
}

int     loop() {
  int   ret;

  rain(cube);
  if ((ret = renderCube(hdlr, cube)) < 0) {
    return ret;
  }

  usleep(hdlr.config.delay);
  return 0;
}

int             renderCube(spi_handler hdlr, cube_t cube) {
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
    if ((ret = transfer(hdlr, tx, NULL, sizeof(tx))) < 0) {
      return ret;
    }
  }

  return 0;
}
