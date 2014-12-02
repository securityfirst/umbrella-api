'use strict';

/* Services */

var secFirstServices = angular.module('secFirstServices', []);

secFirstServices.factory('Segments', ['$http', '$cookieStore',
  function($http, $cookieStore){
    var factory = {};

    factory.getAll = function() {
        return $http.get('/v1/segments_raw', {
            headers: {'token': $cookieStore.get('token')}
        });
    };

    factory.getRaw = function() {
        return $http.get('/v1/segments', {
            headers: {'token': $cookieStore.get('token')}
        });
    };

    factory.getByCat = function(categoryId) {
        if (typeof categoryId === "undefined") {
          return $http.get('/v1/segments', {
              headers: {'token': $cookieStore.get('token')}
          });
        } else{
          return $http.get("/v1/segments/"+categoryId+"/category");
        }
    };

    factory.updateCat = function(categoryId, formData) {
      return $http({method: 'PUT', url: '/v1/segments/'+categoryId+'/category/', headers: {'token': $cookieStore.get('token')}, data: formData
      });
    };

    factory.getId = function(segmentId) {
        return $http.get("/v1/segments_raw/"+segmentId);
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

secFirstServices.factory('Categories', ['$http', '$cookieStore',
  function($http, $cookieStore){
    var factory = {};

    factory.getAll = function() {
        return $http.get('/v1/categories', {
            headers: {'token': $cookieStore.get('token')}
        });
    };

    factory.getId = function(categoryId) {
        return $http.get("/v1/categories/"+categoryId);
    };

    factory.deleteCat = function(deleteId) {
        return $http.delete('/v1/categories/'+deleteId, {
            headers: {'token': $cookieStore.get('token')}
        });
    };

    factory.getSorted = function(response) {
      var sortedCats = [];
      for (var i = 0; i < response.length; i++) {
        if (response[i].parent===0) {
          response[i].parentName = '';
          sortedCats.push(response[i]);
          for (var j = 0; j < response.length; j++) {
            if (response[i].id === response[j].parent) {
              response[j].parentName = response[i].category;
              sortedCats.push(response[j]);
            }
          }
        }
      }
      return sortedCats;
    };

    factory.getParentOnly = function(response) {
      var parentOnly = [];
      for (var i = 0; i < response.length; i++) {
        if (response[i].parent===0) {parentOnly.push(response[i]);}
      }
      return parentOnly;
    };

    factory.update = function(categoryId, formData) {
      return $http({method: 'PUT', url: '/v1/categories/'+categoryId, headers: {'token': $cookieStore.get('token')}, data: formData
      });
    };

    factory.create = function(formData) {
      return $http({method: 'POST', url: '/v1/categories/', headers: {'token': $cookieStore.get('token')}, data: formData
      });
    };

    return factory;
  }]);