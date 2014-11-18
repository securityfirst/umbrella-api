'use strict';

/* Services */

var secFirstServices = angular.module('secFirstServices', []);

secFirstServices.factory('Segment', ['$http', '$cookieStore',
  function($http, $cookieStore){
    var factory = {};

    factory.getAll = function() {
        return $http.get('/v1/segments', {
            headers: {'token': $cookieStore.get('token')}
        });
    };

    factory.getId = function(segmentId) {
        console.log(segmentId);
        return $http.get("/v1/segments/"+segmentId);
    };
    return factory;
  }]);

secFirstServices.factory('Login', ['$http', '$cookieStore',
  function($http, $cookieStore){
    var factory = {};

    factory.postForm = function(formData) {
        return $http({
          method  : 'POST',
          url     : '/v1/account/login',
          headers: {'token': $cookieStore.get('token')},
          data    : formData
      });
    };
    return factory;
  }]);