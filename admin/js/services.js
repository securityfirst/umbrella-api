'use strict';

/* Services */

var phonecatServices = angular.module('secFirstServices', []);

phonecatServices.factory('Segment', ['$http',
  function($http){
    var factory = {};

    factory.getAll = function() {
        return $http.get("/v1/segments");
    };

    factory.getId = function(segmentId) {
        console.log(segmentId);
        return $http.get("http://localhost:3000/v1/segments/"+segmentId);
    };
    return factory;
  }]);