'use strict';

/* Controllers */

var secFirstControllers = angular.module('secFirstControllers', []);

secFirstControllers.controller('SegmentList', ['$scope', '$routeParams', '$http', '$cookieStore', '$location', 'Segments', 'Categories',
  function($scope, $routeParams, $http, $cookieStore, $location, Segments, Categories) {

    $scope.catId = $routeParams.categoryId;

    Segments.getByCat($routeParams.categoryId).success(function(response) {
      $scope.segments = response;
    });
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

secFirstControllers.controller('CheckItemList', ['$scope', '$routeParams', '$http', '$cookieStore', '$location', 'CheckItems', 'Categories',
  function($scope, $routeParams, $http, $cookieStore, $location, CheckItems, Categories) {

    $scope.catId = $routeParams.categoryId;

    CheckItems.getByCat($routeParams.categoryId).success(function(response){
      $scope.checkItems = response;
    });

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

    $scope.deleteItem = function() {
      CheckItems.deleteCat($scope.toDelete).success(function(response) {
        $scope.showModal = false;
        for (var i = 0; i < $scope.checkItems.length; i++) {
          if ($scope.checkItems[i].id==$scope.toDelete) {
            $scope.checkItems.splice(i, 1);
          }
        }
      });
    };

  }]);

secFirstControllers.controller('CheckItemDetail', ['$scope', '$routeParams', '$http', '$cookieStore', '$location', 'Categories', 'CheckItems',
  function($scope, $routeParams, $http, $cookieStore, $location, Categories, CheckItems) {

    $scope.token = $cookieStore.get('token');
    $scope.action = $routeParams.action;

    $scope.showModal = false;
    $scope.toggleModal = function(toDelete){
      $scope.showModal = !$scope.showModal;
      $scope.toDelete = $scope.showModal?toDelete:0;
    };

    Categories.getAll().success(function(response){
      $scope.categories = Categories.getParentOnly(response);
    });

    CheckItems.getId($routeParams.categoryId).success(function(response){
      $scope.checkItem = response;
    });


    $scope.processForm = function() {
      if ($routeParams.action=='create') {
        CheckItems.create($scope.category).success(function(data) {
            $scope.error = '';
            $location.url('/check_items');
        }).error(function(data, status, header, config){
            $scope.error = data.error;
            $scope.message = "";
        });
      } else if ($routeParams.action=='edit'){
        CheckItems.update($routeParams.categoryId, $scope.category).success(function(data) {
            $scope.error = '';
            $location.url('/check_items');
        }).error(function(data, status, header, config){
            $scope.error = data.error;
            $scope.message = "";
        });
      }
    };

  }]);

secFirstControllers.controller('CategoryList', ['$scope', '$routeParams', '$http', '$cookieStore', '$location', 'Categories',
  function($scope, $routeParams, $http, $cookieStore, $location, Categories) {

    $scope.catId = $routeParams.categoryId;
    $scope.token = $cookieStore.get('token');

    $scope.items = ["One", "Two", "Three"];

    $scope.sortableOptions = {
        update: function(e, ui) { console.log($scope.categories); },
        'ui-floating': true,
        axis: 'y'
      };

    $scope.showModal = false;
    $scope.toggleModal = function(krneki){
      $scope.showModal = !$scope.showModal;
      $scope.toDelete = $scope.showModal?krneki:0;
    };

    Categories.getAll().success(function(response){
      $scope.categories = Categories.getSorted(response);
    });

    $scope.deleteCat = function() {
      Categories.deleteCat($scope.toDelete).success(function(response) {
        $scope.showModal = false;
        for (var i = 0; i < $scope.categories.length; i++) {
          if ($scope.categories[i].id==$scope.toDelete) {
            $scope.categories.splice(i, 1);
          }
        }
      });
    };

  }]);

secFirstControllers.controller('CategoryDetail', ['$scope', '$routeParams', '$http', '$cookieStore', '$location', 'Categories',
  function($scope, $routeParams, $http, $cookieStore, $location, Categories) {

    $scope.backLink = 'categories';
    $scope.backLinkName = 'Categories';

    $scope.token = $cookieStore.get('token');
    if (Object.keys($routeParams).length==1) {
      $scope.action = $routeParams.categoryId;
      if (typeof($scope.action) === "undefined") { $scope.action = ''; }
    } else{
      $scope.action = $routeParams.action;
      $scope.categoryId = $routeParams.categoryId;
    }

    $scope.showModal = false;
    $scope.toggleModal = function(toDelete){
      $scope.showModal = !$scope.showModal;
      $scope.toDelete = $scope.showModal?toDelete:0;
    };

    Categories.getAll().success(function(response){
      $scope.categories = Categories.getParentOnly(response);
    });

    if (typeof($scope.categoryId)!=='undefined') {
      Categories.getId($routeParams.categoryId).success(function(response){
        if ($scope.action==='create') {
          $scope.backLink = 'segments/'+$routeParams.categoryId+'/category';
          $scope.backLinkName = 'Segments';
          var category = {};
          category.parent = response.id;
          $scope.category = category;
        } else{
          $scope.category = response;
          if ($scope.category.parent!==0) {
            for (var i = 0; i < $scope.categories.length; i++) {
              if ($scope.categories[i].id==$scope.category.parent) {
                $scope.category.parentName = $scope.categories[i].category;
              }
            }
          }
        }
      });
    }

    $scope.processForm = function() {
      if ($routeParams.action=='create') {
        Categories.create($scope.category).success(function(data) {
            $scope.error = '';
            $location.url('/categories');
        }).error(function(data, status, header, config){
            $scope.error = data.error;
            $scope.message = "";
        });
      } else if ($routeParams.action=='edit'){
        Categories.update($routeParams.categoryId, $scope.category).success(function(data) {
            $scope.error = '';
            $location.url('/categories');
        }).error(function(data, status, header, config){
            $scope.error = data.error;
            $scope.message = "";
        });
      }
    };

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
          $cookieStore.put('loginProfile', data.profile);
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

    $scope.$watch(function() { return $cookieStore.get('token');}, function(newValue) {
        $scope.loginProfile = $cookieStore.get('loginProfile');
        $scope.token = $cookieStore.get('token');
    });
  }]);