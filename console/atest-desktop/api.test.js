const { getPort, getHealthzUrl, getHomePage } = require('./api');

describe('getPort function', () => {
  test('should return the default port number 7788', () => {
    const port = getPort();
    expect(port).toBe(7788);
  });
});

describe('getHealthzUrl function', () => {
    test('should return the default healthz url', () => {
      const url = getHealthzUrl();
      expect(url).toBe('http://localhost:7788/healthz');
    });
});

describe('getHomePage function', () => {
    test('should return the default home page url', () => {
      const url = getHomePage();
      expect(url).toBe('http://localhost:7788');
    });
})
