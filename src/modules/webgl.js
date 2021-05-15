/**
 * Fix WebGL Vendor/Renderer being set to Google in headless mode
 * Example data (Apple Retina MBP 13): {vendor: "Intel Inc.", renderer: "Intel(R) Iris(TM) Graphics 6100"}
 */
const webglVendor = (utils) => {
  const getParameterProxyHandler = {
    apply(target, ctx, args) {
      const param = (args || [])[0];
      const result = utils.cache.Reflect.apply(target, ctx, args);
      if (param === 37445) return 'Intel Inc.'; // default in headless: Google Inc.
      if (param === 37446) return 'Intel Iris OpenGL Engine'; // default in headless: Google SwiftShader
      return result;
    },
  };
  const addProxy = (obj, propName) => {
    utils.replaceWithProxy(obj, propName, getParameterProxyHandler);
  };
  addProxy(WebGLRenderingContext.prototype, 'getParameter');
  addProxy(WebGL2RenderingContext.prototype, 'getParameter');
};

// eslint-disable-next-line import/prefer-default-export
export { webglVendor };
