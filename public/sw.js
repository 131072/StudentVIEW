/** An empty service worker! */
self.addEventListener('fetch', function(event) {
  /** An empty fetch handler! */
});

var cacheName = 'shell-content';
var filesToCache = [
  '/studentview/index.html',
  '/studentview/StudentVIEW.js',
  '/studentview/login.html',
  '/studentview/loginstyle.css',
  '/studentview/main.js',
  '/studentview/styles.css',
  '/studentview/material.cyan-light_blue.min.css',
  'https://fonts.googleapis.com/icon?family=Material+Icons',
  'https://fonts.googleapis.com/css?family=Roboto:regular,bold,italic,thin,light,bolditalic,black,medium&amp;lang=en',
  '/studentview/icons/apple-touch-icon.png',
  '/studentview/icons/favicon-32x32.png',
  '/studentview/icons/favicon-194x194.png',
  '/studentview/icons/android-chrome-192x192.png',
  '/studentview/icons/favicon-16x16.png',
  '/studentview/icons/safari-pinned-tab.svg',
  '/studentview/icons/favicon.ico',
  '/studentview/icons/mstile-144x144.png',
  '/studentview/icons/browserconfig.xml',
  '/studentview/material.min.js',
  'https://cdn.rawgit.com/kimmobrunfeldt/progressbar.js/master/dist/progressbar.min.js',
  'https://cdn.rawgit.com/HubSpot/pace/master/pace.min.js',
  'https://cdn.rawgit.com/HubSpot/pace/master/themes/purple/pace-theme-fill-left.css'
];

self.addEventListener('install', function(e) {
  console.log('[ServiceWorker] Install');
  e.waitUntil(
    caches.open(cacheName).then(function(cache) {
      console.log('[ServiceWorker] Caching app shell');
      return cache.addAll(filesToCache);
    })
  );
});
