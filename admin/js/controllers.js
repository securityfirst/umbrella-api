'use strict';

/* Controllers */

var secFirstControllers = angular.module('secFirstControllers', []);

secFirstControllers.controller('SegmentList', ['$scope', '$http', 'Segment',
  function($scope, $http, Segment) {

    Segment.getAll().success(function(response) {
      $scope.segments = response;
    });
    $scope.orderProp = 'age';
  }]);

secFirstControllers.controller('SegmentDetail', ['$scope', '$routeParams', 'Segment',
  function($scope, $routeParams, Segment) {
    Segment.getId($routeParams.segmentId).success(function(response){
      $scope.segment = response;
    });
  }]);

secFirstControllers.controller('LoginForm', ['$scope', '$http', '$cookieStore', '$location', 'Login',
  function($scope, $http, $cookieStore, $location, Login) {
    $scope.formData = {};

    $scope.processForm = function() {
      Login.postForm($scope.formData).success(function(data) {
          $scope.message = data.token;
          $scope.data = null;
          $cookieStore.put('token', data.token);
          $scope.token = data.token;
          $location.url('/segments');
          $scope.error = "";
      }).error(function(data, status, header, config){
          $scope.error = data.error;
          $scope.message = "";
      });
    };
  }]);