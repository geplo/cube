#define _BSD_SOURCE // For usleep(3) (fix warning on linux).
#include <unistd.h> // usleep(3), close(2).

#include <time.h>   // time(2) (random seed).

#include <stdint.h>  // uint8_t & co.
#include <stdbool.h> // bool support.

#include <stdio.h>     // perror(3), printf(3) & co.
#include <stdlib.h>    // abort(3), srand(3), rand(3).
#include <fcntl.h>     // open(2).
#include <sys/ioctl.h> // ioctl(2).

#include <linux/spi/spidev.h> // SPI structs & ioctls.

typedef struct  spi_config {
  const char*   device; // Device path.
  uint8_t       mode;   // SPI mode.
  uint8_t       bits;   // Bits per words.
  uint32_t      speed;  // Frequency (Hz).
  uint16_t      delay;  // Delay between words (usec).
}               spi_config;

typedef struct  spi_handler {
  spi_config    config;
  int           fd;
}               spi_handler;


void clearCube(uint8_t cube[8][8]);
void shift(uint8_t cube[8][8], uint8_t dir);
void setPlane(uint8_t cube[8][8],uint8_t axis, uint8_t i);
void setVoxel(uint8_t cube[8][8],uint8_t x, uint8_t y, uint8_t z);
void planeBoing(uint8_t cube[8][8]);
void rain(uint8_t cube[8][8]);
void setup_cube();

int renderCube(spi_handler hdlr, uint8_t cube[8][8]);
int loop(spi_handler hdlr, uint8_t cube[8][8]);

int                             transfer(spi_handler hdlr, uint8_t tx[], uint64_t len) {
  uint8_t*                      rx = malloc(sizeof(uint8_t) * len); // TODO: See if we can get rid of this one.
  struct spi_ioc_transfer       tr = {
    .tx_buf        = (unsigned long)tx,
    .rx_buf        = (unsigned long)rx,
    .len           = len,
    .delay_usecs   = hdlr.config.delay,
    .speed_hz      = hdlr.config.speed,
    .bits_per_word = hdlr.config.bits,
  };

  if (ioctl(hdlr.fd, SPI_IOC_MESSAGE(1), &tr) < 1) {
    return -1;
  }
  return 0;
}

int     spi_setup(spi_handler* hdlr) {
  int   ret;

  // TODO: Check if it would work in WR_ONLY, skipping all the RD iotctls and NULL the tx rd.

  // Open the device in read/write.
  if ((ret = open(hdlr->config.device, O_RDWR))  < 0) {
    return ret;
  }
  hdlr->fd = ret;

  // Set the SPI mode.
  if ((ret = ioctl(hdlr->fd, SPI_IOC_WR_MODE, &hdlr->config.mode)) < 0) {
    return ret;
  }
  if ((ret = ioctl(hdlr->fd, SPI_IOC_RD_MODE, &hdlr->config.mode)) < 0) {
    return ret;
  }

  // Set the bits per word.
  if ((ret = ioctl(hdlr->fd, SPI_IOC_WR_BITS_PER_WORD, &hdlr->config.bits)) < 0) {
    return ret;
  }
  if ((ret = ioctl(hdlr->fd, SPI_IOC_RD_BITS_PER_WORD, &hdlr->config.bits)) < 0) {
    return ret;
  }

  // Set max speed.
  if ((ret = ioctl(hdlr->fd, SPI_IOC_WR_MAX_SPEED_HZ, &hdlr->config.speed)) < 0) {
    return ret;
  }
  if ((ret = ioctl(hdlr->fd, SPI_IOC_RD_MAX_SPEED_HZ, &hdlr->config.speed)) < 0) {
    return ret;
  }

  // Success.
  return 0;
}

int             main() {
  uint8_t       cube[8][8];
  spi_handler   hdlr = {
			.config = {
				   .device = "/dev/spidev0.0",
				   .mode   = 0,
				   .bits   = 8,
				   .speed  = 8000000,
				   .delay  = 5,
				   },
  };

  // Initialize SPI.
  if (spi_setup(&hdlr) < 0) {
    perror("error setting up SPI device");
    return -1;
  }

  // TODO: Handle signals for clean exit.

  for (setup_cube(); ; loop(hdlr, cube)) {
    usleep(hdlr.config.delay);
  }

  // Cleanup SPI.
  if (close(hdlr.fd) < 0) {
    perror("error closing SPI device");
    return -1;
  }

  return 0;
}


#define XAXIS 0
#define YAXIS 1
#define ZAXIS 2

#define POS_X 0
#define NEG_X 1
#define POS_Z 2
#define NEG_Z 3
#define POS_Y 4
#define NEG_Y 5

#define RAIN_TIME        260
#define PLANE_BOING_TIME 220

uint16_t timer;
uint64_t randomTimer;
bool     loading;

void setup_cube() {
  loading = true;
  randomTimer = 0;
  srand(time(NULL));
}

int     loop(spi_handler hdlr, uint8_t cube[8][8]) {
  randomTimer++;
  rain(cube);
  return renderCube(hdlr, cube);
}

int             renderCube(spi_handler hdlr, uint8_t cube[8][8]) {
  uint8_t       txx[9];
  int           ret;

  for (uint8_t i = 0; i < 8; i++) {
    txx[0] = 0x01 << i;
    for (uint8_t j = 0; j < 8; j++) {
      txx[j + 1] = cube[7 - i][j + 1];
      txx[j + 2] = cube[7 - i][j];
      j++;
    }
    if ((ret = transfer(hdlr, txx, sizeof(txx))) < 0) {
      return ret;
    }
  }
  return 0;
}

void rain(uint8_t cube[8][8]) {
  if (loading) {
    clearCube(cube);
    loading = false;
  }
  timer++;
  if (timer > RAIN_TIME) {
    timer = 0;
    shift(cube, NEG_Y);
    uint8_t numDrops = rand() % 5;
    for (uint8_t i = 0; i < numDrops; i++) {
      setVoxel(cube, rand() % 8, 7, rand() % 8);
    }
  }
}

uint8_t planePosition  = 0;
uint8_t planeDirection = 0;
bool    looped         = false;

void planeBoing(uint8_t cube[8][8]) {
  if (loading) {
    clearCube(cube);
    uint8_t axis = rand() % 3;
    planePosition = rand() % 2 * 7;
    setPlane(cube, axis, planePosition);
    if (axis == XAXIS) {
      if (planePosition == 0) {
        planeDirection = POS_X;
      } else {
        planeDirection = NEG_X;
      }
    } else if (axis == YAXIS) {
      if (planePosition == 0) {
        planeDirection = POS_Y;
      } else {
        planeDirection = NEG_Y;
      }
    } else if (axis == ZAXIS) {
      if (planePosition == 0) {
        planeDirection = POS_Z;
      } else {
        planeDirection = NEG_Z;
      }
    }
    timer = 0;
    looped = false;
    loading = false;
  }

  timer++;
  if (timer > PLANE_BOING_TIME) {
    timer = 0;
    shift(cube, planeDirection);
    if (planeDirection % 2 == 0) {
      planePosition++;
      if (planePosition == 7) {
        if (looped) {
          loading = true;
        } else {
          planeDirection++;
          looped = true;
        }
      }
    } else {
      planePosition--;
      if (planePosition == 0) {
        if (looped) {
          loading = true;
        } else {
          planeDirection--;
          looped = true;
        }
      }
    }
  }
}

void setVoxel(uint8_t cube[8][8], uint8_t x, uint8_t y, uint8_t z) {
  cube[7 - y][7 - z] |= (0x01 << x);
}

void setPlane(uint8_t cube[8][8], uint8_t axis, uint8_t i) {
  for (uint8_t j = 0; j < 8; j++) {
    for (uint8_t k = 0; k < 8; k++) {
      if (axis == XAXIS) {
        setVoxel(cube, i, j, k);
      } else if (axis == YAXIS) {
        setVoxel(cube, j, i, k);
      } else if (axis == ZAXIS) {
        setVoxel(cube, j, k, i);
      }
    }
  }
}

void shift(uint8_t cube[8][8], uint8_t dir) {
  if (dir == POS_X) {
    for (uint8_t y = 0; y < 8; y++) {
      for (uint8_t z = 0; z < 8; z++) {
        cube[y][z] = cube[y][z] << 1;
      }
    }
  } else if (dir == NEG_X) {
    for (uint8_t y = 0; y < 8; y++) {
      for (uint8_t z = 0; z < 8; z++) {
        cube[y][z] = cube[y][z] >> 1;
      }
    }
  } else if (dir == POS_Y) {
    for (uint8_t y = 1; y < 8; y++) {
      for (uint8_t z = 0; z < 8; z++) {
        cube[y - 1][z] = cube[y][z];
      }
    }
    for (uint8_t i = 0; i < 8; i++) {
      cube[7][i] = 0;
    }
  } else if (dir == NEG_Y) {
    for (uint8_t y = 7; y > 0; y--) {
      for (uint8_t z = 0; z < 8; z++) {
        cube[y][z] = cube[y - 1][z];
      }
    }
    for (uint8_t i = 0; i < 8; i++) {
      cube[0][i] = 0;
    }
  } else if (dir == POS_Z) {
    for (uint8_t y = 0; y < 8; y++) {
      for (uint8_t z = 1; z < 8; z++) {
        cube[y][z - 1] = cube[y][z];
      }
    }
    for (uint8_t i = 0; i < 8; i++) {
      cube[i][7] = 0;
    }
  } else if (dir == NEG_Z) {
    for (uint8_t y = 0; y < 8; y++) {
      for (uint8_t z = 7; z > 0; z--) {
        cube[y][z] = cube[y][z - 1];
      }
    }
    for (uint8_t i = 0; i < 8; i++) {
      cube[i][0] = 0;
    }
  }
}

void clearCube(uint8_t cube[8][8]) {
  for (uint8_t i = 0; i < 8; i++) {
    for (uint8_t j = 0; j < 8; j++) {
      cube[i][j] = 0;
    }
  }
}
