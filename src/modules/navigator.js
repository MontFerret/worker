/* eslint-disable no-proto */

/**
 * Set the hardwareConcurrency to 4 (optionally configurable with `hardwareConcurrency`)
 */
const navigatorHardwareConcurrency = () => {
  Object.defineProperty(navigator, 'hardwareConcurrency', { get: () => 12 });
};

/**
 * Pass the Languages Test.
 */
const navigatorLanguages = () => {
  Object.defineProperty(navigator, 'languages', { get: () => ['en-US', 'en'] });
  Object.defineProperty(navigator, 'language', { get: () => 'en-US' });
};

/**
 * Pass the Platform Test.
 */
const navigatorPlatform = () => {
  Object.defineProperty(navigator, 'platform', { get: () => 'MacIntel' });
};

/**
 * Pass the UserAgent Test.
 */
const navigatorUserAgent = () => {
  Object.defineProperty(navigator, 'userAgent', {
    get: () => 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36',
  });
  Object.defineProperty(navigator, 'appVersion', {
    get: () => '5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36',
  });
};

/**
 * Pass the Vendor Test.
 */
const navigatorVendor = () => {
  Object.defineProperty(navigator, 'vendor', { get: () => 'Google Inc.' });
  Object.defineProperty(navigator, 'vendorFlavors', { get: () => ['chrome'] });
};

/**
 * Pass the Webdriver Test.
 * Will delete `navigator.webdriver` property.
 */
const navigatorWebDriver = () => {
  const newProto = navigator.__proto__;
  delete newProto.webdriver;
  navigator.__proto__ = newProto;
  if (navigator.webdriver !== false && navigator.webdriver !== undefined) {
    delete navigator.webdriver;
  }
};

/**
 * Pass the Device Memory Test.
 */
const navigatorDeviceMemory = () => {
  Object.defineProperty(navigator, 'deviceMemory', { get: () => 8 });
};

export {
  navigatorHardwareConcurrency,
  navigatorLanguages,
  navigatorPlatform,
  navigatorUserAgent,
  navigatorVendor,
  navigatorWebDriver,
  navigatorDeviceMemory,
};
