// Adds a contentWindow proxy to the provided iframe element
const addContentWindowProxy = (iframe) => {
  const contentWindowProxy = {
    get(target, key) {
      if (key === 'self') {
        return this;
      }
      if (key === 'frameElement') {
        return iframe;
      }
      return Reflect.get(target, key);
    },
  };

  if (!iframe.contentWindow) {
    const proxy = new Proxy(window, contentWindowProxy);
    Object.defineProperty(iframe, 'contentWindow', {
      get() {
        return proxy;
      },
      set(newValue) {
        return newValue; // contentWindow is immutable
      },
      enumerable: true,
      configurable: false,
    });
  }
};

// Handles iframe element creation, augments `srcdoc` property so we can intercept further
const handleIframeCreation = (target, thisArg, args) => {
  const iframe = target.apply(thisArg, args);
  const _iframe = iframe;
  const _srcdoc = _iframe.srcdoc;
  Object.defineProperty(iframe, 'srcdoc', {
    configurable: true, // Important, so we can reset this later
    get() {
      return _iframe.srcdoc;
    },
    set(newValue) {
      addContentWindowProxy(this);
      Object.defineProperty(iframe, 'srcdoc', {
        configurable: false,
        writable: false,
        value: _srcdoc,
      });
      _iframe.srcdoc = newValue;
    },
  });
  return iframe;
};

// Adds a hook to intercept iframe creation events
const addIframeCreationSniffer = (utils) => {
  const createElementHandler = {
    get(target, key) {
      return Reflect.get(target, key);
    },
    apply(target, thisArg, args) {
      const isIframe = args && args.length && `${args[0]}`.toLowerCase() === 'iframe';
      if (!isIframe) {
        return target.apply(thisArg, args);
      }
      return handleIframeCreation(target, thisArg, args);
    },
  };
  utils.replaceWithProxy(document, 'createElement', createElementHandler);
};

// eslint-disable-next-line import/prefer-default-export
export { addIframeCreationSniffer };
