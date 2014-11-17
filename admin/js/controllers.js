'use strict';

/* Controllers */

var phonecatControllers = angular.module('secFirstControllers', []);

phonecatControllers.controller('SegmentList', ['$scope', '$http', 'Segment',
  function($scope, $http, Segment) {

    Segment.getAll().success(function(response) {
      $scope.segments = response;
    });
    $scope.orderProp = 'age';
  }]);

phonecatControllers.controller('SegmentDetail', ['$scope', '$routeParams', 'Segment',
  function($scope, $routeParams, Segment) {
    Segment.getId($routeParams.segmentId).success(function(response){
      $scope.segment = response;
    });
  }]);
