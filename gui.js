angular
    .module('App', ['ngMaterial'])
    .config(function($mdThemingProvider) {
        $mdThemingProvider
            .theme('default')
            .primaryPalette('blue')
            .accentPalette('pink');
    })
    .controller('AppCtrl', function($scope, $mdSidenav) {
        $scope.toggleSidenav = function(menuId) {
            $mdSidenav(menuId).toggle();
        };
    });