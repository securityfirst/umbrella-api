'use strict';

/* Controllers */

var secFirstControllers = angular.module('secFirstControllers', []);

secFirstControllers.controller('SegmentList', ['$scope', '$routeParams', '$http', '$cookieStore', '$location', 'Segments', 'Categories',
  function($scope, $routeParams, $http, $cookieStore, $location, Segments, Categories) {

    $scope.catId = $routeParams.categoryId;

    Segments.getByCat($routeParams.categoryId).success(function(response) {
      $scope.segments = response;
    });
    $scope.orderProp = 'age';
    $scope.token = $cookieStore.get('token');

    Categories.getAll().success(function(response){
      var sortedCats = [];
      for (var i = 0; i < response.length; i++) {
        if (response[i].parent===0) {
          sortedCats.push(response[i]);
          for (var j = 0; j < response.length; j++) {
            if (response[i].id === response[j].parent) {
              sortedCats.push(response[j]);
            }
          }
        }
      }
      $scope.categories = sortedCats;
    });

  }]);

secFirstControllers.controller('SegmentDetail', ['$scope', '$routeParams', '$cookieStore', '$location', 'Segments',
  function($scope, $routeParams, $cookieStore, $location, Segments) {
    Segments.getByCat($routeParams.segmentId).success(function(response){
      $scope.segment = response[0];
    });
    $scope.token = $cookieStore.get('token');

    $scope.options = {
      toolbar: [
        ['style', ['bold', 'italic', 'underline', 'clear']]
      ]
    };

    $scope.processForm = function() {
      Segments.updateCat($routeParams.segmentId, $scope.segment).success(function(data) {
          $scope.error = '';
          $location.url('/segments');
      }).error(function(data, status, header, config){
          $scope.error = data.error;
          $scope.message = "";
      });
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

secFirstControllers.controller('LogOut', ['$scope', '$cookieStore', '$location',
  function($scope, $cookieStore, $location) {
    $cookieStore.put('token', '');
    $location.url('/segments');
  }]);

secFirstControllers.controller('TopNav', ['$scope', '$cookieStore', '$location',
  function($scope, $cookieStore, $location) {

    $scope.token = $cookieStore.get('token');

  }]);