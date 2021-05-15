/**
 * Fix missing window.outerWidth/window.outerHeight in headless mode
 * Will also set the viewport to match window size, unless specified by user
 */
const windowOuterDimensions = () => {
  if (window.outerWidth && window.outerHeight) return;
  const windowFrame = 85;
  window.outerWidth = window.innerWidth;
  window.outerHeight = window.innerHeight + windowFrame;
};

/**
 * Fix window permissions / notifications
 */
const windowNotifications = () => {
  window.Notification = {
    permission: 'denied',
  };
  const originalQuery = window.navigator.permissions.query;
  // eslint-disable-next-line no-proto
  window.navigator.permissions.__proto__.query = (parameters) =>
    parameters.name === 'notifications' ? Promise.resolve({ state: Notification.permission }) : originalQuery(parameters);

  // Inspired by: https://github.com/ikarienator/phantomjs_hide_and_seek/blob/master/5.spoofFunctionBind.js
  const oldCall = Function.prototype.call;
  function call(...rest) {
    return oldCall.apply(this, rest);
  }
  // eslint-disable-next-line no-extend-native
  Function.prototype.call = call;

  const nativeToStringFunctionString = Error.toString().replace(/Error/g, 'toString');
  const oldToString = Function.prototype.toString;

  function functionToString() {
    if (this === window.navigator.permissions.query) {
      return 'function query() { [native code] }';
    }
    if (this === functionToString) {
      return nativeToStringFunctionString;
    }
    return oldCall.call(oldToString, this);
  }
  // eslint-disable-next-line no-extend-native
  Function.prototype.toString = functionToString;
};

export { windowOuterDimensions, windowNotifications };
