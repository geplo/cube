#include <sys/signal.h> // signal(2) & co.

int setup();
int loop();
int cleanup();

static volatile int _running = 1;

static void intHandler() {
    _running = 0;
}

int     main() {
   signal(SIGINT, intHandler);
   signal(SIGTERM, intHandler);

  if (setup() < 0) {
    return 1;
  }

  // Main loop.
  while (loop() >= 0 && _running);

  if (cleanup() < 0) {
    return 1;
  }

  return 0;
}
