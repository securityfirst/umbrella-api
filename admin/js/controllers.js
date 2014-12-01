'use strict';

/* Controllers */

var secFirstControllers = angular.module('secFirstControllers', []);

secFirstControllers.controller('SegmentList', ['$scope', '$http', '$cookieStore', 'Segments',
  function($scope, $http, $cookieStore, Segments) {

    Segments.getRaw().success(function(response) {
      $scope.segments = response;
    });
    $scope.orderProp = 'age';
    $scope.token = $cookieStore.get('token');
  }]);

secFirstControllers.controller('SegmentDetail', ['$scope', '$routeParams', '$cookieStore', 'Segments',
  function($scope, $routeParams, $cookieStore, Segments) {
    Segments.getByCat($routeParams.segmentId).success(function(response){
      $scope.segment = response[0];
    });
    $scope.token = $cookieStore.get('token');

    $scope.options = {
        // height: 300,
        focus: true,
        toolbar: [
          ['style', ['bold', 'italic', 'underline', 'clear']]
        ]
      };

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