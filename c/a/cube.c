#include "cube.h" // cube_t, & co.

// setVoxel tuerns on the x*y*z LED.
void setVoxel(cube_t cube, int x, int y, int z) {
  cube[CUBE_SIZE - 1 - y][CUBE_SIZE - 1 - z] |= (0x01 << x);
}

// clearCube turns off the whole cube.
void clearCube(cube_t cube) {
  for (unsigned int i = 0; i < CUBE_SIZE; i++) {
    for (unsigned int j = 0; j < CUBE_SIZE; j++) {
      cube[i][j] = 0x00;
    }
  }
}

// shift slides the whole cube on the given direction.
void shift(cube_t cube, shift_dir_t dir) {
  switch (dir) {
  case shiftPosX:
    for (unsigned int y = 0; y < CUBE_SIZE; y++) {
      for (unsigned int z = 0; z < CUBE_SIZE; z++) {
        cube[y][z] = cube[y][z] << 1;
      }
    }
    break;
  case shiftNegX:
    for (unsigned int y = 0; y < CUBE_SIZE; y++) {
      for (unsigned int z = 0; z < CUBE_SIZE; z++) {
        cube[y][z] = cube[y][z] >> 1;
      }
    }
    break;

  case shiftPosY:
    for (unsigned int y = 1; y < CUBE_SIZE; y++) {
      for (unsigned int z = 0; z < CUBE_SIZE; z++) {
        cube[y - 1][z] = cube[y][z];
      }
    }
    for (unsigned int i = 0; i < CUBE_SIZE; i++) {
      cube[CUBE_SIZE - 1][i] = 0;
    }
    break;
  case shiftNegY:
    for (unsigned int y = CUBE_SIZE - 1; y > 0; y--) {
      for (unsigned int z = 0; z < CUBE_SIZE; z++) {
        cube[y][z] = cube[y - 1][z];
      }
    }
    for (unsigned int i = 0; i < CUBE_SIZE; i++) {
      cube[0][i] = 0;
    }
    break;

  case shiftPosZ:
    for (unsigned int y = 0; y < CUBE_SIZE; y++) {
      for (unsigned int z = 1; z < CUBE_SIZE; z++) {
        cube[y][z - 1] = cube[y][z];
      }
    }
    for (unsigned int i = 0; i < CUBE_SIZE; i++) {
      cube[i][CUBE_SIZE - 1] = 0;
    }
    break;
  case shiftNegZ:
    for (unsigned int y = 0; y < CUBE_SIZE; y++) {
      for (unsigned int z = CUBE_SIZE - 1; z > 0; z--) {
        cube[y][z] = cube[y][z - 1];
      }
    }
    for (unsigned int i = 0; i < CUBE_SIZE; i++) {
      cube[i][0] = 0;
    }
    break;
  }
}

// setPlane turns on the Nth plane from the given axis.
void setPlane(cube_t cube, axis_t axis, int n) {
  for (unsigned int j = 0; j < CUBE_SIZE; j++) {
    for (unsigned int k = 0; k < CUBE_SIZE; k++) {
      switch (axis) {
      case axisX:
        setVoxel(cube, n, j, k);
        break;
      case axisY:
        setVoxel(cube, j, n, k);
        break;
      case axisZ:
        setVoxel(cube, j, k, n);
        break;
      }
    }
  }
}
