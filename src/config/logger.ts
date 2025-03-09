class Logger {
  private formatMessage(message: string | Error, level: string): string {
    const timestamp = new Date().toISOString();
    const msg = message instanceof Error ? message.stack || message.message : message;
    return `${timestamp} ${level.toUpperCase()}: ${msg}`;
  }

  info(message: string | object): void {
    const formatted = this.formatMessage(String(message), 'info');
    console.log(formatted);
  }

  debug(message: string | object): void {
    const formatted = this.formatMessage(String(message), 'debug');
    console.debug(formatted);
  }

  error(message: string | Error ): void {
    const formatted = this.formatMessage(message, 'error');
    console.error(formatted);
  }
}

const logger = new Logger();

export default logger;
