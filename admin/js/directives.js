'use strict';

/* Directives */

var secFirstDirectives = angular.module('secFirstDirectives', []);

angular.module('secFirstDirectives').directive('modal', function () {
    return {
      template: '<div class="modal fade">' +
          '<div class="modal-dialog">' +
            '<div class="modal-content">' +
              '<div class="modal-header">' +
                '<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>' +
                '<h4 class="modal-title">{{ title }}</h4>' +
              '</div>' +
              '<div class="modal-body" ng-transclude></div>' +
            '</div>' +
          '</div>' +
        '</div>',
      restrict: 'E',
      transclude: true,
      replace:true,
      scope:true,
      link: function postLink(scope, element, attrs) {
        scope.title = attrs.title;

        scope.$watch(attrs.visible, function(value){
          if(value == true)
            $(element).modal('show');
          else
            $(element).modal('hide');
        });

        $(element).on('shown.bs.modal', function(){
          scope.$apply(function(){
            scope.$parent[attrs.visible] = true;
          });
        });

        $(element).on('hidden.bs.modal', function(){
          scope.$apply(function(){
            scope.$parent[attrs.visible] = false;
          });
        });
      }
    };
  });

angular.module("secFirstDirectives").directive("navigation", [
    "$location", function($location) {
      return {
        restrict: 'A',
        scope: {},
        link: function(scope, element) {
          var classSelected, navLinks;

          scope.location = $location;

          classSelected = 'active';

          navLinks = element.find('a');

          scope.$watch('location.path()', function(newPath) {
            var segments =  newPath.split('/');
            if (segments.length>1) {newPath = '/'+segments[1];}
            var el;
            el = navLinks.filter('[href="#' + newPath + '"]');

            navLinks.not(el).closest('li').removeClass(classSelected);
            return el.closest('li').addClass(classSelected);
          });
        }
      };
    }
  ]);