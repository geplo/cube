int setup();
int loop();
int cleanup();

int     main() {
  // TODO: Handle signals for clean exit.
  if (setup() < 0) {
    return 1;
  }

  // Main loop.
  // TODO: Use select(2)?
  for (; loop() >= 0;);

  if (cleanup() < 0) {
    return 1;
  }

  return 0;
}
