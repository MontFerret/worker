/**
 * Mock the `chrome` object
 */
const chrome = () => {
  Object.defineProperty(window, 'chrome', {
    writable: true,
    enumerable: true,
    configurable: false, // note!
    value: {}, // We'll extend that later
  });
};

/**
 * Mock the `chrome.app` object
 */
const chromeApp = (utils) => {
  if ('app' in window.chrome) return;
  const makeError = {
    ErrorInInvocation: (fn) => {
      const err = new TypeError(`Error in invocation of app.${fn}()`);
      return utils.stripErrorWithAnchor(err, `at ${fn} (eval at <anonymous>`);
    },
  };
  const STATIC_DATA = {
    isInstalled: false,
    InstallState: {
      DISABLED: 'disabled',
      INSTALLED: 'installed',
      NOT_INSTALLED: 'not_installed',
    },
    RunningState: {
      CANNOT_RUN: 'cannot_run',
      READY_TO_RUN: 'ready_to_run',
      RUNNING: 'running',
    },
  };
  window.chrome.app = {
    ...STATIC_DATA,
    get isInstalled() {
      return false;
    },
    getDetails: function getDetails() {
      if (arguments.length) throw makeError.ErrorInInvocation(`getDetails`);
      return null;
    },
    getIsInstalled: function getDetails() {
      if (arguments.length) throw makeError.ErrorInInvocation(`getIsInstalled`);
      return false;
    },
    runningState: function getDetails() {
      if (arguments.length) throw makeError.ErrorInInvocation(`runningState`);
      return 'cannot_run';
    },
  };
  utils.patchToStringNested(window.chrome.app);
};

/**
 * Mock the `chrome.csi` function if not available (e.g. when running headless).
 * It's a deprecated (but unfortunately still existing) chrome specific API to fetch browser timings.
 */
const chromeCsi = (utils) => {
  if ('csi' in window.chrome) return;
  if (!window.performance || !window.performance.timing) return;
  const { timing } = window.performance;
  window.chrome.csi = () => {
    return {
      onloadT: timing.domContentLoadedEventEnd,
      startE: timing.navigationStart,
      pageT: Date.now() - timing.navigationStart,
      tran: 15,
    };
  };
  utils.patchToString(window.chrome.csi);
};

/**
 * Mock the `chrome.loadTimes` function if not available (e.g. when running headless).
 * It's a deprecated (but unfortunately still existing) chrome specific API to fetch browser timings and connection info.
 */
const chromeLoadTimes = (utils) => {
  if ('loadTimes' in window.chrome) return;
  if (!window.performance || !window.performance.timing || !window.PerformancePaintTiming) return;
  const { performance } = window;
  const ntEntryFallback = {
    nextHopProtocol: 'h2',
    type: 'other',
  };
  const protocolInfo = {
    get connectionInfo() {
      const ntEntry = performance.getEntriesByType('navigation')[0] || ntEntryFallback;
      return ntEntry.nextHopProtocol;
    },
    get npnNegotiatedProtocol() {
      const ntEntry = performance.getEntriesByType('navigation')[0] || ntEntryFallback;
      return ['h2', 'hq'].includes(ntEntry.nextHopProtocol) ? ntEntry.nextHopProtocol : 'unknown';
    },
    get navigationType() {
      const ntEntry = performance.getEntriesByType('navigation')[0] || ntEntryFallback;
      return ntEntry.type;
    },
    get wasAlternateProtocolAvailable() {
      return false;
    },
    get wasFetchedViaSpdy() {
      const ntEntry = performance.getEntriesByType('navigation')[0] || ntEntryFallback;
      return ['h2', 'hq'].includes(ntEntry.nextHopProtocol);
    },
    get wasNpnNegotiated() {
      const ntEntry = performance.getEntriesByType('navigation')[0] || ntEntryFallback;
      return ['h2', 'hq'].includes(ntEntry.nextHopProtocol);
    },
  };

  const { timing } = window.performance;

  function toFixed(num, fixed) {
    const re = new RegExp(`^-?\\d+(?:.\\d{0,${fixed || -1}})?`);
    return num.toString().match(re)[0];
  }
  const timingInfo = {
    get firstPaintAfterLoadTime() {
      return 0;
    },
    get requestTime() {
      return timing.navigationStart / 1000;
    },
    get startLoadTime() {
      return timing.navigationStart / 1000;
    },
    get commitLoadTime() {
      return timing.responseStart / 1000;
    },
    get finishDocumentLoadTime() {
      return timing.domContentLoadedEventEnd / 1000;
    },
    get finishLoadTime() {
      return timing.loadEventEnd / 1000;
    },
    get firstPaintTime() {
      const fpEntry = performance.getEntriesByType('paint')[0] || {
        startTime: timing.loadEventEnd / 1000, // Fallback if no navigation occured (`about:blank`)
      };
      return toFixed((fpEntry.startTime + performance.timeOrigin) / 1000, 3);
    },
  };

  window.chrome.loadTimes = () => {
    return {
      ...protocolInfo,
      ...timingInfo,
    };
  };
  utils.patchToString(window.chrome.loadTimes);
};

/**
 * Mock the `chrome.runtime` object if not available (e.g. when running headless) and on a secure site.
 */
const chromeRuntime = (utils) => {
  const existsAlready = 'runtime' in window.chrome;
  if (existsAlready || !window.location.protocol.startsWith('https')) return;

  window.chrome.runtime = {
    ...{
      OnInstalledReason: {
        CHROME_UPDATE: 'chrome_update',
        INSTALL: 'install',
        SHARED_MODULE_UPDATE: 'shared_module_update',
        UPDATE: 'update',
      },
      OnRestartRequiredReason: {
        APP_UPDATE: 'app_update',
        OS_UPDATE: 'os_update',
        PERIODIC: 'periodic',
      },
      PlatformArch: {
        ARM: 'arm',
        ARM64: 'arm64',
        MIPS: 'mips',
        MIPS64: 'mips64',
        X86_32: 'x86-32',
        X86_64: 'x86-64',
      },
      PlatformNaclArch: {
        ARM: 'arm',
        MIPS: 'mips',
        MIPS64: 'mips64',
        X86_32: 'x86-32',
        X86_64: 'x86-64',
      },
      PlatformOs: {
        ANDROID: 'android',
        CROS: 'cros',
        LINUX: 'linux',
        MAC: 'mac',
        OPENBSD: 'openbsd',
        WIN: 'win',
      },
      RequestUpdateCheckStatus: {
        NO_UPDATE: 'no_update',
        THROTTLED: 'throttled',
        UPDATE_AVAILABLE: 'update_available',
      },
    },
    get id() {
      return undefined;
    },
  };
  utils.patchToString(window.chrome.runtime);
};

/**
 * Fix Chromium not reporting "probably" to codecs like `videoEl.canPlayType('video/mp4; codecs="avc1.42E01E"')`.
 * (Chromium doesn't support proprietary codecs, only Chrome does)
 */
const chromeCodec = (utils) => {
  const parseInput = (arg) => {
    const [mime, codecStr] = arg.trim().split(';');
    let codecs = [];
    if (codecStr && codecStr.includes('codecs="')) {
      codecs = codecStr
        .trim()
        .replace(`codecs="`, '')
        .replace(`"`, '')
        .trim()
        .split(',')
        .filter((x) => !!x)
        .map((x) => x.trim());
    }
    return {
      mime,
      codecStr,
      codecs,
    };
  };
  const canPlayType = {
    apply(target, ctx, args) {
      if (!args || !args.length) return target.apply(ctx, args);
      const { mime, codecs } = parseInput(args[0]);
      if (mime === 'video/mp4' && codecs.includes('avc1.42E01E')) return 'probably';
      if (mime === 'audio/x-m4a' && !codecs.length) return 'maybe';
      if (mime === 'audio/aac' && !codecs.length) return 'probably';
      return target.apply(ctx, args);
    },
  };
  utils.replaceWithProxy(HTMLMediaElement.prototype, 'canPlayType', canPlayType);
};

export { chrome, chromeApp, chromeCsi, chromeLoadTimes, chromeRuntime, chromeCodec };
