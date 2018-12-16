#ifndef __CUBE_H__
# define __CUBE_H__

// Define the cube as 8x8x8, one color.
typedef unsigned char cube_size_t;                  // 8 bits per element.
#define CUBE_SIZE     (sizeof(cube_size_t) * 8)
typedef cube_size_t   cube_t[CUBE_SIZE][CUBE_SIZE]; // 2d array so we have x*y*z.

// Enum for shift direction.
typedef enum {
              shiftPosX,
              shiftNegX,
              shiftPosY,
              shiftNegY,
              shiftPosZ,
              shiftNegZ,
} shift_dir_t;

// Enum for axis select.
typedef enum {
              axisX,
              axisY,
              axisZ,
} axis_t;

void set_voxel(cube_t cube, int x, int y, int z);
void clear_cube(cube_t cube);
void shift(cube_t cube, shift_dir_t dir);
void set_plane(cube_t cube, axis_t axis, int i);

#endif /* !__CUBE_H__ */
