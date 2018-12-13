#include <fcntl.h>     // open(2).
#include <unistd.h>    // close(3).
#include <sys/ioctl.h> // ioctl(2).

#include <linux/spi/spidev.h> // spi_ioc_transfer & ioctls consts.

#include "spi.h"

// spi_transfer uses SPI to send tx and receives in rx. tx and rx must be allocated with len size.
int                             spi_transfer(spi_handler* hdlr, void* tx, void* rx, int len) {
  struct spi_ioc_transfer       tr =
    {
     .tx_buf        = (unsigned long)tx,
     .rx_buf        = (unsigned long)rx,
     .len           = len,
     .speed_hz      = hdlr->config.speed,
     .delay_usecs   = hdlr->config.delay,
     .bits_per_word = hdlr->config.bits,
    };

  if (ioctl(hdlr->fd, SPI_IOC_MESSAGE(1), &tr) < 1) {
    return -1;
  }
  return 0;
}

// spi_setup initializes the SPI with the hdlr->config values.
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

// spi_cleanup closes down the SPI.
int     spi_cleanup(spi_handler* hdlr) {
  int   ret;

  if ((ret = close(hdlr->fd)) < 0) {
    return ret;
  }
  hdlr->fd = -1;

  return 0;
}
