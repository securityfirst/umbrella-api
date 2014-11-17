'use strict';

/* App Module */

var phonecatApp = angular.module('phonecatApp', [
  'ngRoute',
  'secFirstAnimations',

  'secFirstControllers',
  'secFirstFilters',
  'secFirstServices'
]);

phonecatApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/segments', {
        templateUrl: 'partials/segments.html',
        controller: 'SegmentList'
      }).
      when('/segments/:segmentId', {
        templateUrl: 'partials/segment-detail.html',
        controller: 'SegmentDetail'
      }).
      otherwise({
        redirectTo: '/segments'
      });
  }]);
