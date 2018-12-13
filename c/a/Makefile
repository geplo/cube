NAME    = cube

SRCS    = main.c \
          loop.c \
          spi.c \
          cube.c \
          scene_planeShift.c \
          scene_rain.c
OBJS    = ${SRCS:.c=.o}

CC      = gcc
LD      = gcc
CFLAGS  = -W -Wall -Werror -ansi -pedantic -std=c99
LDFLAGS =

.DEFAULT_GOAL = ${NAME}

# Dependency tree.
cube.c:             cube.h
loop.c:             cube.h spi.h scenes.h
scene_planeShift.c: cube.h
scene_rain.c:       cube.h
scenes.h:           cube.h
spi.c:              spi.h


# Main target.
${NAME} : ${OBJS}
	${LD} -o $@ ${LDFLAGS} $+

# Cleanup.
.PHONY  : clean fclean re
clean   :
	${RM} ${OBJS}

fclean  : clean
	${RM} ${NAME}

re      : fclean ${NAME}
