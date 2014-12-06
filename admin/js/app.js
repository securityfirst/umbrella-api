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
  'summernote',
  'ui.sortable'
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
      when('/check_items', {
        templateUrl: 'partials/check_items.html',
        controller: 'CheckItemList'
      }).
      when('/check_items/:categoryId/category', {
        templateUrl: 'partials/check_items.html',
        controller: 'CheckItemList'
      }).
      when('/check_items/:categoryId/:action', {
        templateUrl: 'partials/check_item_detail.html',
        controller: 'CheckItemDetail'
      }).
      when('/check_items/:action', {
        templateUrl: 'partials/check_item_detail.html',
        controller: 'CheckItemDetail'
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
      when('/categories/:categoryId/:action/:parentId', {
        templateUrl: 'partials/category_detail.html',
        controller: 'CategoryDetail'
      }).
      when('/categories/:categoryId', {
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
        templateUrl: 'partials/segment_detail.html',
        controller: 'SegmentDetail'
      }).
      otherwise({
        redirectTo: '/segments'
      });
  }]);
