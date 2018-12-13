#ifndef __SPI_H__
# define __SPI_H__

# include <stdint.h> // uint8_t & co.

/**
   Example:

   spi_handler   hdlr = {
     .config = {
       .device = "/dev/spidev0.0",
       .mode   = 0,
       .bits   = 8,
       .speed  = 8000000,
       .delay  = 5,
     },
   };
*/

typedef struct {
    const char* device; // Device path.
    uint8_t     mode;   // SPI mode.
    uint8_t     bits;   // Bits per words.
    uint32_t    speed;  // Frequency (Hz).
    uint16_t    delay;  // Delay between words (usec).
}               spi_config;

typedef struct {
    spi_config    config;
    int           fd;
}               spi_handler;

int     spi_setup(spi_handler* hdlr);
int     spi_cleanup(spi_handler* hdlr);
int     transfer(spi_handler hdlr, uint8_t tx[], uint8_t rx[], uint64_t len);

#endif /* !__SPI_H__ */
