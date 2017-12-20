var cacheName = 'shell-content'
var filesToCache = [
  '.',
  'index.html',
  'sw.js',
  'StudentVIEW.js',
  'login.html',
  'loginstyle.css',
  'styles.css',
  'material.cyan-light_blue.min.css',
  'https://fonts.googleapis.com/icon?family=Material+Icons',
  'https://fonts.googleapis.com/css?family=Roboto:regular,bold,italic,thin,light,bolditalic,black,medium&lang=en',
  'icons/apple-touch-icon.png',
  'icons/favicon-32x32.png',
  'icons/favicon-194x194.png',
  'icons/android-chrome-192x192.png',
  'icons/favicon-16x16.png',
  'icons/safari-pinned-tab.svg',
  'icons/favicon.ico',
  'icons/mstile-144x144.png',
  'icons/browserconfig.xml',
  'https://cdn.rawgit.com/HubSpot/pace/master/themes/purple/pace-theme-fill-left.css'
]

self.addEventListener('install', function (event) {
  console.log('Attempting to install service worker and cache static assets');
  event.waitUntil(
    caches.open(cacheName)
    .then(function(cache) {
      return cache.addAll(filesToCache);
    })
  );
});

self.addEventListener('fetch', function(event) {
  event.respondWith(
    fetch(event.request).catch(function() {
      return caches.match(event.request);
    })
  );
});
