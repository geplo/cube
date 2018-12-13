#include <stdint.h> // uint8_t & co.

#include <fcntl.h>     // open(2).
#include <unistd.h>    // close(3).
#include <sys/ioctl.h> // ioctl(2).

#include <linux/spi/spidev.h> // SPI structs & ioctls.

#include "spi.h"

int                             transfer(spi_handler hdlr, uint8_t tx[], uint8_t rx[], uint64_t len) {
  struct spi_ioc_transfer       tr =
    {
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

int     spi_cleanup(spi_handler* hdlr) {
  int   ret;

  if ((ret = close(hdlr->fd)) < 0) {
    return ret;
  }
  hdlr->fd = -1;

  return 0;
}
