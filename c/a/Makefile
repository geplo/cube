NAME    = cube

SRCS    = main.c \
          loop.c \
          spi.c \
          cube.c \
          scene_planeshift.c \
          scene_rain.c \
          scene_manual.c
HEADERS = cube.h \
          spi.h \
          scenes.h
OBJS    = ${SRCS:.c=.o}

CC      = gcc
LD      = gcc
CFLAGS  = -W -Wall -Werror -ansi -pedantic -std=c99
LDFLAGS =

.DEFAULT_GOAL = ${NAME}

# Dependency tree.
cube.c:             cube.h
scene_planeshift.c: cube.h
scene_rain.c:       cube.h
scene_manual.c:     cube.h
spi.c:              spi.h
loop.c:             cube.h spi.h scenes.h
scenes.h:           cube.h

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

# Helper.
${SRCS} ${HEADERS}:
	@touch $@
