import utils from './lib/utils';

import { windowOuterDimensions, windowNotifications } from './modules/window';

import { chrome, chromeApp, chromeCsi, chromeLoadTimes, chromeRuntime, chromeCodec } from './modules/chrome';

import {
  navigatorHardwareConcurrency,
  navigatorLanguages,
  navigatorPlatform,
  navigatorUserAgent,
  navigatorVendor,
  navigatorWebDriver,
  navigatorDeviceMemory,
} from './modules/navigator';

import { injectCanva } from './modules/canva';
import { addIframeCreationSniffer } from './modules/iframe';
import { navigatorPlugins } from './modules/plugins';
import { webglVendor } from './modules/webgl';

utils.preloadCache();

// window
windowOuterDimensions();
windowNotifications();

// chrome
chrome();
chromeApp(utils);
chromeCsi(utils);
chromeLoadTimes(utils);
chromeRuntime(utils);
chromeCodec(utils);

// navigator
navigatorHardwareConcurrency();
navigatorLanguages();
navigatorPlatform();
navigatorUserAgent();
navigatorVendor();
navigatorWebDriver();
navigatorDeviceMemory();

// hard
injectCanva();
addIframeCreationSniffer(utils);
navigatorPlugins(utils);
webglVendor(utils);
