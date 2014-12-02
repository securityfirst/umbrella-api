'use strict';

/* App Module */

var secFirstApp = angular.module('secFirstApp', [
  'ngRoute',
  'ngCookies',
  'secFirstAnimations',
  'secFirstControllers',
  'secFirstFilters',
  'secFirstServices',
  'secFirstDirectives',
  'summernote'
]).run(function($rootScope, $location) {
    $rootScope.location = $location;
});

secFirstApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/segments', {
        templateUrl: 'partials/segments.html',
        controller: 'SegmentList'
      }).
      when('/segments/:categoryId/category', {
        templateUrl: 'partials/segments.html',
        controller: 'SegmentList'
      }).
      when('/login', {
        templateUrl: 'partials/login.html',
        controller: 'LoginForm'
      }).
      when('/categories', {
        templateUrl: 'partials/categories.html',
        controller: 'CategoryList'
      }).
      when('/categories/:categoryId/:action', {
        templateUrl: 'partials/category_detail.html',
        controller: 'CategoryDetail'
      }).
      when('/categories/:action', {
        templateUrl: 'partials/category_detail.html',
        controller: 'CategoryDetail'
      }).
      when('/logout', {
        templateUrl: 'partials/about.html',
        controller: 'LogOut'
      }).
      when('/about', {
        templateUrl: 'partials/about.html'
      }).
      when('/segments/:segmentId', {
        templateUrl: 'partials/segment-detail.html',
        controller: 'SegmentDetail'
      }).
      otherwise({
        redirectTo: '/segments'
      });
  }]);
